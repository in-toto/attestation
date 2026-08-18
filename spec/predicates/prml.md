# Predicate type: PRML (Pre-Registered Manifest Lock)

Type URI: https://falsify.dev/prml/v0.1

Version: 0.1

Predicate Name: PRML

Authors: Cüneyt Öztürk ([@sk8ordie84](https://github.com/sk8ordie84))

## Purpose

This predicate expresses the **success criteria of an AI/ML evaluation, locked
before the evaluation is run**. It records the metric, comparator, threshold,
dataset identity (by content hash), random seed, and claimant identity, bound
together by a SHA-256 digest over a canonical serialisation of a [PRML]
manifest.

Every existing evaluation-shaped predicate in this directory describes what
happened _after_ a run (Test Result records outcomes, SVR records a policy
verdict, SCAI records observed attributes). None of them can answer the one
question a results attestation cannot answer about itself: **were the criteria
that decide pass/fail fixed before the result was known, or adjusted after?**
This predicate exists to make that ordering attestable. It is the evaluation
analogue of pre-registration in empirical science.

## Use Cases

An organisation evaluates an ML model against an accuracy bar (e.g. "accuracy
>= 0.85 on dataset X, seed 42"). The team commits the bar as a PRML manifest,
attests it (this predicate), and only then runs the evaluation. The results of
the run are attested separately — for example with the
[Test Result](test-result.md) predicate. A verification policy can then check
both attestations together and answer:

1.  Was the success bar committed **before** the results existed? (Compare the
    PRML attestation's timestamp — ideally countersigned by an RFC 3161 TSA or
    an append-only transparency log — against the time of the test run.)
2.  Did the reported run use the pre-committed criteria — same dataset (by
    content hash), same seed, same threshold — or did any of them drift?
3.  Does the reported result actually clear the pre-committed bar?

The Test Result predicate alone cannot cover this use case: it captures the
configuration _used_ by a run, as reported at result time. A configuration
reported together with the result can have been chosen after seeing the
result. The PRML predicate is deliberately produced at a different point in
the supply chain (before execution), which is the property policies check.

Concrete policy questions this predicate answers:

-   "Reject any evaluation claim whose criteria were not attested before the
    run started."
-   "Reject any result whose dataset hash or seed differs from the
    pre-committed manifest."
-   "Treat a result whose recomputed manifest digest does not match
    `manifest_sha256` as unsupported."

This mirrors requirements now appearing in AI-procurement and regulatory
contexts (e.g. EU AI Act Annex IV documentation of test logs, and the 2026
Code of Practice on Transparency's "documented internal testing"), where the
missing piece is evidence that the bar did not move.

## Prerequisites

The in-toto [Attestation Framework] and the [PRML] v0.1 specification
(openly specified, CC BY 4.0). Producing or verifying the predicate requires
only a SHA-256 implementation; a reference implementation exists but is not
required.

## Model

The producer of this attestation is the **claimant** of a future evaluation
result (a model developer, an eval team, a vendor under a contractual
pre-registration clause) or a neutral registry acting on their behalf. It is
produced **before** the evaluation step it constrains, and is consumed by
whoever must rely on the evaluation result afterwards (an internal reviewer, a
buyer, an auditor, a leaderboard maintainer).

The `subject` of the Statement is the locked claim itself (named by its
`claim_id`, with digest = the SHA-256 of the canonical PRML manifest) together
with the dataset the claim is about (named by its dataset id, with digest =
the dataset content hash). This lets downstream policies match later
result-time attestations about the same dataset artifact.

## Schema

```jsonc
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [
        {
            "name": "<PRML claim_id (UUIDv7)>",
            "digest": { "sha256": "<SHA-256 of the canonical PRML manifest>" }
        },
        {
            "name": "<dataset id>",
            "digest": { "sha256": "<dataset content hash>" }
        }
    ],
    "predicateType": "https://falsify.dev/prml/v0.1",
    "predicate": {
        "claim_id": "<UUIDv7>",
        "created_at": "<RFC 3339 timestamp, UTC>",
        "metric": "<metric name>",
        "comparator": "<one of \">=\", \"<=\", \">\", \"<\", \"==\">",
        "threshold": <number>,
        "seed": <integer>,
        "producer": { "id": "<claimant identity>" },
        "prml_version": "prml/0.1",
        "manifest_sha256": "<SHA-256 of the canonical PRML manifest>"
    }
}
```

### Parsing Rules

This predicate follows the in-toto Attestation Framework's
[standard parsing rules], with one deliberate deviation: predicate field names
use `snake_case`, not lowerCamelCase. The fields mirror, name-for-name, the
fields of the PRML manifest whose canonical serialisation is bound by
`manifest_sha256`; renaming them at the attestation layer would break the
byte-level correspondence between the predicate and the digest it carries.
Consumers MUST recompute the manifest digest from the canonical PRML
serialisation as defined by the [PRML] specification, and MUST treat a
mismatch with `manifest_sha256` (or with the first subject's digest) as a
verification failure.

The predicate is versioned by the Type URI together with `prml_version`.
Unknown fields MUST be ignored (monotonic principle): additional fields may
appear in future PRML versions, but the meaning of the fields defined here
will not change within v0.1.

### Fields

`claim_id` _string (UUIDv7), required_

Unique identifier of the locked claim. UUIDv7 is required by PRML so that the
identifier itself carries a coarse creation time.

`created_at` _string (RFC 3339, timezone Z), required_

Time at which the claimant locked the criteria. This is the claimant's own
assertion; policies that need a non-self-asserted time SHOULD rely on an
independent countersignature of the attestation (RFC 3161 timestamp authority
or an append-only transparency log) rather than on this field.

`metric` _string, required_

Name of the metric that decides the claim (e.g. `accuracy`, `pass_rate`).
Semantics are agreed between producer and consumer, as with test names in the
Test Result predicate.

`comparator` _string (enum), required_

One of `>=`, `<=`, `>`, `<`, `==`. The claim succeeds iff
`observed <comparator> threshold` (for `==`, within the PRML-specified
tolerance).

`threshold` _number, required_

The pre-committed bar the observed metric value is compared against.

`seed` _integer, required_

Random seed pre-committed for the evaluation run.

`producer` _object, required_

Identity of the claimant. Contains at least `id` (string).

`prml_version` _string, required_

The PRML specification version of the locked manifest, e.g. `prml/0.1`.

`manifest_sha256` _string (hex SHA-256), required_

SHA-256 digest over the canonical serialisation of the full PRML manifest, as
defined by the PRML specification. Equals the digest of the first subject.

## Example

Generated by the PRML reference implementation from a valid manifest (the
statement below round-trips through its validator):

```json
{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [
        {
            "name": "01991a2b-7c3d-7e4f-8a5b-6c7d8e9f0a1b",
            "digest": {
                "sha256": "97a4da62e1d22a35d336d58ee98c4d3b22e94e150d9dd6952fb43f186565e511"
            }
        },
        {
            "name": "mmlu-test-v1",
            "digest": {
                "sha256": "2a3548c1e7beb15958c29066826fbdf0ee37bee27de3d28096022e820855f133"
            }
        }
    ],
    "predicateType": "https://falsify.dev/prml/v0.1",
    "predicate": {
        "claim_id": "01991a2b-7c3d-7e4f-8a5b-6c7d8e9f0a1b",
        "created_at": "2026-08-18T21:00:00Z",
        "metric": "accuracy",
        "comparator": ">=",
        "threshold": 0.85,
        "seed": 42,
        "producer": {
            "id": "eval-team.example.org"
        },
        "prml_version": "prml/0.1",
        "manifest_sha256": "97a4da62e1d22a35d336d58ee98c4d3b22e94e150d9dd6952fb43f186565e511"
    }
}
```

A typical verification policy pairs this attestation with a later
[Test Result](test-result.md) attestation over the same dataset subject, and
accepts the result only if the PRML attestation predates the run and the
run's configuration matches the pre-committed fields.

## Changelog and Migrations

Initial version, tracking PRML v0.1.

[Attestation Framework]: https://github.com/in-toto/attestation
[PRML]: https://spec.falsify.dev/v0.1
[standard parsing rules]: ../v1/README.md#parsing-rules
