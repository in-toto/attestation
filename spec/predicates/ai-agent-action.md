# Predicate type: AI Agent Action

Type URI: https://in-toto.io/attestation/ai-agent-action/v0.1

Version: 0.1.0

Authors: Elankumaran Srinivasan (@elang2)

Predicate Name: ai-agent-action

## Purpose

This predicate describes actions performed by AI agents through tool-calling
protocols. As AI agents increasingly execute real-world operations autonomously
(file modifications, API calls, infrastructure changes, financial transactions),
organizations need cryptographically verifiable records of what an agent did,
when, and with what outcome.

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

When an AI agent causes an incident (deleted a production database, sent an
unauthorized email, approved an invalid transaction), investigators need
tamper-evident records. Application logs written by the same process whose
behavior is under investigation create a conflict of interest. Third-party
attestation from a protocol intermediary provides independent evidence.

### Tail truncation detection

A hash-chained log provides prefix integrity (modifying any record breaks the
chain from that point forward), but an adversary with storage access can
silently truncate records from the tail without breaking existing hash linkage.
Checkpoint records address this. The attestor periodically emits a checkpoint
containing the current chain head hash, a checkpoint sequence counter, and the
total record count. The consumer externalizes one checkpoint to a store the
attestor cannot overwrite (transparency log, separate cloud account, HSM). On
verification, any records missing after the externalized checkpoint's position
are detectable.

