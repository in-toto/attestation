# Predicate type: AI Agent Action

Type URI: https://in-toto.io/attestation/ai-agent-action/v0.1

Version: 0.1.0

Authors: Elankumaran Srinivasan (@elang2)

Predicate Name: ai-agent-action

## Purpose

This predicate describes actions performed by AI agents through tool-calling
protocols. As AI agents increasingly execute real-world operations autonomously
(file modifications, API calls, infrastructure changes, financial transactions),
organizations need records of what an agent did, when, and with what outcome
that can be verified by third parties without trusting the agent or its runtime.

Existing in-toto predicates cover software supply chain operations (builds,
scans, deployments). AI agent actions represent a new class of automated
operation that produces side effects in external systems but lacks standardized
attestation. This predicate fills that gap by recording tool invocations made
through protocols like the Model Context Protocol (MCP), enabling the same
verification workflows that in-toto provides for build and release pipelines.

## Use Cases

### Compliance auditing for autonomous AI systems

EU AI Act Article 14 requires human oversight capabilities for high-risk AI
systems, which implies the ability to review and verify historical agent actions.
An AI agent that approves expenses, modifies cloud infrastructure, or sends
communications on behalf of a user needs verifiable records that cannot be
fabricated or selectively deleted by the agent itself.

This predicate enables an intermediary (gateway, proxy, or audit sidecar) to
produce signed attestations of each tool call. A compliance team can later
verify the complete history of agent actions without trusting the agent's own
logs.

### AI supply chain integrity

When AI agents participate in software development workflows (code generation,
PR creation, deployment approval), their actions become part of the software
supply chain. This predicate integrates AI agent actions into the same
attestation framework used for builds, scans, and releases, enabling end-to-end
verification policies that span both human and AI contributors.

### Incident investigation

When an AI agent causes an incident, investigators need tamper-evident records
that no participant in the incident could have altered. Application logs written
by the same process whose behavior is under investigation cannot serve this
purpose. A protocol intermediary producing independent attestations separates
the evidence chain from the system under investigation. The hash-chained
structure additionally proves no records were removed between the triggering
event and the investigation.

## Prerequisites

Familiarity with the in-toto Attestation Framework (v1 Statement format),
AI agent architectures (tool-calling patterns), and the Model Context Protocol
(MCP) is helpful but not required to understand this predicate.

## Model

The predicate models a single tool invocation as observed by an intermediary
sitting between an AI client and a tool server. The intermediary does not
execute the tool; it observes the request and response passing through it and
records the metadata.

Key entities:

- **Agent**: The AI system that initiated the tool call (identified by model,
  session, or principal)
