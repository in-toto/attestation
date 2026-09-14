# Predicate type: Source Review Coverage

Type URI: `https://drvelvetfog.github.io/source-review-coverage/v0.1`

Version: v0.1

Predicate name: `source-review-coverage`

## Purpose

To carry the evidence needed to decide, independently of the system that produced it, whether a source revision was covered by the review it claims.

SLSA v1.2 Source Track Level 4 requires that changes to protected branches "be agreed to by two or more trusted persons prior to submission," and requires this of the *final revision submitted* — if changes occur during review, they need re-review. The requirement is precise. The verification is absent: SLSA "leaves source provenance attestations undefined and up to the SCSs to determine what works best," and does not address how squash merges or rebases interact with approval, nor how post-approval changes are detected.

In practice the Source Control System asserts the property. A consumer receives "two people approved this" from the same platform that performed the merge, and either believes it or does not. This predicate replaces that assertion with recomputable evidence: which tree each approval covered, which tree shipped, and — where those differ because the base moved beneath the change — whether the shipped tree is the reviewed change replayed onto its base, exactly.

## Use Cases

**A reviewed change is not the change that shipped.** A reviewer approves a branch; the base advances; the forge merges. The tree that ships is no longer the tree anyone read. Under this predicate a verifier replays the reviewed change onto the base it landed on and compares: an exact match means the review covers the shipped state, and any difference is emitted as the specific bytes no approval covers. Conflict resolution is the common cause — a hunk is resolved into a third form at merge time, appears in no pull request, and reaches the default branch unreviewed.

**Content enters inside a merge commit.** A merge commit can carry changes present in no parent. Forges do not display these in the pull-request diff, so they are unreviewed as a matter of course. Replaying the parents and comparing against the merge commit's tree surfaces them.

**An approval is claimed but not bound.** Where an approval does not name the tree it covers, this predicate has nothing to check, and a verifier reports that rather than assuming coverage — which is the distinction between evidence and a claim.

Existing predicates do not cover these. SLSA Provenance describes how an artifact was *built*, not how a change was *admitted*. A Source VSA communicates *that* a revision meets a level — a conclusion, not the evidence to reach it independently. Link, SCAI and SBOM predicates describe steps, attributes and contents; none binds an approval to a tree.

## Prerequisites

A Git repository, or a version control system with content-addressed tree objects and a deterministic merge. The `gitCommit`, `gitTree`, `gitBlob` and `gitTag` digest types defined by the attestation framework.

Verification requires the reviewed revision to be reachable. Under a squash workflow it is not an ancestor of the subject; forges that retain pull references (`refs/pull/*/head`) make it recoverable, and where a repository retains none, squash revisions are unverifiable and MUST be reported as such.

## Model

Three functionaries appear, and separating them is the point.

A **reviewer** approves a revision, producing an approval bound to the tree they read. A **merge system** admits the change to a protected branch, possibly transforming it — squashing, rebasing, or merging — and records how. A **verifier**, trusting neither of the other two, recomputes the transform and reports whether the approved tree accounts for the shipped tree.

The predicate is issued at admission, over the revision that shipped. It is expected to be signed by a workflow identity rather than by the party whose work is being attested; an attestation signed only by the audited party is an assertion, and a verifier is expected to treat it accordingly.