Limits: a gap between checkpoints proves records were dropped but cannot
distinguish adversarial truncation from a crash (detection, not attribution).
Records emitted after the most recently externalized checkpoint remain in an
unprotected window until the next checkpoint lands.

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
- **Tool**: The capability being invoked (identified by name and namespace)
- **Intermediary**: The entity producing the attestation (identified in the
  statement's signing metadata)
- **Outcome**: Whether the tool call succeeded or failed, its duration, and
  any error classification

The predicate does not include the tool's input arguments or output content by
default, as these may contain sensitive data. The `contentDigest` field allows
optional binding to the actual request/response payloads without including them
inline.

### Field provenance model

Fields in this predicate have different trust characteristics depending on who
asserts them:

- **Witnessed by intermediary**: `timestamp`, `durationMs`, `success`,
  `errorCode`, `toolName`, `method`, `previousHash` — the intermediary
  directly measures or extracts these from the protocol messages it forwards
  and the chain state it maintains.
- **Asserted by client (declared)**: `model`, `sessionId`, `turnId`,
  `invocationReason` — provided by the AI client. The intermediary cannot
  independently verify these. Consumers SHOULD treat them as claims, not facts.
- **Deployment-dependent**: `principal` — observed when the intermediary
  authenticates the client connection (e.g., from an OAuth token or mTLS
  certificate), but declared under bare stdio transports where no
  authentication layer exists. The `parties` array disambiguates per deployment.

The `parties` array makes provenance machine-readable. Each entry identifies
an entity and the fields it asserts:

```json
[
  { "party": "gateway", "role": "witness", "scope": ["id", "timestamp", "method", "toolName", "..."] },
  { "party": "client", "role": "asserter", "scope": ["aiInvocation"] }
]
```

Scope entries name fields of the underlying signed audit record (the JSONL line
written by the gateway), which may include fields like `id` that do not appear
in the predicate schema. This is intentional: the scope describes what the
party's signature covers at the implementation layer.

A field present in the record but not covered by any party's `scope` MUST be
treated as having unknown provenance. Verification policies SHOULD reject or
flag records where security-relevant fields have no attesting party.

## Schema

The predicate defines three record types that share a chain:

### Action record (tool_call)

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "<agent-session-identifier>",
    "digest": { "sha256": "<chain-hash-of-underlying-audit-record>" }
  }],
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "tool_call",
      "protocol": "mcp",
      "method": "tools/call",
      "toolName": "<string>",
      "namespace": "<string or null>",
      "timestamp": "<RFC 3339 timestamp>",
      "durationMs": "<integer>",
      "success": "<boolean>",
      "errorCode": "<integer or null>",
      "errorClass": "<string or null>"
    },
    "agent": {
      "principal": "<string or null>",
      "model": "<string or null>",
      "sessionId": "<string or null>",
      "turnId": "<string or null>",
      "invocationReason": "<string or null>"
    },
    "parties": [
      { "party": "<string>", "role": "<'witness' | 'asserter'>", "scope": ["<field-name>", "..."] }
    ],
    "upstream": {
      "name": "<string or null>",
      "transport": "<string or null>"
    },
    "chain": {
      "previousHash": "<hex string; 'genesis' for the very first record only>"
    },
    "contentDigest": {
      "request": { "sha256": "<hex>" },
      "response": { "sha256": "<hex>" }
    },
    "extensions": {},
    "metadata": {
      "attestorVersion": "<string>",
      "configHash": "<hex or null>"
    }
  }
}
```

### Checkpoint record

```jsonc
{
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "checkpoint",
      "timestamp": "<RFC 3339 timestamp>"
    },
    "checkpoint": {
      "sequence": "<integer>",
      "recordCount": "<integer>",
      "previousHash": "<hex string>"
    },
    "parties": [...]
  }
}
```

### Chain break record

```jsonc
{
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "chain_break",
      "timestamp": "<RFC 3339 timestamp>"
    },
    "chainBreak": {
      "reason": "<string>",
      "priorHead": "<hex string or null>",
      "priorSequence": "<integer or null>",
      "priorRecordCount": "<integer or null>"
    }
  }
}
```

### Canonicalization

This predicate uses two distinct canonical forms for different purposes.

#### Signing canonical form (chain integrity + attestation signature)

The signing canonical form is a **tuple-array**: an ordered array of
`[field-name, value]` pairs, serialized with `JSON.stringify`. The field order
is fixed by the implementation (not sorted alphabetically — the order matches
the type definition). Nested structured values within optional fields use
**M/L type tagging** to prevent cross-type collisions:

- Objects become `["M", [[key, value], ...]]` with keys sorted by UTF-16
  code unit values (ECMAScript default sort).
- Arrays become `["L", [value, ...]]`.
- Scalars (string, number, boolean, null) pass through untagged.

This form is used to compute the bytes the attestor signs (HMAC-SHA256 or
Ed25519). It covers a deterministic subset of fields in fixed order.

The **chain hash** (`previousHash` for the next record) is different: it is
`SHA-256(JSON.stringify(record))` — the SHA-256 of the complete JSON-serialized
audit record as written to the JSONL log, including the attestation signature
field. This intentionally differs from the signing canonical form: the
signature covers fields in deterministic order; the chain hash covers the full
serialized bytes including the signature, creating a hash-then-sign-then-chain
layering.

The first record in a chain carries `previousHash` set to the literal string
`"genesis"`.

Constraints on both forms:
- Numbers MUST be safe integers (absolute value < 2^53, per RFC 7493 §2.2).
- Strings MUST NOT contain unpaired UTF-16 surrogates.

Reference implementation and 90+ cross-language conformance vectors (JS +
Python): https://github.com/elang2/mcp-audit-gateway/tree/main/test/vectors

#### Content digest form (payload binding)

Content digests (`contentDigest.request` and `contentDigest.response`) bind
arbitrary MCP tool call payloads to the attestation without inlining them.

These use **RFC 8785 (JSON Canonicalization Scheme)** directly: keys sorted
lexicographically, numbers serialized per IEEE 754 / ECMAScript rules, no
whitespace. This is necessary because MCP payloads can contain floating-point
values that the signing canonical form rejects by design.

The digest is `SHA-256(JCS(payload))` where `JCS` is the RFC 8785
canonicalization and `payload` is the JSON-RPC `params` object (for requests)
or `result` object (for responses).

Using RFC 8785 for content digests aligns with the registry-level content
digest scheme proposed in in-toto/attestation#570.

Implementations MUST NOT mix the two forms: signing fields use the tuple-array
form; content payloads use JCS.

### Genesis and chain continuity

The literal string `"genesis"` appears as `previousHash` exactly once in a
chain's lifetime: on the very first record. Its chain hash (SHA-256 of its
full JSON serialization as written) becomes the anchor for the chain.

If a chain is broken and restarted (crash, forced rotation, state corruption),
the attestor MUST emit a signed `chain_break` record. The record after the
break carries the break record's chain hash as its `previousHash` — NOT
`"genesis"`. This binds the discontinuity evidence into the new chain so that
discarding the break record would break hash linkage. A consumer that
encounters a second `"genesis"` anywhere in the log SHOULD treat it as
evidence of adversarial chain replacement.

Chain continuity is preserved across log file rotation: the last record's
chain hash carries forward into the first record of the new file. Rotation
does not reset the chain.

### Chain injectivity and verification order

Exactly one record MAY carry any given `previousHash` value in a chain.
Verifiers MUST reconstruct the chain as a strict walk forward from the
genesis record, following the unique successor at each step; verifying only
that every `previousHash` in a presented record set resolves to some record
in the set is insufficient. Set-membership verification accepts chain forks
in which two records share a predecessor, producing a parallel history that
is internally consistent but differs from the honest chain after the fork
point. The concrete attack this rule prevents is described in Security
Considerations, "Chain fork".

### Parsing Rules

This predicate follows the in-toto attestation parsing rules. Summary:

- Consumers MUST ignore unrecognized fields.
- The `predicateType` URI includes the major version number and will always
  change whenever there is a backwards incompatible change.
- `predicate.action.type` is one of `"tool_call"`, `"checkpoint"`, or
  `"chain_break"`. Consumers MUST reject records with unrecognized types.
- The `predicate.chain` object is OPTIONAL on `tool_call` records. When
  absent, the attestation stands alone without ordering guarantees.
- The `predicate.contentDigest` object is OPTIONAL. When absent, the
  attestation does not bind to specific request/response payloads.
- The `predicate.parties` array is OPTIONAL. When absent, all fields are
  treated as intermediary-witnessed (legacy behavior).
- The `predicate.extensions` object is OPTIONAL. Nesting depth within
  `extensions` MUST NOT exceed 128 levels. Counting rule per
  in-toto/attestation#570: the outermost brace is depth 1; each nested open
  container (`{` or `[`) increments the count by one. Implementations MUST
  reject records exceeding this bound without attempting to parse the excess
  depth.

### Fields

#### `predicate.action`

| Field        | Type              | Required | Description                                               |
| ------------ | ----------------- | -------- | --------------------------------------------------------- |
| `type`       | string            | Yes      | `"tool_call"`, `"checkpoint"`, or `"chain_break"`         |
| `protocol`   | string            | Yes*     | Protocol observed. e.g. `"mcp"`. *Required for tool_call* |
| `method`     | string            | Yes*     | Protocol method. e.g. `"tools/call"`. *Required for tool_call* |
| `toolName`   | string            | Yes*     | Tool name as invoked. *Required for tool_call*            |
| `namespace`  | string or null    | No       | Tool namespace (for multi-server routing)                 |
| `timestamp`  | string (RFC 3339) | Yes      | When the action was observed by the intermediary          |
| `durationMs` | integer           | Yes*     | Wall-clock duration in milliseconds. *Required for tool_call* |
| `success`    | boolean           | Yes*     | Whether the tool call completed without error. *Required for tool_call* |
| `errorCode`  | integer or null   | No       | Protocol-level error code (e.g., JSON-RPC error code)     |
| `errorClass` | string or null    | No       | Human-readable error classification                       |

#### `predicate.agent`

| Field              | Type           | Required | Description                                       |
| ------------------ | -------------- | -------- | ------------------------------------------------- |
| `principal`        | string or null | No       | Identity of the requesting agent or user. Provenance is deployment-dependent (see Field provenance model) |
| `model`            | string or null | No       | AI model identifier. **Provenance: declared** (client-asserted) |
| `sessionId`        | string or null | No       | Session or conversation identifier. **Provenance: declared** |
| `turnId`           | string or null | No       | Turn within the session. **Provenance: declared** |
| `invocationReason` | string or null | No       | Why the agent invoked this tool. **Provenance: declared** |

#### `predicate.parties`

| Field   | Type     | Required | Description                                                    |
| ------- | -------- | -------- | -------------------------------------------------------------- |
| `party` | string   | Yes      | Identifier of the asserting entity (e.g., `"gateway"`, `"client"`, `"policy-engine"`) |
| `role`  | string   | Yes      | `"witness"` (intermediary-observed) or `"asserter"` (declared by that party, unverified by intermediary) |
| `scope` | string[] | Yes      | Field names this party's signature covers. Names refer to the underlying signed audit record and may include implementation-layer fields (e.g., `id`) not present in this predicate schema |

#### `predicate.upstream`

| Field       | Type           | Required | Description                                      |
| ----------- | -------------- | -------- | ------------------------------------------------ |
| `name`      | string or null | No       | Name of the upstream tool server                 |
| `transport` | string or null | No       | Transport type (e.g., `"stdio"`, `"streamable-http"`) |

#### `predicate.chain`

| Field          | Type   | Required | Description                                                    |
| -------------- | ------ | -------- | -------------------------------------------------------------- |
| `previousHash` | string | Yes      | Chain hash of the preceding record (`SHA-256(JSON.stringify(prev_record))`). `"genesis"` for the very first record in the log only; after a chain_break, the successor carries the break record's chain hash |

#### `predicate.checkpoint`

| Field         | Type    | Required | Description                                                    |
| ------------- | ------- | -------- | -------------------------------------------------------------- |
| `sequence`    | integer | Yes      | Monotonically increasing checkpoint ordinal. Safe integer bound applies (< 2^53) |
| `recordCount` | integer | Yes      | Total records emitted since chain genesis or last chain_break  |
| `previousHash`| string  | Yes      | Chain head at checkpoint time: the chain hash of the most recent record, which equals the `previousHash` the next record will carry |

#### `predicate.chainBreak`

| Field             | Type              | Required | Description                                         |
| ----------------- | ----------------- | -------- | --------------------------------------------------- |
| `reason`          | string            | Yes      | Why the chain was broken (e.g., `"crash_recovery"`, `"forced_rotation"`) |
| `priorHead`       | string or null    | No       | Last known chain head before the break              |
| `priorSequence`   | integer or null   | No       | Last checkpoint sequence before the break           |
| `priorRecordCount`| integer or null   | No       | Last known record count before the break            |

#### `predicate.contentDigest`

| Field      | Type                          | Required | Description                                    |
| ---------- | ----------------------------- | -------- | ---------------------------------------------- |
| `request`  | ResourceDescriptor (digests)  | No       | SHA-256 of JCS-canonicalized (RFC 8785) request payload |
| `response` | ResourceDescriptor (digests)  | No       | SHA-256 of JCS-canonicalized (RFC 8785) response payload |

#### `predicate.extensions`

An arbitrary JSON object for implementation-specific metadata. Nesting depth
MUST NOT exceed 128 levels (outermost brace = depth 1; each nested `{` or `[`
increments by one). If present, implementations SHOULD compute an
`extensionsDigest` (SHA-256 of the M/L-tagged canonical form of the extensions
object) and include it in the signing tuple. This allows the extensions content
to be stripped for storage efficiency while remaining signature-bound.

#### `predicate.metadata`

| Field            | Type           | Required | Description                                        |
| ---------------- | -------------- | -------- | -------------------------------------------------- |
| `attestorVersion`| string         | Yes      | Version of the attestation-producing software      |
| `configHash`     | string or null | No       | SHA-256 of the attestor's policy configuration JSON |

## Example

A complete worked example for an AI agent creating a pull request via the
GitHub MCP server. Every digest below has a shown preimage and is independently
recomputable.

### Configuration preimage

The `configHash` is SHA-256 of the gateway's serialized policy configuration:

```json
{"defaultEffect":"allow","rules":[{"effect":"deny","tools":["fs/delete_file"],"principals":["*"]}]}
```

`configHash` = SHA-256 of above = `f9a6ceae76b6df4bb5424fc76a204c443fa85a4384bd5ddf4bd5f532e3610cb3`

### Underlying gateway audit record (JSONL line)

This is the record as written to the gateway's `audit.jsonl` file. It is the
artifact whose chain hash becomes the subject digest:

```json
{"id":"bf7a2f62-4d0f-4cce-afd2-cbfbf7bca2a5","timestamp":"2026-08-18T14:32:01.998Z","method":"tools/call","toolName":"github/create_pull_request","namespace":"github","upstream":"github-server","principal":"user:alice@example.com","durationMs":1247,"success":true,"aiInvocation":{"model":"claude-sonnet-4-20250514","sessionId":"conv_abc123","turnId":"turn_1"},"parties":[{"party":"gateway","role":"witness","scope":["id","timestamp","method","toolName","namespace","upstream","principal","durationMs","success","errorCode","previousHash"]},{"party":"client","role":"asserter","scope":["aiInvocation"]}],"previousHash":"genesis","attestation":"a33a18da2dea9936e327b8cc19ff6c7ca532682be975ad8460466fed4209d03c"}
```

Chain hash = `SHA-256(above)` = `cd26e6c4930f34da6dbbb53988f4920b13eedc7e3354ac51a82efeac9574e664`

The attestation field is `HMAC-SHA256(key, signing_canonical_form)`. For
reproducibility, this example uses the trivial key `00` repeated 32 bytes
(`0000...0000`, 64 hex chars). The signing canonical form of this record is:

```json
[["id","bf7a2f62-4d0f-4cce-afd2-cbfbf7bca2a5"],["timestamp","2026-08-18T14:32:01.998Z"],["method","tools/call"],["toolName","github/create_pull_request"],["namespace","github"],["upstream","github-server"],["principal","user:alice@example.com"],["durationMs",1247],["success",true],["errorCode",null],["previousHash","genesis"],["aiInvocation",["M",[["model","claude-sonnet-4-20250514"],["sessionId","conv_abc123"],["turnId","turn_1"]]]],["parties",[{"party":"gateway","role":"witness","scope":["id","timestamp","method","toolName","namespace","upstream","principal","durationMs","success","errorCode","previousHash"]},{"party":"client","role":"asserter","scope":["aiInvocation"]}]]]
```

### Content digest preimages

Request payload (JCS of the `params` object):
```
{"method":"tools/call","params":{"arguments":{"base":"main","body":"Corrects a spelling mistake in line 42","head":"fix-typo","title":"Fix typo in README"},"name":"github/create_pull_request"}}
```
SHA-256 = `167bd5c6ecde61c67bb42cb8607bd17c39aca91729b302247c67833fe7427815`

Response payload (JCS of the `result` object):
```
{"result":{"content":[{"text":"Pull request #123 created successfully","type":"text"}]}}
```
SHA-256 = `46943f801d4623b020293ed8f31d3603972dd20d9687949a31f13a5f12acdd24`

### In-toto Statement

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "session:agent-workspace-4f2a",
    "digest": {
      "sha256": "cd26e6c4930f34da6dbbb53988f4920b13eedc7e3354ac51a82efeac9574e664"
    }
  }],
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "tool_call",
      "protocol": "mcp",
      "method": "tools/call",
      "toolName": "github/create_pull_request",
      "namespace": "github",
      "timestamp": "2026-08-18T14:32:01.998Z",
      "durationMs": 1247,
      "success": true,
      "errorCode": null,
      "errorClass": null
    },
    "agent": {
      "principal": "user:alice@example.com",
      "model": "claude-sonnet-4-20250514",
      "sessionId": "conv_abc123",
      "turnId": "turn_1",
      "invocationReason": null
    },
    "parties": [
      {
        "party": "gateway",
        "role": "witness",
        "scope": ["id", "timestamp", "method", "toolName", "namespace", "upstream", "principal", "durationMs", "success", "errorCode", "previousHash"]
      },
      {
        "party": "client",
        "role": "asserter",
        "scope": ["aiInvocation"]
      }
    ],
    "upstream": {
      "name": "github-server",
      "transport": "stdio"
    },
    "chain": {
      "previousHash": "genesis"
    },
    "contentDigest": {
      "request": { "sha256": "167bd5c6ecde61c67bb42cb8607bd17c39aca91729b302247c67833fe7427815" },
      "response": { "sha256": "46943f801d4623b020293ed8f31d3603972dd20d9687949a31f13a5f12acdd24" }
    },
    "metadata": {
      "attestorVersion": "mcp-audit-gateway/0.4.0",
      "configHash": "f9a6ceae76b6df4bb5424fc76a204c443fa85a4384bd5ddf4bd5f532e3610cb3"
    }
  }
}
```

The subject digest `cd26e6c4...` is the chain hash of the underlying gateway
audit record shown above. It is the `previousHash` the next record in this
chain will carry.

### Checkpoint example

Emitted after the single record above. This is a toy chain (1 record) for
verifiability; production deployments typically checkpoint every 100 records or
60 seconds.

Underlying checkpoint JSONL line:

```json
{"id":"ckpt_9a1b2c3d-4e5f-6789-abcd-ef0123456789","type":"checkpoint","timestamp":"2026-08-18T14:33:42.101Z","sequence":1,"recordCount":1,"previousHash":"cd26e6c4930f34da6dbbb53988f4920b13eedc7e3354ac51a82efeac9574e664","parties":[{"party":"gateway","role":"witness","scope":["sequence","recordCount","previousHash"]}],"attestation":"aa8ac887a8c3064d610306f8694a44ebe6ba4f4e118e417d268602b47a462f53"}
```

Chain hash = `SHA-256(above)` = `68964a4bb4cf84608fb1a63d1164c49485507f0545cf4fc8c7299361019b4898`

In-toto Statement:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "session:agent-workspace-4f2a",
    "digest": { "sha256": "68964a4bb4cf84608fb1a63d1164c49485507f0545cf4fc8c7299361019b4898" }
  }],
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "checkpoint",
      "timestamp": "2026-08-18T14:33:42.101Z"
    },
    "checkpoint": {
      "sequence": 1,
      "recordCount": 1,
      "previousHash": "cd26e6c4930f34da6dbbb53988f4920b13eedc7e3354ac51a82efeac9574e664"
    },
    "parties": [
      { "party": "gateway", "role": "witness", "scope": ["sequence", "recordCount", "previousHash"] }
    ],
    "metadata": {
      "attestorVersion": "mcp-audit-gateway/0.4.0"
    }
  }
}
```

The consumer externalizes `{"sequence": 1, "recordCount": 1, "previousHash":
"cd26e6c4..."}` to a store the attestor cannot write. On the next verification
pass, if the chain head differs from `cd26e6c4...` at record count 1, or if
fewer records exist, truncation is detected.

### Chain break example

Emitted after a crash recovery. The next record chains from this record's
hash (not from "genesis").

Underlying chain_break JSONL line:

```json
{"id":"brk_7f8e9d0c-1b2a-3456-cdef-0123456789ab","type":"chain_break","timestamp":"2026-08-18T15:01:00.000Z","reason":"crash_recovery","priorHead":"68964a4bb4cf84608fb1a63d1164c49485507f0545cf4fc8c7299361019b4898","priorSequence":1,"priorRecordCount":1,"attestation":"28ca6f7f395cf0028fd2b0290c0a507782a916ad3c1c8c735edea13926b901bb"}
```

Chain hash = `SHA-256(above)` = `2cb2c22007b553470cd0775088dac6be6f3be1594a1caa9437e493916dea07a9`

In-toto Statement:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "session:agent-workspace-4f2a",
    "digest": { "sha256": "2cb2c22007b553470cd0775088dac6be6f3be1594a1caa9437e493916dea07a9" }
  }],
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "chain_break",
      "timestamp": "2026-08-18T15:01:00.000Z"
    },
    "chainBreak": {
      "reason": "crash_recovery",
      "priorHead": "68964a4bb4cf84608fb1a63d1164c49485507f0545cf4fc8c7299361019b4898",
      "priorSequence": 1,
      "priorRecordCount": 1
    },
    "metadata": {
      "attestorVersion": "mcp-audit-gateway/0.4.0"
    }
  }
}
```

