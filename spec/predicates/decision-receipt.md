# Predicate type: Decision Receipt

Type URI: https://veritasacta.com/attestation/decision-receipt/v0.1

Version: 0.1

Predicate Name: DecisionReceipt

## Purpose

AI agents and autonomous systems make consequential decisions: invoking tools,
accessing APIs, modifying state, approving transactions, and controlling
physical actuators. There is currently no standard in-toto predicate for
attesting to these decisions.

This predicate captures what decision was made, by whom, under what policy,
and with what inputs and outputs. Each attestation is independently verifiable
and can be hash-chained to form a tamper-evident audit trail.

The predicate is designed for two domains:

- **Software agent governance**: AI agents invoking tools via the Model Context
  Protocol (MCP), LangChain, Pydantic AI, Vercel AI SDK, or similar frameworks.
  Each tool call is attested with the policy that governed it and the decision
  outcome.
- **Physical sensor attestation**: IoT devices producing signed environmental
  readings (temperature, shock, GPS, light) from hardware secure elements. Each
  reading is attested with the policy that evaluated it.

Both domains produce the same predicate structure, enabling a single verifier
to process attestations from software agents and physical devices
interchangeably.

## Use Cases

**Post-incident audit of an AI agent**: An AI agent ran a deployment pipeline.
After an outage, the operator needs to prove which tools the agent called, what
the inputs were, and whether each call was authorized by the active policy. The
decision receipt chain provides this proof as independently verifiable
attestations anchored in a transparency log.

**Cold chain supply chain verification**: A sensor device attached to a
pharmaceutical shipment records temperature every 5 minutes. At destination,
the inspector verifies the complete chain of readings. A temperature excursion
that was denied by the on-device policy produces a decision receipt proving
the policy gate held.

**Regulatory compliance evidence**: For EU AI Act Article 12 (record-keeping)
and SLSA build provenance, decision receipts provide the cryptographic evidence
layer that standard logging cannot: signed, chained, and verifiable without
trusting the operator.

## Prerequisites

