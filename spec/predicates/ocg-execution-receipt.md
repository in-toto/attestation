# Predicate type: OCG Execution Receipt

Type URI: https://ainumbers.co/attestations/ocg-execution-receipt/v0.1

Version: 0.1

Predicate Name: OCG Execution Receipt

## Purpose

An [OpenChainGraph] (OCG) artifact is a self-contained, hash-anchored record of a single deterministic
computation: the inputs it ran on, the output it produced, and a re-computable `execution_hash` binding
the two. This predicate re-expresses that artifact as the `predicate` of an [in-toto Statement], so that
generic in-toto / SLSA tooling — verifiers, transparency-log ingesters, and policy engines built against
the in-toto ecosystem — can consume an OCG receipt without bespoke parsing.

It does **not** replace the OCG artifact as the source of truth; it is a lossless re-expression layered
on top, the same posture OCG already takes toward its other export profiles (W3C Verifiable Credentials,
SD-JWT, PROV-DM). The normative predicate specification is published at the Type URI above.

## Use Cases

An autonomous agent runs a compliance or payment-routing workflow client-side or server-side, and must
hand a counterparty (often another agent) a verifiable "I computed exactly this output from exactly these
inputs" receipt — including, optionally, a succinct zero-knowledge proof that a specific kernel produced
the output and independent timestamp/transparency evidence of when.

Existing predicates do not cover this:

- **[SLSA Provenance]** attests how a software *artifact was built* by a build system; an OCG receipt
  attests the *result of running* a versioned deterministic function on caller-supplied inputs, with a
  re-computable input→output hash. The subject is a decision output, not a build product.
- **[SCAI]** expresses evidence-based *attributes* of an artifact or supply-chain operation; it does not
  carry the recomputable-hash execution record plus optional zkVM execution proof this predicate needs.
- **[Simple Verification Result]** / **[Test Result]** record a pass/fail evaluation against a policy or
  test; an OCG receipt records the *computation itself* (inputs, output, deterministic binding), from
  which such a verdict could later be derived.

The distinguishing content is the recomputable `executionHash` over the `{policyParameters,
outputPayload}` pair, the multi-parent chain lineage, and the optional zkVM `computeProof` /
`anchorBindings` — none of which the above predicates model.

## Prerequisites

- The [in-toto Attestation Framework] (in-toto Statement v1).
- The [OpenChainGraph specification] (artifact envelope, `execution_hash` derivation, and the optional
  §16 signature / §17 kernel-identity / §18 compute-proof / §20 anchor-binding layers), published at
  `https://ainumbers.co/`. This predicate is a mechanical re-expression of that envelope; it defines no
  computation of its own.

## Model