The `priorHead` is the chain hash of the last known record (the checkpoint
above). The next record in the new chain carries
`previousHash: "2cb2c220..."` — the break record's own chain hash — binding
the discontinuity evidence into the successor chain.

## Security Considerations

### Declared vs. witnessed fields

The `agent.model`, `agent.sessionId`, `agent.turnId`, and
`agent.invocationReason` fields are client-asserted. A compromised or malicious
client can lie about its model identity, session, or reasoning. Verification
policies SHOULD NOT rely on these fields for security decisions without
corroborating evidence (e.g., an OAuth token binding the session to a known
principal at the transport layer).

The `agent.principal` field has deployment-dependent provenance. Under
authenticated transports (OAuth, mTLS), the gateway witnesses it directly.
Under bare stdio, it is declared by the client. The `parties` array
disambiguates which case applies for a given deployment.

### Tail truncation window

Records emitted after the most recently externalized checkpoint are
unprotected. The checkpoint interval represents a security/performance
trade-off: shorter intervals reduce the window but increase storage I/O.
Implementations SHOULD document their checkpoint interval in operational
guidance.

### Chain fork

An adversary with write access to the record store can construct a chain
fork: two records that share a common `previousHash` value, producing
parallel histories that share a prefix through the fork point. Every
`previousHash` in a presented sub-log resolves and no hash is broken; the
fork is detectable only by enforcing injectivity of the predecessor
relation (see "Chain injectivity and verification order") and by
reconstructing the chain by strict walk from genesis rather than by
set-membership check of `previousHash` values.

