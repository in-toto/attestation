# Predicate type: AI Agent Decision

Type URI: https://contextpassport.com/attestation/agent-decision/v0.1

Version: v0.1

Predicate Name: AI Agent Decision

## Purpose

Software supply chains increasingly contain decisions no human made. An agent
triages a vulnerability report and closes it as not exploitable. An agent
approves a dependency bump. An agent decides a failing test is flaky and re-runs
it. An agent drafts and merges a patch.

Each of those changes what ships. None of them currently produces an
attestation, so a policy engine that can verify *how* an artifact was built
cannot ask *who or what decided it should be built that way*.

This predicate expresses one such decision, together with enough information to
verify that the record of the decision has not been altered since it was made.
It embeds a [Context Passport] record, a CC0 format whose integrity block is a
SHA-256 chain computed over the [RFC 8785] (JCS) canonical form of the payload.

## Use Cases

A policy engine or human reviewer needs to answer questions that no existing
predicate covers:

1.  Was this artifact produced or approved by an automated agent, and which
    model and provider?
2.  Was a human in the loop, and did they override the agent's conclusion?
3.  Has the decision record been altered since it was made? Each record commits
    to the hash of the decision before it, so an edit anywhere in a sequence is
    detectable.
4.  Was a decision taken automatically at a confidence below the threshold that
    should have required escalation?
5.  Across a release, were there agent decisions with no corresponding human
    review where policy required one?

### Why the existing predicates do not cover this

-   [SLSA Provenance] records how an artifact was built: builder, source,
    parameters. It says nothing about a judgement made during the process.
    Agent metadata could be placed in `internalParameters`, but that field is
    explicitly unstructured and not intended to carry claims a verifier acts on.
-   [Test Result] records what a test harness concluded. An agent decision is
    not a test with a pass or fail; it has a rationale, a confidence, and
    sometimes a human who overruled it.
-   [SCAI Report] is the closest fit, and is deliberately general enough to
    express much of this. Two things argue for a distinct type: decisions have a
    recurring shape worth standardising rather than re-deriving per producer,
    and SCAI has no equivalent of the `record` field, which binds an attestation
    to a tamper-evident sequence of *earlier* decisions. Maintainer input on
    whether that is sufficient justification is explicitly welcome.
-   SPDX and CycloneDX describe composition, not judgement.

## Prerequisites

The in-toto Attestation Framework.

Producing the `record` field requires an implementation of the [Context
Passport] specification v2.0 or later, which defines the canonicalization and
hashing rules. Verifying it requires only SHA-256 and an RFC 8785
canonicalizer: no network access, account, or trust in the attestation's issuer
is needed.

## Model

This predicate applies to any supply chain step where an automated agent
reaches a conclusion that affects what is produced or released, and to the
functionaries who review those conclusions.

The subject is the artifact the decision concerns. A predicate of this type
SHOULD NOT be used with the decision record itself as the subject: the record is
the evidence, not the thing being attested about.

