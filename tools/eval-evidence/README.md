# eval-evidence validator

Companion to [`spec/predicates/eval-result-evidence-profile-v0.1.md`](../../spec/predicates/eval-result-evidence-profile-v0.1.md).

```bash
# no third-party dependencies: standard library only
python3 validate_profile.py            # 23 vectors against the committed examples
python3 validate_profile.py --build    # rebuild examples from the cited bundle, then run
```

**The cited bundle is not committed here.** The four files are third-party evidence from
in-toto/attestation#565; `out-of-band/MANIFEST.md` gives a direct URL and the digest for each, so
a reviewer with only public access can fetch them and get the same bytes.
Point `EVAL_EVIDENCE_SRC` at a directory holding them if they are not in `out-of-band/`.

Exit codes are deliberate:

| exit | meaning |
|---|---|
| 0 | all 23 vectors behaved as specified |
| 1 | a vector deviated — a real finding |
| 2 | the bundle is absent, so nothing was verified. **Not a failure, and not a pass.** |

A vector's expected answer is one of `accept`, `reject`, or `indeterminate`. The third one exists
for the case where the verifier was not given an artifact it would need (the dataset behind
`total`, or the per-task rows behind `numerator`): the honest answer is that it cannot decide, and
reporting `accept` there would be a pass the verifier did not earn. Those two vectors (I1, I2) are
the only ones whose conditions are about the *verifier's* inputs rather than the artifact's content.

`validate_profile_js.mjs` is a second implementation of the same rules (Node, no dependencies,
written as an independent rewrite of `validate_profile.py`). It runs the same 23 vectors and
writes `profile_vector_results_js.json`; exit codes are identical. Two implementations agreeing
is not an interoperability proof — they share one author — but it does show the verdicts do not
depend on one language's JSON or decimal handling.

**A fourth round of vectors found a defect in this implementation, not in the vectors.**
PV-4 required three fields to be present; it never checked that `excludedRowsDigest` *resolves*
to a row list, which the profile text demands. Vector **N13** builds a self-consistent artifact
with five rows excluded and a well-formed 64-hex digest pointing at nothing — the validator
returned `accept`, with no codes. The fix makes that case indeterminate
(`cannot_verify:excluded rows`), on the same principle as PV-9: an artifact the verifier was not
given is not one it may call satisfied.

The other half of that clause needed its own cases, because a suite that only shows refusals
cannot tell "refuses opaque exclusion" from "refuses all exclusion". **P6** supplies a row list
whose digest matches — `accept`; **N15** supplies one that resolves but carries four rows for
`excluded: 5` — `exclusion_row_count_mismatch`, since a digest will happily match a list of the
wrong length and only counting catches it. The profile text now names both branches
(`exclusion_opaque` when the verifier holds the list and it disagrees;
`cannot_verify:excluded rows` when it holds no list at all), which is what the code had been
doing before the prose said so.

**P5 / N14** pin the other boundary: the seven-day window is inclusive, so exactly 7 days before
passes and one second outside rejects.

Three vectors were added after the first external review of this suite:

- **N11** closes a gap the reviewer named: PV-5 requires the cannot-determine policy to be a
  *pre-registered* rule whose digest is on the artifact, but no vector exercised that clause
  (N5 covers a missing policy, not an unresolvable one).
- **N12** pins the reading of PV-7: the digest is over the four inputs' 32-byte values, not over
  the text of their hex strings. Both readings can be written to look self-consistent.
- **P4** pins the comparison base: at 18 significant digits, a binary-float comparison reads
  `9 >= 9.0000000000000001` as true and rejects an honest `passed=false`; the decimal comparison
  PV-3 requires accepts it.

`examples/qwen2.5-coder-3b-q2K.statement.json` is deliberately **non-conformant**: P1 asserts it
fails on `template_not_run_bound`, because its chat-template digest was captured about three months
after `evaluatedAt`. P2 builds the corrected variant in memory (template `servedAt` set to the run
time) and asserts it verifies clean — there is no separate `-fixed` file on disk.
