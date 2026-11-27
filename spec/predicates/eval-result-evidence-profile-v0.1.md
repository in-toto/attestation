# Evidence profile for `eval-result` v0.1 — draft v0.1.2

Profile URI: `urn:lgd:in-toto:profile:eval-evidence:v0.1`

The profile identifier is a `urn:`, not an `https:` URI, on purpose: there is no registry for
profiles yet, and an `https:` identifier that returns 503 is a defect a reviewer can see in one
`curl`. When one exists, an `https:` URI will be minted alongside, and this URN stays as an alias.
Applies to: `https://in-toto.io/attestation/eval-result/v0.1` (in-toto/attestation#575).
Note: that type URI resolves only once #575 merges; as of 2026-11-17 it returns HTTP 404.
Checked, not assumed — a reader running `curl` today will see the same 404, and it belongs to
the predicate's own landing, not to this profile.
Status: draft for review in in-toto/attestation#565. **Changes nothing in #575.**

---

## 1. What problem this closes

#575 settles the predicate: who signed which claims, at which threshold, with which verdict. Three
things stay open, and #575 leaves them to an evidence profile by design:

1. **The verdict is signed but not recomputable.** `claims[].passed` is the issuer's verdict; a
   consumer with the predicate alone can authenticate it and nothing more. b7n0de's 16 Nov check of
   Aditi's statement is exactly this gap: the predicate passes, but "was it really 0 of 15?" can only
   be answered from the referenced evidence.
2. **Total count vs evaluated count.** `sampleSize` belongs to the Statement. A 15-task screen and a
   164-task confirmation are two Statements over one subject. Nothing in #575 says how the *evaluated*
   count relates to the *total*, nor what happens to rows that could not be scored.
3. **How the outputs were generated.** `harness` names and can bind the harness artifact, but says
   nothing about the chat template, the decode settings, or the runtime. In the published census the
   same two files scored 5/15 and 6/15 without the model's chat template and 13/15 and 15/15 with it.
   That is a 2× swing carried by a field the predicate does not have.

This profile is the artifact that `evidence[]` points at. It does not add a field to `eval-result`.

## 2. Where it sits

```
Statement (signed, #575)                 Evidence artifact (this profile)
├─ subject: model digest                 ├─ counts {total, evaluated, excluded, scored, cannotDetermine}
├─ predicate.evaluatedAt                 ├─ claims[] {numerator, denominator}     ← recomputes passed
├─ predicate.claims[].passed             ├─ exclusionPolicy {ruleId, ruleDigest, excludedRowsDigest}
├─ predicate.sampleSize                  ├─ generation {chatTemplate, decode, runtime}
├─ predicate.subjectProfile (snake_case) ├─ aggregation {ruleId, ruleDigest}
└─ predicate.evidence[0].digest ─────────┤─ reproducibility {grade, digest, inputs}
                                         └─ drift[] {impl, field, expected, observed}
```

A consumer that only needs authenticity verifies the Statement and stops. A consumer that wants to
recompute reads the evidence artifact, whose digest is pinned in the Statement, and applies §4.

## 3. Evidence artifact schema

```jsonc
{
  "profile": "urn:lgd:in-toto:profile:eval-evidence:v0.1",
  "statement": { "sha256": "<digest of the signed Statement, one direction only>" },
  "counts": {
    "total": 15,             // rows in the pinned dataset (MUST equal its row count)
    "evaluated": 15,         // rows run and scored-or-judged
    "excluded": 0,           // rows committed but outside both counts
    "scored": 15,            // rows that produced a verdict
    "cannotDetermine": 0     // judged but undecidable, kept apart from failed
  },
  "claims": [
    { "metric": "tasks_solved", "comparator": ">=", "threshold": "9",
      "numerator": 0, "denominator": 15, "passed": false, "metricKind": "capability_positive",
      "cannotDeterminePolicy": "<ruleId, REQUIRED iff metricKind=failure_absence and passed=true and cannotDetermine>0>" }
  ],
  "exclusionPolicy": null,   // REQUIRED non-null iff counts.excluded > 0
  "generation": {
    "chatTemplate": { "digest": { "sha256": "<hex>" }, "servedAt": "<RFC3339>", "source": "<uri>" },
    "decode": { "temperature": 0, "seed": 42, "numCtx": 2048, "numPredict": 768 },
    "runtime": { "name": "ollama", "version": "<string>", "digest": { "sha256": "<hex>" } }
  },
  "aggregation": { "ruleId": "pass-rate", "ruleDigest": { "sha256": "<hex>" }, "description": "<how numerator is counted>" },
  "reproducibility": { "grade": "reproduced|third_party|self|none", "algorithm": "sha256-concat-v1",
                       "inputs": ["sha256:<hex>", "sha256:<hex>", "sha256:<hex>", "sha256:<hex>"],
                       "digest": { "sha256": "<hex>" } },
  "drift": [ { "impl": "proofbundle v6.5.0", "field": "subjectProfile",
               "expected": "public_model", "observed": "public-model", "note": "<informational>" } ]
}
```

**Binding direction (normative).** The Statement pins the evidence digest
(`predicate.evidence[0].digest`); the evidence artifact MAY pin the statement
(`statement.sha256`). Because each would otherwise hash the other, the statement digest is defined
as §6 of the ACS family defines it: `sha256` over the JCS serialization of the Statement with the
**profile evidence entry's digest blanked** (`""`) before hashing. That makes the pair constructible
and checkable in both directions without a fixed point, and it is the only definition this profile
recognises. A verifier MUST recompute both sides; an artifact that omits `statement.sha256` is
allowed, an artifact whose `statement.sha256` does not match is rejected.

## 4. Conformance checks (normative)

**PV-1 count arithmetic.** `total = evaluated + excluded` and `evaluated = scored + cannotDetermine`.
`total` MUST equal the row count of the pinned dataset — it is not a free integer. Reject
`count_arithmetic`, `total_unbound`.

**PV-1b sample size binding.** `predicate.sampleSize` MUST equal `counts.evaluated`. Reject
`sample_size_mismatch`.

**PV-1c claim binding.** The Statement's single `predicate.claims[0]` and this artifact's single
`claims[0]` MUST agree on `metric`, `comparator`, `threshold` and `passed`. The `denominator`
is bound to `counts.scored` by PV-2, not to `sampleSize`; when a row cannot be determined,
`denominator < sampleSize` is conformant. Otherwise the artifact computes a verdict the
signed predicate never asserted. Reject `claim_binding`.

**PV-2 denominator binding.** Every `claims[].denominator` MUST equal `counts.scored`, not
`counts.evaluated`, and `numerator <= denominator`. The denominator is the number of rows that
actually produced a verdict: taking `evaluated` makes PV-2 and PV-5 mutually unsatisfiable the
moment any row is undecidable (this draft's own v0.1 text did exactly that). Reject
`denominator_mismatch`, `numerator_out_of_range`.

**PV-2b numerator provenance.** When the per-task artifact is available, `numerator` MUST equal the
count of passing rows for the attested tag inside it. A numerator that cannot be recomputed row by
row is an assertion, not a measurement. Reject `numerator_unbound`.

**PV-3 verdict recomputability.** `claims[].passed` MUST equal the comparator applied to
`numerator` and `threshold` (decimal comparison, no float). Reject `verdict_not_recomputable`.
This is the check #575 cannot perform and the reason this profile exists.

**PV-4 exclusion transparency.** If `counts.excluded > 0`, `exclusionPolicy.ruleDigest`,
`exclusionPolicy.excludedRowsDigest` and a non-empty `description` MUST be present, so that the
excluded rows are committed to rather than merely dropped. A digest field that points at nothing
verifiable is not transparency: the digest MUST resolve to a row list. Reject `exclusion_opaque`.

**PV-5 cannot-determine separation.** `cannotDetermine` MUST NOT appear in `numerator` or
`denominator`. For `metricKind: "failure_absence"` (a pass that depends on not finding failures),
`passed: true` with `cannotDetermine > 0` MUST be rejected unless
`claims[].cannotDeterminePolicy` names a pre-registered rule (a `ruleId` whose digest also appears in
`exclusionPolicy.ruleDigest` or `aggregation.ruleDigest`) that says how undecidable rows are treated.
The field is part of the `claims` schema above — an exception condition that lives only in prose
cannot be checked, and a validator will reject a claim that cites a field the schema does not define.
Reject `cannot_determine_folded` when no policy is named; reject `cannot_determine_policy_unbound`
when a policy is named but the rule it names resolves to neither digest — a `ruleId` alone is an
assertion, not a pre-registration. Vector N11 is the case: undecidable rows present, a policy named,
and no registered rule behind it.

**PV-6 generation disclosure and run binding.** For a generative evaluation,
`generation.chatTemplate.digest`, `generation.decode`, `generation.runtime.version` and
`generation.runtime.digest` MUST be present. `chatTemplate.servedAt` MUST satisfy
`0 <= evaluatedAt - servedAt <= 7 days`: at or before the run, and at most a week before it.
A template digest captured *after* the run cannot explain it — that is a post-hoc reconstruction of
the run, not a record of it. Reject `generation_undisclosed`, `template_not_run_bound`.

**PV-7 reproducibility grade binding.** If `reproducibility.grade != "none"`,
`reproducibility.digest` MUST be recomputable as `sha256-concat-v1`, defined exactly as:
the four `inputs` entries are `sha256:<64 lowercase hex>` strings naming (1) dataset, (2) harness,
(3) aggregation rule, (4) generations; the digest is `sha256` over the concatenation of the four
32-byte values decoded from those hex strings, in that order, with no separators and no re-hashing
of the text form. `third_party` and `reproduced` additionally require an `evidence[]` entry named
`reproduction*` **and** a `reproducibility.reproducer` block carrying `key_id` and `signature`, so
that the higher grade is attributable to someone other than the issuer. Reject
`reproducibility_unbound`, `grade_unbound`.

**PV-8 drift disclosure is informational.** A listed mismatch MUST NOT fail conformance; an
*unlisted* known mismatch is the maintainer's call, not this profile's. No reason code.

**PV-9 indeterminacy is a verdict, not a pass.** Three checks need an artifact the verifier may
not hold: PV-1c needs the pinned dataset, PV-2b needs the per-task rows, and any digest comparison
needs the file itself. When such an artifact is absent, the verifier MUST report the case as
*indeterminate*, naming the missing artifact as `cannot_verify:<evidence name>`. It MUST NOT
report acceptance: a rule that was not evaluated has not been satisfied, and silence in place of a
verdict is the fail-open shape this profile exists to close. A verifier that has every artifact
MUST never return this verdict.

This rule is normative because of what a reader does without it: given PV-1c and PV-2b as written,
two reviewers of this draft concluded that a verifier lacking the dataset should simply accept —
there is no sentence telling it otherwise. The behaviour was in the validator before it was in this
text, which is the wrong order.

## 5. Non-claims

This profile does not make the metric correct, the dataset representative, or the harness fair. It
makes the *arithmetic* checkable and the *generation settings* attributable. A recomputed verdict is a
recomputation of the issuer's counting rule, not an endorsement of it. It also does not make the
issuer honest about time: PV-6 checks that the recorded template predates the run, not that the
recorded time was the real one.

## 6. Worked example (deliberately non-conformant, plus its fixed form)

`examples/qwen2.5-coder-3b-q2K.statement.json` and `...evidence.json` are built from Aditi's
published example in #565, with all four referenced digests **re-downloaded and recomputed** by this
package (see `07-外人实测日志.md`) and the per-task scores re-counted from `census_rescore.jsonl`
(0 of 15 for this tag). The verdict is recomputed by PV-3 rather than trusted; `subjectProfile`
uses `public_model` (snake_case, per #575); and `counts.total` is the real 15-row count of
`smoke_suite_v1.json`, not the 164-row figure quoted in the thread.

**This example fails PV-6, on purpose.** Its chat-template digest was captured 2026-11-16, about
96 days *after* `evaluatedAt` — exactly the post-hoc reconstruction PV-6 exists to catch.
`validate_profile.py` reports it as `template_not_run_bound`; vector P1 asserts that failure, and
vector P2 shows the same numbers with a run-bound `servedAt` verifying clean. Publishing a
non-conformant example as if it were conformant would be the mistake this profile is meant to
prevent, so the example carries its own failure and this section says so.

## 7. Changelog

- v0.1.2 (this revision), after three external model reviews (DeepSeek-V3.2, GLM-5.3 via
  SiliconFlow): denominator changed from `evaluated` to `scored` (PV-2 and PV-5 were jointly
  unsatisfiable once any row was undecidable); added `sampleSize` and Statement↔evidence claim
  binding; bound `total` to the dataset row count; bound `numerator` to the per-task artifact; added `claims[].cannotDeterminePolicy` to the schema
  (PV-5 cited it while §3 never defined it — flagged by the third external model);
  defined `inputs` and `sha256-concat-v1` byte-exactly; made PV-6 one-directional
  (`servedAt <= evaluatedAt`) and added `runtime.digest`; required a signed `reproducer` block
  above grade `self`; stated the blanked-entry digest definition that removes the circular
  dependency; corrected the worked example's `total` from 164 to 15 and declared it non-conformant.
- v0.1: initial draft. Supersedes nothing; requested by b7n0de's 16 Nov comment
  (`#565` comment 6928971370) and by the 25 Oct commitment (comment 6398799756).
