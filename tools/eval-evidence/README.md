# eval-evidence validator

Companion to [`spec/predicates/eval-result-evidence-profile-v0.1.md`](../../spec/predicates/eval-result-evidence-profile-v0.1.md).

```bash
pip install jsonschema
python3 validate_profile.py            # 14 vectors against the committed examples
python3 validate_profile.py --build    # rebuild examples from the cited bundle, then run
```

**The cited bundle is not committed here.** The four files are third-party evidence from
in-toto/attestation#565; `out-of-band/MANIFEST.md` lists their digests and how to place them.
Point `EVAL_EVIDENCE_SRC` at a directory holding them if they are not in `out-of-band/`.

Exit codes are deliberate:

| exit | meaning |
|---|---|
| 0 | all 14 vectors behaved as specified |
| 1 | a vector deviated — a real finding |
| 2 | the bundle is absent, so nothing was verified. **Not a failure, and not a pass.** |

A vector's expected answer is one of `accept`, `reject`, or `indeterminate`. The third one exists
for the case where the verifier was not given an artifact it would need (the dataset behind
`total`, or the per-task rows behind `numerator`): the honest answer is that it cannot decide, and
reporting `accept` there would be a pass the verifier did not earn. Those two vectors (I1, I2) are
the only ones whose conditions are about the *verifier's* inputs rather than the artifact's content.

`examples/qwen2.5-coder-3b-q2K.statement.json` is deliberately **non-conformant**: P1 asserts it
fails on `template_not_run_bound`, because its chat-template digest was captured about three months
after `evaluatedAt`. P2 builds the corrected variant in memory (template `servedAt` set to the run
time) and asserts it verifies clean — there is no separate `-fixed` file on disk.