The producing functionary is an OCG compute node — a versioned, deterministic function identified by
`toolId` / `toolVersion`, executed in one of three modes (`server`, `browser`, `wasm-vm`). Each run
emits one artifact; this predicate expresses one such artifact. Nodes may be chained into a DAG, in
which case a child receipt records its parents' hashes (`chain`), giving downstream verifiers the full
lineage without a separate materials model. The predicate is relevant to any consumer that must verify a
computation's integrity — agentic-commerce counterparties, compliance-decision auditors, and
transparency-log operators.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    // "<mandateType>/<toolId>", e.g. "payment_mandate/ap2-fee-route-optimizer"
    "name": "...",
    // the OCG execution_hash, sha256: prefix stripped (in-toto digest maps are algorithm-keyed)
    "digest": { "sha256": "..." }
  }],

  // Predicate:
  "predicateType": "https://ainumbers.co/attestations/ocg-execution-receipt/v0.1",
  "predicate": {
    "chaingraphVersion": "0.4.0",
    "computeMode": "server",
    "mandateType": "payment_mandate",
    "toolId": "ap2-fee-route-optimizer",
    "toolVersion": "1.4.0",
    "generatedAt": "2026-07-10T14:32:07.118Z",
    "buildType": "https://ainumbers.co/chaingraph/context/v0.2#WebCryptoSHA256",
    "executionHash": "2255b6e8...",
    "chain": { "parent_hashes": [], "parent_tool_ids": [], "chain_depth": 0 },
    "policyParameters": { /* the run's inputs — half the executionHash preimage */ },
    "outputPayload": { /* the run's output — the other half */ },
    "complianceFlags": [],
    // OPTIONAL fields, omitted when absent on the source artifact:
    "dcp1Profile": "ocg-deterministic-compute",
    "kernelIdentity": { "kernelDigest": "...", "buildType": "...", "sourceRef": "..." },
    "computeProof": null,
    "auditSignature": { "type": "DataIntegrityProof", "cryptosuite": "eddsa-jcs-2022", /* ... */ },
    "anchorBindings": [ { "type": "rfc3161-tst", /* ... */ } ]
  }
}
```

### Parsing Rules

This predicate opts in to the framework's [standard parsing rules], including the monotonic principle:
consumers MUST ignore unrecognized fields, and absence of an OPTIONAL field MUST NOT be read as a
negative assertion (an absent `computeProof` means "no proof provided," not "proof failed"). The
predicate is versioned solely through the `predicateType` URI path (`.../v0.1`); a breaking change to
field semantics ships as a new version URI.

`executionHash` is authoritative and self-checking: a verifier MUST be able to recompute it as the
SHA-256 of the JCS-canonical (RFC 8785) serialization of `{policy_parameters, outputPayload}` and get
the value in both `predicate.executionHash` and `subject[0].digest.sha256`. All predicate fields are
verbatim or renamed passthroughs of the source OCG artifact; the predicate mints no new claims.

### Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `chaingraphVersion` | string | Yes | OCG schema version; `"0.4.0"` under the frozen v0.4 schema. |
| `computeMode` | string | Yes | `"server"` \| `"browser"` \| `"wasm-vm"`. |
| `mandateType` | string | Yes | Producer's decision-domain label (open taxonomy, not a fixed enum). |
| `toolId` | string | Yes | Kebab-case identifier of the deterministic node. |
| `toolVersion` | string | Yes | Node version. |
| `generatedAt` | string (RFC 3339) | Yes | When the receipt was generated. |
| `buildType` | string (URI) | Yes | Hash-algorithm identifier for `executionHash`. |
| `executionHash` | string (hex) | Yes | SHA-256 over JCS-canonical `{policyParameters, outputPayload}`. Equals `subject[0].digest.sha256`. |
| `chain` | object | Yes | `{ parent_hashes[], parent_tool_ids[], chain_depth }` — DAG lineage. Empty/zero for a root run. |
| `policyParameters` | object | Yes | The run's inputs; half the `executionHash` preimage. |
| `outputPayload` | object | Yes | The run's output; the other half of the preimage. |
| `complianceFlags` | array | Yes | Producer decision metadata (not regulatory certification). |
| `dcp1Profile` | string | No | `"ocg-deterministic-compute"` when the node is in the OCG Deterministic Compute profile scope; omitted otherwise. |
| `kernelIdentity` | object | No | `{ kernelDigest, buildType, sourceRef? }` — an advisory published claim of which kernel source ran; **not** a proof of execution. |
| `computeProof` | object \| null | No | zkVM proof that the *named* kernel produced this output; enables succinct verification without re-execution. |
| `auditSignature` | object | No | A `DataIntegrityProof` (`eddsa-jcs-2022`) over the artifact when signed. Signing is opt-in and de-anonymizing. |
| `anchorBindings` | array | No | Transparency-log / timestamp evidence (`rfc3161-tst`, `opentimestamps`, `c2sp-tlog-proof-v1`, `scitt-receipt-rfc9942`). |

**What the attestation asserts** (from `executionHash` alone): the recorded `outputPayload` is exactly
what a deterministic function produced from the recorded `policyParameters`; post-hoc tampering with
either is detectable. **Only if the corresponding OPTIONAL field is present** does it additionally assert
a named-kernel claim (`kernelIdentity`, advisory), a proof of which kernel ran (`computeProof`), a
key binding (`auditSignature`), or independent time evidence (`anchorBindings`). It **never** asserts
that the inputs are true or non-malicious, that the producer is authorized or compliant with any
external regime, the real-world legal identity of a signer, or that the kernel is bug-free.

## Example

A complete in-toto Statement (server-mode payment-routing run, signed, RFC 3161 timestamped):

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    { "name": "payment_mandate/ap2-fee-route-optimizer",
      "digest": { "sha256": "2255b6e807e739e23073c256e039f83f009c17be3fd69bfce58994a0cfe29b49" } }
  ],
  "predicateType": "https://ainumbers.co/attestations/ocg-execution-receipt/v0.1",
  "predicate": {
    "chaingraphVersion": "0.4.0",
    "computeMode": "server",
    "mandateType": "payment_mandate",
    "toolId": "ap2-fee-route-optimizer",
    "toolVersion": "1.4.0",
    "generatedAt": "2026-07-10T14:32:07.118Z",
    "buildType": "https://ainumbers.co/chaingraph/context/v0.2#WebCryptoSHA256",
    "executionHash": "2255b6e807e739e23073c256e039f83f009c17be3fd69bfce58994a0cfe29b49",
    "chain": { "parent_hashes": [], "parent_tool_ids": [], "chain_depth": 0 },
    "policyParameters": {
      "execution_backend": "server",
      "input_parameters": {
        "route_candidates": [ { "rail": "a2a", "fee_bps": 12 }, { "rail": "rtp", "fee_bps": 18 } ],
        "settlement_amount_usd": 250000
      }
    },
    "outputPayload": { "selected_rail": "a2a", "effective_fee_bps": 12, "estimated_fee_usd": 300 },
    "complianceFlags": [],
    "dcp1Profile": "ocg-deterministic-compute",
    "kernelIdentity": {
      "kernelDigest": "sha256:4a1f9e3d6c2b8e7a0f5d9c1b3a6e8f2d7c4b9a1e6f3d8c2b5a9e1f4d7c3b8a6e",
      "buildType": "https://ainumbers.co/chaingraph/context/v0.6#SourceDigestSHA256",
      "sourceRef": "https://github.com/PostOakLabs/ainumbers/blob/a1b2c3d/chaingraph/kernels/ap2-fee-route-optimizer.kernel.mjs"
    },
    "computeProof": null,
    "auditSignature": {
      "type": "DataIntegrityProof",
      "cryptosuite": "eddsa-jcs-2022",
      "proofPurpose": "assertionMethod",
      "verificationMethod": "did:key:z6MkrJVnaZkeFzdQyMZu1cgjaeW7MoFaNKg4tLCf4g3kEhWh",
      "created": "2026-07-10T14:32:07.201Z",
      "proofValue": "z3fRQx8mNcVYh2pQd9sT4kLwJb6nXe1oZaCf5uHgYtRvKm9DwLp2VjXhNqE7BsGtMcRz"
    },
    "anchorBindings": [
      { "type": "rfc3161-tst",
        "anchored_hash": "sha256:2255b6e807e739e23073c256e039f83f009c17be3fd69bfce58994a0cfe29b49",
        "log_origin": "https://freetsa.org/tsr",
        "gen_time": "2026-07-10T14:32:09Z" }
    ]
  }
}
```

<!-- Reference links -->
[OpenChainGraph]: https://ainumbers.co/
[OpenChainGraph specification]: https://ainumbers.co/
[in-toto Statement]: ../v1/statement.md
[in-toto Attestation Framework]: ../v1/README.md
[standard parsing rules]: ../v1/README.md#parsing-rules
[SLSA Provenance]: https://slsa.dev/provenance
[SCAI]: scai.md
[Simple Verification Result]: svr.md
[Test Result]: test-result.md
