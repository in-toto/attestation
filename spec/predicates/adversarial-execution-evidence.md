# Predicate type: Adversarial Execution Evidence

Type URI: https://in-toto.io/attestation/adversarial-execution-evidence/v0.5

Version: 0.5.0

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
  "predicateType": "https://in-toto.io/attestation/adversarial-execution-evidence/v0.5",
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
        "basis": "substrate",
        "method": "intercepted",
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

A design invariant follows from the recompute: any per-observation property
that the recompute or the documented consumer gating reads travels on the
row itself, as a required member with a closed vocabulary, fail-closed on
missing or unknown values. Run-level pins in `observationEnvironment` never
substitute for a row-level property, because the recompute reads rows.

### Fields

`result` _string, required_

One of `pass`, `degraded`, `fail` (lowercase). Defined as a total,
deterministic, severity-independent function of the predicate: `fail` iff any
`attackResults` row carries a containment-observed label from the published
caught set, an out-of-vocabulary label (fail-closed), or a missing or
out-of-vocabulary `basis` or `method` (fail-closed, same rule); otherwise `degraded`
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
`basis` (required; the observation's vantage, see below), `method`
(required; the observation's directness, see below), `actualLayer`
(required; which enforcement layer acted, see below), and `interceptRefs`
(indexes into `interceptRecords` binding this attack to its signed
interceptions).
Coverage integrity is checked at attack granularity: the union of attackIds
for the assessed classes must exactly equal the manifest's. That granularity
is what stops a failing attack from being quietly omitted inside a class the
producer still reports as assessed.

`basis` states each observation's vantage, with a closed two-value
vocabulary:

-   `substrate`: every input the row's claim depends on was obtained at a
    vantage the executed artifact could neither forge nor suppress (a
    network boundary, syscall supervision, a hypervisor's read of guest
    state). This names a class of vantage, cooperation-independence, not
    the identity of the enforcing substrate: a passive tap, an inline gate,
    and an adversarial corpus endpoint logging the connection it received
    all sit at this basis and differ only in enforcement role, which is
    `actualLayer`'s question.
-   `artifact`: at least one input the claim depends on derives from output
    the executed artifact itself produced (its stdout/stderr, exit status,
    or self-emitted logs).

An input is artifact-sourced when the claim relies on it as a channel
whose content the artifact can populate arbitrarily without performing
the claimed containment event: testimony about an event rather than the
event itself. An egress capture is not artifact-sourced even though the
packet bytes were artifact-authored, because the artifact cannot cause
the boundary to record an egress without performing one; its stdout is
artifact-sourced because the artifact can write anything there at no
cost.

`basis` is the vantage of the claim's weakest input. A derived observation
inherits `artifact` from any artifact-sourced input it consumed: a state
diff computed by substrate machinery over the artifact's own logs is
`artifact`, however trusted the machinery, because the artifact could have
populated what the machinery read without performing any containment
event.

`method` states how the row's claim was established, with a closed
two-value vocabulary:

-   `intercepted`: the claim rests on events captured as they occurred. On
    a clean row, a live capture vantage was armed for the attack and no
    capture was attributed to it.
-   `reconstructed`: the claim derives from state examined after the fact,
    e.g. a snapshot-to-snapshot diff. A reconstruction can miss a transient
    raised and undone between the states it compares; on a clean row that
    tolerance is part of the claim.

Like `basis`, `method` composes by weakest input: a claim inherits
`reconstructed` from any state-derived input it depends on. Post-hoc
decode of an event stream captured as it occurred (a packet capture
parsed later, a hardware trace decoded after the run) does not demote a
row, provided the capture channel was armed for the claimed event class
before the event; a row that fuses a live capture with after-the-fact
state examination is `reconstructed`.

`method` describes how the observation was made, not how the row was
attributed to its attack: attribution strength (a hash-pinned payload
versus a time-window pairing) stays in the producer's
`containmentObserved` vocabulary, and so does attribution tolerance on
clean rows, including window bleed between same-layer siblings.

Both fields are REQUIRED on every row and both vocabularies are closed: a
missing value, or any value outside them, is fail-closed exactly as an
out-of-vocabulary `containmentObserved` label is — the row forces
`result` to `fail` and can support nothing stronger. The 0.4 values
`substrate_observed`, `artifact_reported`, and `inferred` are
out-of-vocabulary in 0.5, with no alias, for the same
single-canonicalization reason the old `does_not_assert` spelling is
rejected. `inferred` has no successor value because it conflated the two
axes: a hypervisor snapshot diff and a diff parsed from the artifact's own
logs were both "derived indirectly" while carrying opposite vantages.
Under 0.5 the first is (`substrate`, `reconstructed`), the second is
(`artifact`, `reconstructed`), and no row can read as more independent
than its weakest input.

The two axes bind a consumer's confidence on opposite sides. A `fail`'s
supporting rows are the rows that force `result` to `fail`: rows whose
`containmentObserved` label is in the published caught set, plus rows
fail-closed on a missing or out-of-vocabulary member. `basis` bounds a
`fail` over that supporting set: it answers whether the artifact could
have manufactured the observation. A `fail` whose supporting rows are all
`artifact` SHOULD be treated as a weaker claim than one whose supporting
set carries a `substrate` row, and a consumer MAY reject it; a `fail`
from a (`substrate`, `reconstructed`) row is still an observation the
artifact could not manufacture without performing a containment-relevant
event, weaker than an interception only in that it was derived after the
fact. A row fail-closed on `basis` or `method` sits at the bottom of both
orderings: its vantage is unknown, so it can strengthen nothing. `method`
bounds a `pass`: it answers whether a real event could have slipped past
the observation. A `pass` whose clean rows are all (`substrate`,
`intercepted`) makes the strongest absence claim this predicate can
carry; a `pass` resting on any `reconstructed` clean row SHOULD be read
as tolerating transients between the observed states; a `pass` resting on
any `artifact` clean row is self-reported absence, the weakest. The
clean-row ordering applies equally to the clean rows of a `degraded`
result's assessed classes, with `degraded` additionally bounded by its
disclosed coverage gap. Both orderings are consumer guidance; the
predicate still carries no verdict.

A caught row carrying `method: intercepted` SHOULD reference its
interception through `interceptRefs`, and a consumer MAY reject, never
downgrade, an attestation whose intercepted caught rows reference no
verifiable intercept record, as incoherent with this predicate's model
that each interception travels as an independently signed record. A clean
row's `intercepted` has no record by definition: an armed vantage that
captured nothing produces nothing to sign.

`basis` and `method` are producer claims, not cryptographically bound;
they defend against honest mislabeling and make the claim explicit,
uniform, and auditable. A substrate operator who signs false evidence is
outside this predicate's threat model, as for every self-asserted field:
that operator holds the signing key (see Model). Consumers MAY
coherence-check both fields against the pinned `observationEnvironment`:
a `substrate` row claiming a network-boundary observation under a
`networkPosture` that provides no interception path at that boundary is
incoherent, and a consumer MAY reject the attestation on that ground.

`actualLayer` names the enforcement layer that acted on the row's
containment event. It is required on every row; a row missing the member
is malformed under the framework's standard parsing rules and the
attestation is invalid, rather than the row forcing `fail`. That altitude
is deliberate and follows from the design invariant under Parsing Rules:
fail-closed-row semantics are reserved for members the recompute or the
documented consumer gating reads (`containmentObserved`, `basis`,
`method`); `actualLayer` is read by neither, so its absence is a
malformed statement, not weak evidence. On a row whose
`containmentObserved` label is from the published vocabulary but not in
the caught set (a clean row: nothing acted), the producer MUST emit the
literal string `none`. `none` is explicit rather than the field being
omitted so that "no layer needed to act" is distinguishable from an
accidental omission. `none` is also valid on a caught row, and there
states that the containment event was observed but no enforcement layer
acted: the observing vantage was positioned to see, not to act (a passive
tap, a monitor-only deployment). This is deliberate: enforcement role
travels here and only here, so `basis` never has to encode who could act.
Whether anything was positioned to see is answered by `basis` and
`method`: a clean row carrying (`basis: substrate`, `method:
intercepted`) states a live substrate vantage was armed and no capture
was attributed, which is the strongest claim a `pass` can rest on; a
`reconstructed` clean row makes the bounded version of that claim
described under `method`.

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
is the single canonical spelling: earlier internal versions used a
snake_case spelling (see the changelog), which is not accepted as an alias,
since two accepted spellings would mean two canonicalizations for the same
content. Migrating old producer output to the new name is a producer
concern, not something the wire format carries.

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

0.5 incorporates review feedback on 0.4:

-   Split the per-row `basis` field into two orthogonal required fields:
    `basis` (closed vocabulary `substrate` / `artifact`) now names only the
    observation's vantage, defined by its weakest input with a stated
    artifact-sourcing criterion, and the new `method` (closed vocabulary
    `intercepted` / `reconstructed`) names its directness, with the same
    weakest-input composition rule. The 0.4 values `substrate_observed`,
    `artifact_reported`, and `inferred` are rejected, not aliased, under
    the same single-canonicalization rule as the `does_not_assert` rename;
    `inferred` has no successor because it conflated the two axes.
-   Made `actualLayer` required on every row (missing member: malformed
    statement, deliberately a different altitude than the fail-closed row
    members, per the stated design invariant) and extended its literal
    `none` to caught rows, where it states observed-but-not-enforced, so
    enforcement role never leaks into `basis`.
-   Added consumer strength orderings on both sides with a defined
    supporting set (`basis` bounds a `fail`, `method` bounds a `pass`,
    fail-closed rows at the lattice bottom, clean-row ordering extended to
    `degraded`), a coherence check of row claims against the pinned
    `observationEnvironment`, and a row-internal check that intercepted
    caught rows reference verifiable intercept records — all consumer-side.
-   Stated the row-travel design invariant under Parsing Rules and the
    producer-claim trust boundary for `basis`/`method`.

[Runtime Traces]: runtime-trace.md
[SCAI]: scai.md
[SVR]: svr.md
[Test Result]: test-result.md
[VSA]: vsa.md
