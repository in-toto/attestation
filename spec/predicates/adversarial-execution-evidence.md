# Predicate type: Adversarial Execution Evidence

<a id="req-type-uri"></a>
Type URI: https://in-toto.io/attestation/adversarial-execution-evidence/v0.7

Version: 0.7.0

Predicate Name: Adversarial Execution Evidence

> Status: DRAFT submission for vetting. The schema matches a production
> implementation and a public conformance suite.

## Purpose

Records what a containment substrate observed while a known corpus of adversarial
inputs was thrown at an under-trusted artifact -- an agent tool, an MCP server, a
plugin, a build step -- executing inside that substrate. The predicate carries
what was thrown, what the substrate was configured to catch, what it observed
(each observation an independently signed record), and the coverage bounds the
observation holds under. A consumer recomputes `result` from the carried
predicate alone. The conformance [suite] holds executable vectors, including the
required refusals, and a reference verifier, and the [long form] carries what this
page states without arguing, in `rationale.md`, `wire-profile.md`, `fields.md`,
`verification.md`, `changelog.md` and `settled-readings.md`, which states this
text's reading of seven places a second implementation found underdetermined.

## Use Cases

-   An admission controller gating a third-party MCP server or agent tool image
    on evidence that it was executed against a named attack corpus under an
    enforcing catch policy, with the policy digest and network posture pinned.
-   An auditor re-verifying offline, without trusting the producer's
    infrastructure, that a specific interception happened: the signed record
    binds the destination, the payload commitment and the substrate context.
-   A security team checking that two runs of one artifact assessed the same
    attacks, because the manifest is digest-committed at attack granularity.

Adjacent predicates cover different ground: [Runtime Traces] carries monitor
activity with no corpus binding and no coverage denominator, [SCAI] carries
attribute assertions with no adversarial corpus, [VSA] and [SVR] carry verdicts
computed downstream of evidence like this, and [Test Result] carries outcomes with
no cryptographic binding of the inputs. This predicate is the evidence layer those
predicates consume.

## Prerequisites

The in-toto Attestation Framework, [DSSE] (each observation record is a
DSSE-shaped envelope) and [RFC 8785] canonical JSON, which every digest binding
is defined over. Statements and record payloads are parsed under a pinned encoding profile:
[RFC 7493] safe integers, no duplicate member at any depth, well-formed Unicode
scalar values checked on the raw bytes, a nesting bound of 128, and BMP-only
strings on every signed canonical surface, and canonical [RFC 4648] base64
wherever a signed surface carries base64. A violation is malformed, fail-closed.

**Run binding.** A statement carrying a `basis: substrate` row derives a run
binding digest: the lowercase 64-hex SHA-256 of the RFC 8785 canonicalization of
`{"aeeBindingVersion": "2", "catchPolicy": <digest>, "corpus": <digest>,
"networkPosture": <digest of the carried object>, "observationVocabulary":
<digest>, "runEntropy": <digest>, "subject": <subject[0] digest>, "substrate":
<digest>}`. Every input is configuration fixed before corpus injection: the arming record
carries this digest inside its own signature and is signed before injection, so an
outcome of the run can never be an input. No field carries it, and every
substrate-signed record repeats it inside its signed payload, which is what makes
a record signed under a different configuration unspliceable. `subject` MUST carry
exactly one entry, and a verifier MUST reject, fail-closed, an
`aeeBindingVersion` it does not implement. The run binding's
`networkPosture` input is a digest over the whole carried object; **the pinned
`networkPosture` digest**, which `arming` and `sealed` records repeat as
`aeePostureDigest`, is the distinct value
`observationEnvironment.networkPosture.digest.sha256`. Reading either for the
other refuses every substrate-row statement the other accepts.

## Model

The producer is a containment substrate operator: a functionary that runs the
subject artifact inside an isolated, instrumented environment (a microVM, a
sandbox, an eBPF-supervised process), injects the corpus, and signs what the
substrate observed -- each interception, the armed vantage it occurred under, and
the seal that the vantage stayed armed. The subject is the executed artifact, by
digest, and the substrate is referenced by subject so the evidence inherits a
trust chain rather than a bare name. Verdicts against an organization's policy
belong in a downstream predicate such as [VSA], computed over this evidence.

## Schema

