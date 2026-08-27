# Predicate type: AI Agent Action

Type URI: https://in-toto.io/attestation/ai-agent-action/v0.1

Version: 0.1.0

Authors: Elankumaran Srinivasan (@elang2)

Conformance corpus: Sankalp Gilda (@astrogilda), see the "Conformance corpus"
subsection under Canonicalization.

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

## Requirements notation

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to
be interpreted as described in BCP 14
([RFC 2119](https://www.rfc-editor.org/rfc/rfc2119),
[RFC 8174](https://www.rfc-editor.org/rfc/rfc8174)) when, and only when,
they appear in all capitals.

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

-   **Agent**: The AI system that initiated the tool call (identified by model,
  session, or principal)
-   **Tool**: The capability being invoked (identified by name and namespace)
-   **Intermediary**: The entity producing the attestation (identified in the
  statement's signing metadata)
-   **Outcome**: Whether the tool call succeeded or failed, its duration, and
  any error classification

The predicate does not include the tool's input arguments or output content by
default, as these may contain sensitive data. The `contentDigest` field allows
optional binding to the actual request/response payloads without including them
inline.

### Field provenance model

Fields in this predicate have different trust characteristics depending on who
asserts them:

-   **Witnessed by intermediary**: `timestamp`, `durationMs`, `success`,
  `errorCode`, `errorClass`, `toolName`, `method`, `previousHash`, and the
  content digests it computes from the payloads it forwards. The intermediary
  directly measures or extracts these from the protocol messages it forwards
  and the chain state it maintains.
-   **Asserted by client (declared)**: `model`, `sessionId`, `turnId`,
  `invocationReason`. Provided by the AI client, the intermediary cannot
  independently verify these. Consumers SHOULD treat them as claims, not facts.
-   **Deployment-dependent**: `principal`. Observed when the intermediary
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

Scope arrays describe witnessing, not signing. The emitter already signs
every byte it emits, whether or not a field appears in any scope: the
attestation member covers the signing tuple, the chain hash covers the
record, and the DSSE envelope covers the Statement. What a scope adds is
the provenance claim, naming the party that stands behind a field's value
and in what role. The gateway's scope in the worked example excludes
`aiInvocation` not because the gateway leaves those bytes unsigned but
because it only relays them; the client asserts them, and the scope split
is what makes that distinction machine-readable.

The structural members `id`, `type`, and `timestamp` are witnessed by the
emitter by construction, because no other party could originate them, and
they need not appear in any scope. Every other field present in the record
but not covered by any party's `scope` MUST be treated as having unknown
provenance: signed by the emitter, vouched for by nobody. Verification
policies SHOULD reject or flag records where security-relevant fields have
no attesting party. When the `parties` array is absent entirely, the
legacy rule in Parsing Rules applies and every field is treated as
intermediary-witnessed.

## Schema

The predicate defines three record types that share a chain:

### Action record (tool_call)

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "<agent-session-identifier>",
    "digest": { "sha256": "<chain-hash-of-the-chain's-genesis-record (stable chain identifier)>" }
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
      "previousHash": "<lowercase 64-hex chain hash; the literal 'genesis' for the very first record only>"
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
    "chain": {
      "previousHash": "<lowercase 64-hex; the chain linkage, equals checkpoint.previousHash>"
    },
    "checkpoint": {
      "sequence": "<integer>",
      "recordCount": "<integer>",
      "previousHash": "<lowercase 64-hex; the chain head being externalized>"
    },
    "parties": [...],
    "metadata": {
      "attestorVersion": "<string; Statement-layer annotation, see Underlying record shape>"
    }
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
    "chain": {
      "previousHash": "<lowercase 64-hex; equals chainBreak.priorHead; omitted when the prior head is unrecoverable>"
    },
    "chainBreak": {
      "reason": "<string>",
      "priorHead": "<hex string or null>",
      "priorSequence": "<integer or null>",
      "priorRecordCount": "<integer or null>"
    },
    "metadata": {
      "attestorVersion": "<string; Statement-layer annotation, see Underlying record shape>"
    }
  }
}
```

### Underlying record shape

The predicate is a view over an underlying signed audit record: the JSONL
line the intermediary writes, whose canonical bytes the chain hash commits
to and whose members the signing tuple enumerates. Which members live in
the record and which exist only in the Statement wrapper is normative,
because each layer has a different protection mechanism (see Signature
coverage boundary).

The `id` member, present as the first tuple element on every record type
and signed by every attestation, is an opaque producer-chosen string of at
most 128 bytes that MUST be unique within a chain. Consumers MUST treat it
as an opaque identifier and MUST NOT parse structure from it; a v0.2
registry may profile particular producers (for example, UUIDv7 for
timestamp-ordered ids) without changing this rule. A verifier MUST reject,
fail-closed, a chain in which two records share the same `id`, because two
records with the same identifier and different chain hashes are a fork by
another name.

A `tool_call` record carries exactly these members. Required: `id`, `type`
(the string `"tool_call"`), `timestamp`, `method`, `toolName`, `durationMs`,
`success`, `attestorVersion`, `previousHash` (when the record is part of a
chain), and `attestation`. Optional, omitted from the record when absent:
`namespace`, `upstream`, `principal`, `errorCode`, `errorClass`,
`contentDigestRequest`, `contentDigestResponse`, `configHash`,
`decisionContextDigest`, `extensionsDigest`, `extensions`, `aiInvocation`,
and `parties`.

A `checkpoint` record carries `id`, `type` (`"checkpoint"`), `timestamp`,
`sequence`, `recordCount`, `previousHash`, and `attestation`, plus
optionally `parties`. A `chain_break` record carries `id`, `type`
(`"chain_break"`), `timestamp`, `reason`, and `attestation`, and MUST carry
all three of `priorHead`, `priorSequence`, and `priorRecordCount`. When the
prior head is unknown (the recovery case that motivates the record), the
three members MUST be encoded as explicit `null` values in the record, not
omitted; without this rule two producers with identical knowledge would
disagree on whether the members appear at all, and JCS would then produce
two different chain-hash preimages for the same observation.

The Statement is derived from the record, never the reverse.
`contentDigestRequest` and `contentDigestResponse` surface as
`predicate.contentDigest.request` and `.response`; `attestorVersion` and
`configHash` surface under `predicate.metadata`; `aiInvocation` members
surface under `predicate.agent` alongside `principal`; the remaining action
members surface under `predicate.action`. On `checkpoint` and `chain_break`
statements, `predicate.metadata` is a Statement-layer annotation with no
record-layer counterpart.

`predicate.action.protocol` on `tool_call` statements is a REQUIRED
Statement-only annotation identifying the wire protocol the intermediary
observed. At v0.1 it takes the fixed lowercase value `"mcp"` for the Model
Context Protocol; future protocols will extend the vocabulary in a
backward-compatible v0.2 registry. The value is intermediary-declared and
sits outside the record-layer protections (the signing tuple and the chain
hash); like every Statement-layer member it is covered by the DSSE envelope
signature, so its integrity rests on the envelope and its truthfulness on
the intermediary that signed it.

The record-layer `upstream` member is a bare string naming the upstream
tool server. The Statement-layer `predicate.upstream` is an object whose
`name` field surfaces the record's `upstream` value verbatim and whose
optional `transport` field (allowed values `"stdio"` and `"streamable-http"`
at v0.1, matching the transports MCP defines, extensible in a v0.2
registry) is a Statement-only annotation declared by the intermediary. Like
`protocol`, `transport` sits outside the record-layer protections at v0.1
and is covered only by the DSSE envelope, so consumers rely on the
intermediary's identity for its truthfulness. A v0.2 profile MAY promote
`transport` to a record-layer member and add it to the signing tuple; the
present layout keeps v0.1 example digests stable while making the trust
boundary explicit for readers.

On a `checkpoint` record, the single `previousHash` member populates two
Statement fields: `predicate.chain.previousHash` (the linkage the chain
walker follows for every record type) and `predicate.checkpoint.previousHash`
(the checkpoint's own restatement of the head at the moment the checkpoint
was cut). Both derive from the same record member and MUST be byte-equal in
the Statement; a verifier MUST reject, fail-closed, a checkpoint statement
whose two `previousHash` values differ.

Three protection layers apply. The signing tuple covers the record members
enumerated in the signing canonical form, under the attestor's key. The
chain hash covers the complete record bytes, including the `attestation`
member, binding each record into history. The DSSE envelope signature
covers the whole Statement, including Statement-layer members the first two
layers never see. A field's protection is determined by the layers that
contain it, and the Signature coverage boundary subsection below states
what each layer can and cannot promise.

### Canonicalization

This predicate names three canonical byte-string derivations, and the boundary
between them is drawn by what a party signs or hashes rather than by what a
party stores. Every byte string any rule below hashes is named exactly once,
so that two implementations reading only this document derive identical bytes
from identical observations.

#### The record canonical form (chain hash preimage)

The record canonical form is the RFC 8785 (JCS) canonicalization of the audit
record object, including its `attestation` member, encoded as UTF-8. Object
members are sorted by UTF-16 code unit, numbers take the shortest form that
round-trips under IEEE 754, and no insignificant whitespace appears. A
producer MUST write each JSONL line as exactly these bytes, and a verifier
MUST recompute the form from the parsed record and MUST reject, fail-closed,
any line whose bytes differ from the recomputation. The line terminator is
not part of the record: the chain hash preimage ends at the closing brace.

That last obligation is the one that cannot be dropped. Without it the
preimage has two readings, the bytes on disk and the re-serialization of the
parsed object, and they coincide only by accident. Any log shipper that
reparses and re-emits JSON moves a verifier from one reading to the other
while every field value stays identical, so a chain that verified before the
shipper ran fails after it.

`JSON.stringify` MUST NOT be used to derive the record canonical form. It is
not a canonical form and it is not portable. ECMAScript orders canonical
numeric property names ahead of every other name and in ascending numeric
order regardless of insertion order, so an `extensions` object whose members
are `10`, `2`, `aa`, `zz` serializes as `2, 10, zz, aa` in a JavaScript
gateway, as `10, 2, zz, aa` in a Python gateway that preserves insertion
order, and as `10, 2, aa, zz` in a Go gateway, whose `encoding/json` sorts
map keys. That is three chain hashes for one observation, and none of the
three implementations has done anything wrong. JCS fixes the order to
`10, 2, aa, zz` for all of them.

String escaping is fixed by the same rule. JCS emits a character rather than
an escape wherever the character is permitted, so a serializer whose default
is ASCII-escaping output, and a serializer that escapes `<`, `>` and `&`,
both produce non-canonical bytes and both are rejected. Neither behaviour is
exotic: the first is Python's default and the second is Go's.

The `attestation` member's value serializes as a lowercase hex string of
fixed length for its declared signature scheme: 64 hex characters (32 bytes)
for HMAC-SHA256, and 128 hex characters (64 bytes) for Ed25519. This
encoding is normative because the chain hash covers the byte sequence that
includes the `attestation` member, so any deviation (uppercase hex,
base64url, base64, differing length) changes the chain-hash preimage and
desynchronizes the chain across producers. A producer MUST NOT emit any
other encoding, and a verifier MUST reject, fail-closed, a record whose
`attestation` value is not lowercase hex of the expected length for the
declared scheme.

#### The content digest form (payload binding)

Content digests bind MCP payloads without inlining them. They use RFC 8785
over the payload, and the digest is the lowercase 64-hex SHA-256 of those
bytes.

`contentDigest.request` is the digest of the JSON-RPC `params` member. For
`contentDigest.response` the preimage depends on which JSON-RPC response
arrived. A success response carries a `result` member and the digest is over
that member. An error response carries no `result` member at all, and the
digest is over the `error` member. A record with `action.success` false is by
construction an error response, so this is not an edge case: it is every
failed tool call, which is the half of an audit trail an investigator reaches
for first. A producer MUST NOT digest `null`, an empty object, or the whole
response envelope in place of the named member.

A JSON-RPC response that carries neither `result` nor `error` is malformed,
as is one that carries both. In either case the `action.success`
classification is unreliable and a producer digesting the "convenient"
member silently would let the trace look clean while the transport observed
something wrong. A producer MUST still record the call, because suppressing
the record would hand a misbehaving tool server an audit-evasion channel.
It MUST record it with `success: false`, SHOULD set `errorClass` to a value
identifying the protocol violation, and MUST NOT emit a
`contentDigest.response` member, because a malformed response has no
well-defined preimage to digest. A verifier that holds the response payload
MUST reject, fail-closed, a record carrying a `contentDigest.response`
computed from a malformed response.

Floats are permitted here and only here. MCP tool payloads are arbitrary JSON
and routinely carry them.

Using RFC 8785 for content digests aligns with the registry-level content
digest scheme proposed in in-toto/attestation#570.

#### The signing canonical form (attestation signature)

The signing canonical form is a **tuple-array**: an ordered array of
`[field-name, value]` pairs, serialized as compact JSON with no insignificant
whitespace. The field list and its order are part of this specification, per
record type. A form whose input list lives in an implementation's type
definition is not a specification of anything: a second implementer can
reproduce the tagging rules, the sort order and the tags exactly and still
reproduce no signature at all, and no conformance vector can be written for
it, because a vector needs a preimage the text determines.

Fields referenced by any party's `scope` remain in the tuple even when the
intermediary is not their sole witness. Dropping them would leave scope
entries pointing at bytes the attestor never signed, which hollows the
provenance mechanism the `parties` array exists to express. This is why
`method`, `upstream`, and `principal` stay in the `tool_call` tuple even
though the client, not the intermediary, originates each of them.

Every string encoded anywhere in the tuple form is serialized under the JCS
string rules: JCS-preferred character-rather-than-escape for every codepoint
JCS permits unescaped, no ASCII-escaping default, no HTML-escape of `<`,
`>`, `&`, and no additional implementation escapes. This applies to
`toolName`, `method`, `id`, `errorClass`, `principal`, every string inside
`parties`, every field name in the outer tuple, and every string value
including those inside M-tagged and L-tagged structures. Numeric
serialization follows the shortest-round-trip rule already stated for the
record canonical form. Two implementations otherwise reproducing the field
list and tagging exactly would still diverge on any string containing a
codepoint their default escape strategies disagree on; this paragraph closes
that gap universally.

For a `tool_call` record the tuple contains, in exactly this order:

```text
id, type, timestamp, method, toolName, namespace, upstream, principal,
durationMs, success, errorCode, errorClass, previousHash,
contentDigestRequest, contentDigestResponse, attestorVersion, configHash
```

An absent member of this base list is encoded as `null`, not omitted. Four
conditional members join the tuple only when present on the record, appended
after `configHash` in this order and omitted entirely when absent:
`decisionContextDigest` (digest of policy-engine decision context, when a
policy engine participates; when present it is a lowercase 64-hex SHA-256
whose preimage is deployment-defined at v0.1, and a v0.2 registry may pin
the preimage), then `extensionsDigest`, then `aiInvocation`, then
`parties`.

For a `checkpoint` record the tuple contains, in exactly this order:

```text
id, type, timestamp, sequence, recordCount, previousHash
```

followed by `parties` when present. For a `chain_break` record the tuple
contains, in exactly this order:

```text
id, type, timestamp, reason, priorHead, priorSequence, priorRecordCount
```

with the three `prior*` members encoded as `null` when unknown.

Structured values inside the tuple are encoded as follows. The `aiInvocation`
value (and the `extensions` object when computing `extensionsDigest`) uses
**M/L type tagging** to prevent cross-type collisions:

-   Objects become `["M", [[key, value], ...]]` with keys sorted by UTF-16
  code unit values (ECMAScript default sort).
-   Arrays become `["L", [value, ...]]`.
-   Scalars (string, number, boolean, null) pass through untagged.

The `parties` value is **not** M/L-tagged: it is serialized as a plain JSON
array whose entries are objects with exactly the members `party`, `role`,
`scope`, in exactly that order. The member order is fixed by this paragraph
precisely because plain objects otherwise reintroduce the serialization
ambiguity the tuple form exists to remove. Every string inside `parties`,
including every entry of the `scope` array, is serialized under the JCS
string rules that govern the record canonical form: JCS-preferred
character-rather-than-escape for every codepoint JCS permits unescaped, no
ASCII-escaping default, no HTML-escape of `<`, `>`, `&`, and no additional
implementation escapes. Numeric serialization inside `parties`, if any,
follows the shortest-round-trip rule already stated for the record canonical
form. This paragraph closes the escape-ambiguity gap that a non-M/L-tagged
carve-out would otherwise reopen.

This form is used to compute the bytes the attestor signs (HMAC-SHA256 or
Ed25519). Numbers in the signing form MUST be safe integers (absolute value
< 2^53, per RFC 7493 §2.2); floats are rejected by design. This constraint
binds the signing form and the record canonical form. It does not bind
content payloads, which is why the content digest form exists.

Implementations MUST NOT substitute one form for another: chain hashes use
the record canonical form, attestation signatures use the signing canonical
form, and payload digests use the content digest form.

#### Signature coverage boundary

Three mechanisms protect different byte ranges, and a verifier needs to
know which one it is trusting for any given field. The attestor's HMAC or
signature covers the signing tuple; for `tool_call` records at v0.1 the
tuple spans the complete record shape, including the content digests,
`errorClass`, `attestorVersion`, and `configHash`. The chain hash covers
the complete record bytes including the `attestation` member, which is what
binds each record into history and makes deletion and reordering
detectable. The DSSE envelope signature covers the whole Statement, and it
is the only mechanism covering Statement-layer members the record never
carries: the subject, the predicate type, and `predicate.metadata` on
`checkpoint` and `chain_break` statements.

In statement-at-a-time policy evaluation (single-record verification
outside a chain walk), the tuple signature and the DSSE envelope are the
integrity guarantees available; the chain hash contributes nothing until a
walk traverses it. In the tail-truncation window (see Security
Considerations), a record emitted after the most recently externalized
checkpoint has no downstream chain link yet, so its history binding is
still one-sided; its field integrity still holds under the tuple signature
and the envelope. A policy that treats `configHash` as compliance evidence
or the content digests as payload binding can therefore rely on the
attestor's signature alone. A policy that needs completeness or ordering
MUST walk the chain from a consumer-externalized checkpoint.

Reference implementation and cross-language conformance vectors (JS +
Python): https://github.com/elang2/mcp-audit-gateway/tree/main/test/vectors

#### Strict I-JSON, statement-wide

The whole audit record and the whole Statement are parsed as strict I-JSON.

A duplicate member anywhere, at any depth, makes the record malformed, and a
verifier MUST reject it fail-closed. A lenient parser that keeps the last of
a repeated member lets a record carry `"toolName": "read_file"` and
`"toolName": "delete_repository"` at once: a first-wins reader shows the
auditor the harmless call while the hash commits to the destructive one, and
both readers believe the chain intact.

Every string literal MUST be a well-formed sequence of Unicode scalar values,
in member-name and value position alike. The record MUST be valid UTF-8 with
no overlong form and no surrogate encoded directly in UTF-8; a `\u` escape
naming a high surrogate MUST be immediately followed by one naming a low
surrogate, and an unpaired escape of either half is malformed; a `\u` escape
MUST be exactly four hexadecimal digits with no sign, whitespace or radix
prefix; a string MUST NOT carry a raw unescaped character below U+0020; and
the Unicode noncharacters, U+FDD0 through U+FDEF and U+nFFFE and U+nFFFF in
every plane, are excluded. A verifier MUST apply this to the raw bytes before
any decoded string is read, because a lenient decoder does not fail on
ill-formed input, it substitutes U+FFFD, and every check after that point
reads a string the producer never wrote.

A valid surrogate pair is one supplementary-plane character and is well
formed. A verifier that rejects it is over-rejecting.

A verifier MUST reject, fail-closed, a record whose JSON nesting depth
exceeds 128. Depth is the number of arrays and objects open at a point,
counting the outermost brace as depth 1; scalars do not increase it. The
counting rule follows in-toto/attestation#570. The bound is normative rather
than a resource limit, because with no bound stated two conforming verifiers
disagree about whether identical bytes are evidence at all across the whole
range between their private choices.

#### Conformance corpus

The externally authored conformance corpus for this predicate is maintained
at [`astrogilda/aee-conformance`](https://github.com/astrogilda/aee-conformance)
under `vectors-ai-agent-action/`. The corpus at commit
`92a4340e557615418b1d59f7482f73259beb0b75` certifies the predecessor of
this revision: its `spec-vendored/` directory freezes the specification
text at commit `639ec56`, and regeneration against the present text is
scheduled as a follow-up, at which point the pin advances with it. Once a
corpus pin regenerated against this text lands, an implementation claiming
conformance to this predicate MUST accept every vector in the `accept/`
directory and MUST reject every vector in the `reject/` directory at that
pin. The corpus's `check_vectors.py` self-check refuses to pass a corpus in
which any reject condition lacks an accepting twin carrying the same
condition id, so a verifier that trivially rejects everything does not
satisfy the corpus.

The corpus is authored by Sankalp Gilda (@astrogilda).
Corpus-versus-specification drift is detectable by comparing the vendored
spec against the pinned specification commit, and either version is
authoritative for its own artifact.

This corpus follows the external-corpus door in
[in-toto/ITE#63](https://github.com/in-toto/ITE/pull/63): a conformance
suite authored by someone other than the specification's author, with
negative controls enforced by its self-check. A published run against the
reference implementation will follow the regeneration.

### Chain shape

#### The chain hash

The chain hash of a record is the lowercase 64-hex SHA-256 of that record's
record canonical form as defined above, including its `attestation` member.
The next record carries it as `previousHash`.

`previousHash` is either lowercase 64-hex or the literal lowercase string
`genesis`, and nothing else. Uppercase hex is not canonical. A digest of any
other length is not admissible, and there is no algorithm agility in this
field at v0.1; a future version that needs another algorithm adds an
algorithm member rather than widening what this one accepts. Without this, a
verifier that folds hex case treats two distinct byte strings as one link
while the two successors carrying them hash differently, so the same logical
chain has two identities.

Chain continuity is preserved across log file rotation: the last record's
chain hash carries forward into the first record of the new file. Rotation
does not reset the chain.

#### One head, and one predecessor

Exactly one record in a chain MUST carry any given `previousHash` value. A
verifier MUST reject, fail-closed, a log in which two records share one, and
MUST reconstruct the chain as a strict walk from the genesis record rather
than by checking that each record's `previousHash` appears somewhere in the
presented set.

This is the rule that makes the chain a chain. Nothing in a hash link forbids
a fork: two records may each chain from record one, and every hash in both
branches verifies. Since the subject digest is the genesis hash, both
branches carry the same subject digest and a policy targeting the chain
cannot distinguish them. A presenter holding a four-record chain can
therefore present a two-record branch that omits whichever calls it prefers
an auditor not see, with no hash broken, no second genesis, and no gap in
sequence for a checkpoint to catch. Set-membership verification accepts it; a
strict walk plus the injectivity rule does not. The concrete attack this rule
prevents is described in Security Considerations, "Chain fork".

#### Every record type is on the chain

`checkpoint` and `chain_break` records are chained exactly as `tool_call`
records are: the record following any record of any type carries that
record's chain hash. At the predicate layer, every chained record of every
type carries the linkage in `predicate.chain.previousHash`.
`predicate.checkpoint.previousHash` restates the chain head for the consumer
that externalizes it and is not the linkage; the linkage is always
`predicate.chain.previousHash`. On a checkpoint record the two values are
equal by construction. The checkpoint's predecessor is the record whose
chain hash is the head it snapshots, and a verifier MUST reject a checkpoint
record on which they differ.

At the underlying record layer, `tool_call` and `checkpoint` records carry
the linkage in their `previousHash` member; a `chain_break` record carries it
in `priorHead`, and `predicate.chain.previousHash` on a break record equals
`chainBreak.priorHead`. When the prior head is unrecoverable (`priorHead`
null, meaning the crash also lost or corrupted the chain state), the break record
carries no `predicate.chain` and roots a new chain segment; the missing
linkage is precisely the discontinuity the break attests, and the record
following the break still carries the break record's chain hash, so the scar
binds into the successor chain either way.

Stating the linkage for `chain_break` alone would leave the checkpoint off
the chain. A verifier walking the documented linkage field would step past
every checkpoint, so a checkpoint could be removed without breaking any
documented link, and the checkpoint is the entire mechanism standing between
this predicate and silent tail truncation.

`predicate.chain` is REQUIRED on every `tool_call` and `checkpoint` record,
and on every `chain_break` record whose `chainBreak.priorHead` is non-null.
A verifier MUST reject, fail-closed, a `tool_call` or `checkpoint` record
that lacks `predicate.chain`. There is no standalone mode at v0.1: a
missing chain object is malformed, not merely unanchored.

A `chain_break` record with `chainBreak.priorHead: null` is the single
legitimate case in which a chained record omits `predicate.chain`, because
the prior head is by definition unrecoverable and there is no truthful
value to place there. Such a record does not attest to any prior chain
state, and it MUST NOT be interpreted as evidence about the pre-break
segment; its role is exactly to mark the discontinuity. The record
following the break chains from the break's own record chain hash,
computed as for any other record, so the successor segment remains
verifiable from its genesis. A verifier MUST NOT confuse this case with a
`tool_call` or `checkpoint` record that lacks `predicate.chain`, which
remains rejected fail-closed by the paragraph above.

#### Genesis and breaks

The literal `genesis` appears as `previousHash` exactly once in a chain's
lifetime, on the first record. Its chain hash becomes the anchor for the
chain and the stable chain identifier: every statement in the chain,
`tool_call`, `checkpoint`, and `chain_break` alike, carries the genesis
record's chain hash as its subject digest, so verification policies can
target entire audit chains.

If a chain is broken and restarted (crash, forced rotation, state
corruption), the attestor MUST emit a signed `chain_break` record. After a
`chain_break`, the successor carries the break record's chain hash, never
`genesis`, which binds the discontinuity into the successor chain so that
discarding the break record breaks linkage.

A verifier MUST reject, fail-closed, a log in which `genesis` appears more
than once. This is a MUST and not a SHOULD. A detection obligation a
conformant verifier may decline is not a defence against an adversary who is
choosing which verifier to present to.

### Parsing Rules

This predicate follows the in-toto attestation parsing rules. Summary:

-   Consumers MUST ignore unrecognized fields.
-   The `predicateType` URI includes the major version number and will always
  change whenever there is a backwards incompatible change.
-   `predicate.action.type` is one of `"tool_call"`, `"checkpoint"`, or
  `"chain_break"`. Consumers MUST reject records with unrecognized types.
-   The `predicate.chain` object is REQUIRED on every `tool_call` record,
  every `checkpoint` record, and every `chain_break` record whose
  `chainBreak.priorHead` is non-null. It is OPTIONAL only on `chain_break`
  records with `chainBreak.priorHead: null`, which is the single legitimate
  case in which a chained record omits `predicate.chain` (see Chain shape).
  A verifier MUST reject, fail-closed, any `tool_call` or `checkpoint`
  record that lacks `predicate.chain`, and MUST NOT treat a `chain_break`
  with `priorHead: null` as evidence about the pre-break segment.
-   The `predicate.contentDigest` object is OPTIONAL. When absent, the
  attestation does not bind to specific request/response payloads.
-   The `predicate.parties` array is OPTIONAL. When absent, all fields are
  treated as intermediary-witnessed (legacy behavior).
-   The `predicate.extensions` object is OPTIONAL.
-   Records and Statements are strict I-JSON (see Canonicalization): a
  duplicate member at any depth, an ill-formed string, or nesting depth
  beyond 128 levels (for the whole record, not only `extensions`) MUST be
  rejected fail-closed, without attempting to parse the excess depth.
  Counting rule per in-toto/attestation#570: the outermost brace is depth 1;
  each nested open container (`{` or `[`) increments the count by one.

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
| `previousHash` | string | Yes      | Chain hash of the preceding record: the lowercase 64-hex SHA-256 of the predecessor's record canonical form (RFC 8785, including its `attestation` member), or the literal lowercase `"genesis"` for the very first record in the log only. No other case, length, or algorithm is admissible. After a chain_break, the successor carries the break record's chain hash. Exactly one record may carry any given value (see Chain shape) |

#### `predicate.checkpoint`

| Field         | Type    | Required | Description                                                    |
| ------------- | ------- | -------- | -------------------------------------------------------------- |
| `sequence`    | integer | Yes      | Monotonically increasing checkpoint ordinal. Safe integer bound applies (< 2^53) |
| `recordCount` | integer | Yes      | Total records emitted since chain genesis or last chain_break  |
| `previousHash`| string  | Yes      | Chain head at checkpoint time: the chain hash of the most recent record. Restates the head for externalization; the chain linkage is `predicate.chain.previousHash`, which MUST carry the same value on a checkpoint record (a mismatch is malformed) |

#### `predicate.chainBreak`

| Field             | Type              | Required | Description                                         |
| ----------------- | ----------------- | -------- | --------------------------------------------------- |
| `reason`          | string            | Yes      | Why the chain was broken (e.g., `"crash_recovery"`, `"forced_rotation"`) |
| `priorHead`       | string or null    | Yes (may be null) | Last known chain head before the break. When non-null, `predicate.chain.previousHash` carries the same value; when null, the break record carries no `predicate.chain` and roots a new chain segment (see Chain shape). MUST appear in the record even when the value is null; MUST NOT be omitted |
| `priorSequence`   | integer or null   | Yes (may be null) | Last checkpoint sequence before the break. MUST appear in the record even when null; MUST NOT be omitted           |
| `priorRecordCount`| integer or null   | Yes (may be null) | Last known record count before the break. MUST appear in the record even when null; MUST NOT be omitted            |

#### `predicate.contentDigest`

| Field      | Type                          | Required | Description                                    |
| ---------- | ----------------------------- | -------- | ---------------------------------------------- |
| `request`  | ResourceDescriptor (digests)  | No       | SHA-256 of the JCS-canonicalized (RFC 8785) JSON-RPC `params` member |
| `response` | ResourceDescriptor (digests)  | No       | SHA-256 of the JCS-canonicalized (RFC 8785) JSON-RPC `result` member for success responses, or `error` member for error responses (see Canonicalization) |

#### `predicate.extensions`

An arbitrary JSON object for implementation-specific metadata. The
statement-wide strict I-JSON rules apply, including the 128-level depth bound
(outermost brace = depth 1; each nested `{` or `[` increments by one). If
present, implementations SHOULD compute an `extensionsDigest` (SHA-256 of the
M/L-tagged canonical form of the extensions object) and include it in the
signing tuple.

Retention behaviour is normative: `extensionsDigest` binds the extensions
content to the attestor's signature, but a producer that has recorded and
signed a record MUST NOT strip `extensions` from the stored record or from
the DSSE-enveloped Statement afterwards. The record canonical form
committed to by the chain hash covers the full record including
`extensions`, and the DSSE envelope covers the full Statement including
`predicate.extensions`; stripping either post-signing corrupts the chain
hash for every downstream record and invalidates the envelope signature. A
producer that wants to omit the extensions payload from the storage layer
altogether MUST do so at emission time, before the chain hash is computed:
emit the record with `extensionsDigest` set and `extensions` absent,
carrying the extensions payload out-of-band if needed. The
`extensionsDigest` protects a never-inlined-payload workflow, not an
inline-then-strip workflow.

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

This is the record as written to the gateway's `audit.jsonl` file, in the
record canonical form (RFC 8785: members sorted by UTF-16 code unit, no
insignificant whitespace), including the `attestation` member. It is the
artifact whose chain hash becomes the subject digest:

```json
{"aiInvocation":{"model":"claude-sonnet-4-20250514","sessionId":"conv_abc123","turnId":"turn_1"},"attestation":"d1618968e997bd6423e2951a047a9bf6746bf2a56023d11c3cf81be984701027","attestorVersion":"mcp-audit-gateway/0.4.0","configHash":"f9a6ceae76b6df4bb5424fc76a204c443fa85a4384bd5ddf4bd5f532e3610cb3","contentDigestRequest":"b47ce84f9c856142547199358dd70203c504fbe3c82f362f46248bb05cf133cc","contentDigestResponse":"12fcbdcf920251bd7596b9e255eca20a314a179d5c33a70505d764706469825a","durationMs":1247,"id":"bf7a2f62-4d0f-4cce-afd2-cbfbf7bca2a5","method":"tools/call","namespace":"github","parties":[{"party":"gateway","role":"witness","scope":["id","type","timestamp","method","toolName","namespace","upstream","principal","durationMs","success","errorCode","errorClass","previousHash","contentDigestRequest","contentDigestResponse","attestorVersion","configHash"]},{"party":"client","role":"asserter","scope":["aiInvocation"]}],"previousHash":"genesis","principal":"user:alice@example.com","success":true,"timestamp":"2026-08-18T14:32:01.998Z","toolName":"github/create_pull_request","type":"tool_call","upstream":"github-server"}
```

Chain hash = SHA-256 of exactly these bytes (the line terminator is not part
of the preimage) = `6b1ded60d7b715d4bd97117ef5f5e24687b30517482c225b1ac768a14309ea8a`

The record carries the content digests, `attestorVersion`, and `configHash`
at the record layer, so both the signing tuple and the chain hash cover
them; their preimages are shown in the next subsection.

The attestation field is `HMAC-SHA256(key, signing_canonical_form)`. For
reproducibility, this example uses the trivial key `00` repeated 32 bytes
(`0000...0000`, 64 hex chars); a published key proves nothing, and this
choice is what makes every value on this page recomputable. Deployments
claiming attestor accountability use Ed25519 instead (see Security
Considerations, "Signature scheme and trust model"). The signing canonical
form of this record, built from the field list in Canonicalization, is:

```json
[["id","bf7a2f62-4d0f-4cce-afd2-cbfbf7bca2a5"],["type","tool_call"],["timestamp","2026-08-18T14:32:01.998Z"],["method","tools/call"],["toolName","github/create_pull_request"],["namespace","github"],["upstream","github-server"],["principal","user:alice@example.com"],["durationMs",1247],["success",true],["errorCode",null],["errorClass",null],["previousHash","genesis"],["contentDigestRequest","b47ce84f9c856142547199358dd70203c504fbe3c82f362f46248bb05cf133cc"],["contentDigestResponse","12fcbdcf920251bd7596b9e255eca20a314a179d5c33a70505d764706469825a"],["attestorVersion","mcp-audit-gateway/0.4.0"],["configHash","f9a6ceae76b6df4bb5424fc76a204c443fa85a4384bd5ddf4bd5f532e3610cb3"],["aiInvocation",["M",[["model","claude-sonnet-4-20250514"],["sessionId","conv_abc123"],["turnId","turn_1"]]]],["parties",[{"party":"gateway","role":"witness","scope":["id","type","timestamp","method","toolName","namespace","upstream","principal","durationMs","success","errorCode","errorClass","previousHash","contentDigestRequest","contentDigestResponse","attestorVersion","configHash"]},{"party":"client","role":"asserter","scope":["aiInvocation"]}]]]
```

### Content digest preimages

Request digest preimage. JCS of the JSON-RPC `params` member:

```json
{"arguments":{"base":"main","body":"Corrects a spelling mistake in line 42","head":"fix-typo","title":"Fix typo in README"},"name":"github/create_pull_request"}
```

SHA-256 = `b47ce84f9c856142547199358dd70203c504fbe3c82f362f46248bb05cf133cc`

Response digest preimage. This call succeeded, so the preimage is JCS of the
JSON-RPC `result` member:

```json
{"content":[{"text":"Pull request #123 created successfully","type":"text"}]}
```

SHA-256 = `12fcbdcf920251bd7596b9e255eca20a314a179d5c33a70505d764706469825a`

Had the call failed, the response would carry no `result` member and the
preimage would be JCS of the `error` member instead. For example, for the
error member:

```json
{"code":-32602,"message":"Invalid params"}
```

SHA-256 = `0bf47b3590aca56f86def7fc8fd3445aa813579274a6cfd24b0d7f6b9c442d8f`

### In-toto Statement

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "session:agent-workspace-4f2a",
    "digest": {
      "sha256": "6b1ded60d7b715d4bd97117ef5f5e24687b30517482c225b1ac768a14309ea8a"
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
        "scope": ["id", "type", "timestamp", "method", "toolName", "namespace", "upstream", "principal", "durationMs", "success", "errorCode", "errorClass", "previousHash", "contentDigestRequest", "contentDigestResponse", "attestorVersion", "configHash"]
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
      "request": { "sha256": "b47ce84f9c856142547199358dd70203c504fbe3c82f362f46248bb05cf133cc" },
      "response": { "sha256": "12fcbdcf920251bd7596b9e255eca20a314a179d5c33a70505d764706469825a" }
    },
    "metadata": {
      "attestorVersion": "mcp-audit-gateway/0.4.0",
      "configHash": "f9a6ceae76b6df4bb5424fc76a204c443fa85a4384bd5ddf4bd5f532e3610cb3"
    }
  }
}
```

This record is the genesis record, so the subject digest `6b1ded60...` is its
own chain hash: the stable chain identifier every subsequent statement in
this chain also carries as its subject digest. It is the `previousHash` the
next record in this chain will carry.

### Checkpoint example

Emitted after the single record above. This is a toy chain (1 record) for
verifiability; production deployments typically checkpoint every 100 records or
60 seconds.

Underlying checkpoint JSONL line (record canonical form):

```json
{"attestation":"0a8e65aaffc5c0351891b160359ac425e9b543f9fab06ed925f91fbe35db8989","id":"ckpt_9a1b2c3d-4e5f-6789-abcd-ef0123456789","parties":[{"party":"gateway","role":"witness","scope":["sequence","recordCount","previousHash"]}],"previousHash":"6b1ded60d7b715d4bd97117ef5f5e24687b30517482c225b1ac768a14309ea8a","recordCount":1,"sequence":1,"timestamp":"2026-08-18T14:33:42.101Z","type":"checkpoint"}
```

Chain hash = `SHA-256(above)` = `8df1e3359ce163bccae2ed5e998e51440bc1874d871cb65e9332a5e4f105f7cd`

In-toto Statement:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "session:agent-workspace-4f2a",
    "digest": { "sha256": "6b1ded60d7b715d4bd97117ef5f5e24687b30517482c225b1ac768a14309ea8a" }
  }],
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "checkpoint",
      "timestamp": "2026-08-18T14:33:42.101Z"
    },
    "chain": {
      "previousHash": "6b1ded60d7b715d4bd97117ef5f5e24687b30517482c225b1ac768a14309ea8a"
    },
    "checkpoint": {
      "sequence": 1,
      "recordCount": 1,
      "previousHash": "6b1ded60d7b715d4bd97117ef5f5e24687b30517482c225b1ac768a14309ea8a"
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

The subject digest is the chain identifier (the genesis record's chain hash),
not the checkpoint's own chain hash. `chain.previousHash` is the linkage;
`checkpoint.previousHash` restates the same value as the head being
externalized. On a checkpoint record the two are equal by construction.

The consumer externalizes `{"sequence": 1, "recordCount": 1, "previousHash":
"6b1ded60..."}` to a store the attestor cannot write. On the next verification
pass, if the chain head differs from `6b1ded60...` at record count 1, or if
fewer records exist, truncation is detected. The next record in the chain
carries `previousHash: "8df1e335..."`, the checkpoint's own chain hash, so
deleting the checkpoint breaks linkage.

### Chain break example

Emitted after a crash recovery. The next record chains from this record's
hash (not from "genesis").

Underlying chain_break JSONL line (record canonical form):

```json
{"attestation":"b4a0b7507f8284fa56f70adb3a432d4a85753ff54893a2615ff23f96a8cbf30e","id":"brk_7f8e9d0c-1b2a-3456-cdef-0123456789ab","priorHead":"8df1e3359ce163bccae2ed5e998e51440bc1874d871cb65e9332a5e4f105f7cd","priorRecordCount":1,"priorSequence":1,"reason":"crash_recovery","timestamp":"2026-08-18T15:01:00.000Z","type":"chain_break"}
```

Chain hash = `SHA-256(above)` = `290c37cb1b132b13c55fa7f5222e4adfe2324f1a07f06f526d13fdaa744a5952`

In-toto Statement:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "session:agent-workspace-4f2a",
    "digest": { "sha256": "6b1ded60d7b715d4bd97117ef5f5e24687b30517482c225b1ac768a14309ea8a" }
  }],
  "predicateType": "https://in-toto.io/attestation/ai-agent-action/v0.1",
  "predicate": {
    "action": {
      "type": "chain_break",
      "timestamp": "2026-08-18T15:01:00.000Z"
    },
    "chain": {
      "previousHash": "8df1e3359ce163bccae2ed5e998e51440bc1874d871cb65e9332a5e4f105f7cd"
    },
    "chainBreak": {
      "reason": "crash_recovery",
      "priorHead": "8df1e3359ce163bccae2ed5e998e51440bc1874d871cb65e9332a5e4f105f7cd",
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
above), and `chain.previousHash` carries the same value: the break links
backward when the prior head is known. The next record in the new chain
carries `previousHash: "290c37cb..."`, the break record's own chain hash,
binding the discontinuity evidence into the successor chain.

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
guidance. Deleting a checkpoint record does not evade detection: checkpoints
are chained like every other record (see Chain shape), so removing one breaks
linkage.

### Chain fork

An adversary with write access to the record store can construct a chain
fork: two records that share a common `previousHash` value, producing
parallel histories that share a prefix through the fork point. Every
`previousHash` in a presented sub-log resolves and no hash is broken; the
fork is detectable only by enforcing injectivity of the predecessor
relation (see "Chain shape") and by
reconstructing the chain by strict walk from genesis rather than by
set-membership check of `previousHash` values.

The concrete attack shape: a consumer targets the chain by subject digest,
where the subject digest equals the genesis record's chain hash. The
adversary constructs a second record that also chains from genesis (its
`previousHash` equals the genesis chain hash, exactly as an honest
successor's would). The adversary presents a two-record sub-log containing
the genesis record and the parallel-branch successor. Every hash link is
valid; the second-`"genesis"` rule (see Chain shape) sees only one
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
record breaks hash linkage. An adversary cannot silently remove the scar.
The literal string `"genesis"` appears exactly once, at the true start of the
log; a verifier MUST reject, fail-closed, a log in which it appears more than
once (see Chain shape).
Implementations that persist chain state to disk SHOULD validate the state
file's integrity on startup (e.g., separate MAC over the state file, or
storing the last hash in a second location).

The `chain_break` record type with `priorHead: null` is a residual risk that
implementations MUST mitigate. When the prior chain head is unrecoverable,
the break carries no chain linkage backward, so an attestor with intent to
hide records A-B-C can emit a `chain_break` with `reason: "crash_recovery"`
and `priorHead: null` to root a new chain segment that has no cryptographic
tie to the abandoned prefix. No hash link is broken, the second-`"genesis"`
rule is not tripped, and the presenter chooses whether the abandoned prefix
is discoverable. A conformant deployment claiming resistance against a
compromised attestor MUST mitigate this by one of: (1) requiring that any
`chain_break` with `priorHead: null` be preceded by an externalized
checkpoint whose head a consumer can compare against any produced-elsewhere
successor chain segment; (2) requiring an out-of-band co-attestation over
the `chain_break` record's `id`, `timestamp`, `reason`, and its own record
chain hash, signed by a distinct party whose key the compromised attestor
does not hold, and published to a verifier-reachable channel independent of
the attestor (v0.1 does not admit `parties` on `chain_break` records, so
the co-attestation is a sibling in-toto statement rather than a member of
the break record itself); or (3) forbidding `priorHead: null` in the
deployment profile. Absent one of these, a planted break is an accepted
residual risk that a verifier should surface as an unattested discontinuity
rather than silently accepting.

### Signature scheme and trust model

The signing form is used with either HMAC-SHA256 (symmetric) or Ed25519
(asymmetric). Deployments claiming attestor accountability (the compliance
use case in the Use Cases section) MUST use Ed25519 or another asymmetric
scheme; HMAC-SHA256 is permitted only for internal single-trust-boundary
deployments where a verifier holding the shared key is also trusted with
the ability to forge records. HMAC-SHA256's symmetry means anyone holding
the key, including the attestor, any provisioned verifier, and any operator
with root access on the gateway host, can produce records that verify
identically to authentic ones. The compliance-auditing framing rests on the
attestor being a distinct trust boundary from the verifier; that boundary
is preserved by Ed25519 (the attestor holds a private signing key the
verifier does not) and is collapsed by HMAC (a shared secret). A deployment
that uses HMAC MUST document the trust boundary explicitly and MUST NOT
claim non-repudiation against the attestor.

### Safe integer bound

All integer fields (`durationMs`, `errorCode`, checkpoint `sequence` and
`recordCount`) MUST remain within the I-JSON safe integer range (absolute
value < 2^53, per RFC 7493 Section 2.2). Implementations MUST reject records
with integers at or above this bound. At one checkpoint per 100 records and
one record per millisecond, the checkpoint sequence provides approximately
29 million years of headroom before overflow.

### Three canonical forms

The predicate names three canonicalization procedures: the record canonical
form (RFC 8785 over the complete audit record; the chain hash preimage), the
signing canonical form (tuple-array; the signature preimage), and the content
digest form (RFC 8785 over MCP payload members). The signing form rejects
floats and uses M/L type tags: properties that make it injective and safe
for structured attestation fields. Content digests use RFC 8785 because MCP
tool payloads are arbitrary JSON that may contain floats, and the registry
should converge on a single content-digest scheme (aligning with
in-toto/attestation#570). Implementations MUST NOT substitute one form for
another: chain hashes use the record canonical form, signatures the signing
form, payload digests the content digest form.

## Changelog and Migrations

### v0.1.0 (current)

Initial version. Changes from the original submission based on the first
review round:

-   Added `checkpoint` and `chain_break` as `action.type` values with their
  own sub-schemas, addressing the truncation detection gap
-   Added `parties` array with `witness`/`asserter` roles to distinguish
  field provenance (intermediary-observed vs. client-declared)
-   Documented genesis convention (`previousHash: "genesis"`) and chain_break
  requirement for crash recovery
-   Added I-JSON safe integer bound on all integer fields (< 2^53, RFC 7493)
-   Clarified deployment-dependent provenance of `principal`
-   Replaced placeholder example with fully recomputable worked example
-   Added Security Considerations section
-   Moved listing from "Vetted Predicates" to community contributions per
  ITE-63 process

Canonicalization hardening from the second review round:

-   Replaced the `JSON.stringify` chain-hash definition with the record
  canonical form: RFC 8785 (JCS) over the complete audit record including its
  `attestation` member. Producers MUST write the JSONL line as exactly those
  bytes; verifiers MUST recompute and reject on byte mismatch (fail-closed),
  so a log shipper that reparses and re-emits can no longer silently move the
  preimage
-   Enumerated the signing canonical form's complete field list and order for
  all three record types in the specification text, including null-vs-omitted
  behavior and the fixed member order of `parties` entries
-   Made strict I-JSON statement-wide: duplicate members at any depth reject
  fail-closed; well-formed-string rules applied to raw bytes and `\u`
  escapes; the 128-level depth bound now covers the whole record, not only
  `extensions`
-   Constrained `previousHash` to lowercase 64-hex or the literal `"genesis"`;
  no case folding, no other digest lengths, no algorithm agility at v0.1
-   Added the injective-predecessor and strict-walk verification rules
  (exactly one record per `previousHash` value; reject forks fail-closed)
-   Wired checkpoint linkage: every record type carries
  `predicate.chain.previousHash`; `checkpoint.previousHash` restates the head
  and MUST equal it; deleting a checkpoint now breaks documented linkage
-   Upgraded second-genesis detection from SHOULD to MUST reject (fail-closed)
-   Pinned content digest preimages: the JSON-RPC `params` member for requests,
  the `result` member for success responses, the `error` member for error
  responses; floats are permitted in content payloads only
-   Fixed the worked example's subject digests: every statement in a chain
  carries the genesis record's chain hash (the stable chain identifier), not
  the record's own chain hash
-   Added the Underlying record shape section pinning which members the JSONL
  record carries per record type and how they map to Statement fields
-   Widened the `tool_call` signing tuple to cover `type`, `errorClass`, the
  content digests, `attestorVersion`, and `configHash` (adopting the
  reviewer's proposed additions), so payload binding and configuration are
  signature-covered at v0.1
-   Documented the three protection layers (signing tuple, chain hash, DSSE
  envelope) and the signature-scheme trust model: Ed25519 for attestor
  accountability, HMAC-SHA256 only inside a single trust boundary
-   Named the planted-break residual risk (`chain_break` with `priorHead`
  null) and its required mitigations
-   Recomputed the worked example under the record canonical form and the
  widened signing tuple

## Acknowledgments

The record canonical form, the strict I-JSON profile, the chain-shape rules,
and much of the canonicalization language in this document are adapted, with
permission, from replacement prose contributed by Sankalp Gilda (@astrogilda)
during review of in-toto/attestation#588, itself adapted from the
canonicalization text of in-toto/attestation#570. An externally authored
conformance corpus for this predicate is maintained at
https://github.com/astrogilda/aee-conformance.

### Deviations from the vendored prose

The following four points depart from the replacement prose deliberately.
Each is documented in the normative sections referenced.

- **`chain_break` with `priorHead: null`** carries no `predicate.chain`
  object and roots a new segment. The vendored text has `chain_break` carry
  `predicate.chain.previousHash` in every case; when the crash also lost the
  chain state, there is nothing truthful to place there. Security
  Considerations names the resulting planted-break residual risk and three
  implementable mitigations.
- **`tool_call` signing tuple retains `method`, `upstream`, `principal`** at
  positions 4, 7, 8. The vendored list drops all three. They stay because
  the `parties` scopes reference them and dropping would hollow provenance;
  the retention rule is stated in the signing canonical form subsection.
- **Line terminator is not part of the chain-hash preimage.** The chain
  hash preimage ends at the closing brace. Vendored prose left this
  implicit.
- **On a checkpoint record, `predicate.chain.previousHash` MUST equal
  `predicate.checkpoint.previousHash` byte-for-byte.** Both derive from the
  single record `previousHash`. Vendored prose left the equality implicit;
  it is now a MUST-reject on mismatch.
