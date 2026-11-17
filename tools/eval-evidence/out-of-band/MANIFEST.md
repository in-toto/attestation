# out-of-band evidence — manifest, not data

Four files cited in in-toto/attestation#565 (Aditi's worked example) are the ground truth this
profile was validated against. They are **not committed here**: redistributing someone else's
evaluation data in this repository is not mine to decide, and a third of a megabyte of it would sit
in the tree forever.

What is committed is the digest of each file, so a reviewer can fetch the originals and prove they
are the same bytes:

| file | sha256 |
|---|---|
| `smoke_suite_v1.json` | `0f50c6f29e2dfd2afe6153a60d9e806d7b8292bbebe23fe6276340a2180abf96` |
| `qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl` | `8d4ad7e0ded869f685a614379159894ccb2156c743ebbf35731afa39ba51809b` |
| `census_rescore.jsonl` | `fd448d2d6f94e877a09e5b8b444e0d148d7ffa984b8680c7ac383bfefd208349` |
| `census_digests.json` | `210dfbe1fc7a627f47c0a3720c662a341ea34d98d0fa1d4cf7908c49bbd003b4` |

(`census_rescore.jsonl` and `census_digests.json` live under `rerun-2026-10/` in the original
bundle; this table keeps the flat name so a single directory works.)

To rebuild and re-verify the examples:

```bash
mkdir -p tools/eval-evidence/out-of-band
# place the four originals here with the names above, then
python3 tools/eval-evidence/validate_profile.py --build
python3 tools/eval-evidence/validate_profile.py        # 12/12
```

Without the files present, the two rebuild vectors report `skipped`, never `ok`: a missing
out-of-band bundle means the claim was not re-derived, and the run says so instead of passing.
