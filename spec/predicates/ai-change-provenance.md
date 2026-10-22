# Predicate type: AI Change Provenance

Type URI: https://noru.tech/spec/ai-change-provenance/v0.3

Version: 0.3

Predicate Name: AI Change Provenance

## Purpose

Record, for one or more software changes (pull requests or merge requests),
who the effective author was, meaning a human, or a coding agent together with
the human who directed it. Also record who reviewed and merged the change, and
whether a human independent of that effective author approved the reviewed
head before merge. The predicate carries the normalized facts, the policy they
were evaluated under, and the result: an assessment per rule per change, with
findings for the failures. A consumer can re-evaluate it offline and confirm
that the result follows from the facts.

The account-based four-eyes check, "the author is not the approver", stops
working once coding agents open changes under their own accounts. An engineer
directs an agent, the agent opens the pull request from a bot account, and the
engineer approves it. The platform sees two accounts and one human judgment.
AI Change Provenance (ACP) replaces the account check with "the effective human
is not an approving human, on the current head, before merge", or an explicit
`unknown` when the record cannot tell. The normative definition is the
[ACP specification](https://github.com/noru-tech/agent-change-control/blob/v0.5.2/spec/ai-change-provenance.md)
("the specification" below).

## Use Cases

-   **Four-eyes change control for agent-authored changes.** Change-control
    frameworks require approval by someone other than the person who made the
    change: SOC 2 CC8.1, ISO/IEC 27001:2022 A.8.32 and A.5.3, PCI DSS 6.5.1,
    and NIST SP 800-53 CM-3 and AC-5. A Statement per change is reproducible
    evidence for that control. It records whether the approver was independent
    of the human behind an agent, not only of the account that opened the
    change.
-   **Merge or deploy gate.** A CI job produces the Statement for a pull
    request's head. A policy engine refuses to merge or deploy a head without a
    Statement whose findings are clean and whose collection was complete.
-   **Release evidence.** The Statement for each change in a release is signed
    and stored next to the release's build provenance. The merge commit is a
    subject of both, so the change-control evidence and the build provenance
    are found together.
-   **Audit population.** One Statement per change over a reporting window is
    the population an auditor samples from. Each item answers the independence
    question, and any gap is marked incomplete rather than read as clean.

Existing and proposed predicates cover adjacent ground, but none of them models
the agent, the human who directed it, and the reviewer as separate parties with
an independence test between them:

-   [SLSA Provenance](https://slsa.dev/provenance) describes how an artifact
    was built, not how its source change was approved.
-   The [SLSA Source track](https://slsa.dev/spec/v1.2/source-requirements)
    defines two-party review at Source Level 4: changes to protected branches
    "MUST be agreed to by two or more trusted persons". The source control
    system asserts it, and it is summarized in a Source VSA
    (`SLSA_SOURCE_LEVEL_4`). A trusted person is an authorized human, but the
    track does not say what happens when one of the two accounts is an agent
    that a trusted person directed. ACP answers exactly that question, per
    change, and can back an organization-defined `ORG_SOURCE_` property in a
    Source VSA.
-   [`source-review-coverage`](https://github.com/in-toto/attestation/pull/581)
    (proposed) records which reviews cover a source tree. It has no notion of
    agent authorship or of an operator.
-   [AI Agent Action](https://github.com/in-toto/attestation/pull/588)
    (proposed) records what an agent did through a gateway, as a chain of tool
    calls. ACP records who is accountable for a change and whether someone
    independent approved it. The two are complementary: a gateway-signed chain
    that names its principal could become strong operator evidence for ACP once
    a binding from a tool-call chain to a git commit exists.
-   The human-review predicates proposed in
    [#151](https://github.com/in-toto/attestation/pull/151) stalled while the
    SLSA Source track settled.

## Prerequisites

-   The in-toto Attestation Framework: the
    [Statement v1](https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md)
    layer and a DSSE envelope for signing.
-   The terminology of the specification, section 2: change, actor, effective
    author, effective human, operator, qualifying approval, evidence,
    complete.
-   [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785) (JSON Canonicalization
    Scheme) and [RFC 7493](https://www.rfc-editor.org/rfc/rfc7493) (I-JSON),
    for digests.

## Model

The predicate belongs to the change-control step between writing source and
building it. The functionaries are:

-   **The forge** (GitHub, GitLab, …), which records the change, its commits,
    its review history and its merge. Facts read from it are `observed`
    evidence.
-   **The producer of an authorship claim**: the agent's integration or its
    operator. It declares which agent wrote the change and which human directed
    it, as an inline declaration (`declared`) or a signed provenance document
    (`signed` once verified). Records that agents write at authoring time, such
    as vendor commit trailers and Agent Trace records, give `derived`
    evidence.
-   **Reviewers**, human or agent. An agent reviewer's own signed review
    document can add who operated it, who owns its instructions, and its
    verified signing identity.
-   **The evaluator**, which collects the facts read-only, evaluates the rules
    deterministically and offline, and emits the Statement unsigned. The
    reference implementation is
    [`acc`](https://github.com/noru-tech/agent-change-control).
-   **The signer**, typically the CI workflow identity or an organization key,
    wraps the Statement in DSSE. Trust in the result comes from the signer's
    identity, which the verifier pins, and never from the predicate itself.

```text
Statement
├── subject: for each change
│   ├── <change id>        gitCommit = head commit (what the approvals are bound to)
│   └── <change id>:merge  gitCommit = merge commit (when merged, known and distinct)
├── predicateType: https://noru.tech/spec/ai-change-provenance/v0.3
└── predicate (manifest)
    ├── events      repository, window, actors, attestations, changes (facts with evidence)
    ├── policy      the resolved rules, severities, evidence minimums, agent_review
    ├── findings    one per failed rule per change, stable identifier, disposition
    ├── assessments pass | fail | unknown | not_applicable per rule per change
    ├── summary     counts: clean, with findings, indeterminate
    └── generated   tool, version, source digest over the events
```

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {"name": "github:acme/api:pr:421", "digest": {"gitCommit": "<head commit>"}},
    {"name": "github:acme/api:pr:421:merge", "digest": {"gitCommit": "<merge commit>"}}
  ],

  // Predicate:
  "predicateType": "https://noru.tech/spec/ai-change-provenance/v0.3",
  "predicate": {
    "version": "0.3",
    "events": {
      "version": "0.3",
      "repository": "<OWNER/REPO>",
      "window": {"from": "<RFC 3339>", "to": "<RFC 3339>", "complete": true, "reason": null},
      "actors": {"<namespace:name>": {"kind": "human | agent | bot | service | unknown", ...}},
      "attestations": {"attestation:<16 hex>": {...}},  // optional
      "changes": [{
        "id": "<change id>", "head_sha": "<commit>", "merge_commit_sha": "<commit> | null",
        "author": {...}, "agent_operator": {...} | null, "reviews": [...],
        "reviews_complete": true, ...
      }]
    },
    "policy": {...},
    "findings": [{"id": "acc-<16 hex>", "rule_id": "ACC001", ...}],
    "assessments": [{"change_id": "...", "rule_id": "ACC001", "status": "pass", "reason": "..."}],
    "summary": {...},
    "generated": {"tool": "<evaluator>", "version": "<version>", "source_digest": "sha256:<hex>"}
  }
}
```

A JSON Schema (2020-12) for the predicate is published with the specification
as
[`schemas/manifest.schema.json`](https://github.com/noru-tech/agent-change-control/blob/v0.5.2/schemas/manifest.schema.json).
The Statement shape `acc` emits is
[`schemas/statement.schema.json`](https://github.com/noru-tech/agent-change-control/blob/v0.5.2/schemas/statement.schema.json),
a strict subset of Statement v1. It validates the example below and the
manifest of every accept vector in the conformance corpus.

### Parsing Rules

-   **Subjects.** Each change contributes a subject named by its change
    identifier, with `digest.gitCommit` set to the head commit. The approvals
    are bound to that commit. A merged change whose merge commit is known and
    differs from the head also contributes `<change id>:merge`. That is the
    commit reachable from the target branch after a squash or rebase merge. A
    merge commit the forge did not report is never invented.
-   **Alternative subject for SHA-2-only signers.** A signer that accepts only
    SHA-2 subjects (GitHub artifact attestations, `cosign attest-blob`) may
    produce a single subject whose `sha256` is the digest of the predicate's
    RFC 8785 serialization. The commits are then found inside the predicate,
    in `head_sha` and `merge_commit_sha`.
-   **Verification.** A consumer:
    -   MUST check that the subjects cover the commits it is asking about.
    -   MUST check that the subjects are exactly those the predicate's changes
        produce, or that the single `sha256` subject recomputes from the
        predicate.
    -   SHOULD re-evaluate the predicate before trusting its findings. The
        embedded events evaluated under the embedded policy must reproduce the
        findings, assessments and summary exactly. The findings are derived;
        the events are the input, and a consumer recomputes derived data rather
        than trusting it.
-   **Incomplete and unknown.** A predicate with `events.window.complete:
    false`, or with a change whose `reviews_complete` is `false`, is an
    incomplete record and MUST NOT be read as clean, whatever its findings say.
    An assessment of `unknown` is never a pass.
-   **Canonicalization.** Every digest in the predicate, and every byte
    sequence that is signed, is SHA-256 over, or equal to, the
    [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785) serialization of the
    value. This is section 8.2 of the specification. A preimage never includes
    a trailing newline. The finding identifier is `acc-` followed by the first
    16 hexadecimal characters of SHA-256 over the JCS array
    `[repository, change_id, rule_id, actor_ids]`. `generated.source_digest`
    is over the JCS serialization of `events`.
-   **I-JSON.** The predicate is I-JSON (RFC 7493) with these constraints:
    -   numbers are integer literals, with no fraction or exponent, in
        ±(2^53 − 1);
    -   strings have no unpaired surrogates;
    -   member names are unique within an object;
    -   nesting is at most 128 deep, counting the outermost container as
        depth 1.

    A consumer rejects anything else. For the depth bound, it rejects without
    parsing further.
-   **Unrecognized fields: a deviation from the framework's parsing rules.**
    Consumers MUST reject unrecognized fields inside the predicate, rather than
    ignore them. The predicate must re-evaluate to exactly the same content, so
    a field the consumer does not understand would make that check
    meaningless. The one exception is `merge_commit_sha`, which may be absent
    in records made before it existed, and then means `null`. Unrecognized
    fields at the Statement layer follow the framework's rules.
-   **Timestamps** are RFC 3339 in UTC with `Z`. Every timestamp is a source
    fact. None of them is the time of evaluation, which is deliberately not
    recorded so that output is deterministic.
-   **Versioning.** The type URI carries the version, and the predicate's
    `version` member must agree with it. Additive changes increment the minor
    version. Before 1.0, a change to the meaning of a rule, identifier or
    format also increments the minor version, and the specification defines how
    the previous version's documents are still validated. Version 0.2 is
    validated with its legacy serialization, whose preimages ended in a newline
    (see Changelog). Rule identifiers are never reused.

### Fields

All fields are required unless marked optional. "Evidence" is a list of
`{source, ref, kind}`, where `kind` is one of `derived`, `declared`,
`observed` or `signed`, ordered by trust in that order. "Identity" is
`{actor_id, provenance: Evidence}`.

| Field | Type | Description |
| --- | --- | --- |
| `version` | string | Predicate version, `0.3`. Must match the type URI. |
| `events.version` | string | Version of the event export the facts were collected as: `0.1`, `0.2` or `0.3`. |
| `events.repository` | string | The forge repository, `OWNER/REPO`. |
| `events.window` | object | `from` and `to` (RFC 3339): the inclusive window of merge times the changes were selected from. `complete` (boolean) says whether collection was complete, and `reason` (string or null) says why not. |
| `events.actors` | map | Actor registry keyed by `namespace:name` (`github:alice`, `agent:claude-code`). Each entry has `kind` (`human`, `agent`, `bot`, `service`, `unknown`), `display_name` (string or null) and optional `vendor` (string or null, for agents). |
| `events.attestations` | map, optional | Attestations the collector drew evidence from, keyed `attestation:<16 hex>`. Each has `file`, `predicate_type`, `payload_digest`, `signed`, `verified_by` (who the caller says verified the signature), optional `signer` (the identity the verifier established) and optional `matched`. |
| `events.changes[].id` | string | Change identifier, for example `github:acme/api:pr:421`. |
| `events.changes[].repository`, `.forge`, `.title`, `.url` | string | Change identity. `forge` is `github` in this version. |
| `events.changes[].opened_at` | string | RFC 3339, UTC. |
| `events.changes[].merged_at` | string or null | RFC 3339, UTC, when merged. |
| `events.changes[].head_sha` | string | The current head commit, which the approvals are bound to. The first subject. |
| `events.changes[].merge_commit_sha` | string or null, optional | The commit the merge produced, when merged and reported. The second subject. |
| `events.changes[].commits[]` | array | `sha`, and `author` as an Identity or null. |
| `events.changes[].forge_author` | Identity | The account that opened the change. |
| `events.changes[].author` | Identity | The effective author: a human, or an agent when agent authorship is established (specification section 3). |
| `events.changes[].agent_operator` | object or null | For an agent author: `actor_id` (a human, or null), `confidence` (`explicit`, `derived` or `unknown`) and `provenance`. |
| `events.changes[].reviews[]` | array | The complete review history, each entry with `id`, `actor_id`, `state` (`approved`, `changes_requested`, `commented`, `dismissed`), `at`, `commit_sha` and `provenance`. `agent` is optional, or null: for an agent reviewer with a matched review document it holds `operator`, `identity`, `instructions_owner` and `model`. |
| `events.changes[].reviews_complete` | boolean | Whether the review history was fully retrieved. |
| `events.changes[].merger` | Identity or null | Who merged, when merged. |
| `events.changes[].provenance` | Evidence | Where the change record came from. |
| `events.changes[].labels` | array of string | Forge labels, sorted. |
| `policy.version` | string | Policy format version. |
| `policy.fail_on` | string | Severity threshold (`info` < `warning` < `medium` < `high`). |
| `policy.minimum_authorship_evidence`, `policy.minimum_review_evidence` | string | The weakest evidence kind an operator claim or an approval may rest on. |
| `policy.agent_review` | object | Opt-in mode in which a signed, independent agent approval may satisfy independence: `satisfies_independence`, `require` (`operator`, `provider`, `identity`, `instructions`), `minimum_evidence` (always `signed`) and `labels`. |
| `policy.rules` | map | Per rule name: `enabled` (boolean) and `severity`. |
| `findings[]` | array | One per failed rule per change. Each has `id` (`acc-<16 hex>`), `rule_id`, `rule`, `severity`, `change_id`, `actor_ids`, `explanation`, `provenance` and `disposition`. `legacy_ids` is optional: the identifier the same finding had under 0.2. |
| `findings[].disposition` | object | A human decision: `status` (`open`, `accepted`, `remediated`, `false_positive`), `owner`, `decided_at`, `expires_at`, `rationale` and `remediated_at`. Editing it does not change the finding. |
| `assessments[]` | array | Per rule per change: `change_id`, `rule_id`, `status` (`pass`, `fail`, `unknown`, `not_applicable`) and `reason`. |
| `summary` | object | Integer counts: `changes`, `human_authored`, `agent_authored`, `agent_operator_unknown`, `clean`, `with_findings`, `indeterminate` and `findings`. |
| `generated` | object | `tool`, `version` and `source_digest` (`sha256:<hex>`). |

The rules are defined in section 6.2 of the specification:

| Rule | Fails when |
| --- | --- |
| ACC001 | An agent authored the change, its effective human is known, reviews are complete, and no qualifying independent approval exists. |
| ACC002 | An `approved` review, at any point, is by the effective human. |
| ACC003 | The change is merged, its effective human is known, reviews are complete, and no qualifying independent approval exists. |
| ACC006 | An agent authored the change and no human operator is recorded (advisory). |
| ACC007 | An agent approved the current head (an observation, `info`). |
| ACC008 | Every approval of the head is by an agent of the author agent's vendor. |
| ACC009 | Under `agent_review`, every agent approval is dependent on a required dimension. |
| ACC010 | Under `agent_review`, an agent approval lacks signed evidence (`warning`). |

ACC004 and ACC005 are reserved.

## Conformance

Evaluator conformance is defined by a published corpus, not by agreement with
the reference implementation (specification section 9):

-   **Suite:** `acp-evaluator-conformance`, `suiteRevision` 2, for
    specification version 0.3.
-   **Pinned by:** repository `noru-tech/agent-change-control`, directory
    `conformance/`, release tag `v0.5.2`, commit
    `1fad4bb75fd9ebb60bd5196c89a203eb4de0c5d0`.
-   **Digest list:** `CORPUS-DIGESTS.txt`, SHA-256
    `ffa39c97c2a7247bc1873749c73c620884df1e17c52bb74b0604d70b5c730c52`.
    It is signed at the release tag with a GitHub artifact attestation;
    verify with `gh attestation verify conformance/CORPUS-DIGESTS.txt
    --repo noru-tech/agent-change-control --signer-workflow
    noru-tech/agent-change-control/.github/workflows/conformance-release.yml
    --source-ref refs/tags/v0.5.2`.
-   **Contents:** 35 accept, 18 reject and 4 incomplete vectors. Every reject
    vector is one mutation of an accept vector.
-   **Contract:** the evaluator runs as `<cmd> <vector-file>`. Its exit status
    is `0` evaluated, `3` invalid or `4` incomplete. The last line of stdout is
    `{"verdict", "codes", "assessments", "manifestDigest"}`. `codes` are
    compared as a set, `assessments` exactly, and `manifestDigest`, when
    present, must match. That last check tests canonicalization as well as rule
    semantics.
-   **Harness:** `python3 conformance/run.py --verifier "<cmd>"`, standard
    library only. It is also available as a GitHub Action.

The corpus is a public answer key. A pass shows that an implementation produced
the expected outputs for these inputs. It does not prove that the
implementation implements the rules. Independent implementations, and runs that
disagree with the corpus, are the most useful evidence.

## Security Considerations

-   **Declarations are claims.** An inline declaration of agent authorship is a
    claim by whoever wrote the change description, which is mutable and can be
    forged or omitted. ACP makes the claim explicit and auditable; it does not
    authenticate it. `signed` evidence is only as good as the signature check
    behind it. The collector records who verified each attestation
    (`verified_by`), but it does not verify signatures itself.
-   **Operators may be unknown.** When no record names the human who directed
    an agent, the operator is `unknown`. ACC006 reports it, and independence is
    `unknown`, never `pass`. The merger, the opener or a reviewer is never
    substituted.
-   **Forge data is mutable.** Descriptions are edited, accounts renamed or
    deleted, and reviews dismissed. The predicate is a snapshot at collection
    time. Anything that could not be fully retrieved is marked incomplete, and
    an incomplete record can never be clean.
-   **Scope.** A clean result is a statement about the recorded scope (one
    repository and one window) and the trusted inputs. It is not a compliance
    certification.
-   **Where trust comes from.** Re-evaluation detects tampering with the
    derived parts (findings, assessments, summary). It cannot detect false
    facts. Trust in the facts comes from the signer, the evaluator identity a
    consumer chose to trust, and verifiers MUST pin signer identities
    themselves. A predicate never vouches for its own signer.
-   **Agent reviewers.** By default an agent's approval never satisfies
    independence. The `agent_review` mode is opt-in, requires `signed`
    evidence, and reports `unknown` when a required dimension cannot be
    established.

## Privacy Considerations

The predicate records employee activity: account identifiers, display names,
change titles, the full review history with timestamps, and who directed which
agent. It is personal data in most jurisdictions.

-   Before publishing an attestation to a public transparency log, for example
    by keyless signing through the public Sigstore instance, consider that the
    entry is permanent and cannot be withdrawn.
-   An attestation store is readable by everyone who can read the store. For a
    public repository the facts are already visible on the forge. For a private
    one they are not, so sign with a key you control and keep attestations in a
    store with the same access as the repository.
-   The per-change form carries only the actors its change refers to.
    Pseudonymization for sharing outside the organization is planned but not
    part of this version.

## Example

A merged pull request written by Claude Code under the direction of
`github:alice`, who then approved it. The platform sees two accounts,
the agent's and Alice's. ACP sees one human. The Statement has two subjects,
the head and the merge commit. It has three findings: ACC001 (no independent
human approved the agent's change), ACC002 (the operator approved it) and
ACC003 (merged without independent approval). ACC006 passes, because the
operator is recorded, and ACC007 to ACC010 do not apply.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "digest": {
        "gitCommit": "3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e"
      },
      "name": "github:acme/api:pr:421"
    },
    {
      "digest": {
        "gitCommit": "8b1e0c7d6f5a4e3b2c1d0e9f8a7b6c5d4e3f2a1b"
      },
      "name": "github:acme/api:pr:421:merge"
    }
  ],
  "predicateType": "https://noru.tech/spec/ai-change-provenance/v0.3",
  "predicate": {
    "assessments": [
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "No qualifying independent human approval of the current head exists before merge.",
        "rule_id": "ACC001",
        "status": "fail"
      },
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "The effective human author or agent operator approved this change.",
        "rule_id": "ACC002",
        "status": "fail"
      },
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "No qualifying independent human approval of the current head exists before merge.",
        "rule_id": "ACC003",
        "status": "fail"
      },
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "A human agent operator is recorded.",
        "rule_id": "ACC006",
        "status": "pass"
      },
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "No agent reviewed the current head.",
        "rule_id": "ACC007",
        "status": "not_applicable"
      },
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "A human approved the current head.",
        "rule_id": "ACC008",
        "status": "not_applicable"
      },
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "Agent approvals do not satisfy independence under this policy.",
        "rule_id": "ACC009",
        "status": "not_applicable"
      },
      {
        "change_id": "github:acme/api:pr:421",
        "reason": "Agent approvals do not satisfy independence under this policy.",
        "rule_id": "ACC010",
        "status": "not_applicable"
      }
    ],
    "events": {
      "actors": {
        "agent:claude-code": {
          "display_name": "Claude Code",
          "kind": "agent",
          "vendor": "anthropic"
        },
        "agent:codex": {
          "display_name": "Codex",
          "kind": "agent",
          "vendor": "openai"
        },
        "github:alice": {
          "display_name": "Alice",
          "kind": "human",
          "vendor": null
        },
        "github:bob": {
          "display_name": "Bob",
          "kind": "human",
          "vendor": null
        },
        "github:ci": {
          "display_name": "CI",
          "kind": "bot",
          "vendor": null
        }
      },
      "attestations": {},
      "changes": [
        {
          "agent_operator": {
            "actor_id": "github:alice",
            "confidence": "explicit",
            "provenance": [
              {
                "kind": "declared",
                "ref": "https://github.com/acme/api/pull/421",
                "source": "pr_metadata"
              }
            ]
          },
          "author": {
            "actor_id": "agent:claude-code",
            "provenance": [
              {
                "kind": "observed",
                "ref": "https://api.github.com/repos/acme/api/pulls/421",
                "source": "github_api"
              }
            ]
          },
          "commits": [
            {
              "author": {
                "actor_id": "github:alice",
                "provenance": [
                  {
                    "kind": "observed",
                    "ref": "https://api.github.com/repos/acme/api/pulls/421",
                    "source": "github_api"
                  }
                ]
              },
              "sha": "3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e"
            }
          ],
          "forge": "github",
          "forge_author": {
            "actor_id": "github:alice",
            "provenance": [
              {
                "kind": "observed",
                "ref": "https://api.github.com/repos/acme/api/pulls/421",
                "source": "github_api"
              }
            ]
          },
          "head_sha": "3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e",
          "id": "github:acme/api:pr:421",
          "labels": [],
          "merge_commit_sha": "8b1e0c7d6f5a4e3b2c1d0e9f8a7b6c5d4e3f2a1b",
          "merged_at": "2026-08-03T12:00:00Z",
          "merger": {
            "actor_id": "github:alice",
            "provenance": [
              {
                "kind": "observed",
                "ref": "https://api.github.com/repos/acme/api/pulls/421",
                "source": "github_api"
              }
            ]
          },
          "opened_at": "2026-08-01T12:00:00Z",
          "provenance": [
            {
              "kind": "observed",
              "ref": "https://api.github.com/repos/acme/api/pulls/421",
              "source": "github_api"
            }
          ],
          "repository": "acme/api",
          "reviews": [
            {
              "actor_id": "github:alice",
              "agent": null,
              "at": "2026-08-02T12:00:00Z",
              "commit_sha": "3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e",
              "id": "1",
              "provenance": [
                {
                  "kind": "observed",
                  "ref": "https://api.github.com/repos/acme/api/pulls/421/reviews/1",
                  "source": "github_api"
                }
              ],
              "state": "approved"
            }
          ],
          "reviews_complete": true,
          "title": "Example change",
          "url": "https://github.com/acme/api/pull/421"
        }
      ],
      "repository": "acme/api",
      "version": "0.2",
      "window": {
        "complete": true,
        "from": "2026-08-01T00:00:00Z",
        "reason": null,
        "to": "2026-08-31T23:59:59Z"
      }
    },
    "findings": [
      {
        "actor_ids": [
          "agent:claude-code",
          "github:alice"
        ],
        "change_id": "github:acme/api:pr:421",
        "disposition": {
          "decided_at": null,
          "expires_at": null,
          "owner": null,
          "rationale": null,
          "remediated_at": null,
          "status": "open"
        },
        "explanation": "No qualifying independent human approval of the current head exists before merge.",
        "id": "acc-76345141bc169388",
        "legacy_ids": [
          "acc-e59f983b2bb043de"
        ],
        "provenance": [
          {
            "kind": "observed",
            "ref": "https://api.github.com/repos/acme/api/pulls/421",
            "source": "github_api"
          },
          {
            "kind": "observed",
            "ref": "https://api.github.com/repos/acme/api/pulls/421/reviews/1",
            "source": "github_api"
          },
          {
            "kind": "declared",
            "ref": "https://github.com/acme/api/pull/421",
            "source": "pr_metadata"
          }
        ],
        "rule": "agent_change_without_independent_human",
        "rule_id": "ACC001",
        "severity": "high"
      },
      {
        "actor_ids": [
          "agent:claude-code",
          "github:alice"
        ],
        "change_id": "github:acme/api:pr:421",
        "disposition": {
          "decided_at": null,
          "expires_at": null,
          "owner": null,
          "rationale": null,
          "remediated_at": null,
          "status": "open"
        },
        "explanation": "The effective human author or agent operator approved this change.",
        "id": "acc-1e2fe2ca19119412",
        "legacy_ids": [
          "acc-25b6499b1ba489be"
        ],
        "provenance": [
          {
            "kind": "observed",
            "ref": "https://api.github.com/repos/acme/api/pulls/421",
            "source": "github_api"
          },
          {
            "kind": "observed",
            "ref": "https://api.github.com/repos/acme/api/pulls/421/reviews/1",
            "source": "github_api"
          },
          {
            "kind": "declared",
            "ref": "https://github.com/acme/api/pull/421",
            "source": "pr_metadata"
          }
        ],
        "rule": "approver_is_author",
        "rule_id": "ACC002",
        "severity": "high"
      },
      {
        "actor_ids": [
          "agent:claude-code",
          "github:alice"
        ],
        "change_id": "github:acme/api:pr:421",
        "disposition": {
          "decided_at": null,
          "expires_at": null,
          "owner": null,
          "rationale": null,
          "remediated_at": null,
          "status": "open"
        },
        "explanation": "No qualifying independent human approval of the current head exists before merge.",
        "id": "acc-6e1239e4e65c437c",
        "legacy_ids": [
          "acc-843391ac0a886679"
        ],
        "provenance": [
          {
            "kind": "observed",
            "ref": "https://api.github.com/repos/acme/api/pulls/421",
            "source": "github_api"
          },
          {
            "kind": "observed",
            "ref": "https://api.github.com/repos/acme/api/pulls/421/reviews/1",
            "source": "github_api"
          },
          {
            "kind": "declared",
            "ref": "https://github.com/acme/api/pull/421",
            "source": "pr_metadata"
          }
        ],
        "rule": "merged_without_independent_approval",
        "rule_id": "ACC003",
        "severity": "high"
      }
    ],
    "generated": {
      "source_digest": "sha256:5beced00260c88dad83d6cbc00de4c2c828820e0c5ccf0bed0fa6ee04df70cab",
      "tool": "agent-change-control",
      "version": "0.5.2"
    },
    "policy": {
      "agent_review": {
        "labels": [],
        "minimum_evidence": "signed",
        "require": [
          "operator",
          "provider",
          "identity"
        ],
        "satisfies_independence": false
      },
      "fail_on": "medium",
      "minimum_authorship_evidence": "derived",
      "minimum_review_evidence": "observed",
      "rules": {
        "agent_approval_not_independent": {
          "enabled": true,
          "severity": "high"
        },
        "agent_approval_recorded": {
          "enabled": true,
          "severity": "info"
        },
        "agent_approval_unsigned": {
          "enabled": true,
          "severity": "warning"
        },
        "agent_change_without_independent_human": {
          "enabled": true,
          "severity": "high"
        },
        "approver_is_author": {
          "enabled": true,
          "severity": "high"
        },
        "merged_without_independent_approval": {
          "enabled": true,
          "severity": "high"
        },
        "same_vendor_write_and_review": {
          "enabled": true,
          "severity": "high"
        },
        "unknown_agent_operator": {
          "enabled": true,
          "severity": "warning"
        }
      },
      "version": "0.3"
    },
    "summary": {
      "agent_authored": 1,
      "agent_operator_unknown": 0,
      "changes": 1,
      "clean": 0,
      "findings": 3,
      "human_authored": 0,
      "indeterminate": 0,
      "with_findings": 1
    },
    "version": "0.3"
  }
}
```

### Digests in this example

Each preimage below is shown exactly: the bytes inside the code block, without the newline that ends the block. To recompute one, save the block's content without that final newline and run `sha256sum`; or serialize the corresponding part of the Statement above with any RFC 8785 implementation.

**Finding `acc-76345141bc169388`** (ACC001): `acc-` and the first 16 hexadecimal characters of the SHA-256 of the JCS array `[repository, change_id, rule_id, actor_ids]`.

```json
["acme/api","github:acme/api:pr:421","ACC001",["agent:claude-code","github:alice"]]
```

SHA-256: `76345141bc169388fb37455fa7aa6118acb598d7e418add2f9359871365055ef`

**Finding `acc-1e2fe2ca19119412`** (ACC002): `acc-` and the first 16 hexadecimal characters of the SHA-256 of the JCS array `[repository, change_id, rule_id, actor_ids]`.

```json
["acme/api","github:acme/api:pr:421","ACC002",["agent:claude-code","github:alice"]]
```

SHA-256: `1e2fe2ca191194124d6444015e27e5eb84ab651519fa8d7919e8ef3cb673c92d`

**Finding `acc-6e1239e4e65c437c`** (ACC003): `acc-` and the first 16 hexadecimal characters of the SHA-256 of the JCS array `[repository, change_id, rule_id, actor_ids]`.

```json
["acme/api","github:acme/api:pr:421","ACC003",["agent:claude-code","github:alice"]]
```

SHA-256: `6e1239e4e65c437ca293e2bf174dd5c69e5d4220fb10b4ba1af353d30d364787`

**`generated.source_digest`**: the SHA-256 of the JCS serialization of `predicate.events`.

<details><summary>Preimage</summary>

```json
{"actors":{"agent:claude-code":{"display_name":"Claude Code","kind":"agent","vendor":"anthropic"},"agent:codex":{"display_name":"Codex","kind":"agent","vendor":"openai"},"github:alice":{"display_name":"Alice","kind":"human","vendor":null},"github:bob":{"display_name":"Bob","kind":"human","vendor":null},"github:ci":{"display_name":"CI","kind":"bot","vendor":null}},"attestations":{},"changes":[{"agent_operator":{"actor_id":"github:alice","confidence":"explicit","provenance":[{"kind":"declared","ref":"https://github.com/acme/api/pull/421","source":"pr_metadata"}]},"author":{"actor_id":"agent:claude-code","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"commits":[{"author":{"actor_id":"github:alice","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"sha":"3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e"}],"forge":"github","forge_author":{"actor_id":"github:alice","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"head_sha":"3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e","id":"github:acme/api:pr:421","labels":[],"merge_commit_sha":"8b1e0c7d6f5a4e3b2c1d0e9f8a7b6c5d4e3f2a1b","merged_at":"2026-08-03T12:00:00Z","merger":{"actor_id":"github:alice","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"opened_at":"2026-08-01T12:00:00Z","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}],"repository":"acme/api","reviews":[{"actor_id":"github:alice","agent":null,"at":"2026-08-02T12:00:00Z","commit_sha":"3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e","id":"1","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421/reviews/1","source":"github_api"}],"state":"approved"}],"reviews_complete":true,"title":"Example change","url":"https://github.com/acme/api/pull/421"}],"repository":"acme/api","version":"0.2","window":{"complete":true,"from":"2026-08-01T00:00:00Z","reason":null,"to":"2026-08-31T23:59:59Z"}}
```

</details>

SHA-256: `5beced00260c88dad83d6cbc00de4c2c828820e0c5ccf0bed0fa6ee04df70cab`

**The `sha256` subject form**: a signer that only accepts SHA-2 subjects uses a single subject whose `sha256` is the digest of the JCS serialization of `predicate`, which is also the exact content of the manifest file `acc --format json` writes.

<details><summary>Preimage</summary>

```json
{"assessments":[{"change_id":"github:acme/api:pr:421","reason":"No qualifying independent human approval of the current head exists before merge.","rule_id":"ACC001","status":"fail"},{"change_id":"github:acme/api:pr:421","reason":"The effective human author or agent operator approved this change.","rule_id":"ACC002","status":"fail"},{"change_id":"github:acme/api:pr:421","reason":"No qualifying independent human approval of the current head exists before merge.","rule_id":"ACC003","status":"fail"},{"change_id":"github:acme/api:pr:421","reason":"A human agent operator is recorded.","rule_id":"ACC006","status":"pass"},{"change_id":"github:acme/api:pr:421","reason":"No agent reviewed the current head.","rule_id":"ACC007","status":"not_applicable"},{"change_id":"github:acme/api:pr:421","reason":"A human approved the current head.","rule_id":"ACC008","status":"not_applicable"},{"change_id":"github:acme/api:pr:421","reason":"Agent approvals do not satisfy independence under this policy.","rule_id":"ACC009","status":"not_applicable"},{"change_id":"github:acme/api:pr:421","reason":"Agent approvals do not satisfy independence under this policy.","rule_id":"ACC010","status":"not_applicable"}],"events":{"actors":{"agent:claude-code":{"display_name":"Claude Code","kind":"agent","vendor":"anthropic"},"agent:codex":{"display_name":"Codex","kind":"agent","vendor":"openai"},"github:alice":{"display_name":"Alice","kind":"human","vendor":null},"github:bob":{"display_name":"Bob","kind":"human","vendor":null},"github:ci":{"display_name":"CI","kind":"bot","vendor":null}},"attestations":{},"changes":[{"agent_operator":{"actor_id":"github:alice","confidence":"explicit","provenance":[{"kind":"declared","ref":"https://github.com/acme/api/pull/421","source":"pr_metadata"}]},"author":{"actor_id":"agent:claude-code","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"commits":[{"author":{"actor_id":"github:alice","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"sha":"3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e"}],"forge":"github","forge_author":{"actor_id":"github:alice","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"head_sha":"3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e","id":"github:acme/api:pr:421","labels":[],"merge_commit_sha":"8b1e0c7d6f5a4e3b2c1d0e9f8a7b6c5d4e3f2a1b","merged_at":"2026-08-03T12:00:00Z","merger":{"actor_id":"github:alice","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}]},"opened_at":"2026-08-01T12:00:00Z","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"}],"repository":"acme/api","reviews":[{"actor_id":"github:alice","agent":null,"at":"2026-08-02T12:00:00Z","commit_sha":"3f9c2e4b7a1d8c6e5f0a2b4c6d8e0f1a3b5c7d9e","id":"1","provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421/reviews/1","source":"github_api"}],"state":"approved"}],"reviews_complete":true,"title":"Example change","url":"https://github.com/acme/api/pull/421"}],"repository":"acme/api","version":"0.2","window":{"complete":true,"from":"2026-08-01T00:00:00Z","reason":null,"to":"2026-08-31T23:59:59Z"}},"findings":[{"actor_ids":["agent:claude-code","github:alice"],"change_id":"github:acme/api:pr:421","disposition":{"decided_at":null,"expires_at":null,"owner":null,"rationale":null,"remediated_at":null,"status":"open"},"explanation":"No qualifying independent human approval of the current head exists before merge.","id":"acc-76345141bc169388","legacy_ids":["acc-e59f983b2bb043de"],"provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"},{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421/reviews/1","source":"github_api"},{"kind":"declared","ref":"https://github.com/acme/api/pull/421","source":"pr_metadata"}],"rule":"agent_change_without_independent_human","rule_id":"ACC001","severity":"high"},{"actor_ids":["agent:claude-code","github:alice"],"change_id":"github:acme/api:pr:421","disposition":{"decided_at":null,"expires_at":null,"owner":null,"rationale":null,"remediated_at":null,"status":"open"},"explanation":"The effective human author or agent operator approved this change.","id":"acc-1e2fe2ca19119412","legacy_ids":["acc-25b6499b1ba489be"],"provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"},{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421/reviews/1","source":"github_api"},{"kind":"declared","ref":"https://github.com/acme/api/pull/421","source":"pr_metadata"}],"rule":"approver_is_author","rule_id":"ACC002","severity":"high"},{"actor_ids":["agent:claude-code","github:alice"],"change_id":"github:acme/api:pr:421","disposition":{"decided_at":null,"expires_at":null,"owner":null,"rationale":null,"remediated_at":null,"status":"open"},"explanation":"No qualifying independent human approval of the current head exists before merge.","id":"acc-6e1239e4e65c437c","legacy_ids":["acc-843391ac0a886679"],"provenance":[{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421","source":"github_api"},{"kind":"observed","ref":"https://api.github.com/repos/acme/api/pulls/421/reviews/1","source":"github_api"},{"kind":"declared","ref":"https://github.com/acme/api/pull/421","source":"pr_metadata"}],"rule":"merged_without_independent_approval","rule_id":"ACC003","severity":"high"}],"generated":{"source_digest":"sha256:5beced00260c88dad83d6cbc00de4c2c828820e0c5ccf0bed0fa6ee04df70cab","tool":"agent-change-control","version":"0.5.2"},"policy":{"agent_review":{"labels":[],"minimum_evidence":"signed","require":["operator","provider","identity"],"satisfies_independence":false},"fail_on":"medium","minimum_authorship_evidence":"derived","minimum_review_evidence":"observed","rules":{"agent_approval_not_independent":{"enabled":true,"severity":"high"},"agent_approval_recorded":{"enabled":true,"severity":"info"},"agent_approval_unsigned":{"enabled":true,"severity":"warning"},"agent_change_without_independent_human":{"enabled":true,"severity":"high"},"approver_is_author":{"enabled":true,"severity":"high"},"merged_without_independent_approval":{"enabled":true,"severity":"high"},"same_vendor_write_and_review":{"enabled":true,"severity":"high"},"unknown_agent_operator":{"enabled":true,"severity":"warning"}},"version":"0.3"},"summary":{"agent_authored":1,"agent_operator_unknown":0,"changes":1,"clean":0,"findings":3,"human_authored":0,"indeterminate":0,"with_findings":1},"version":"0.3"}
```

</details>

SHA-256: `0be85acde6dd0805c7fd360e9dab2b469a1d1ae36cd95b3f8048969e8f936e31`

**The Statement itself**: the bytes a DSSE signer wraps are the JCS serialization of the Statement above, 7078 bytes with SHA-256 `734c3d9a4aa63e73c1e386ff00aa8e91c60e07077b609b6d791429cc5351d62b`.

## Related Predicate Types

The specification defines two more predicate types. They will be proposed
separately, once the naming conventions for this one are settled:

-   `https://noru.tech/spec/ai-change-provenance/provenance/v0.1`: an agent's
    integration or operator states which agent wrote a change and which human
    directed it, with the head commit as subject (specification section 3.2).
-   `https://noru.tech/spec/ai-change-provenance/review/v0.1`: a reviewer's
    tooling states who reviewed which head and what they decided, and for an
    agent reviewer who operated it and who owns its instructions (specification
    section 3.7).

### Composition with AI code-generation provenance

This section is non-normative.

The OpenFab Generation predicate,
`https://open-fab.ai/attestation/generation/v0.1`, proposed in
[in-toto/attestation#604](https://github.com/in-toto/attestation/issues/604),
records which line ranges inside files a model generated and which a human
wrote. ACP stops at the change boundary: it says who owned a change and whether
an independent human approved it, not which lines inside it came from a model.
When both statements exist for the same change, a verifier can evaluate them
together, from "who owned this change" down to "which lines came from which
model".

-   **No shared subject.** ACP subjects are commits (`digest.gitCommit`).
    Generation subjects and `generated[]` entries are SHA-256 digests of file
    content. The two statements never share a subject digest, so a verifier
    cannot join them on the subject alone.
-   **Joining through the commit tree.** Given an ACP change, a verifier reads
    each `generated[].path` from the tree of the change's head commit and
    recomputes the Generation statement's digests from those bytes, as that
    specification defines them. If every digest matches, the Generation
    statement describes content of the change ACP attests. Git object
    identifiers hash a header together with the content and are not the same
    value, so the verifier hashes the file bytes, not the blob identifier.
-   **Which commit.** Join against `head_sha`, which the approvals are bound
    to. The tree of a merge commit can include other changes. A file edited
    after generation no longer matches, and that is the right result: the
    Generation statement no longer describes what was approved.
-   **Finding the statement.** Neither predicate points to the other. ACP 0.3
    records in `events.attestations` only the authorship and review attestations
    its collector drew evidence from, and a Generation statement is neither. A
    verifier obtains both statements from its own attestation store, and the
    join above establishes that they describe the same content.
-   **Approvals are counted once.** Generation `signoffs[]` are human
    cryptographic approvals of the generated artifact. ACP does not count them
    as reviews, and a Generation verifier does not count ACP reviews as
    sign-offs. A policy that wants both reads both statements.
-   **Claimed and observed.** The Generation predicate separates a producer's
    self-reported acceptance result from a verifier's re-execution of the
    contract. ACP makes the same split between its embedded findings, which a
    producer claims, and the findings a consumer obtains by re-evaluating the
    embedded events under the embedded policy.

## Changelog and Migrations

-   **0.3.**
    -   Every digest is SHA-256 over RFC 8785 bytes, and input must be I-JSON
        with the constraints above.
    -   Version 0.2 used a project canonicalization whose preimages ended in a
        newline. Every finding identifier, `source_digest` and `sha256` subject
        digest therefore changed, while findings, assessments and verdicts did
        not. Findings carry their 0.2 identifier in `legacy_ids` for one minor
        version.
    -   A `v0.2` Statement is still validated, by recomputing its digests with
        the legacy canonicalization.
-   **0.2.** Agent reviewers: rules ACC007 to ACC010 and the opt-in
    `agent_review` policy.
-   **0.1.** The initial version: head-commit subjects, then merge-commit
    subjects and the JSON Lines form, and signed authorship evidence.
