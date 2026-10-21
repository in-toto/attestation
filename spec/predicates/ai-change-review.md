# Predicate type: AI Change Review

Type URI: https://noru.tech/spec/ai-change-provenance/review/v0.1

Version: 0.1

Predicate Name: AI Change Review

## Purpose

State, for one software change, who reviewed its head commit and what they
decided, signed by the reviewer's tooling. When the reviewer is a coding agent,
the predicate also records the facts only its tooling can state:

-   the human who directed it;
-   who owns the instructions it was given;
-   the model it ran.

The predicate is the "review document" of the
[AI Change Provenance specification](https://noru.tech/spec/ai-change-provenance/v0.3),
section 3.7 ("the specification" below). It is consumed by the
[AI Change Provenance](ai-change-provenance.md) evaluation predicate, next to
[AI Change Authorship](ai-change-authorship.md), which states who directed the
agent that wrote the change.

Review predicates proposed so far carry the reviewer's identity in the
signature, not in the predicate. That works when the consumer verifies the
signature itself. It fails for a consumer that receives already-verified
attestations and cannot see who signed. This predicate names the reviewer in
the predicate, and binds that name to a review the forge already recorded.

## Use Cases

-   **Agent reviewers in four-eyes change control.** An organization may let a
    reviewing agent's approval count as the second pair of eyes, under strict
    conditions. The agent's approval must be signed, and it must be independent
    of the change's author. That means its operator is not the author's
    operator, its vendor differs from the authoring agent's vendor, its signing
    identity differs from the authorship claim's, and, optionally, its
    instructions are not owned by the author. Only the reviewer's tooling can
    state the operator and the instructions owner; the forge records neither.
    This predicate carries them.
-   **Signed human reviews.** A human reviewer's tooling can sign what they
    decided. The forge's record of the review then gains `signed` evidence,
    which a policy can require.
-   **Accountability for automated review.** An organization can show which
    human directed each reviewing agent, and who controls what it was told to
    check.

Adjacent work:

-   The human-review predicates proposed in
    [#151](https://github.com/in-toto/attestation/pull/151) stalled. They carry
    the reviewer in the signature and do not model agent reviewers.
-   [`source-review-coverage`](https://github.com/in-toto/attestation/pull/581)
    (proposed) records which reviews cover a source tree. It has no operator or
    instructions model.
-   The [SLSA Source track](https://slsa.dev/spec/v1.2/source-requirements)
    requires review by trusted persons at Level 4. The source control system
    asserts it and it is summarized in a Source VSA. It does not record agent
    reviewers.

## Prerequisites

-   The in-toto Attestation Framework: the
    [Statement v1](https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md)
    layer and a DSSE envelope for signing.
-   Section 2 of the specification defines the terms. Section 6.5 defines the
    independence dimensions that the agent-reviewer fields feed.

## Model

The predicate belongs to the review step, after authoring and before merge.
There are three functionaries:

-   **The reviewer**, a human or a coding agent, reviews the change on the
    forge. The forge records the review, including the reviewer's account, the
    commit reviewed and the decision.
-   **The producer and signer**: the reviewer's tooling (the reviewing agent's
    integration, or a human reviewer's client). It emits and signs a Statement
    about that same review.
-   **The consumer**, typically a change-control evaluator, verifies the
    signature against identities it trusts for review claims, matches the
    Statement to the forge's review, and adds its evidence and facts to that
    review.

The subject is the head commit the review applies to.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {"name": "<change identifier>", "digest": {"gitCommit": "<head commit>"}}
  ],

  // Predicate:
  "predicateType": "https://noru.tech/spec/ai-change-provenance/review/v0.1",
  "predicate": {
    "spec_version": "0.1",
    "reviewer": {"kind": "human | agent", "id": "<namespace:account>", "agent": "<agent> | null"},
    "decision": "approved | changes_requested | commented",
    "change": {"head_commit": "<commit>"},
    "submitted_at": "<RFC 3339>",
    // Agent reviewers only; a human reviewer's document MUST NOT carry them:
    "operator": {"id": "<namespace:name>"} | null,
    "instructions": {"owner": "<namespace:name>", "digest": "sha256:<hex> | null"} | null,
    "model": "<model> | null",
    "session": {"id": "<session id>", "started_at": "<RFC 3339>"} | null
  }
}
```

A JSON Schema (2020-12) for the predicate is published with the specification
as
[`schemas/review.schema.json`](https://github.com/noru-tech/agent-change-control/blob/v0.5.2/schemas/review.schema.json).

### Parsing Rules

-   **Binding to the change.** At least one subject MUST carry a `gitCommit`
    digest equal to `predicate.change.head_commit`, compared as hexadecimal
    without regard to case. A consumer MUST reject a Statement whose head
    differs from the change's head.
-   **It upgrades a forge review and never creates one.** A consumer MUST match
    the Statement to a review the forge recorded, on three things: the
    reviewer account (`reviewer.id`), the head commit, and the decision.
    -   A Statement that matches no recorded review is kept as *unmatched* and
        yields no review. A signed file must not approve a change that the
        platform never showed as approved.
    -   A matched Statement adds its evidence to that review. Its facts (the
        operator, the instructions owner, the model, and the signing identity
        the verifier established) are recorded on the review.
-   **Kind and agent must agree.** `reviewer.kind` MUST agree with the kind of
    the forge account, human or not. For an agent reviewer, `reviewer.agent`
    MUST agree with any mapping the consumer has verified from that account to
    an agent. Either disagreement is an error for the change, never a choice.
-   **Human reviewers.** A document whose reviewer is human MUST NOT name an
    `operator` or `instructions`. A consumer MUST treat one that does as an
    error.
-   **People resolve only to humans.** `operator.id` and `instructions.owner`
    resolve only to accounts the forge confirms are human. Anything else
    becomes null, meaning unknown. An agent reviewer whose operator is unknown
    can never count as independent on that dimension.
-   **Evidence strength.** The review's added evidence is `signed` only when
    the Statement arrived in a DSSE envelope with a signature that the
    consumer's operator verified against identities trusted for review claims.
    It records who verified it. Otherwise the evidence is `declared`. The
    signing identity recorded on the review is the one the verifier
    established, never one the predicate states.
-   **Unrecognized fields: a deviation from the framework's parsing rules.**
    Consumers MUST reject unrecognized fields in the predicate. A field a
    consumer does not understand could change what the review asserts. New
    fields come with a new type URI version. Unrecognized fields at the
    Statement layer follow the framework's rules.
-   **Digests and serialization.** A signed Statement's digest is over the DSSE
    payload bytes as signed. A Statement handed over without an envelope has
    its digest over its [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785)
    serialization. The predicate is I-JSON (RFC 7493) with these constraints,
    and a consumer rejects anything else:
    -   integer-only numbers in ±(2^53 − 1);
    -   no unpaired surrogates;
    -   unique member names;
    -   nesting of at most 128.

    Section 8.2 of the specification defines them.
    `instructions.digest` identifies the instructions the agent reviewer was
    given. The producer defines its preimage, and consumers compare it without
    interpreting it.
-   **Timestamps** are RFC 3339.
-   **Versioning.** `spec_version` MUST be `0.1` for this type URI. The
    predicate is versioned on its own.

### Fields

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `spec_version` | string | yes | `0.1`. |
| `reviewer.kind` | string | yes | `human` or `agent`. Must agree with the forge account's kind. |
| `reviewer.id` | string | yes | The forge account that submitted the review, `namespace:account`. |
| `reviewer.agent` | string or null | no | For an agent reviewer, the reviewing agent named by tool, for example `claude-code-review`. |
| `decision` | string | yes | `approved`, `changes_requested` or `commented`. Must agree with the forge's record of the review. |
| `change.head_commit` | string | yes | The head commit reviewed. It must appear as a `gitCommit` subject. |
| `submitted_at` | string (RFC 3339) | yes | When the review was submitted, as the producer recorded it. |
| `operator.id` | string, or `operator` null | no | Agent reviewers only: the human who directed the reviewing agent. |
| `instructions.owner` | string, or `instructions` null | no | Agent reviewers only: the human or team that owns the instructions the reviewing agent was given. |
| `instructions.digest` | string (`sha256:<64 hex>`) or null | no | Agent reviewers only: a digest identifying those instructions. |
| `model` | string or null | no | Agent reviewers only: the model the reviewing agent ran. It is recorded, and not used for independence. |
| `session` | object or null | no | Agent reviewers only: `id` and `started_at` (RFC 3339) of the review session. |

## How consumers use it

The AI Change Provenance evaluation predicate tests an agent's approval on four
independence dimensions, each one `independent`, `dependent` or `unknown`:

| Dimension | Independent when |
| --- | --- |
| `operator` | The reviewing agent's operator differs from the human behind the change. |
| `provider` | The reviewing agent's vendor differs from the authoring agent's vendor. |
| `identity` | This Statement's verified signer differs from every verified signer behind the authorship claim. |
| `instructions` | The instructions owner differs from the human behind the change. |

An agent's approval counts as independent only when three things hold: a policy
opts in, its evidence is `signed`, and every required dimension is
`independent`. `unknown` on a required dimension never counts. Rules ACC007 to
ACC010 of the specification report agent approvals and why they did or did not
qualify.

## Testing

This predicate has no dedicated producer-side corpus yet. What exists:

-   The evaluation predicate's conformance corpus exercises the agent-review
    rules. Its vectors carry agent review facts that are signed or only
    declared, and together they make each of the four dimensions decide an
    outcome. Each vector's expected outcome is written down.
    -   Suite revision 1, published at `v0.5.0`, covers the `operator`,
        `provider` and `identity` dimensions.
    -   Suite revision 2, published at `v0.5.2`, adds the opt-in
        `instructions` dimension. Its digest list (SHA-256
        `ffa39c97c2a7247bc1873749c73c620884df1e17c52bb74b0604d70b5c730c52`)
        is signed at that tag.
-   The reference consumer,
    [`acc`](https://github.com/noru-tech/agent-change-control), tests:
    -   matching, and unmatched Statements;
    -   agreement between the reviewer's kind and the forge account;
    -   agreement between `reviewer.agent` and a verified account-to-agent
        mapping;
    -   the rule that a human reviewer's document names no operator or
        instructions;
    -   the signer read from a verifier's output.

    Its example Statement is validated against the schema in CI.

A producer-side corpus, with Statements paired by one mutation, is planned.
Independent consumers are welcome.

## Security Considerations

-   **The forge record is authoritative for the review itself.** The Statement
    can only add evidence and facts to a review the forge recorded. It cannot
    create an approval, change a decision, or apply to another commit.
-   **The added facts are claims.** The operator and the instructions owner are
    stated by the reviewer's tooling. A dishonest integration could name
    someone else to make an agent reviewer look independent. Consumers limit
    this in three ways. They trust only specific signer identities for review
    claims. They resolve people only to forge-confirmed humans. And they test
    the `identity` dimension, the one dimension resting on a verified fact
    rather than a claim, so the same signer cannot vouch for both authorship
    and review.
-   **Agent approval is opt-in.** Under the evaluation predicate's default
    policy, an agent's approval never counts as independent, whatever this
    Statement says.

## Privacy Considerations

The predicate names reviewers, the humans directing reviewing agents, and the
owners of review instructions, with timestamps. That is employee activity
data. Before a Statement is published to a public transparency log, consider
that the entry is permanent. For private repositories, keep the attestations
in a store with the same access as the repository.

## Example

A reviewing agent, `claude-code-review` running as `github:acme-review[bot]`,
approved head `7c4a8d09…`. It was directed by `github:carol`, and its
instructions are owned by `github:security-team`:

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
  "predicateType": "https://noru.tech/spec/ai-change-provenance/review/v0.1",
  "predicate": {
    "spec_version": "0.1",
    "reviewer": {
      "kind": "agent",
      "id": "github:acme-review[bot]",
      "agent": "claude-code-review"
    },
    "decision": "approved",
    "change": {
      "head_commit": "7c4a8d09ca3762af61e59520943dc26494f8941b"
    },
    "submitted_at": "2026-08-02T12:00:00Z",
    "operator": {
      "id": "github:carol"
    },
    "instructions": {
      "owner": "github:security-team",
      "digest": null
    },
    "model": "claude-opus-5"
  }
}
```

Handed over without an envelope, the Statement's digest is the SHA-256 of its
RFC 8785 serialization. The preimage is exactly the bytes below, without the
newline that ends the code block:

```json
{"_type":"https://in-toto.io/Statement/v1","predicate":{"change":{"head_commit":"7c4a8d09ca3762af61e59520943dc26494f8941b"},"decision":"approved","instructions":{"digest":null,"owner":"github:security-team"},"model":"claude-opus-5","operator":{"id":"github:carol"},"reviewer":{"agent":"claude-code-review","id":"github:acme-review[bot]","kind":"agent"},"spec_version":"0.1","submitted_at":"2026-08-02T12:00:00Z"},"predicateType":"https://noru.tech/spec/ai-change-provenance/review/v0.1","subject":[{"digest":{"gitCommit":"7c4a8d09ca3762af61e59520943dc26494f8941b"},"name":"github:acme/api:pr:421"}]}
```

SHA-256: `59e6658506512ead7b7960f4914e656b38e88b533501648d1b3760b5f7c07f97`

## Changelog and Migrations

-   **0.1.** Initial version. It was introduced in revision 5 of AI Change
    Provenance 0.1, and the rules that read it arrived in specification 0.2.
    It is unchanged through 0.3.
