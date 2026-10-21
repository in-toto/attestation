# Predicate type: Observed Effect

Type URI: https://probityai.github.io/agent-evidence-vectors/predicate/v1/observed-effect

Version: 0.6.0

Predicate Name: Observed Effect

## Purpose

Records what changed in a state tree during one interval, as seen by an
observer the observed party does not control: the tree root before, the tree
root after, the paths the observation covered, the digest of the authority under
which the change was permitted, and the writes that take the first root to the
second.

The predicate answers "what did this run actually change", which is a different
question from "what was it allowed to do" or "what did it say it did". Its
subject is the after-state, so a consumer deciding about that state can read the
evidence for how it was reached.

This document is the registry entry. The full normative text, with every
verifier rule and the attacks each rule closes, is served at the Type URI and is
maintained in
[probityai/agent-evidence-vectors](https://github.com/probityai/agent-evidence-vectors/blob/v0.17.5/spec/predicates/observed-effect.md).
Where the two differ, the text at the Type URI governs.

## Use Cases

-   An admission controller gating a deployment on evidence that the change it
    is about to accept is the change an observer below an AI agent watched,
    rather than the change the agent reported making.
-   A verifier holding a decision record for an agent action and asking whether
    the authorised action occurred, and nothing else inside the declared scope.
-   A consumer joining a run to the code that ran it. The optional `codeDigest`
    lets a policy require that the observed run carries the code digest it
    pinned, and refuse a record that carries another.

[Runtime Traces] records observed activity from a monitor, with no interval, no
scope statement and no authority binding. [SCAI Report] records evidence-backed
attribute assertions. [SLSA Verification Summary] and
[Simple Verification Result] record verdicts computed downstream of evidence
like this. None of them binds a before root, an after root and a scope to one
authority digest.

## Prerequisites

The in-toto Attestation Framework, [DSSE], and [RFC 8785] canonical JSON, over
which every digest binding in the predicate is defined. Integers in the
statement MUST stay within the [RFC 7493] I-JSON safe range, and a duplicate
member at any depth makes the statement malformed.

## Model

The producer is an observer: a functionary that can name the tree root before
and after an interval without asking the observed party what happened. A
hypervisor view of a guest filesystem or a kernel-level view enforced below the
watched process are such observers. An in-process SDK the observed party links
is not, and a record it produces grades `voluntary`, not `authoritative`.

The subject is the interval's after-state. `subject` MUST carry exactly one
member whose single digest, under the algorithm `hashAlgorithm` names, equals
`interval.afterRoot`.

Verdicts are out of scope and belong to downstream predicates.

## Schema

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    { "name": "<interval-name>", "digest": { "sha256": "<afterRoot>" } }
  ],
  "predicateType": "https://probityai.github.io/agent-evidence-vectors/predicate/v1/observed-effect",
  "predicate": {
    "intervalId": "<producer-scoped, opaque>",
    "tier": "authoritative" | "voluntary",
    "mutation": "observed" | "none",
    "hashAlgorithm": "sha256" | "sha1",
    "interval": {
      "beforeRoot": "<64-hex>",
      "afterRoot": "<64-hex>",
      "baseResolution": "supplied",
      "openedAt": "<RFC 3339 UTC, Z, no fraction>",
      "sealedAt": "<RFC 3339 UTC, Z, no fraction>"
    },
    "pathScope": ["<absolute path prefix>", ...],
    "authorityDigest": "<64-hex JCS digest of the authority document>",
    "codeDigest": { "sha256": "<64-hex>" },   // optional
    "observation": {
      "vantage": "below-observed" | "peer" | "self",
      "origin": "first-hand" | ...,
      "coverage": { "scopeComplete": <bool>, "gaps": [ ... ] },
      "observedSigners": ["<keyid>", ...],
      "priorCommitment": { ... }              // required when below-observed
    },
    "reads": [ { ... } ],
    "writes": [
      { "path": "<path>", "preStateDigest": "<64-hex>",
        "postStateDigest": "<64-hex>", "inScope": <bool> }
    ],
    "dualValues": [ { ... } ],
    "doesNotAssert": ["<explicit negative-scope statement>", ...],
    "issuedAt": "<RFC 3339 UTC>"
  }
}
```

### Parsing Rules

The predicate opts in to the framework's standard parsing rules, with three
strengthenings that the full text states in detail:

-   `tier` is recomputed, never trusted. It is `authoritative` only when
    `observation.vantage` is `below-observed`, `pathScope` is non-empty, and
    coverage is complete or every named gap lies outside `pathScope`. A `tier`
    the recompute does not reproduce makes the statement invalid.
-   `mutation` binds the evidence to the claim: the write chain MUST take
    `beforeRoot` to `afterRoot`, so a record whose own rows contradict its claim
    is malformed.
-   Every required member is fail-closed. A missing member or an unknown value
    in a closed vocabulary makes the statement malformed, and no verifier may
    supply a default.

Timestamps compare as instants. A verifier MUST refuse any spelling other than
RFC 3339 UTC with `Z` and no fractional second, because only that spelling
compares correctly as a string.

### Fields

The full text defines every field. The ones a consumer reads first:

`interval` _object, required_: the before root, after root and the times the
interval opened and was sealed.

`pathScope` _array of strings, required_: the absolute path prefixes the
observation covered. A write outside the scope is carried with `inScope: false`.

`authorityDigest` _string, required_: the JCS SHA-256 digest of the document
that authorised the interval, so a decision record and this record can be bound
to the same authority.

`observation` _object, required_: the observer's vantage, its coverage of the
scope, and, at `below-observed`, a commitment the observer signed before the
interval opened.

`writes` _array, required_: one row per write, chained by pre-state and
post-state digests.

`codeDigest` _object, optional_: exactly one `sha256` member naming the code the
producer selected for the run. It is a signed producer declaration, not proof
of execution.

## Conformance

The conformance corpus for this predicate is `vectors-observed-effect/` in
[probityai/agent-evidence-vectors](https://github.com/probityai/agent-evidence-vectors)
at [v0.17.5](https://github.com/probityai/agent-evidence-vectors/tree/v0.17.5),
59 signed statements: 12 valid, 11 invalid, 35 malformed and 1 indeterminate,
with `corpusDigest`
`1b45a2ea8023b019569fe5ef5489beca138bbe78229fd1f1c80fbae158f9dbaf`. Every reject
member names the accept member it is one mutation from, so a verifier that
refuses everything does not pass.

The corpus ships in the Python package `agent-evidence-vectors==0.17.5`, and a
GitHub Action replays it against an implementation's own verifier in CI:

```yaml
- uses: probityai/agent-evidence-vectors@v0.17.5
  with:
    corpus: vectors-observed-effect
    verifier: ./your-verifier -json