Two independent checks apply, and they answer different questions. Envelope
signature verification establishes who issued the attestation. Recomputing the
hashes in `record` establishes that the decision has not been rewritten since it
was made, including by whoever signed the attestation. The second is the reason
this predicate carries hashes at all: an attestation signed by the party that
made the decision proves they said it, while the hash chain proves they have not
since changed what they said.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ "name": "...", "digest": { "sha256": "..." } }],

  // Predicate:
  "predicateType": "https://contextpassport.com/attestation/agent-decision/v0.1",
  "predicate": {
    "decision": {
      "agentId": "string",
      "agentName": "string",
      "provider": "string",
      "model": "string",
      "eventType": "string",
      "decidedAt": "<TIMESTAMP>",
      "input": "string",
      "output": "string",
      "confidence": 0.0
    },
    "humanOversight": {
      "reviewed": true,
      "overrodeAgent": false,
      "reviewerRef": "string",
      "authority": "string",
      "reviewedAt": "<TIMESTAMP>"
    },
    "record": {
      "id": "string",
      "schemaVersion": "string",
      "payloadHash": "sha256:...",
      "parentHash": "sha256:...",
      "integrityHash": "sha256:..."
    },
    "chain": {
      "rootId": "string",
      "position": 0,
      "length": 0
    }
  }
}
```

### Parsing Rules

This predicate follows the standard [parsing rules]. In particular, consumers
MUST apply the monotonic principle, and one consequence of it is load-bearing
here:

An absent `humanOversight` object means **no claim is made** about whether a
human was involved. It does not assert that none was. A policy requiring human
review MUST require `humanOversight.reviewed` to be present and `true`, and MUST
NOT infer anything from the object's absence.

Similarly, an absent `confidence` means the agent did not report one. It does
not mean zero.

Versioning follows the framework convention: the version is encoded in the
`predicateType` URI, and any change to the meaning or requiredness of an
existing field requires a new version.

### Fields

`decision` _object, required_

> The decision itself.

`decision.agentId` _string, required_

> Stable identifier for the deciding agent. Self-asserted by the producer and
> not verified by this predicate. Producers needing verifiable agent identity
> SHOULD use a resolvable identifier such as a W3C DID.

`decision.agentName` _string, required_

> Human-readable name of the agent.

`decision.provider` _string, optional_

> The organization or platform serving the model, for example `anthropic`,
> `openai`, or `self-hosted`.

`decision.model` _string, optional_

> The model that produced the decision.

`decision.eventType` _string, required_

> The kind of decision, drawn from the Context Passport event type registry.
> Custom types are permitted and SHOULD be namespaced.

`decision.decidedAt` _[Timestamp], required_

> When the agent reached the decision. Not the time the attestation was
> generated.

`decision.input` _string, required_

> What the agent was asked. See the privacy note below.

`decision.output` _string, required_

> What the agent concluded.

`decision.confidence` _number, optional_

> Confidence as reported by the agent, in the range 0.0 to 1.0 inclusive.
> Absent when the agent reported none.

`humanOversight` _object, optional_

> Present only when a human reviewed the decision. See the parsing rules: its
> absence asserts nothing.

`humanOversight.reviewed` _boolean, required_

> Whether a human reviewed the decision.

`humanOversight.overrodeAgent` _boolean, required_

> Whether the human's conclusion differed from the agent's.

`humanOversight.reviewerRef` _string, optional_

> An opaque handle for the reviewer. See the privacy note below.

`humanOversight.authority` _string, optional_

> The delegation, policy, or role under which the reviewer acted.

`humanOversight.reviewedAt` _[Timestamp], optional_

> When the human acted.

`record` _object, required_

> The Context Passport integrity block for this decision.

`record.id` _string, required_

> The Context Passport record identifier.

`record.schemaVersion` _string, required_

> The Context Passport specification version the record was written against,
> for example `2.0`. Determines which canonicalization rules apply.

`record.payloadHash` _string, required_

> `sha256:` followed by the hex SHA-256 of the RFC 8785 canonical payload.

`record.parentHash` _string or null, required_

> The `integrityHash` of the preceding record, or `null` for the first record
> in a chain.

`record.integrityHash` _string, required_

> `sha256:` followed by the hex SHA-256 over `payloadHash` and `parentHash`, as
> defined in the Context Passport specification.

`chain` _object, optional_

> Where this record sits in a sequence.

`chain.rootId` _string, optional_

> Identifier of the first record in the sequence.

`chain.position` _integer, optional_

> Zero-based index of this record within the sequence.

`chain.length` _integer, optional_

> Number of records in the sequence at the time the attestation was generated.

#### Privacy

`decision.input` and `decision.output` are free text and will frequently
contain the substance of a decision about a person. Producers SHOULD NOT place
identifying information in them, and SHOULD refer to data subjects and
reviewers by opaque handles, as `reviewerRef` does.

An attestation is intended to be shown to third parties. Anything placed in one
should be assumed readable by everyone it is ever shown to, permanently,
because the hash makes later redaction detectable rather than silent.

## Example

An agent triages a CVE reported against a bundled dependency and concludes the
vulnerable path is unreachable. A human security engineer disagrees and
escalates. Both facts are recorded.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "acme-api-server:2.4.1",
      "digest": {
        "sha256": "5f2c9e41a7b03d8c6e15f9a2d47b8c0e3f6a1d9b4c8e2f7a0d3b6c9e1f4a7d2c"
      }
    }
  ],
  "predicateType": "https://contextpassport.com/attestation/agent-decision/v0.1",
  "predicate": {
    "decision": {
      "agentId": "agent-secops-01",
      "agentName": "Vulnerability Triage Agent",
      "provider": "anthropic",
      "model": "claude-sonnet-5",
      "eventType": "audit",
      "decidedAt": "2026-03-29T10:01:00Z",
      "input": "Assess CVE-2026-1471 in transitive dependency parse-yaml@3.2.0 against this service.",
      "output": "not_exploitable: the vulnerable parser path is unreachable from any exported entrypoint",
      "confidence": 0.83
    },
    "humanOversight": {
      "reviewed": true,
      "overrodeAgent": true,
      "reviewerRef": "reviewer_2c9d",
      "authority": "Security Engineering on-call, policy SEC-2026-03",
      "reviewedAt": "2026-03-29T14:22:00Z"
    },
    "record": {
      "id": "ctx_1774778460000_d0c1de01a1b2",
      "schemaVersion": "2.0",
      "payloadHash": "sha256:9f1c7d4e2a8b05c3f6d9e1a4b7c0d3e6f9a2b5c8d1e4f7a0b3c6d9e2f5a8b1c4",
      "parentHash": null,
      "integrityHash": "sha256:31ab6f9c2d5e8a1b4c7d0e3f6a9b2c5d8e1f4a7b0c3d6e9f2a5b8c1d4e7f0a3d"
    },
    "chain": {
      "rootId": "ctx_1774778460000_d0c1de01a1b2",
      "position": 0,
      "length": 2
    }
  }
}
```

The agent was wrong, a human caught it, and the record says so. That is the
shape of evidence an auditor asks for, and it is what the current predicate set
cannot express.

## Changelog and Migrations

Initial version.

[Context Passport]: https://github.com/contextpassport/spec
[RFC 8785]: https://www.rfc-editor.org/rfc/rfc8785
[SCAI Report]: scai.md
[SLSA Provenance]: https://slsa.dev/provenance
[Test Result]: test-result.md
[Timestamp]: ../v1/field_types.md#timestamp
[parsing rules]: ../v1/README.md#parsing-rules