- [in-toto Attestation Framework v1](https://github.com/in-toto/attestation/tree/main/spec/v1)
- Ed25519 (RFC 8032) or ES256 (ECDSA P-256) signing capability
- JCS canonicalization (RFC 8785) for deterministic serialization of the
  predicate payload before signing

## Model

The decision receipt predicate models a single access control decision:

- **Principal**: The agent or device making the decision (the signer).
- **Action**: The tool invoked or sensor reading captured.
- **Policy**: The authorization policy that governed the decision.
- **Decision**: The outcome (allow, deny, alert).
- **Evidence**: Input/output hashes and chain links.

Steps:

1. An agent or device is about to perform an action (tool call, sensor reading).
2. The active policy is evaluated against the action context.
3. The decision (allow/deny/alert) is recorded.
4. The decision is signed by the principal's key.
5. The signed receipt is chained to the previous receipt by hash.
6. The attestation may be anchored in a transparency log (e.g., Sigstore Rekor).

Functionaries:

- **Signer**: The agent host (e.g., protect-mcp) or device secure element
  (e.g., ATECC608B).
- **Verifier**: Any party with the signer's public key. Verification is offline.
- **Log**: Optional transparency log (Rekor) for temporal anchoring.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    // The tool or action being attested
    "name": "<tool_name or sensor_type>",
    "digest": {
      "sha256": "<hex digest of the canonicalized action input>"
    }
  }],

  // Predicate:
  "predicateType": "https://veritasacta.com/attestation/decision-receipt/v0.1",
  "predicate": {
    // Decision outcome
    "decision": "<allow | deny | alert>",
    "reason": "<human-readable reason for the decision>",

    // Policy that governed the decision
    "policyId": "<policy identifier>",
    "policyDigest": {
      "sha256": "<hex digest of the policy file at evaluation time>"
    },

    // Temporal and ordering fields
    "issuedAt": "<RFC 3339 timestamp>",
    "sequence": 1,
    "previousReceiptDigest": {
      "sha256": "<hex digest of the previous receipt's canonical form>"
    },

    // Signer identity
    "issuerId": "<issuer identifier (key fingerprint or DID)>",

    // Optional output evidence
    "outputDigest": {
      "sha256": "<hex digest of the action output>"
    },

    // Optional metadata
    "metadata": {
      "deviceId": "<hardware device identifier, if physical>",
      "sessionId": "<agent session identifier, if software>",
      "framework": "<protect-mcp | pydantic-ai | llamaindex | vercel-ai | ...>"
    }
  }
}
```

### Parsing Rules

The decision receipt predicate follows the in-toto Attestation Framework
[standard parsing rules](/spec/v1/README.md#parsing-rules) with the following
additions:

- **Versioning**: The `predicateType` URI includes a version suffix (`/v0.1`).
  Consumers MUST check the version and reject unknown major versions.
- **Unknown fields**: Consumers SHOULD ignore unknown fields in the `predicate`
  and `metadata` objects. This allows forward-compatible extensions.
- **Chain validation**: If `previousReceiptDigest` is present and non-null,
  the consumer SHOULD verify it matches the digest of the chronologically
  preceding attestation.
- **Canonicalization**: The `predicate` object is canonicalized with JCS
  (RFC 8785) before digest computation. All digest values in `DigestSet`
  fields use lowercase hexadecimal encoding.

### Fields

**Statement-level fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `subject[0].name` | string | Yes | The tool name (e.g., "Bash", "read_file") or sensor type (e.g., "temperature") |
| `subject[0].digest` | DigestSet | Yes | Digest of the JCS-canonicalized action input |

**Predicate fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `decision` | string | Yes | One of: "allow", "deny", "alert" |
| `reason` | string | No | Human-readable explanation of the decision |
| `policyId` | string | Yes | Identifier for the policy that governed the decision |
| `policyDigest` | DigestSet | No | Digest of the policy file at evaluation time |
| `issuedAt` | Timestamp | Yes | RFC 3339 UTC timestamp of the decision |
| `sequence` | integer | Yes | Monotonically increasing sequence number within a chain |
| `previousReceiptDigest` | DigestSet | No | Digest of the previous receipt. Null for the first receipt in a chain. |
| `issuerId` | string | Yes | Identifier for the signing entity |
| `outputDigest` | DigestSet | No | Digest of the action output, if available |
| `metadata` | object | No | Additional context. Known fields: `deviceId`, `sessionId`, `framework` |

## Example

### Software agent tool call (allow)

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "Bash",
    "digest": {
      "sha256": "a3f8c91e4b7d2f6a1e8b3c5d9f0a2b4e6c8d0f1a3b5c7d9e1f0a2b4c6d8e0f"
    }
  }],
  "predicateType": "https://veritasacta.com/attestation/decision-receipt/v0.1",
  "predicate": {
    "decision": "allow",
    "reason": "all parameters within range",
    "policyId": "development-safe",
    "policyDigest": {
      "sha256": "b7e2f4a6c8d0e1f3a5b7c9d1e3f5a7b9c1d3e5f7a9b1c3d5e7f9a1b3c5d7e9"
    },
    "issuedAt": "2026-04-16T10:30:00.000Z",
    "sequence": 1,
    "previousReceiptDigest": null,
    "issuerId": "sb:agent:4437ca56815c",
    "outputDigest": {
      "sha256": "c9d1e3f5a7b9c1d3e5f7a9b1c3d5e7f9a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9"
    },
    "metadata": {
      "sessionId": "sess_abc123",
      "framework": "protect-mcp"
    }
  }
}
```

### Physical sensor reading (deny due to temperature excursion)

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "temperature",
    "digest": {
      "sha256": "d1e3f5a7b9c1d3e5f7a9b1c3d5e7f9a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1"
    }
  }],
  "predicateType": "https://veritasacta.com/attestation/decision-receipt/v0.1",
  "predicate": {
    "decision": "deny",
    "reason": "temp 22.4C > 18.0C limit",
    "policyId": "cold-chain-wine-premium",
    "policyDigest": {
      "sha256": "e3f5a7b9c1d3e5f7a9b1c3d5e7f9a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3"
    },
    "issuedAt": "2026-04-10T18:00:00.000Z",
    "sequence": 5,
    "previousReceiptDigest": {
      "sha256": "f5a7b9c1d3e5f7a9b1c3d5e7f9a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5"
    },
    "issuerId": "seal:SB-SEAL-001",
    "metadata": {
      "deviceId": "SB-SEAL-001",
      "framework": "scopeblind-seal"
    }
  }
}
```

### Session summary (anchored in Rekor)

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "agent-session",
    "digest": {
      "sha256": "a7b9c1d3e5f7a9b1c3d5e7f9a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7"
    }
  }],
  "predicateType": "https://veritasacta.com/attestation/decision-receipt/v0.1",
  "predicate": {
    "decision": "allow",
    "reason": "session complete: 47 tool calls, 0 denials",
    "policyId": "development-safe",
    "issuedAt": "2026-04-16T11:45:00.000Z",
    "sequence": 47,
    "previousReceiptDigest": {
      "sha256": "b9c1d3e5f7a9b1c3d5e7f9a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9"
    },
    "issuerId": "sb:agent:4437ca56815c",
    "metadata": {
      "sessionId": "sess_abc123",
      "framework": "protect-mcp",
      "rekorLogIndex": 1315193837,
      "rekorUuid": "108e9186e8c5677a..."
    }
  }
}
```