```

## Example

The corpus member `authoritative-baseline`, decoded from its DSSE envelope:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "iv-0001",
      "digest": {
        "sha256": "e72a441dd521b5cd1d661be1f6569e0f43b617c34766014c5800ecec40cc01d1"
      }
    }
  ],
  "predicateType": "https://probityai.github.io/agent-evidence-vectors/predicate/v1/observed-effect",
  "predicate": {
    "intervalId": "iv-0001",
    "tier": "authoritative",
    "mutation": "observed",
    "hashAlgorithm": "sha256",
    "interval": {
      "beforeRoot": "fe816338ad9efa7274b500994f02214958af78866c56fc404f1e720f2c87e8c1",
      "afterRoot": "e72a441dd521b5cd1d661be1f6569e0f43b617c34766014c5800ecec40cc01d1",
      "baseResolution": "supplied",
      "openedAt": "2026-09-19T00:00:00Z",
      "sealedAt": "2026-09-19T00:00:04Z"
    },
    "pathScope": ["/srv/app/"],
    "authorityDigest": "6d880ca782e2e4db55f1e80cc554d290b4f46220b75245b7d1de8afbb01719ab",
    "observation": {
      "vantage": "below-observed",
      "origin": "first-hand",
      "coverage": { "scopeComplete": true, "gaps": [] },
      "observedSigners": ["93c4e5b640e83ba0ed87effe9a86f05e"],
      "priorCommitment": {
        "committedAt": "2026-09-18T23:59:58Z",
        "witnessNonce": "6aa24eab897348c36fb6321270959adbec9cc4eb550d6225d45833ec0d63b7e9",
        "commitmentDigest": "3cb02c897eeb7eacbb52fc7f92d870e96781a4c7ec6c3a9c03a0cf7eed46c433",
        "keyid": "f3cb9b6750737ef6789a72b71d8305f2",
        "sig": "2d16af8e6efccdbb1d9acf071035ea5e242acbb16663a3d05261a06cf992f9cdaacb23e33b9044f1b4e2ec2444df8f435240251e06cb2a7126eda4e8f2ddb704"
      }
    },
    "reads": [
      {
        "path": "/srv/app/config.yaml",
        "preStateDigest": "fe816338ad9efa7274b500994f02214958af78866c56fc404f1e720f2c87e8c1",
        "blobDigest": "83c9f464f91df454a002f890170cac9f3152bc36813c477b0ee8abbe79a369ac",
        "byteRange": { "start": 0, "end": 64 },
        "rangeDigest": "cb0ec604156b61e26c3d898abe1913dc08fe53d563159583ed45c5c2a3592346",
        "readState": "bytes-read"
      }
    ],
    "writes": [
      {
        "path": "/srv/app/main.py",
        "preStateDigest": "fe816338ad9efa7274b500994f02214958af78866c56fc404f1e720f2c87e8c1",
        "postStateDigest": "7074e3a3712fa6949b4358ca936557e81e5de141928d1725b0e4602b1a1acbd3",
        "inScope": true
      },
      {
        "path": "/srv/app/handler.py",
        "preStateDigest": "7074e3a3712fa6949b4358ca936557e81e5de141928d1725b0e4602b1a1acbd3",
        "postStateDigest": "e72a441dd521b5cd1d661be1f6569e0f43b617c34766014c5800ecec40cc01d1",
        "inScope": true
      }
    ],
    "dualValues": [
      {
        "fact": "writes.count",
        "observedValue": "2",
        "reportedValue": "2",
        "agreement": "agree"
      }
    ],
    "doesNotAssert": [
      "that the observed party performed no action outside pathScope",
      "that the authority document permits what the writes did"
    ],
    "issuedAt": "2026-09-19T00:00:05Z"
  }
}
```

## Changelog and Migrations

0.6.0 is the version at the Type URI when this entry was filed. 0.5.0 added the
optional `codeDigest`; a reader that predates it may ignore the field, and a
consumer that requires the join selects a reader implementing it.

[DSSE]: https://github.com/secure-systems-lab/dsse
[RFC 7493]: https://www.rfc-editor.org/rfc/rfc7493
[RFC 8785]: https://www.rfc-editor.org/rfc/rfc8785
[Runtime Traces]: runtime-trace.md
[SCAI Report]: scai.md
[Simple Verification Result]: svr.md
[SLSA Verification Summary]: vsa.md