## Schema

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "<REF>",
    "digest": { "gitCommit": "<COMMIT>", "gitTree": "<TREE>" }
  }],
  "predicateType": "https://drvelvetfog.github.io/source-review-coverage/v0.1",
  "predicate": {
    "reviewCoverage": {
      "result": "identity | replay | residual | unverifiable",
      "residualBase": { "digest": { "gitTree": "<TREE>" } }
    },
    "approvals": [{
      "overTree": { "digest": { "gitTree": "<TREE>" } },
      "approver": "<STRING>",
      "approvedAt": "<TIMESTAMP>"
    }],
    "checks": [{
      "name": "<STRING>",
      "overTree": { "digest": { "gitTree": "<TREE>" } },
      "outcome": "pass | fail",
      "runnerIdentity": "<STRING>"
    }],
    "mergeTransform": {
      "kind": "identity | squash | mergeCommit",
      "reviewedHead": { "digest": { "gitCommit": "<COMMIT>" } },
      "baseAtMerge": { "digest": { "gitCommit": "<COMMIT>" } },
      "parents": [{ "digest": { "gitCommit": "<COMMIT>" } }],
      "expectedTree": { "digest": { "gitTree": "<TREE>" } },
      "replayClean": <BOOL>,
      "strategy": "<STRING>",
      "gitVersion": "<STRING>"
    },
    "authorship": {
      "declaredBy": "<STRING>",
      "agents": [{ "tool": "<STRING>", "model": "<STRING>", "sessionRef": "<STRING>" }],
      "humanOperator": "<STRING>"
    },
    "objectFormat": "sha1 | sha256"
  }
}
```

### Parsing Rules

This predicate follows the standard in-toto attestation parsing rules, with these additions:

-   **`result` is a string, not an enumeration.** A verifier encountering an unrecognised value MUST report it as unrecognised rather than coerce it to a known one.
-   **Digests are wrapped.** Each digest-bearing field contains a `digest` object, as elsewhere in the framework. `parents` is repeated and so cannot be a bare map.
-   **Unset is not false.** An absent `mergeTransform` means no transform was recorded, not that the subject is the reviewed tree. A verifier MUST NOT infer `identity` from absence.
-   **Replay inputs are object IDs.** Fields naming revisions carry digests, never references. Where a merge conflicts, the conflict markers written by the merge embed the names the merge was given, so a replay computed from reference names produces a different tree than one computed from the commits they point at.
-   **Recorded values are checkable, not authoritative.** `expectedTree` is what the issuer computed. A verifier recomputes and, on divergence, reports drift — including `strategy` and `gitVersion`, since merge results may vary across both.
-   **`approvedAt` bounds nothing.** It is asserted by the signer. Absent a transparency-log inclusion proof, a verifier MUST NOT treat it as evidence of when.

### Fields

**`reviewCoverage`** *object, required*

The outcome of measuring the review against what shipped.

**`reviewCoverage.result`** *string, required*

One of: `identity`, the approved tree is the shipped tree; `replay`, the shipped tree is the reviewed change replayed onto its base, exactly; `residual`, content shipped that no approval covers; `unverifiable`, the reviewed state could not be recovered.

**`reviewCoverage.residualBase`** *object, optional*

The replay result the residual was measured against. Present if and only if `result` is `residual`. The difference between this tree and the subject's tree is the content no approval covers.

**`approvals`** *array of objects, required*

**`approvals[*].overTree`** *object, required*

The tree the approver signed over. An approval that does not name its tree cannot be checked and MUST be rejected rather than assumed to cover the subject.

**`approvals[*].approver`** *string, required* — identity of the approver.

**`approvals[*].approvedAt`** *timestamp, optional* — RFC 3339, timezone `Z`. Asserted; see Parsing Rules.

**`checks`** *array of objects, optional*

**`checks[*].name`** *string, required* — the check's identifier.

**`checks[*].overTree`** *object, required* — the tree the check ran over. A check over a tree no approval covers does not speak for the subject.

**`checks[*].outcome`** *string, required* — `pass` or `fail`.

**`checks[*].runnerIdentity`** *string, optional* — identity of the executing workflow, an OIDC subject where one exists.

**`mergeTransform`** *object, optional*

How the reviewed state became the shipped state. Required where they differ.

**`mergeTransform.kind`** *string, required* — `identity`, `squash`, or `mergeCommit`.

**`mergeTransform.reviewedHead`** *object, required* — the revision the approval covered.

**`mergeTransform.baseAtMerge`** *object, required* — the revision the change landed on.

**`mergeTransform.parents`** *array of objects, optional* — for `mergeCommit` only, in the order the commit records them. Order is normative: a conflicted replay is not symmetric in parent order.

**`mergeTransform.expectedTree`** *object, optional* — the replay result computed at issuance.

**`mergeTransform.replayClean`** *boolean, optional* — whether the replay merged without conflict. A conflicted replay is not itself a failure; it means any residual is resolution work.

**`mergeTransform.strategy`**, **`mergeTransform.gitVersion`** *string, optional* — the merge strategy and version the replay was computed under.

**`authorship`** *object, optional*

Declared, never recomputable. Carried so that policy can act on it, and labelled so that policy cannot mistake it for evidence. A declaration that no agent was involved is exactly as verifiable as one naming a specific model.

**`authorship.declaredBy`** *string, required* — who makes the declaration.

**`authorship.agents`** *array of objects, optional* — `tool`, `model`, `sessionRef`.

**`authorship.humanOperator`** *string, optional* — credential for the human operating the agent, where a deployment binds one.

**`objectFormat`** *string, optional* — `sha1` or `sha256`. Git repositories are `sha1` unless created otherwise; a verifier SHOULD report that tree hashes under `sha1` are not collision-resistant against a motivated attacker.

## Example

A squash merge whose base advanced during review. The shipped tree is not the reviewed tree, and the replay accounts for it exactly.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "refs/heads/main",
      "digest": {
        "gitCommit": "91af93d011e678bb2b884ca26cbff7505d192fdd",
        "gitTree": "a3e3f68f368544c10121912dd580bc55ec721bd0"
      }
    }
  ],
  "predicateType": "https://drvelvetfog.github.io/source-review-coverage/v0.1",
  "predicate": {
    "reviewCoverage": { "result": "replay" },
    "approvals": [
      {
        "overTree": { "digest": { "gitTree": "7a6a0ffba52d712f36edf65df194a56508c70f51" } },
        "approver": "reviewer@example.com"
      }
    ],
    "checks": [
      {
        "name": "unit-tests",
        "overTree": { "digest": { "gitTree": "7a6a0ffba52d712f36edf65df194a56508c70f51" } },
        "outcome": "pass"
      }
    ],
    "mergeTransform": {
      "kind": "squash",
      "reviewedHead": { "digest": { "gitCommit": "64b5ff0c6c298b9e4ff25015b774a2b2d5ac64b1" } },
      "baseAtMerge": { "digest": { "gitCommit": "377c6ab6e170afd390bbcea5f90e9fac19c5c1d2" } },
      "parents": [],
      "expectedTree": { "digest": { "gitTree": "a3e3f68f368544c10121912dd580bc55ec721bd0" } },
      "replayClean": true,
      "strategy": "ort",
      "gitVersion": "2.50.1"
    },
    "authorship": {
      "declaredBy": "ci",
      "agents": [{ "tool": "claude-code" }]
    },
    "objectFormat": "sha1"
  }
}
```

The reviewed tree `7a6a0ffb…` is not the shipped tree `a3e3f68f…`. Replaying `reviewedHead` onto `baseAtMerge` produces `a3e3f68f…`, which is the subject's tree, so the approval covers what shipped and `result` is `replay`. Had the replay produced anything else, `result` would be `residual` and `residualBase` would name the tree the difference is measured against.

## Changelog and Migrations

Initial version.