- **Tool**: The capability being invoked (identified by name and server)
- **Intermediary**: The entity producing the attestation (identified in the
  statement's signing metadata)
- **Outcome**: Whether the tool call succeeded or failed, its duration, and
  any error classification

The predicate does not include the tool's input arguments or output content by
default, as these may contain sensitive data. The `contentDigest` field allows
optional binding to the actual request/response payloads without including them
inline.

### Subject semantics

The statement `subject` identifies the audit chain to which this attestation
belongs, not a traditional build artifact. The `name` field identifies the
agent session or deployment context. The `digest` contains the SHA-256 hash of
the chain's genesis record. This hash is the same for every attestation in a
given chain, serving as a stable chain identifier that verification policies
can match against. A policy can target "all attestations where subject digest
equals X" to collect the complete history of a specific agent session.

This usage follows the precedent that subjects need not be filesystem artifacts.
The DigestSet specification allows any immutable identifier that a verifier can
independently compute. The genesis hash satisfies this: it is computed once at
chain creation and remains constant for the chain's lifetime.

### Trust boundary

The security of attestations produced by this predicate depends entirely on the
integrity of the intermediary's signing key. A compromised intermediary can
fabricate attestation chains that are indistinguishable from legitimate ones.
Deployments SHOULD use hardware-backed signing keys or threshold signing
schemes to reduce this risk.

The `principal` field in `predicate.agent` is asserted by the intermediary based
on whatever authentication context is available. It is not independently verified
by the attestation framework itself. Consumers SHOULD NOT treat the principal
field as proof of identity without external verification.

This predicate provides tamper evidence against third parties who do not hold
the signing key. It does not provide non-repudiation against the signing entity
itself.

### Chain initialization

The first attestation in a chain uses `"genesis"` as its `previousHash`. To
prevent chain-splicing (an attacker with key access starting a parallel chain),
deployments SHOULD bind the genesis record to a unique deployment identifier
and SHOULD publish chain heads to an append-only transparency log or similar
external witness.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "<agent-session-identifier>",
    "digest": { "sha256": "<genesis-record-hash>" }
  }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "tool_call",
      "protocol": "mcp",
      "method": "tools/call",
      "toolName": "<string>",
      "observedAt": "<RFC 3339 timestamp>",
      "durationMs": "<integer>",
      "success": "<boolean>",
      "errorCode": "<integer or null>",
      "errorClass": "<string or null>"
    },
    "agent": {
      "principal": "<string or null>",
      "model": "<string or null>",
      "sessionId": "<string or null>"
    },
    "upstream": {
      "name": "<string or null>",
      "transport": "<string or null>"
    },
    "chain": {
      "previousHash": "<hex string or 'genesis'>",
      "sequenceNumber": "<integer>"
    },
    "contentDigest": {
      "request": { "sha256": "<hex>" },
      "response": { "sha256": "<hex>" }
    },
    "metadata": {
      "attestorVersion": "<string>",
      "configHash": "<hex or null>"
    }
  }
}
```

### Parsing Rules

This predicate follows the in-toto attestation parsing rules. Summary:

- Consumers MUST ignore unrecognized fields.
- The `predicateType` URI includes the major version number and will always
  change whenever there is a backwards incompatible change.
- Minor version changes are always backwards compatible and monotonic. Such
  changes do not update the `predicateType`.
- Producers MAY add extension fields using field names that are unlikely to
  collide with names used by other producers. Field names SHOULD avoid using
  characters like `.` and `$`.
- Fields marked optional MAY be unset or null, and should be treated
  equivalently. Both are equivalent to empty for object or array values.
- The `predicate.chain` object is OPTIONAL. When absent, the attestation
  stands alone without ordering guarantees.
- The `predicate.contentDigest` object is OPTIONAL. When absent, the
  attestation does not bind to specific request/response payloads.

### Fields

#### `predicate.action`

| Field        | Type              | Required | Description                                               |
| ------------ | ----------------- | -------- | --------------------------------------------------------- |
| `type`       | string            | Yes      | Action type. Currently only `"tool_call"` is defined.     |
| `protocol`   | string            | Yes      | Protocol through which the action was observed. e.g. `"mcp"` |
| `method`     | string            | Yes      | Protocol method. e.g. `"tools/call"`                      |
| `toolName`   | string            | Yes      | Name of the tool as invoked by the agent                  |
| `observedAt` | string (RFC 3339) | Yes      | When the action was observed by the intermediary          |
| `durationMs` | integer           | Yes      | Wall-clock duration of the tool execution in milliseconds |
| `success`    | boolean           | Yes      | Whether the tool call completed without error             |
| `errorCode`  | integer or null   | No       | Protocol-level error code (e.g., JSON-RPC error code)     |
| `errorClass` | string or null    | No       | Human-readable error classification                       |

#### `predicate.agent`

| Field       | Type           | Required | Description                                         |
| ----------- | -------------- | -------- | --------------------------------------------------- |
| `principal` | string or null | No       | Asserted identity of the requesting agent or user (not independently verified by the attestation) |
| `model`     | string or null | No       | AI model identifier (e.g., `"gpt-4o-2025-03-15"`)  |
| `sessionId` | string or null | No       | Session or conversation identifier                  |

#### `predicate.upstream`

| Field       | Type           | Required | Description                                      |
| ----------- | -------------- | -------- | ------------------------------------------------ |
| `name`      | string or null | No       | Name of the upstream tool server (assigned by intermediary configuration for multi-server disambiguation) |
| `transport` | string or null | No       | Transport type (e.g., `"stdio"`, `"streamable-http"`) |

#### `predicate.chain`

| Field            | Type              | Required | Description                                                    |
| ---------------- | ----------------- | -------- | -------------------------------------------------------------- |
| `previousHash`   | string            | Yes      | SHA-256 hex of the preceding attestation. First = `"genesis"`  |
| `sequenceNumber` | integer           | Yes      | Monotonically increasing counter within this chain. Verifiers SHOULD use this for ordering; timestamps are not guaranteed monotonic due to clock drift. |

#### `predicate.contentDigest`

| Field      | Type                          | Required | Description                                    |
| ---------- | ----------------------------- | -------- | ---------------------------------------------- |
| `request`  | DigestSet                     | No       | Digest of the canonical tool call request      |
| `response` | DigestSet                     | No       | Digest of the tool call response content       |

#### `predicate.metadata`

| Field            | Type           | Required | Description                                        |
| ---------------- | -------------- | -------- | -------------------------------------------------- |
| `attestorVersion`| string         | Yes      | Version of the attestation-producing software      |
| `configHash`     | string or null | No       | Hash of the attestor's policy configuration        |

## Example

A complete attestation for an AI agent invoking the `create_pull_request` tool
on a GitHub MCP server, as observed by an intermediary gateway:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "session:agent-workspace-4f2a",
    "digest": {
      "sha256": "b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c"
    }
  }],
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "tool_call",
      "protocol": "mcp",
      "method": "tools/call",
      "toolName": "create_pull_request",
      "observedAt": "2026-08-18T14:32:01.998Z",
      "durationMs": 1247,
      "success": true,
      "errorCode": null,
      "errorClass": null
    },
    "agent": {
      "principal": "user:alice@example.com",
      "model": "gpt-4o-2025-03-15",
      "sessionId": "conv_abc123"
    },
    "upstream": {
      "name": "github-server",
      "transport": "stdio"
    },
    "chain": {
      "previousHash": "3fdba35f04dc8c462986c992bcf875546257113072a909c162f7e470e581e278",
      "sequenceNumber": 42
    },
    "contentDigest": {
      "request": { "sha256": "7d865e959b2466918c9863afca942d0fb89d7c9ac0c99bafc3749504ded97730" },
      "response": { "sha256": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" }
    },
    "metadata": {
      "attestorVersion": "mcp-audit-gateway/0.1.0",
      "configHash": "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
    }
  }
}
```

## Changelog and Migrations

This is the initial version (v0.1) of the AI Agent Action predicate. No
migrations are required.
