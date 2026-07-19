# Predicate type: Adversarial Execution Evidence

Type URI: https://in-toto.io/attestation/adversarial-execution-evidence/v0.4

Version: 0.4.0

Predicate Name: Adversarial Execution Evidence

> Status: DRAFT submission for vetting. The schema below matches a production
> implementation; field-level feedback is the point of this PR, and naming or
> structural changes from review are expected before vetting completes.

## Purpose

Records the evidence produced by deliberately executing an untrusted or
under-trusted software artifact (an agent tool, an MCP server, a plugin, a
build step) inside an instrumented containment substrate while a known corpus
of adversarial inputs is thrown at it. The predicate carries what was thrown,
what the substrate was configured to catch, what was actually intercepted
(each interception as an independently signed record), and an explicit
statement of the observation's coverage bounds.

The design goal is that a consumer can recompute the outcome from the
attestation plus the producer's published taxonomy, with no call back to the
producer's infrastructure. The outcome is a deterministic function of the
carried evidence; the coverage denominator is committed by digest, so it is
not something the producer can quietly assert; and each intercept record
verifies on its own before it is read. A producer cannot claim more than the
evidence supports, and a producer claiming less is detectable, since dropping
an inconvenient interception changes the committed batch root.

## Use Cases

-   An admission controller (e.g. a Kubernetes policy engine) gating a
    third-party MCP server or agent tool image on evidence that it was
    detonated against a named attack corpus under an enforcing catch policy,
    with the policy digest and network posture pinned in the evidence.
-   An auditor re-verifying, offline and without trusting the producer's
    infrastructure, that a specific interception happened: the signed
    intercept record binds the destination, the payload commitment, and the
    substrate context.
-   A security team comparing two runs of the same artifact: because the
    corpus manifest is digest-committed at attack granularity, "both runs
    assessed the same attacks" is checkable, not asserted.

Existing predicates cover adjacent but different ground. [Runtime Traces]
carries raw observed activity from a monitor, with no corpus binding, no
coverage denominator, and no per-event signature. [SCAI] carries
evidence-backed attribute assertions but does not model an adversarial corpus
or recomputable outcomes. [VSA] and [SVR] carry policy verdicts computed at
verification time, downstream of evidence like this. [Test Result] carries
test outcomes without cryptographic binding of the inputs or the
interceptions. This predicate is the evidence layer those verdict predicates
can consume.

## Prerequisites

