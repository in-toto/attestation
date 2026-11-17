# eval-evidence validator

Companion to [`spec/predicates/eval-result-evidence-profile-v0.1.md`](../../spec/predicates/eval-result-evidence-profile-v0.1.md).

```bash
pip install jsonschema
python3 validate_profile.py            # 12 vectors against the committed examples
python3 validate_profile.py --build    # rebuild examples from the cited bundle, then run
```

**The cited bundle is not committed here.** The four files are third-party evidence from
in-toto/attestation#565; `out-of-band/MANIFEST.md` lists their digests and how to place them.
Point `EVAL_EVIDENCE_SRC` at a directory holding them if they are not in `out-of-band/`.

Exit codes are deliberate:

| exit | meaning |
|---|---|
| 0 | all 12 vectors behaved as specified |
| 1 | a vector deviated — a real finding |
| 2 | the bundle is absent, so nothing was verified. **Not a failure, and not a pass.** |

`examples/qwen2.5-coder-3b-q2K.statement.json` is deliberately **non-conformant**: P1 asserts it
fails on `template_not_run_bound`, because its chat-template digest was captured about three months
after `evaluatedAt`. P2 builds the corrected variant in memory (template `servedAt` set to the run
time) and asserts it verifies clean — there is no separate `-fixed` file on disk.