The concrete attack shape: a consumer targets the chain by subject digest,
where the subject digest equals the genesis record's chain hash. The
adversary constructs a second record that also chains from genesis (its
`previousHash` equals the genesis chain hash, exactly as an honest
successor's would). The adversary presents a two-record sub-log containing
the genesis record and the parallel-branch successor. Every hash link is
valid; the second-`"genesis"` rule (previous subsection) sees only one
`"genesis"` because it appears only in the genesis record. A verifier that
checks `previousHash` membership rather than performing the strict walk
accepts the branch as valid and cannot distinguish it from the honest
chain.

The strict-walk-plus-injectivity rule prevents this: after reading the
genesis record, the walk asks "what is the unique successor?"; the presence
of two candidates carrying the same `previousHash` is itself the detection.

An externalized checkpoint (see "Tail truncation window") mitigates this
for records preceding the checkpoint: only one branch can carry the
checkpoint's chain head, so any parallel branch that continues past the
checkpoint's position is exposed. Records emitted after the most recently
externalized checkpoint remain in the same unprotected window described
under tail truncation, subject to the same detection-not-attribution
limit.

### Chain break vs. planted state

An attestor that crashes and restarts MUST emit a `chain_break` record. The
successor record chains from the break record's hash, binding the
discontinuity evidence into the new chain. This means discarding the break
record breaks hash linkage — an adversary cannot silently remove the scar.
The literal string `"genesis"` appears exactly once, at the true start of the
log; a second `"genesis"` anywhere is evidence of chain replacement.
Implementations that persist chain state to disk SHOULD validate the state
file's integrity on startup (e.g., separate MAC over the state file, or
storing the last hash in a second location).

