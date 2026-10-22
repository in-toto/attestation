# Predicate type: AI Change Authorship

Type URI: https://noru.tech/spec/ai-change-provenance/provenance/v0.1

Version: 0.1

Predicate Name: AI Change Authorship

## Purpose

State, for one software change, which coding agent produced it and which human
directed it, bound to the change's head commit and signed by the party that can
know: the agent's integration, or the human operator.

Development platforms record the account that opened a change. They do not
record whether a coding agent wrote it, or who directed that agent. When an
agent opens a change under its own account, the human whose judgment the change
embodies is missing from the record. A reviewer who is independent of the bot
account may be the very person who directed the agent. This predicate makes
that human explicit and binds the claim to the exact commit it covers, so that a
consumer can test the independence of approvals against the right person.

The predicate is the "provenance document" of the
[AI Change Provenance specification](https://noru.tech/spec/ai-change-provenance/v0.3),
section 3.2 ("the specification" below). It is named *authorship* here because
it records who authored a change, not how an artifact was built, and should not
be confused with SLSA Provenance. The type URI keeps the specification's path.
The evaluation predicate that consumes it,
[AI Change Provenance](ai-change-provenance.md), is proposed alongside.

## Use Cases

-   **Authorship evidence for four-eyes change control.** An agent integration
    emits this predicate when it opens a pull request. A change-control
    evaluator reads it back and tests every approval against the named operator
    rather than against the bot account. With the claim signed and verified, a
    policy can require `signed` authorship evidence before it accepts that an
    approval was independent.
-   **Accountability for agent-written changes.** An organization can show,
    per change, which human directed which agent. This is what AI governance
    frameworks such as ISO/IEC 42001 ask for as human oversight of AI-assisted
    output.
-   **Agent-opened changes without a human account in the loop.** Agents that
    commit and open pull requests under an app or bot account have no other
    place to record their operator that is bound to a revision and signed.

Adjacent records do not cover this:

-   Vendor `Co-Authored-By` commit trailers name the agent, but not the human
    who directed it. They are unsigned and live in mutable commit messages.
-   [Agent Trace](https://agent-trace.dev) records which model produced which
    lines. It does not name the directing human, and its storage is left to the
    implementation.
-   [AI Agent Action](https://github.com/in-toto/attestation/pull/588)
    (proposed) records an agent's tool calls through a gateway. It is not bound
    to a git commit.
-   [SLSA Provenance](https://slsa.dev/provenance) describes builds.

## Prerequisites

-   The in-toto Attestation Framework: the
    [Statement v1](https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md)
    layer and a DSSE envelope for signing.
-   Section 2 of the specification defines the terms agent, operator and
    change.

## Model

The predicate belongs to the authoring step, before review. There are two
functionaries:

-   **The producer and signer**: the agent's integration (for example the CI
    workflow or app that runs the agent), or the operator. It knows which agent
    ran and on whose instruction, and it signs the Statement.
-   **The consumer**: typically a change-control evaluator. It verifies the
    signature against identities it trusts for authorship claims, binds the
    claim to a change by its head commit, and treats the operator as the human
    whose independence reviewers are tested against.

The subject is the head commit of the change the agent produced. A new push to
the change needs a new Statement for the new head.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {"name": "<change identifier>", "digest": {"gitCommit": "<head commit>"}}
  ],

  // Predicate:
  "predicateType": "https://noru.tech/spec/ai-change-provenance/provenance/v0.1",
  "predicate": {
    "spec_version": "0.1",
    "agent": {"name": "<agent>", "version": "<agent version> | null"},
    "operator": {"id": "<namespace:name | name>"},
    "session": {"id": "<session id>", "started_at": "<RFC 3339>"},
    "change": {"base_commit": "<commit>", "head_commit": "<commit>"}
  }
}
```

A JSON Schema (2020-12) for the predicate is published with the specification
as
[`schemas/provenance.schema.json`](https://github.com/noru-tech/agent-change-control/blob/v0.5.1/schemas/provenance.schema.json).

### Parsing Rules

-   **Binding to the change.** At least one subject MUST carry a `gitCommit`
    digest equal to `predicate.change.head_commit`, compared as hexadecimal
    without regard to case. A consumer MUST reject a Statement whose
    `change.head_commit` differs from the head of the change it is applied to.
    Binding to the head is what stops a signed claim being replayed against a
    later, different change. The subject `name` is informative: by convention
    it is the change identifier, such as `github:acme/api:pr:421`.
-   **Agent name.** `agent.name` names the tool, never the account it used.
    Examples are `claude-code`, `codex` and `cursor`. It SHOULD use only
    `[A-Za-z0-9._-]`, the vocabulary of the specification's inline declaration,
    and the same spelling as any other claim for the change, because claims are
    compared exactly.
-   **Operator.** `operator.id` is either namespaced (`github:alice`) or bare
    (`alice`). A bare name resolves in the namespace of the forge that holds
    the change. A consumer MUST resolve it only to an account the forge
    confirms is a human. An email address or a display name never resolves, and
    an unresolved operator is recorded as unknown, never guessed.
-   **Evidence strength.** A consumer counts the claims as `signed` evidence
    only when the Statement arrived in a DSSE envelope with a signature that
    its operator verified against identities trusted to make authorship claims.
    It records who verified it. Otherwise the claims are `declared` evidence, a
    party's unverified assertion. Trust comes from the signer the consumer pins,
    never from anything in the predicate.
-   **Agreement.** When the same change also carries another authorship claim,
    such as a second Statement or the specification's inline declaration, the
    claims MUST name the same agent, and the same operator when both name one.
    Disagreement is an error for that change. A consumer never chooses between
    claims.
-   **Unrecognized fields: a deviation from the framework's parsing rules.**
    Consumers MUST reject unrecognized fields in the predicate. The claim is
    small and closed. A field a consumer does not understand could qualify the
    authorship (a second operator, a partial scope) in a way that silently
    ignoring it would misread. New fields come with a new type URI version.
    Unrecognized fields at the Statement layer follow the framework's rules.
-   **Digests and serialization.** The digest of a signed Statement is over the
    DSSE payload bytes exactly as signed. For a Statement handed over without an
    envelope, it is over the [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785)
    serialization of the Statement. The predicate is I-JSON (RFC 7493) with
    these constraints, and a consumer rejects anything else:
    -   integer-only numbers in ±(2^53 − 1);
    -   no unpaired surrogates;
    -   unique member names;
    -   nesting of at most 128.

    Section 8.2 of the specification defines them.
-   **Timestamps** are RFC 3339. `session.started_at` is when the agent session
    that produced the change started, as the producer recorded it.
-   **Versioning.** `spec_version` MUST be `0.1` for this type URI. The
    predicate is versioned on its own, apart from the evaluation predicate's
    version.

### Fields

All fields are required.

| Field | Type | Description |
| --- | --- | --- |
| `spec_version` | string | `0.1`. |
| `agent.name` | string, non-empty | The coding agent, named by tool, for example `codex`. |
| `agent.version` | string or null | The agent's version, as its integration reports it. `null` when unknown. |
| `operator.id` | string, non-empty | The human who directed the agent for this change: `namespace:name` or a bare name in the forge's namespace. |
| `session.id` | string, non-empty | The producer's identifier for the agent session. It is opaque to consumers and useful for cross-reference to the producer's own logs. |
| `session.started_at` | string (RFC 3339) | When that session started. |
| `change.base_commit` | string, non-empty | The commit the change is based on. |
| `change.head_commit` | string, non-empty | The head commit the claim covers. It must appear as a `gitCommit` subject. |

## Testing

This predicate has no dedicated conformance corpus yet. What exists:

-   The evaluator conformance corpus of the
    [evaluation predicate](ai-change-provenance.md) exercises signed authorship
    evidence. It has an accept vector in which a verified claim satisfies a
    policy that requires `signed` authorship evidence, and a reject vector for
    `signed` evidence that does not resolve to a verified attestation.
-   The reference consumer,
    [`acc`](https://github.com/noru-tech/agent-change-control), tests:
    -   binding by head commit, and rejecting a document for another head;
    -   agreement with inline declarations, where disagreement is an error;
    -   `signed` versus `declared` evidence;
    -   the three containers: a bare Statement, a DSSE envelope and a Sigstore
        bundle.

    Its example Statement is validated against the schema in CI.

A producer-side corpus is planned. It will pair accept and reject Statements
by one mutation, for example a mismatched head, an unknown field or a missing
operator. Independent consumers are welcome.

## Security Considerations

-   **A claim, not a fact.** The Statement proves who signed the claim, not
    that the claim is true. A compromised or dishonest integration could name
    the wrong operator. Naming a colleague would manufacture an appearance of
    independence. Consumers limit this by trusting only specific signer
    identities for authorship claims, such as the integration's workflow
    identity, and not every signer they accept for other predicates. They also
    resolve the operator only to forge-confirmed human accounts, and treat
    disagreement with any other claim as an error.
-   **Replay.** Binding to the head commit stops a claim being reused for a
    different change or a later push. A consumer must not accept a Statement
    whose head differs from the change's head.
-   **Missing claims.** The absence of this predicate proves nothing. A consumer
    that cannot establish the operator records it as unknown and never
    substitutes the merger, the opener or a reviewer.

## Privacy Considerations

The predicate names a person, the operator, and records when they ran an agent
session, which is employee activity data. `session.id` may let someone
correlate the claim with the producer's logs. Before a Statement is published
to a public transparency log, for example through keyless Sigstore signing,
consider that the entry is permanent. For private repositories, sign with a key
you control and keep the attestations in a store with the same access as the
repository.

## Example

An agent integration states that Codex produced the change whose head is
`7c4a8d09…`, directed by `github:alice`:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "github:acme/api:pr:421",
      "digest": {
        "gitCommit": "7c4a8d09ca3762af61e59520943dc26494f8941b"
      }
    }
  ],
  "predicateType": "https://noru.tech/spec/ai-change-provenance/provenance/v0.1",
  "predicate": {
    "spec_version": "0.1",
    "agent": {
      "name": "codex",
      "version": "2026.9.1"
    },
    "operator": {
      "id": "github:alice"
    },
    "session": {
      "id": "sess_7f3a9c",
      "started_at": "2026-09-18T09:12:00Z"
    },
    "change": {
      "base_commit": "5d41402abc4b2a76b9719d911017c592e2c1a5f0",
      "head_commit": "7c4a8d09ca3762af61e59520943dc26494f8941b"
    }
  }
}
```

Handed over without an envelope, the Statement's digest is the SHA-256 of its
RFC 8785 serialization. The preimage is exactly the bytes below, without the
newline that ends the code block:

```json
{"_type":"https://in-toto.io/Statement/v1","predicate":{"agent":{"name":"codex","version":"2026.9.1"},"change":{"base_commit":"5d41402abc4b2a76b9719d911017c592e2c1a5f0","head_commit":"7c4a8d09ca3762af61e59520943dc26494f8941b"},"operator":{"id":"github:alice"},"session":{"id":"sess_7f3a9c","started_at":"2026-09-18T09:12:00Z"},"spec_version":"0.1"},"predicateType":"https://noru.tech/spec/ai-change-provenance/provenance/v0.1","subject":[{"digest":{"gitCommit":"7c4a8d09ca3762af61e59520943dc26494f8941b"},"name":"github:acme/api:pr:421"}]}
```

SHA-256: `77a9aaed7f82b1501b52ea71918f248fd7c06155594607bf42135fda7d4bee04`

Signed, the same Statement is the payload of a DSSE envelope with payload type
`application/vnd.in-toto+json`. Its digest is then over the payload bytes as
the signer produced them.

## Changelog and Migrations

-   **0.1.** Initial version. It was introduced in revision 4 of AI Change
    Provenance 0.1 and is unchanged through specification 0.3.