## Relationship to SLSA Provenance

Decision receipts and SLSA Provenance attest to different properties of the
same subject. SLSA Provenance attests to *how an artifact was produced* (build
inputs, steps, and observed runtime behavior) and is signed by the builder
platform identity. A Decision Receipt attests to *what the policy-enforcement
layer authorized at a specific call* and is signed by the supervisor identity
that runs the policy gate. These are distinct trust domains. Cross-signing
under the builder's key would obscure the supervisor's identity in downstream
verification.

The two compose via the `ResourceDescriptor` reference pattern. A SLSA
Provenance attestation records that Decision Receipt attestations exist for
the same subject by including a byproduct entry with the receipt attestation's
digest, URI, and predicate type. The builder does not cross-sign the receipt
content; it records its existence.

Example byproduct entry in a SLSA provenance (adapted from the
[`agent-commit/v0.2` build type](https://refs.arewm.com/agent-commit/v0.2)):

```json
{
  "name": "decision-receipts",
  "digest": { "sha256": "a8f3c9d2e1b7465f..." },
  "uri": "oci://registry/org/agent-session/run-xyz/receipts:sha256-a8f3c9d2",
  "annotations": {
    "predicateType": "https://veritasacta.com/attestation/decision-receipt/v0.1",
    "signerRole": "supervisor-hook",
    "chainLength": 47,
    "genesisReceiptHash": "sha256:a8f3c9d2e1b7465f",
    "finalReceiptHash":   "sha256:e4d61f7a09b8cd34"
  }
}
```

A consumer fetching both:

1. Verifies the SLSA provenance DSSE signature against the builder identity.
2. Fetches the receipt attestation at the referenced URI, checks its digest
   matches the byproduct entry, then verifies its DSSE signature against the
   supervisor identity named by `issuerId`.
3. Cross-references the receipt's subject against the SLSA provenance subject
   and interprets the chain using this predicate's semantics.

The `issuerId` in this predicate and the `signerRole` annotation in the SLSA
byproduct are complementary: `issuerId` is the concrete identity (key
fingerprint or DID) that signed the receipt, while `signerRole` is the logical
role of that identity relative to the builder.

## Changelog and Migrations

### v0.1 (initial)

- Initial predicate definition.
- Supports software agent tool calls and physical sensor readings.
- Chain integrity via `previousReceiptDigest`.
- Compatible with Sigstore Rekor anchoring via DSSE envelope.
- Composes with SLSA Provenance via `ResourceDescriptor` references in
  byproducts; the builder records the receipt attestation's digest and URI
  without cross-signing its content.

## References

- [IETF draft-farley-acta-signed-receipts](https://datatracker.ietf.org/doc/draft-farley-acta-signed-receipts/) -- Receipt wire format
- [RFC 8032](https://datatracker.ietf.org/doc/html/rfc8032) -- Ed25519 digital signatures
- [RFC 8785](https://datatracker.ietf.org/doc/html/rfc8785) -- JCS canonicalization
- [agent-commit build type](https://refs.arewm.com/agent-commit/v0.2) -- SLSA Provenance build type for AI-agent-produced commits; references this predicate via `ResourceDescriptor` in byproducts
- [protect-mcp](https://www.npmjs.com/package/protect-mcp) -- Reference implementation (npm, 10K+ monthly downloads)
- [@veritasacta/verify](https://www.npmjs.com/package/@veritasacta/verify) -- Offline verification CLI
- [Sigstore Rekor integration](https://github.com/sigstore/rekor/issues/2798) -- Transparency log anchoring (working PoC)
- [SLSA-for-agents discussion](https://github.com/slsa-framework/slsa/issues/1594) -- Composition of build provenance, agent identity, and decision receipts
- [Microsoft Agent Governance Toolkit](https://github.com/microsoft/agent-governance-toolkit/pull/667) -- Enterprise consumer (merged)
- [AWS Cedar for Agents](https://github.com/cedar-policy/cedar-for-agents/pull/64) -- Policy engine WASM bindings (merged)