### Safe integer bound

All integer fields (`durationMs`, `errorCode`, checkpoint `sequence` and
`recordCount`) MUST remain within the I-JSON safe integer range (absolute
value < 2^53, per RFC 7493 Section 2.2). Implementations MUST reject records
with integers at or above this bound. At one checkpoint per 100 records and
one record per millisecond, the checkpoint sequence provides approximately
29 million years of headroom before overflow.

### Two canonical forms

The predicate intentionally uses two different canonicalization procedures
(signing tuple-array vs. RFC 8785 JCS for content digests). The signing form
rejects floats and uses M/L type tags — properties that make it injective and
safe for structured attestation fields. Content digests use RFC 8785 because
MCP tool payloads are arbitrary JSON that may contain floats, and the registry
should converge on a single content-digest scheme (aligning with
in-toto/attestation#570). Implementations MUST NOT mix the two: signing fields
use the tuple-array form; content payloads use JCS.

## Changelog and Migrations

### v0.1.0 (current)

Initial version. Changes from the original submission based on review feedback:

- Added `checkpoint` and `chain_break` as `action.type` values with their
  own sub-schemas, addressing the truncation detection gap
- Added `parties` array with `witness`/`asserter` roles to distinguish
  field provenance (intermediary-observed vs. client-declared)
- Specified two canonical forms: signing tuple-array (M/L-tagged, integer-only,
  surrogate-rejecting) for chain integrity, RFC 8785 JCS for content digests
- Documented genesis convention (`previousHash: "genesis"`) and chain_break
  requirement for crash recovery
- Added I-JSON safe integer bound on all integer fields (< 2^53, RFC 7493)
- Added 128-level depth bound on `extensions` with counting rule per
  in-toto/attestation#570 (outermost brace = depth 1)
- Clarified deployment-dependent provenance of `principal`
- Replaced placeholder example with fully recomputable worked example:
  every digest has a shown preimage, verifiable against the reference
  implementation's conformance vectors
- Added Security Considerations section
- Moved listing from "Vetted Predicates" to community contributions per
  ITE-63 process
