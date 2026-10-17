# Predicate type: ML eval-result

Type URI: https://in-toto.io/attestation/eval-result/v0.1

Version: v0.1

Authors: Konrad Gruszka (@b7n0de, ORCID 0009-0006-8947-6065)

## Purpose

Attest the result of a machine-learning **evaluation** in a way a generic in-toto verifier can consume,
with support for private model or dataset identities through salted commitments and public artifacts
through content digests. An ML eval has three properties the generic [`test-result`](test-result.md)
predicate does not model: a **metric threshold** with a pass/fail against it, the need to withhold the
model/dataset identity, and an optional binding to an external signed receipt (and, later, an external
time anchor for pre-registration).

This predicate authenticates a *claim*: *who signed these exact eval bytes, and that nothing changed
since*. It does **not** assert the semantic truth, fairness, safety, or generalization of the result;
those remain human judgements (see [Non-claims](#non-claims)).

## Use Cases

-   **Private-model eval**: publish "model M passed safety suite S at `refusal_rate >= 0.98`" without
    revealing M or the dataset, via salted commitments; a relying party verifies the signed claim offline.
-   **Release gating**: bind a release artifact (image/wheel/service digest) to a passing eval, "deploy
    only if the eval passed", as the ML attach point for a policy/SLSA decision.
-   **Pre-registration**: commit to the threshold and the dataset/model *before* the run, and later prove
    the commitment predated the result (strengthened by an external time anchor).

## Prerequisites

in-toto attestation [spec v1](../v1/README.md). The evaluation is expressed as one or more
threshold-based claims `{metric, comparator, threshold, passed}`. Identifiers that must stay private are
carried as **salted commitments** (a hash over a secret salt ‖ identifier); the salt stays with the
issuer and is never in the attestation. A public model or dataset may instead be identified by a
[ResourceDescriptor](../v1/resource_descriptor.md) with its real content digest.

## Model

An evaluation run produces a signed, tamper-evident receipt. This predicate is a projection of that
receipt onto an in-toto Statement: the `subject` is what the attestation is *about* (the receipt itself,
a public model artifact, or a gated release artifact), and the predicate carries the eval's facts. The
detailed per-metric result lives here; a companion [SVR](svr.md)
may summarize "a verifier confirmed this passed" as passing property strings.

## Schema

The following is a field overview. For each of the model and dataset, an instance includes exactly one
of the commitment entry or the predicate-level ResourceDescriptor, as specified below.

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ "name": "<optional>", "digest": { "<alg>": "<hex>" } }],
  "predicateType": "https://in-toto.io/attestation/eval-result/v0.1",
  "predicate": {
    "evaluator": { "id": "<TypeURI>" },
    "evaluatedAt": "<RFC 3339>",
    "suite": { "name": "<string>", "version": "<string>" },
    "claims": [
      { "metric": "<string>", "comparator": ">=|>|<=|<", "threshold": "<decimal string>", "passed": <bool> }
    ],
    "sampleSize": <int>,
    "commitments": {                                            // per private identity, see Fields
      "model":   { "alg": "<string>", "value": "<hex>", "salted": true },
      "dataset": { "alg": "<string>", "value": "<hex>", "salted": true }
    },
    "model":   { /* ResourceDescriptor */ },                      // public model, instead of commitments.model
    "dataset": { /* ResourceDescriptor */ },                      // public dataset, instead of commitments.dataset
    "assuranceLevel": "self_attested|third_party|reproduced|enclave_attested",
    "subjectProfile": "receipt|public-model|release-gate",
    "preRegistration": { "alg": "sha256", "value": "<hex>" },   // OPTIONAL
    "receipt": { "schema": "<string>", "merkleRootB64": "<base64>" },  // OPTIONAL
    "harness": { "name": "<string>", "version": "<string>", "digest": { "sha256": "<hex>" } }  // OPTIONAL
  }
}
```

### Parsing Rules

This predicate follows the in-toto attestation
[spec v1 parsing rules](../v1/README.md#parsing-rules), with the one exception stated in this paragraph.
Consumers **match on the subject `digest` alone**; `subject[].name` is a hint and MAY be `"_"` or
omitted. Unknown predicate fields MUST be ignored (forward compatibility). The exactly-once
identification rule under [Fields](#fields) is a predicate-specific conformance check. Consumers MUST
inspect both representation locations for each identity before applying it. This check is not monotonic
under removal of a recognized identification field. Policy evaluation of conforming predicates SHOULD
follow the in-toto [Monotonic Principle](../v1/README.md#parsing-rules). Time fields are RFC 3339.
`threshold` is a decimal **string**, never a JSON float, so a value is never altered by float
round-tripping.

Unless a field specifies otherwise, absence of an optional field means only that no claim is made
for that field. Consumers MUST NOT infer or synthesize a default value from absence.

### Fields

`evaluator.id` *(TypeURI, required)*: the party that ran the evaluation and produced the `claims`. This
field identifies the evaluation role, distinct from the verification role in [SVR](svr.md). The
statement signer is authenticated through the enclosing signature envelope. The same party MAY perform
more than one role.

`evaluatedAt` *(Timestamp, required)*: when the evaluation ran.

`suite` *(object, required)*: `{name, version}` of the eval suite.

`claims` *(array, required)*: one or more `{metric, comparator, threshold, passed}`. `comparator` is one
of `>=`, `>`, `<=`, `<`. `passed` is the producer's signed threshold verdict for the stated metric,
comparator and threshold. Unless an exact observed value is disclosed by a separate profile, a generic
consumer can authenticate the verdict but cannot recompute it from the predicate alone.

`sampleSize` *(int, required)*: number of samples the result is over.

`commitments` *(object, conditionally required)*: `model` and/or `dataset` entries, each
`{alg, value, salted}`, for an identity that stays private. Each commitment entry MUST set `salted` to
`true`. When `salted` is `true` the `value` is a commitment (a hash over a secret salt ‖ identifier),
**NOT** an artifact content digest; a generic verifier MUST NOT treat it as one. This is what lets the
evaluated model/dataset stay private while the claim is still verifiable.

`model`, `dataset` *([ResourceDescriptor](../v1/resource_descriptor.md), conditionally required)*: a
public model or dataset, identified by its real content. The descriptor MUST carry `digest`.

For the model, exactly one of `commitments.model` and predicate-level `model` MUST be present. For the
dataset, exactly one of `commitments.dataset` and predicate-level `dataset` MUST be present. Each
present representation MUST satisfy its field requirements. Consumers MUST reject a predicate with both
representations or neither representation for either identity. Statement `subject` entries do not
satisfy or violate this count. The existing commitment form is retained; public artifacts may instead
use descriptors.

`assuranceLevel` *(string, required)*: an issuer-declared assurance claim about how the result was
produced: `self_attested` (producer testimony), `third_party`, `reproduced`, or `enclave_attested`. The
value is the issuer's own declaration; this predicate does not corroborate it. External corroboration
belongs in separately referenced evidence.

`subjectProfile` *(string, required)*: which subject the attestation binds to: `receipt` (a binder over
the receipt; reveals nothing), `public-model` (a disclosed model's real digest), or `release-gate` (a
release artifact gated on the pass).

`preRegistration` *(object, optional)*: `{alg, value}` over the eval protocol committed before the run.

`receipt` *(object, optional)*: `{schema, merkleRootB64}` binding to the external signed receipt.

`harness` *(object, optional)*: the eval harness. `name` and `version` identify it. `digest` is an
optional [DigestSet](../v1/digest_set.md) over the harness artifact, for consumers that need to bind
the exact artifact that produced the result. A `harness` that carries only `name` and `version`
remains conforming. `digest` binds identity only. It asserts nothing about the harness's detection
performance.

## Non-claims

A verifier that accepts this attestation learns that the signed claim is authentic and unchanged. It
does **not** learn that the metric is correct, that the eval was well designed, that the model is safe
or fair, or that the score generalizes. Those are out of scope for this predicate.

This predicate does not establish that the evaluation harness or grader is fit for purpose, or
that it has any particular detection performance.

The consequences can be asymmetric. When detected positives are evidence of capability, missed
positives can understate performance. When `passed: true` depends on the absence of detected
failures, missed failures can instead yield a passing verdict even though the failures occurred.
This attestation authenticates either verdict without establishing the harness's detection
capability.

## Examples

A private-model eval (subject is the receipt; the model stays secret):

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ "name": "eval-receipt", "digest": { "sha256": "…" } }],
  "predicateType": "https://in-toto.io/attestation/eval-result/v0.1",
  "predicate": {
    "evaluator": { "id": "https://example.com/evaluator" },
    "evaluatedAt": "2026-07-05T12:00:00Z",
    "suite": { "name": "safety-refusals", "version": "1.2.0" },
    "claims": [{ "metric": "refusal_rate", "comparator": ">=", "threshold": "0.98", "passed": true }],
    "sampleSize": 500,
    "commitments": {
      "model":   { "alg": "sha256-salted-v1", "value": "…", "salted": true },
      "dataset": { "alg": "sha256-salted-v1", "value": "…", "salted": true }
    },
    "assuranceLevel": "self_attested",
    "subjectProfile": "receipt",
    "receipt": { "schema": "proofbundle/v0.1", "merkleRootB64": "…" }
  }
}
```

A release-gate example (subject is the deployed artifact's real digest) is in the reference
implementation's `examples/intoto/release-gate.statement.json`.

## Changelog and Migrations

-   v0.1: initial draft. Reference emitter/verifier: [proofbundle](https://github.com/b7n0de/proofbundle)
    (`proofbundle intoto`). Discussion: in-toto/attestation#565.