<a id="req-statement-schema"></a>

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    { "name": "<artifact-name>", "digest": { "sha256": "<64-hex>" } }
  ],
  "predicateType": "https://in-toto.io/attestation/adversarial-execution-evidence/v0.7",
  "predicate": {
    "result": "fail",
    "observationEnvironment": {
      "substrate": { "name": "<substrate-subject>", "digest": { "sha256": "<64-hex>" } },
      "corpus": {
        "name": "<corpus-name>", "uri": "pkg:<producer>/<corpus>@<version>",
        "digest": { "sha256": "<JCS digest of manifest>" },
        "manifest": {
          "classes": { "CO": ["CO-EXFIL-1"] },
          "expectedPayloads": { "CO-EXFIL-1": ["<64-hex>"] }
        }
      },
      "catchPolicy": { "digest": { "sha256": "<JCS digest of catch policy>" } },
      "networkPosture": { "posture": "sinkhole", "digest": { "sha256": "<64-hex>" } },
      "observationVocabulary": {
        "digest": { "sha256": "<JCS digest of {caught, labels}>" },
        "labels": ["egress_captured", "no_egress"], "caught": ["egress_captured"]
      },
      "runEntropy": { "digest": { "sha256": "<run-start value>" } }
    },
    "coverage": { "assessedClasses": ["CO"], "outOfScope": {}, "routedElsewhere": {} },
    "attackResults": [
      {
        "attackId": "CO-EXFIL-1", "containmentObserved": "egress_captured",
        "basis": "substrate", "method": "intercepted", "attribution": "pinned",
        "actualLayer": "policy.egress_sinkhole", "observationRefs": [0]
      }
    ],
    "observationRecords": [
      { "payload": "<base64 canonical +json>", "payloadType": "<... +json>",
        "signatures": [ { "keyid": "<hex>", "sig": "<base64>" } ] }
    ],
    "batchRoot": "<RFC 6962 root over observation records>",
    "doesNotAssert": [ "<negative-scope statements>" ],
    "issuedAt": "2026-06-23T16:08:07Z"
  }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules),
with two strengthenings. `result` is not an independent claim: a consumer MUST be
able to recompute it from the rest of the predicate, and a `result` the recompute
does not reproduce makes the attestation invalid. The recompute reads the carried
predicate alone, never `observationRecords`, a signature outcome, or a consumer
trust decision. And observation records are verify-then-read: a field inside a
payload means nothing until that payload's signature verifies, DSSE PAE over
`(payloadType, payload)`, against a key the consumer trusts.

Four checks are consumption preconditions, each a function of the carried
statement alone: well-formedness including run-binding derivability; eleven
coverage validity requirements; the `result` recompute; and manifest and
vocabulary digest integrity. A consumer that consumes `result` or credits any row
MUST evaluate them first, and on failure the attestation is invalid. None of them
costs a secret against a party able to author a record, so they are security
properties only together with record signature verification.

### Fields

`result` _string, required_

One of `fail`, `degraded`, `pass_indirect`, `pass`, recomputed as the minimum
over three conditions: a caught or fail-closed row gives `fail`, a disclosed
coverage gap gives `degraded`, and a clean row that is not (`substrate`,
`intercepted`) gives `pass_indirect`.

`observationEnvironment` _object, required_

The digest-pinned run context: `substrate`; `corpus`, whose `manifest` carries the
required `classes` denominator and optional per-attack `expectedPayloads`;
`catchPolicy`; `networkPosture`, a `posture` from the closed set `allowlist`,
`no_network`, `sinkhole`, `unsafe_bypass_egress`; `observationVocabulary`, being
`labels`, `caught` and a digest over both; and `runEntropy`, required exactly when
a row declares `basis: substrate`. The manifest and vocabulary pre-images travel,
so a verifier re-derives their digests offline and an edit to either set fails.

`coverage` _object, required_

`assessedClasses`, `outOfScope` and `routedElsewhere`: a disjoint partition of
the manifest's classes, where a class in two of them or a manifest class in none
is malformed.

`attackResults` _array of objects, required_

One row per executed attack, all seven members required and every closed
vocabulary fail-closed on an unknown value: `attackId`, `containmentObserved`,
`basis` (`substrate` or `artifact`, the vantage of the weakest input), `method`
(`intercepted` or `reconstructed`), `attribution` (`pinned` or `paired`, checked
against `expectedPayloads`), `actualLayer` (or the literal `none`), and
`observationRefs`.

`observationRecords` _array of objects, optional_

