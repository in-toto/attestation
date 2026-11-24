# out-of-band evidence — manifest, not data

Four files cited in in-toto/attestation#565 (Aditi's worked example) are the ground truth this
profile was validated against. They are **not committed here**: redistributing someone else's
evaluation data in this repository is not mine to decide, and a third of a megabyte of it would sit
in the tree forever.

What is committed is where to get each file and the digest it must have, so a reviewer can fetch
the originals and prove they are the same bytes. A clean clone plus these four downloads is enough
to reproduce the run; nothing here needs our credentials.

## Where to get them (verified 2026-11-24)

Each URL below was fetched anonymously on 2026-11-24: HTTP 200, and the sha256 of the bytes
received equals the digest in the table.

| file | direct URL | sha256 | size |
|---|---|---|---|
| `smoke_suite_v1.json` | <https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/smoke_suite_v1.json> | `0f50c6f29e2dfd2afe6153a60d9e806d7b8292bbebe23fe6276340a2180abf96` | 456 B |
| `qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl` | <https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/results/transcripts/qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl> | `8d4ad7e0ded869f685a614379159894ccb2156c743ebbf35731afa39ba51809b` | 53,791 B |
| `census_rescore.jsonl` | <https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/rerun-2026-10/census_rescore.jsonl> | `fd448d2d6f94e877a09e5b8b444e0d148d7ffa984b8680c7ac383bfefd208349` | 156,328 B |
| `census_digests.json` | <https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/rerun-2026-10/serial/census_digests.json> | `210dfbe1fc7a627f47c0a3720c662a341ea34d98d0fa1d4cf7908c49bbd003b4` | 140,355 B |

```bash
mkdir -p tools/eval-evidence/out-of-band && cd tools/eval-evidence/out-of-band
curl -fLO https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/smoke_suite_v1.json
curl -fL -o qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl \
  https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/results/transcripts/qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl
curl -fLO https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/rerun-2026-10/census_rescore.jsonl
curl -fL -o census_digests.json \
  https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/rerun-2026-10/serial/census_digests.json
sha256sum smoke_suite_v1.json qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl census_rescore.jsonl census_digests.json
python3 ../../validate_profile.py --build   # from tools/eval-evidence
python3 ../../validate_profile.py           # 14 vectors, with the bundle present
```

(`census_rescore.jsonl` and `census_digests.json` live under `rerun-2026-10/` in the original
bundle; the flat names above are what a single directory needs.)

The links point at the upstream author's tagged release (`v3.4.1`), not at a copy we host. If that
tag is ever moved or deleted, the digests above are what a replacement must match; a bundle whose
digest differs is a different bundle, not an update.

Without the files present, the run exits 2 and the two rebuild vectors report `skipped`, never
`ok`: a missing out-of-band bundle means the claim was not re-derived, and the run says so
instead of passing. Exit 2 is "not verified" — it is neither a pass nor a finding.