The in-toto Attestation Framework, plus an understanding of
[DSSE](https://github.com/secure-systems-lab/dsse) (each intercept record is a
DSSE-shaped envelope) and [RFC 8785 (JCS)](https://www.rfc-editor.org/rfc/rfc8785)
canonical JSON, which the digest bindings are defined over. Producers MUST
enforce the RFC 7493 (I-JSON) safe-integer profile on canonicalized content:
integers with magnitude at or above 2^53 MUST be rejected, so every rail
(producer and verifier, in any language) derives identical bytes.

## Model

The producer is a containment substrate operator: a functionary that runs the
subject artifact inside an isolated, instrumented environment (a microVM, a
sandbox, an eBPF-supervised process), injects the corpus, and signs what the
substrate intercepted. The subject is the executed artifact, by digest. Every
attestation references a substrate by subject (the `substrate` field is
required); that substrate SHOULD in turn carry its own attestation, e.g. build
provenance for the substrate image, so the evidence can inherit a substrate
trust chain rather than a bare name. Verdicts (pass/fail against an
organization's policy) are deliberately out of scope; they belong in a
downstream summary predicate such as [VSA], computed over this evidence.

## Schema

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    { "name": "<executed-artifact-name>", "digest": { "sha256": "<64-hex>" } }
  ],
  "predicateType": "https://in-toto.io/attestation/adversarial-execution-evidence/v0.4",
  "predicate": {
    "result": "fail",
    "observationEnvironment": {
      "substrate": { "name": "<substrate-attestation-subject>", "digest": { "sha256": "<64-hex>" } },
      "corpus": {
        "name": "<corpus-name>",
        "uri": "pkg:<producer>/<corpus>@<version>",
        "digest": { "sha256": "<64-hex-JCS-digest-of-manifest>" },
        "manifest": { "classes": { "CO": ["CO-EXFIL-1"] } }
      },
      "catchPolicy": { "digest": { "sha256": "<64-hex-JCS-digest-of-catch-policy>" } },
      "networkPosture": { "posture": "sinkhole", "digest": { "sha256": "<64-hex>" } }
    },
    "coverage": {
      "assessedClasses": ["CO"],
      "outOfScope": {},
      "routedElsewhere": {}
    },
    "attackResults": [
      {
        "attackId": "CO-EXFIL-1",
        "containmentObserved": "egress_captured",
        "basis": "substrate_observed",
        "actualLayer": "policy.egress_sinkhole",
        "interceptRefs": [0]
      }
    ],
    "interceptRecords": [
      {
        "payload": "<base64(canonical bytes of the signed intercept record)>",
        "payloadType": "<producer-defined media type>",
        "signatures": [ { "keyid": "<hex>", "sig": "<base64>" } ]
      }
    ],
    "batchRoot": "<64-hex-merkle-root-over-intercept-records>",
    "doesNotAssert": [ "<explicit negative-scope statements>" ],
    "issuedAt": "2026-06-23T16:08:07Z"
  }
}
```

### Parsing Rules

The predicate opts in to the framework's standard parsing rules, including the
monotonic principle, with one deliberate strengthening: `result` is not an
independent claim. A consumer MUST be able to recompute it from the rest of
the predicate (rules under `result` below), and a `result` the recompute does
not reproduce makes the attestation invalid. Intercept record `payload`s
follow the same verify-then-read discipline: the fields inside a payload mean
nothing until its signature verifies against a key the consumer trusts.

### Fields

`result` _string, required_

One of `pass`, `degraded`, `fail` (lowercase). Defined as a total,
deterministic, severity-independent function of the predicate: `fail` iff any
`attackResults` row carries a containment-observed label from the published
caught set, an out-of-vocabulary label (fail-closed), or a missing or
out-of-vocabulary `basis` (fail-closed, same rule); otherwise `degraded`
iff `coverage.outOfScope` or `coverage.routedElsewhere` is non-empty;
otherwise `pass`. A `pass` is coverage-bounded-observed: a statement about
what was assessed, never a general safety claim. There is intentionally no
severity threshold, policy ruleset, or free-text reason here; policy belongs
downstream.

`observationEnvironment` _object, required_

The digest-pinned context the evidence was earned under. Four required
members: `substrate` (the subject reference of the substrate's own
attestation), `corpus` (name, uri, RECOMMENDED as a purl; the JCS digest of
the embedded `manifest`, and the `manifest` itself: a map from assessment
class to the complete array of attack identifiers it defines; an attackId
MUST NOT appear under more than one class), `catchPolicy` (JCS digest of the
parsed catch-policy document, so an empty or permissive policy is
distinguishable from an enforcing one), and `networkPosture` (the
substrate-authoritative egress posture, e.g. `no_network`, `allowlist`,
`sinkhole`, with its configuration digest). The manifest pre-image travels in
the attestation, so a verifier re-derives `corpus.digest` offline and any edit
to the assessed set (a dropped attack, a renamed class) fails that check.

`coverage` _object, required_

The coverage bound: `assessedClasses` (array of class codes actually
assessed), `outOfScope` and `routedElsewhere` (maps from class code to a
reason string; empty objects when complete). Disclosing a gap moves the class
into one of these maps and forces `result` to `degraded`, which is the honest
alternative to leaving it out and quietly reporting a narrower run as a full
one.

`attackResults` _array of objects, required_

One row per executed attack: `attackId` (must appear in the manifest),
`containmentObserved` (a label from the producer's published, versioned
observation vocabulary; consumers treat unknown labels as fail-closed),
`basis` (required; the observation's vantage, see below), `actualLayer`
(which enforcement layer acted, see below), and `interceptRefs` (indexes
into `interceptRecords` binding this attack to its signed interceptions).
Coverage integrity is checked at attack granularity: the union of attackIds
for the assessed classes must exactly equal the manifest's. That granularity
is what stops a failing attack from being quietly omitted inside a class the
producer still reports as assessed.

`basis` states where each observation came from, with a closed three-value
vocabulary:

-   `substrate_observed`: the substrate itself observed the event at its own
    vantage (network boundary, syscall supervision, VM introspection),
    independent of the executed artifact's cooperation.
-   `artifact_reported`: the observation derives from output the executed
    artifact itself produced (its stdout/stderr, exit status, or self-emitted
    logs).
-   `inferred`: the observation was derived indirectly from neither of the
    above, e.g. a post-hoc state diff.

`basis` is REQUIRED on every row and the vocabulary is closed: a missing
`basis`, or any value outside these three, is fail-closed exactly as an
out-of-vocabulary `containmentObserved` label is — the row forces `result`
to `fail` and can support nothing stronger. The rationale is the same as for
`result` itself: the recompute reads the row, so the row must carry its own
vantage; `observationEnvironment` pins what the substrate was configured to
do per run, but basis is a per-observation property, and without it a
substrate-intercepted `egress_captured` and one transcribed from the
artifact's own stderr recompute identically. Consumers MUST be able to gate
on `basis`, and a `fail` supported only by `artifact_reported` rows SHOULD
be treated as a weaker claim than one carrying a `substrate_observed` row: a
consumer MAY reject it, since a self-reported observation that lands on
`fail` through a correct recompute would otherwise read downstream as
substrate-checked.

`actualLayer` names the enforcement layer that acted on the row's
containment event. On a row whose `containmentObserved` label is from the
published vocabulary but not in the caught set (a clean row: nothing acted),
the producer MUST emit the literal string `none`. `none` is explicit rather than the
field being omitted so that "the substrate was positioned and no layer
needed to act" is distinguishable from an accidental omission; whether
anything was positioned to see on that clean row is answered by `basis` — a
clean row carrying `basis: substrate_observed` states the substrate had
vantage and observed no containment event, which is the claim a `pass`
actually rests on.

`interceptRecords` _array of objects, optional_

One DSSE-shaped envelope per interception: `payload` (base64 of the exact
canonical bytes the substrate signed at interception time), `payloadType`,
and `signatures`. A consumer verifies each record's signature before reading
any field inside the payload; the record's content (timestamps, destination,
payload commitment digests) is producer-defined and documented with the
producer's vocabulary. What a record carries is a commitment to an intercepted
payload rather than the payload itself, which keeps the attestation
publishable instead of turning it into a sensitive-data store.

`batchRoot` _string, required when `interceptRecords` is non-empty_

64-hex Merkle root over the intercept records, carried once at the predicate
level. Binds the record SET: dropping or reordering a record changes the
root. Omitted on a clean run with no interceptions.

`doesNotAssert` _array of strings, optional_

Explicit negative scope: statements the producer declares this evidence makes
no claim about (e.g. behavior outside the thrown corpus, host integrity
beyond the substrate's own attestation). Advisory: a verifier MUST NOT
require it, and nothing in it weakens the required checks. `doesNotAssert`
is the single canonical spelling: earlier internal versions used snake_case
`does_not_assert`, and that spelling is not accepted as an alias, since two
accepted spellings would mean two canonicalizations for the same content.
Migrating old producer output to the new name is a producer concern, not
something the wire format carries.

`issuedAt` _string (RFC 3339 timestamp), required_

When the producer signed the evidence bundle.

## Example

The schema block above is a complete statement whose `corpus.digest` is
re-derivable from the embedded manifest (`{"classes":{"CO":["CO-EXFIL-1"]}}`
canonicalized under RFC 8785 and hashed); the remaining digests are
placeholders.

## Changelog and Migrations

Versions 0.1–0.2 were internal producer iterations; 0.3 was the first shape
proposed for vetting. Relative to those internal versions, 0.3 removed all
verdict/policy semantics (moved downstream), moved intercepted-payload bytes
out in favor of commitments, moved `batchRoot` from per-record to
predicate-level, and adopted the I-JSON safe-integer profile on every rail.

0.4 incorporates review feedback on 0.3:

-   Added the required per-row `basis` field on `attackResults` (closed
    vocabulary `substrate_observed` / `artifact_reported` / `inferred`,
    fail-closed on unknown values), so each observation carries its own
    vantage and consumers can gate on it.
-   Pinned `actualLayer` clean-run behavior: rows with no containment event
    carry the literal `none` rather than omitting the field.
-   Renamed `does_not_assert` to `doesNotAssert` to match the lowerCamelCase
    convention. The rename is in place with no alias: the old spelling is
    rejected, keeping a single canonicalization per content.

[Runtime Traces]: runtime-trace.md
[SCAI]: scai.md
[SVR]: svr.md
[Test Result]: test-result.md
[VSA]: vsa.md