One DSSE envelope per observation. A payload covering a substrate row MUST be
canonical I-JSON under the profile above and MUST carry `aeeRunBinding`,
`aeeMethod` and `aeeKind` (`interception`, `arming`, `sealed`, `examination`,
`moat-drop`, `uncommitted-observation`) plus the members that kind requires, or
it covers nothing.

`batchRoot` _string, required when `observationRecords` is non-empty_

An RFC 6962 SHA-256 root over the records' DSSE PAE bytes. It binds the carried
set against a party who cannot re-sign the enclosing envelope; it never
establishes completeness against the run.

`doesNotAssert` _array of strings, optional_

Explicit negative scope. Advisory: a verifier MUST NOT require it.

`issuedAt` _Timestamp, required_

RFC 3339 in UTC, separator and zone designator uppercase, zone one of `Z`,
`+00:00` or `-00:00`. `armedAt` carries the same profile.

**Evidence tier (derived, never carried).** Before crediting a
`basis: substrate` row a consumer MUST derive a per-row tier -- `declared` for an
`artifact` row, `attested` when every covering record's signature verifies against
a key its policy names as a substrate observation key, `unattested` otherwise --
and MUST NOT infer that key from the predicate. The tier never alters `result`.
A predicate-level member named `evidenceTier`, or any beginning with the reserved
prefix `aee`, MUST be ignored. The [long form] carries the derivation.

## Examples

### Caught egress, whole Statement

The Schema block above is a complete Statement: one caught row under a sinkhole
posture, `attribution: pinned` against the commitment the manifest declared, and
three records -- the interception the row resolves, the arming record, and the
seal. Its `corpus.digest` is re-derivable from the embedded manifest, the whole
object and not `classes` alone, canonicalized under RFC 8785 and hashed.

### Clean run, coverage complete

`result` is `pass` when the one row is clean -- its `containmentObserved` is in
`labels` and not in `caught` -- and is (`substrate`, `intercepted`). Such a row
carries the literal `none` for `actualLayer` and resolves an `arming` record and a
covering `sealed` record rather than an `interception`, which a clean row MUST NOT
resolve:

```jsonc
{ "attackId": "CO-EXFIL-1", "containmentObserved": "no_egress",
  "basis": "substrate", "method": "intercepted", "attribution": "paired",
  "actualLayer": "none", "observationRefs": [0, 1] }
```

## Changelog and Migrations

A member is born exactly when a normative reader consumes it, not retroactively.
Renames land with no alias, since two accepted spellings would mean two
canonicalizations for one content.

-   **0.1-0.2** internal producer iterations.
-   **0.3** first shape proposed for vetting: verdicts moved downstream, payload
    bytes replaced by commitments, `batchRoot` moved to predicate level.
-   **0.4** required per-row vantage field; `does_not_assert` renamed.
-   **0.5** vantage split into the orthogonal `basis` and `method`;
    `actualLayer` required on every row; strength orderings added.
-   **0.6** coverage denominator made checkable: `classes` in the manifest, the
    vocabulary carried in the attestation, coverage validity as preconditions.
-   **0.7** breaking. Makes checkable four things 0.6 could only describe.
    `expectedPayloads` gives attribution strength its normative reader, so
    `attribution` is required here and at no earlier version;
    `aeeAssessedAttacks` bounds coverage inflation; and `aeeObservedSet` with
    `aeeObservedAttacks` bind record deletion and relabelling on a seal now
    required wherever a substrate row appears. It also adds `pass_indirect`,
    closes the posture vocabulary, and binds the carried `networkPosture` object
    and the vocabulary into run binding version 2. Revised in review: canonical
    base64 required, the pinned posture digest defined, seven readings pinned.

[DSSE]: https://github.com/secure-systems-lab/dsse
[RFC 4648]: https://www.rfc-editor.org/rfc/rfc4648
[RFC 7493]: https://www.rfc-editor.org/rfc/rfc7493
[RFC 8785]: https://www.rfc-editor.org/rfc/rfc8785
[Runtime Traces]: runtime-trace.md
[SCAI]: scai.md
[SVR]: svr.md
[Test Result]: test-result.md
[VSA]: vsa.md
[suite]: https://github.com/probityai/agent-evidence-vectors
[long form]: https://github.com/probityai/agent-evidence-vectors/tree/main/spec/predicates/adversarial-execution-evidence
