# Predicate type: Computation Receipt

Type URI: https://in-toto.io/attestation/computation-receipt/v0.1

Version: v0.1

Predicate Name: Computation Receipt

Authors: Anomly, Inc. (https://github.com/anomly-labs)

## Purpose

To attest that a numeric computation — an HPC reduction, a training step, a model
inference — produced a specific output, in a form checkable by re-execution. The
predicate binds the output, inputs, model weights, computation identity, and the
arithmetic profile by digest into a canonical manifest whose SHA-256 is the claim's
certificate. Under an order-independent exact arithmetic profile, every honest
re-execution reproduces the output bit-identically on any hardware, so verification
is: re-execute, compare digests.

## Use Cases

-   A client verifies on its own (different) hardware that a provider's inference or
    reduction result is exactly correct.
-   An auditor re-derives a recorded model output months later from the pinned
    weights, inputs, and arithmetic profile.
-   A mixed-hardware fleet proves two nodes computed the same answer to the bit.

## Prerequisites

The in-toto Attestation Framework and the Computation Receipt v0.1 specification
(https://github.com/anomly-labs/computation-receipts — spec, conformance vectors
including required refusals, and reference implementations). Section references
below (§2, §3.2, §6, §6.1, §7, §11) are to that specification.

## Model

The subject is the certified output artifact. `subject[0].digest.sha256` MUST equal
the hex digest in the predicate's `manifest.output.digest`. The predicate is a
complete CR receipt object, carried verbatim.

A receipt is a *claim*, not a verdict. It says: "this computation, on these inputs
and this model, under this arithmetic profile, produced this output". Whether the
claim is true is established only by re-executing the computation and comparing the
output digest; the attestation carries what a verifier needs to do that, never the
result of having done it.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "<NAME>",
    "digest": {"sha256": "<HEX>"}          // MUST equal predicate.manifest.output.digest
  }],

  // Predicate: a CR v0.1 receipt, verbatim.
  "predicateType": "https://in-toto.io/attestation/computation-receipt/v0.1",
  "predicate": {
    "manifest": {
      "cr": "0.1",
      "digest_alg": "sha256",
      "arithmetic": {
        "profile": "<REGISTERED_PROFILE_NAME>",
        "accumulation": "exact" | "float",
        "order_independent": true | false,
        "params": { ... }
      },
      "computation": {"id": "<STRING>", "version": "<STRING>", "graph_digest": "<ALG:HEX>" /* optional */},
      "model":  {"digest": "<ALG:HEX>", "n_tensors": <INT>},   // optional
      "input":  {"digest": "<ALG:HEX>", "n_tensors": <INT>},   // optional
      "output": {"digest": "<ALG:HEX>", "shape": [<INT>, ...]}
    },
    "certificate": "<ALG:HEX>",            // digest(canonical(manifest))
    "meta": { ... }                        // optional, untrusted, NOT covered by the certificate
  }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).
In addition:

-   **The predicate is carried verbatim, and `manifest` is closed.** Consumers MUST NOT
    re-serialize the receipt with their own canonicalization before checking it: the
    certificate is `digest(canonical(manifest))` under CR's canonical form (spec §2), and a
    carrier that re-encodes the manifest breaks the join between the bytes and the
    certificate. Because `certificate` covers the whole `manifest` object, the framework's
    general rules are qualified here as the framework permits ("unless otherwise noted in
    the predicate specification"): consumers MUST NOT drop unrecognized fields inside
    `manifest` before recomputing the certificate (doing so turns a well-formed receipt into
    `MALFORMED`), and producers MUST NOT add extension fields inside `manifest` (doing so
    changes `certificate`, which the framework's own rule that an extension must not
    influence any other field already forbids). Extension fields belong in `meta`, which is
    outside the certificate.
-   **Subject join.** `subject` MUST contain exactly one entry, and `subject[0].digest.sha256`
    MUST equal the hex part of `predicate.manifest.output.digest`; `predicate.manifest.digest_alg`
    MUST be `sha256` for that join to be meaningful. A Statement with more than one subject, or
    whose subject and manifest output digests differ, is invalid.
-   **The integrity check is arithmetic-free.** Any consumer can recompute
    `digest(canonical(manifest))` and compare it to `certificate` without any numeric
    capability. A mismatch means the receipt is `MALFORMED` and MUST be refused before
    anything else is considered.
-   **The presence of a receipt is not an ACCEPT.** A well-formed, correctly signed
    Statement carrying a receipt establishes only that someone claimed this output.
    The only path to `ACCEPT` is the verifier's own re-execution of the computation
    reproducing the **certificate** — i.e. the verifier's independently built manifest
    (model, input, computation, arithmetic profile and output digests) canonicalizes to the
    same certificate (spec §6). Comparing the output digest alone is not sufficient: a
    manifest field the verifier does not model (for example an unexpected
    `computation.graph_digest`) yields a different certificate and MUST be refused, which is
    what makes a forged field detectable (spec §13.1). Policies MUST NOT treat a valid
    signature, a transparency-log inclusion, or a certificate-integrity match as verification
    of the output.
-   **Verdict discipline (spec §6).** A verifier that re-executes returns exactly one
    of `MALFORMED`, `REJECT`, `UNVERIFIABLE`, `ACCEPT`. `UNVERIFIABLE` is returned for
    any receipt whose profile is not `order_independent` — **including when the
    re-executed output happens to match**. Under order-dependent accumulation,
    agreement is coincidence and disagreement is expected; neither is evidence.
    `UNVERIFIABLE` is therefore not a soft `REJECT` and MUST NOT be mapped to either
    `ACCEPT` or `REJECT` by a policy.
-   **The registry is authoritative, not the receipt (spec §6.1).** `arithmetic.profile`,
    `arithmetic.accumulation`, `arithmetic.order_independent` and `arithmetic.params`
    are all written by the prover. For a profile the verifier has in its registry
    (spec §7), the declared `accumulation`, `order_independent` and `params` MUST match
    the registry entry, or the receipt is `MALFORMED`. For a profile the verifier does
    not have in its registry, the verdict MUST be `UNVERIFIABLE`: there is nothing to
    check the order-independence claim against, and `ACCEPT` would be trusting the
    prover rather than verifying.
-   **Policy shape.** The recommended policy against this predicate type is *deny unless
    the profile is registered and order-independent, and the verifier's own
    re-execution returned ACCEPT*. A policy that admits `UNVERIFIABLE` receipts is
    admitting unverified claims and SHOULD say so explicitly.
-   **`meta` is not evidence.** `meta` (operator, host, timestamp, hardware) is excluded
    from the certificate by design so that a receipt verifies identically on any
    machine. Consumers MAY log it and MUST NOT use it in a verification decision.
-   **Version refusal.** A consumer that does not implement the `manifest.cr` version
    MUST refuse the receipt (`MALFORMED`) rather than interpret it.

### Fields

**`manifest`, required** object

> The certified claim. Every field in it is bound by `certificate`.

**`manifest.cr`, required** string

> Computation Receipt format version. `"0.1"` for this predicate version.

**`manifest.digest_alg`, required** string

> Digest algorithm used for every digest in the receipt. `"sha256"` is required for
> v0.1 and for the subject join.

**`manifest.arithmetic`, required** object

> The arithmetic profile the computation ran under. This is what decides whether
> the receipt can ever be verified: only an `order_independent` profile admits
> `ACCEPT`.

**`manifest.arithmetic.profile`, required** string

> Registered profile name (spec §7), e.g. `bposit16-quire256`, `float64`.

**`manifest.arithmetic.accumulation`, required** string

> `"exact"` or `"float"`. MUST match the registry entry for `profile`.

**`manifest.arithmetic.order_independent`, required** boolean

> Whether an honest re-execution anywhere reproduces the output bits. MUST match the
> registry entry for `profile`; a verifier never trusts this field on its own.

**`manifest.arithmetic.params`, required** object

> Profile parameters (for a posit/quire profile: `n`, `es`, `quire_bits`,
> `frac_bits`). MUST match the registry entry for `profile`.

**`manifest.computation`, required** object

> Identifies what was run: `id` (string) and `version` (string). An optional
> `graph_digest` (`<alg>:<hex>`) binds a canonical computation graph so that two
> different computations cannot share a name (spec §11, draft).

**`manifest.model`, optional** object

> Named-collection digest of the model weights (`digest`, `n_tensors`), spec §3.2.
> Absent for computations without weights (e.g. a reduction).

**`manifest.input`, optional** object

> Named-collection digest of the inputs (`digest`, `n_tensors`), spec §3.2.

**`manifest.output`, required** object

> Tensor digest (`digest`, `<alg>:<hex>`) and `shape` (array of int) of the certified
> output. Its digest is the Statement's subject.

**`certificate`, required** string

> `digest(canonical(manifest))`, formatted `<alg>:<hex>`. Binds every manifest field
> at once: a change to the model, inputs, computation, arithmetic or output changes
> the certificate.

**`meta`, optional** object

> Untrusted provenance for humans and logs (operator, timestamp, host, hardware).
> Excluded from the certificate; never used in verification.

## Example

A receipt for one linear layer of a model under the exact `bposit16-quire256`
profile, produced by the reference implementation, wrapped in a Statement whose
subject is the certified output. A verifier with the weights and input (fetched by
digest from an artifact store — the Statement carries claims, not tensors) re-executes
the layer, builds its own manifest from what it ran, and compares the resulting
certificate to `certificate`; `ACCEPT` only on a match, and only because the registry
confirms `bposit16-quire256` is order-independent.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "layer0.output",
      "digest": {
        "sha256": "3bbac7c80c16ed99e71e2baab3f19f4a5ac1f3de67543810c3b7909bd4f6c71b"
      }
    }
  ],
  "predicateType": "https://in-toto.io/attestation/computation-receipt/v0.1",
  "predicate": {
    "manifest": {
      "cr": "0.1",
      "digest_alg": "sha256",
      "arithmetic": {
        "profile": "bposit16-quire256",
        "accumulation": "exact",
        "order_independent": true,
        "params": {"n": 16, "es": 3, "quire_bits": 256, "frac_bits": 96}
      },
      "computation": {"id": "llm.linear.example", "version": "1"},
      "model": {
        "digest": "sha256:acfa461cbb55b60371dc6762bc6971f01fa7e5a94f7de722772d53c613056c8f",
        "n_tensors": 1
      },
      "input": {
        "digest": "sha256:88b952bd5a874452b65edef15e37e3bf5f4e9d186cc8d779c9ed7e5d5bd1251b",
        "n_tensors": 1
      },
      "output": {
        "digest": "sha256:3bbac7c80c16ed99e71e2baab3f19f4a5ac1f3de67543810c3b7909bd4f6c71b",
        "shape": [8]
      }
    },
    "certificate": "sha256:e38d2abf9d876188c61cdb49d46082a6803965afe29265fa6982679681a2828d",
    "meta": {
      "producer": "example-runtime 1.0",
      "host": "node-7",
      "timestamp": "2026-09-07T08:00:00Z"
    }
  }
}
```

The same Statement with `"profile": "float64"` and `"order_independent": false` is
equally well-formed and equally signable; a verifier's answer for it is
`UNVERIFIABLE`, whatever its own re-execution produced.

## Changelog and Migrations

-   v0.1: initial version.
