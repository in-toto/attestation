#!/usr/bin/env python3
"""Validator + example builder for the eval-evidence-profile v0.1 draft.

Checks PV-1..PV-8 from eval-evidence-profile-v0.1.md, recomputes referenced digests
from the local copies in out-of-band/, and runs mutation vectors.

Usage:
  python3 validate_profile.py --build     # write examples/*.json from the real data
  python3 validate_profile.py             # run conformance + vectors, print table
"""
import json, os, hashlib, sys, shutil, re
from datetime import datetime, timezone

HERE = os.path.dirname(os.path.abspath(__file__))
OOB = os.path.join(HERE, "out-of-band")
EX = os.path.join(HERE, "examples")
SRC = os.environ.get("EVAL_EVIDENCE_SRC", OOB)  # where the four cited artifacts live;
# default: out-of-band/, or point EVAL_EVIDENCE_SRC at a directory holding them

PROFILE_URI = "urn:lgd:in-toto:profile:eval-evidence:v0.1"


def sha(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()


def sha_text(s):
    return hashlib.sha256(s.encode("utf-8")).hexdigest()


def jcs(obj):
    if obj is None or isinstance(obj, bool):
        return json.dumps(obj)
    if isinstance(obj, int):
        return str(obj)
    if isinstance(obj, str):
        return json.dumps(obj, ensure_ascii=False)
    if isinstance(obj, list):
        return "[" + ",".join(jcs(x) for x in obj) + "]"
    if isinstance(obj, dict):
        return "{" + ",".join(jcs(k) + ":" + jcs(v) for k, v in sorted(obj.items())) + "}"


def stmt_digest(stmt):
    """Statement digest with the profile's own evidence entry blanked, so the
    evidence artifact can bind the statement that points at it (no circularity)."""
    norm = json.loads(json.dumps({k: v for k, v in stmt.items() if k != "signatures"}))
    for e in norm.get("predicate", {}).get("evidence", []):
        if e.get("name") == "eval-evidence-profile":
            e["digest"]["sha256"] = ""
    return sha_text(jcs(norm))


# ---------------------------------------------------------------- checks
def cmp_decimal(a, op, b):
    from decimal import Decimal
    A, B = Decimal(str(a)), Decimal(str(b))
    return {">=": A >= B, ">": A > B, "<=": A <= B, "<": A < B}[op]


def _policy_rule_resolves(rule_id, ev):
    """PV-5: the named rule must be one whose digest is recorded on the artifact.

    A ruleId alone says nothing a verifier can resolve; the profile requires the policy
    to be pre-registered, which on this artifact means its digest sits in
    exclusionPolicy.ruleDigest or aggregation.ruleDigest. The digest form accepted here
    is the one the schema and the examples use: {"sha256": "<64 lowercase hex>"}.
    """
    if not isinstance(rule_id, str) or not rule_id:
        return False
    for section in ("exclusionPolicy", "aggregation"):
        block = ev.get(section) or {}
        if not isinstance(block, dict) or block.get("ruleId") != rule_id:
            continue
        dig = block.get("ruleDigest")
        if isinstance(dig, dict):
            dig = dig.get("sha256")
        if isinstance(dig, str) and re.fullmatch(r"[0-9a-f]{64}", dig):
            return True
    return False


def verify(stmt, ev, files):
    """files: dict name -> path on disk for the artifacts referenced by evidence[]"""
    codes = []
    p = stmt.get("predicate", {})
    ev_at = p.get("evaluatedAt", "")

    # enum checks (added after red-team probe R4: a homoglyph passed silently)
    if p.get("subjectProfile") not in ("receipt", "public_model", "release_gate"):
        codes.append("invalid_enum:subjectProfile")
    if p.get("assuranceLevel") not in ("self_attested", "third_party", "reproduced", "enclave_attested"):
        codes.append("invalid_enum:assuranceLevel")

    # PV-1
    c = ev.get("counts", {})
    total, evaluated, excluded = c.get("total"), c.get("evaluated"), c.get("excluded")
    scored, cd = c.get("scored"), c.get("cannotDetermine")
    if None in (total, evaluated, excluded, scored, cd) or total != evaluated + excluded or evaluated != scored + cd:
        codes.append("count_arithmetic")
    # PV-1b: sampleSize in the statement must equal the evaluated count
    if p.get("sampleSize") != evaluated:
        codes.append("sample_size_mismatch")

    # PV-1c: total must equal the row count of the pinned dataset, not a free integer
    ds_path = files.get("dataset")
    if ds_path and os.path.isfile(ds_path):
        ds = json.load(open(ds_path))
        rows_in_ds = None
        for k in ("tasks", "suite_tasks", "rows", "items"):
            if isinstance(ds, dict) and isinstance(ds.get(k), list):
                rows_in_ds = len(ds[k]); break
        if rows_in_ds is None and isinstance(ds, list):
            rows_in_ds = len(ds)
        if rows_in_ds is None or total != rows_in_ds:
            codes.append("total_unbound")
    else:
        # A verifier without the dataset cannot check this rule. Saying "no problem"
        # would report a pass it did not earn; the honest verdict is indeterminate.
        codes.append("cannot_verify:dataset")

    # PV-2 / PV-3 / PV-5 — denominator is the *scored* count; cannot-determine rows are
    # reported separately by PV-5, so a claim can never be unsatisfiable merely because
    # some rows were undecidable.
    s_claims = p.get("claims") or []
    e_claims = ev.get("claims") or []
    if len(s_claims) != 1 or len(e_claims) != 1:
        codes.append("claim_binding")
    else:
        for key in ("metric", "comparator", "threshold", "passed"):
            if s_claims[0].get(key) != e_claims[0].get(key):
                codes.append("claim_binding")
    # NOTE: denominator is bound to counts.scored by PV-2 below, not to sampleSize.
    # When a row cannot be determined, denominator < sampleSize is the conformant case.
    for cl in e_claims:
        den, num = cl.get("denominator"), cl.get("numerator")
        if den != scored:
            codes.append("denominator_mismatch")
        if not isinstance(num, int) or num < 0 or (isinstance(den, int) and num > den):
            codes.append("numerator_out_of_range")
        if cl.get("comparator") not in (">=", ">", "<=", "<"):
            codes.append("invalid_enum:comparator")
        try:
            if cl.get("passed") != cmp_decimal(num, cl.get("comparator"), cl.get("threshold")):
                codes.append("verdict_not_recomputable")
        except Exception:
            codes.append("verdict_not_recomputable")
        if cl.get("metricKind") == "failure_absence" and cd and cl.get("passed") is True:
            policy = cl.get("cannotDeterminePolicy")
            if not policy:
                codes.append("cannot_determine_folded")
            elif not _policy_rule_resolves(policy, ev):
                # PV-5 names a pre-registered rule and requires its digest to appear in
                # exclusionPolicy.ruleDigest or aggregation.ruleDigest. A bare string that
                # resolves to neither is an assertion, not a registered rule.
                codes.append("cannot_determine_policy_unbound")

    # PV-2b: the numerator must equal the count of solved rows in the per-task artifact
    pt_path = files.get("per-task scores")
    if pt_path and os.path.isfile(pt_path) and e_claims:
        tag = (p.get("model") or {}).get("name")
        rec = None
        for line in open(pt_path, encoding="utf-8"):
            r = json.loads(line)
            if r.get("tag") == tag:
                rec = r; break
        if rec is None:
            codes.append("numerator_unbound")
        else:
            n_true = sum(1 for v in (rec.get("per_task") or {}).values() if v)
            if e_claims[0].get("numerator") != n_true:
                codes.append("numerator_unbound")
    elif e_claims:
        # same reasoning as PV-1c: without the per-task artifact the numerator is an
        # assertion, and a verifier that cannot recompute it must not call it fine.
        codes.append("cannot_verify:per-task scores")


    # PV-4
    xp = ev.get("exclusionPolicy") or {}
    if (excluded or 0) > 0 and not (xp.get("ruleDigest") and xp.get("excludedRowsDigest") and xp.get("description")):
        codes.append("exclusion_opaque")

    # PV-6 — the template digest must be the one that was actually served *at or before*
    # evaluation, and at most 7 days before it. A digest captured later cannot retro-explain
    # the run (the #565 worked example failed exactly here).
    g = ev.get("generation") or {}
    if not (g.get("chatTemplate", {}).get("digest") and g.get("decode") and g.get("runtime", {}).get("version")
            and g.get("runtime", {}).get("digest")):
        codes.append("generation_undisclosed")
    else:
        try:
            served = datetime.fromisoformat(g["chatTemplate"]["servedAt"].replace("Z", "+00:00"))
            evalat = datetime.fromisoformat(ev_at.replace("Z", "+00:00"))
            lag = (evalat - served).total_seconds()
            if lag < 0 or lag > 7 * 86400:
                codes.append("template_not_run_bound")
        except Exception:
            codes.append("template_not_run_bound")

    # PV-7 — inputs are four `sha256:<64 hex>` strings; algorithm sha256-concat-v1 is the
    # sha256 over the concatenation of the four 32-byte digest values (no separators, no
    # re-hashing of the text form).
    r = ev.get("reproducibility") or {}
    if r.get("grade", "none") != "none":
        inputs = r.get("inputs") or []
        hexes = []
        if isinstance(inputs, list) and len(inputs) == 4:
            for it in inputs:
                if isinstance(it, str) and it.startswith("sha256:") and len(it) == 71 \
                        and all(ch in "0123456789abcdef" for ch in it[7:]):
                    hexes.append(it[7:])
        if len(hexes) != 4 or r.get("algorithm") != "sha256-concat-v1":
            codes.append("reproducibility_unbound")
        else:
            want = hashlib.sha256(bytes.fromhex("".join(hexes))).hexdigest()
            if want != (r.get("digest") or {}).get("sha256"):
                codes.append("reproducibility_unbound")
        # grade above `self` needs a separately signed reproduction artifact
        if r.get("grade") in ("third_party", "reproduced"):
            repro = [e for e in p.get("evidence", []) if e.get("name", "").startswith("reproduction")]
            rb = r.get("reproducer") or {}
            if not repro or not (rb.get("key_id") and rb.get("signature")):
                codes.append("grade_unbound")

    # evidence digests must match files on disk (the outsider's only real guarantee)
    for e in p.get("evidence", []):
        name = e.get("name", "")
        ref = files.get(name)
        if ref is None:
            codes.append(f"cannot_verify:{name}")
        else:
            d = ref if re.fullmatch(r"[0-9a-f]{64}", ref) else sha(ref)
            if d != (e.get("digest") or {}).get("sha256"):
                codes.append(f"evidence_digest_mismatch:{name}")

    # statement binding declared by the evidence artifact (optional but checked when present)
    if ev.get("statement", {}).get("sha256") and ev["statement"]["sha256"] != stmt_digest(stmt):
        codes.append("statement_binding_mismatch")

    return sorted(set(codes))


def build(write=False):
    """Recompute the example. Only write=True touches the tree, so a plain run
    leaves the working copy clean and verifies the committed example rather than
    a rebuilt guess. A plain run still needs the out-of-band bundle to recompute
    digests, so without it the script exits 2 (not verified) rather than passing."""
    os.makedirs(OOB, exist_ok=True); os.makedirs(EX, exist_ok=True)
    missing = [n for n in ("smoke_suite_v1.json", "qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl",
                             "census_rescore.jsonl", "census_digests.json")
               if not os.path.exists(os.path.join(SRC, n))]
    if missing:
        sys.stderr.write(
            "out-of-band bundle missing: " + ", ".join(missing) + "\n"
            "These are third-party files and are not committed here.\n"
            "Place them in " + SRC + " (or set EVAL_EVIDENCE_SRC to a directory holding them).\n"
            "Digests for all four are in out-of-band/MANIFEST.md.\n"
            "Exit 2 means 'not verified', not 'specification violated'.\n")
        sys.exit(2)
    names = {
        "smoke_suite_v1.json": "dataset",
        "qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl": "generations",
        "census_rescore.jsonl": "per-task scores",
        "census_digests.json": "run record",
    }
    if write:
        for n in names:
            src, dst = os.path.join(SRC, n), os.path.join(OOB, n)
            if os.path.abspath(src) != os.path.abspath(dst):
                shutil.copy(src, dst)

    d = {n: sha(os.path.join(SRC, n)) for n in names}
    # real per-task count for this tag
    tag = "qwen2.5-coder:3b-instruct-q2_K"
    rows = [json.loads(l) for l in open(os.path.join(SRC, "census_rescore.jsonl"))]
    row = next(r for r in rows if r["tag"] == tag)
    per_task = row["per_task"]
    evaluated = len(per_task); solved = sum(1 for v in per_task.values() if v)
    suite = json.load(open(os.path.join(SRC, "smoke_suite_v1.json")))
    total_rows = len(suite["suite_tasks"]) if isinstance(suite, dict) and isinstance(suite.get("suite_tasks"), list) else 15
    rubric = sha_text("quantcheck v3.4.1: task counts as solved iff all unit tests pass; suite score = solved/15")
    agg = sha_text("pass-rate: numerator = tasks meeting the suite criterion, denominator = evaluated rows")
    commit = "cb09ecb1923f09deec128d9b276d8a83de981dfe"
    harness_d = sha_text(commit)
    template_d = "1e65450c30670713aa47fe23e8b9662bdf4065e81cc8e3cbfaa98924fcc0d320"  # served 2026-11-16, per #565

    inputs = [f"sha256:{d['smoke_suite_v1.json']}", f"sha256:{harness_d}", f"sha256:{agg}", f"sha256:{d['qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl']}"]
    repro_d = hashlib.sha256(bytes.fromhex("".join(x[7:] for x in inputs))).hexdigest()

    stmt = {
        "_type": "https://in-toto.io/Statement/v1",
        "subject": [{"name": tag, "digest": {"sha256": d["smoke_suite_v1.json"] if False else "14c4247b8ce438992cf4fb41d5461ff9d3730a561345c47a224d21ea90f8545e"}}],
        "predicateType": "https://in-toto.io/attestation/eval-result/v0.1",
        "predicate": {
            "evaluator": {"id": "https://github.com/aditi-p31"},
            "evaluatedAt": "2026-08-12T01:19:14Z",
            "suite": {"name": "quantcheck HumanEval+ screen", "version": "smoke_suite_v1"},
            "claims": [{"metric": "tasks_solved", "comparator": ">=", "threshold": "9", "passed": solved >= 9}],
            "sampleSize": evaluated,
            "model": {"name": tag,
                      "digest": {"sha256": "14c4247b8ce438992cf4fb41d5461ff9d3730a561345c47a224d21ea90f8545e"},
                      "mediaType": "application/vnd.ollama.image.model",
                      "uri": "https://registry.ollama.ai/v2/library/qwen2.5-coder/blobs/sha256:14c4247b8ce438992cf4fb41d5461ff9d3730a561345c47a224d21ea90f8545e"},
            "dataset": {"name": "smoke_suite_v1.json", "digest": {"sha256": d["smoke_suite_v1.json"]},
                        "mediaType": "application/json",
                        "downloadLocation": "https://raw.githubusercontent.com/aditi-p31/quantcheck/v3.4.1/smoke_suite_v1.json"},
            "assuranceLevel": "self_attested",
            "subjectProfile": "public_model",
            "evidence": [
                {"name": "eval-evidence-profile", "digest": {"sha256": "<computed>"}, "mediaType": "application/json"},
                {"name": "generations", "digest": {"sha256": d["qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl"]}, "mediaType": "application/jsonl"},
                {"name": "per-task scores", "digest": {"sha256": d["census_rescore.jsonl"]}, "mediaType": "application/jsonl"},
                {"name": "run record", "digest": {"sha256": d["census_digests.json"]}, "mediaType": "application/json"},
            ],
            "harness": {"name": "quantcheck", "version": "3.4.1", "digest": {"sha256": harness_d}},
        },
    }

    ev = {
        "profile": PROFILE_URI,
        "statement": {"sha256": ""},
        "counts": {"total": total_rows, "evaluated": evaluated, "excluded": total_rows - evaluated,
                   "scored": evaluated, "cannotDetermine": 0},
        "claims": [{"metric": "tasks_solved", "comparator": ">=", "threshold": "9",
                    "numerator": solved, "denominator": evaluated, "passed": solved >= 9,
                    "metricKind": "capability_positive"}],
        "exclusionPolicy": None,
        "generation": {"chatTemplate": {"digest": {"sha256": template_d},
                                        "servedAt": "2026-11-16T00:00:00Z",
                                        "source": "https://registry.ollama.ai/v2/library/qwen2.5-coder/blobs/sha256:1e65450c30670713aa47fe23e8b9662bdf4065e81cc8e3cbfaa98924fcc0d320"},
                       "decode": {"temperature": 0, "seed": 42, "numCtx": 2048, "numPredict": 768},
                       "runtime": {"name": "ollama", "version": "0.12.3", "digest": {"sha256": sha_text("ollama 0.12.3")}}},
        "aggregation": {"ruleId": "pass-rate", "ruleDigest": {"sha256": agg},
                        "description": "numerator = tasks whose unit tests all pass; denominator = evaluated rows"},
        "reproducibility": {"grade": "self", "algorithm": "sha256-concat-v1", "inputs": inputs,
                            "digest": {"sha256": repro_d}},
        "drift": [{"impl": "proofbundle v6.5.0", "field": "subjectProfile", "expected": "public_model",
                   "observed": "public-model", "note": "b7n0de, #565 2026-11-16: pre-#575 hyphen spelling"}],
    }
    # bind: statement digest into evidence, evidence digest into statement
    ev["statement"]["sha256"] = stmt_digest(stmt)
    ev_path = os.path.join(EX, "qwen2.5-coder-3b-q2K.evidence.json")
    if write:
        with open(ev_path, "w", newline="\n") as f:
            json.dump(ev, f, indent=2, ensure_ascii=False); f.write("\n")
    elif os.path.exists(ev_path):
        # a plain run must bind the digest of the *committed* evidence file, not one
        # rebuilt on the fly, otherwise verification and the marker drift apart
        pass
    stmt["predicate"]["evidence"][0]["digest"]["sha256"] = sha(ev_path)
    if write:
        with open(os.path.join(EX, "qwen2.5-coder-3b-q2K.statement.json"), "w", newline="\n") as f:
            json.dump(stmt, f, indent=2, ensure_ascii=False); f.write("\n")
    name_map = {"eval-evidence-profile": ev_path,
                "generations": os.path.join(SRC, "qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl"),
                "per-task scores": os.path.join(SRC, "census_rescore.jsonl"),
                "run record": os.path.join(SRC, "census_digests.json"),
                "dataset": os.path.join(SRC, "smoke_suite_v1.json")}
    if not write:
        # verify exactly the committed example, so what a reviewer reads is what runs
        with open(ev_path) as f:
            ev = json.load(f)
        with open(os.path.join(EX, "qwen2.5-coder-3b-q2K.statement.json")) as f:
            stmt = json.load(f)
    return stmt, ev, name_map


def vectors():
    stmt, ev, files = build()

    # Self-consistency gate. If the committed example and the evidence file it names do not
    # agree, the only honest thing to do is stop: every vector below would still "pass"
    # against a mutated copy while the published pair is broken, and the extra evidence
    # codes would be read as a property of some negative case rather than of the release.
    base = [c for c in verify(stmt, ev, files) if not c.startswith("cannot_verify:")]
    if base != ["template_not_run_bound"]:
        sys.stderr.write("published example is self-inconsistent: %s\n" % ",".join(base))
        sys.exit(3)

    rows = []

    def run(name, s, e, f, expect_ok, expect_codes=()):
        codes = verify(s, e, f)
        hard = [c for c in codes if not c.startswith("cannot_verify:")]
        verdict = "reject" if hard else ("indeterminate" if codes else "accept")
        want = "accept" if expect_ok is True else ("reject" if expect_ok is False else "indeterminate")
        ok = (verdict == want) and all(c in codes for c in expect_codes)
        rows.append({"name": name, "expect_ok": expect_ok, "verdict": verdict,
                     "codes": codes, "pass": ok})

    # P1: the *published* example as it stands. It cannot satisfy PV-6 — its chat-template
    # digest was captured 2026-11-16, three months after evaluatedAt. The profile is meant
    # to catch exactly that, so P1 is expected to fail with template_not_run_bound.
    run("P1 published example as-is (template captured late)", stmt, ev, files, False, ["template_not_run_bound"])

    # P2: same example with the template digest recorded at run time -> conformant
    e2 = json.loads(json.dumps(ev)); e2["generation"]["chatTemplate"]["servedAt"] = "2026-08-12T01:19:14Z"
    e2["statement"]["sha256"] = stmt_digest(stmt)
    run("P2 same data, template digest run-bound", stmt, e2, files, True)

    m = json.loads(json.dumps(ev)); m["counts"]["total"] = 160
    run("N1 arithmetic broken (total 160)", stmt, m, files, False, ["count_arithmetic"])
    m = json.loads(json.dumps(ev)); m["claims"][0]["denominator"] = 164
    run("N2 denominator = total instead of evaluated", stmt, m, files, False, ["denominator_mismatch"])
    m = json.loads(json.dumps(ev)); m["claims"][0]["passed"] = True
    run("N3 verdict flipped to pass with 0/15", stmt, m, files, False, ["verdict_not_recomputable"])
    m = json.loads(json.dumps(ev)); m["counts"].update({"total": 20, "excluded": 5})
    m["exclusionPolicy"] = {"ruleId": "stage-2", "description": "synthetic: 5 rows held back, no rule digest"}
    run("N4 excluded rows without rule digest", stmt, m, files, False, ["exclusion_opaque"])
    m = json.loads(json.dumps(ev)); m["counts"]["cannotDetermine"] = 2; m["counts"]["scored"] = 13
    m["claims"][0]["denominator"] = 13
    m["claims"][0]["metricKind"] = "failure_absence"; m["claims"][0]["passed"] = True
    run("N5 cannot-determine folded into a passing verdict", stmt, m, files, False, ["cannot_determine_folded"])
    # N6 must be distinct from P1: P1 already covers "captured after the run". This one
    # covers the other side of the same one-directional rule — recorded eleven days before
    # the run, outside the seven-day window. Found by diffing the vector set against the
    # base pair, which showed N6 and P1 carrying the identical change.
    m = json.loads(json.dumps(ev)); m["generation"]["chatTemplate"]["servedAt"] = "2026-08-01T00:00:00Z"
    run("N6 template recorded eleven days before the run (outside the 7-day window)", stmt, m, files, False, ["template_not_run_bound"])
    m = json.loads(json.dumps(ev)); del m["generation"]["decode"]
    run("N7 decode settings omitted", stmt, m, files, False, ["generation_undisclosed"])
    m = json.loads(json.dumps(ev)); m["reproducibility"]["digest"]["sha256"] = "0" * 64
    run("N8 reproducibility digest not recomputable", stmt, m, files, False, ["reproducibility_unbound"])
    s2 = json.loads(json.dumps(stmt)); s2["predicate"]["evidence"][2]["digest"]["sha256"] = "9" * 64
    run("N9 evidence digest tampered", s2, ev, files, False, ["evidence_digest_mismatch:per-task scores"])
    s3 = json.loads(json.dumps(stmt)); s3["predicate"]["evidence"] = s3["predicate"]["evidence"][:1]
    run("N10 evidence list emptied except profile", s3, ev, {**files}, False, ["statement_binding_mismatch"])

    # N11 — PV-5 says the cannot-determine policy must name a *pre-registered* rule whose digest
    # sits on the artifact. Until this vector existed nothing exercised that clause: N5 covers a
    # missing policy, not an unresolvable one. A bare ruleId matching neither exclusionPolicy nor
    # aggregation with a 64-hex digest is an assertion, not a registration.
    m = json.loads(json.dumps(ev)); m["counts"]["cannotDetermine"] = 2; m["counts"]["scored"] = 13
    m["claims"][0]["denominator"] = 13
    m["claims"][0]["metricKind"] = "failure_absence"; m["claims"][0]["passed"] = True
    m["claims"][0]["cannotDeterminePolicy"] = "stage-9"
    run("N11 cannot-determine policy names an unregistered rule", stmt, m, files, False,
        ["cannot_determine_policy_unbound"])

    # N12 — PV-7 read the wrong way round: the reproducibility digest computed over the *text*
    # of the four digests instead of over their 32-byte values. Both readings can be written to
    # look self-consistent and both can claim to be conformant, so the set rejects one by name.
    m = json.loads(json.dumps(e2))
    hexes = [x[len("sha256:"):] for x in m["reproducibility"]["inputs"]]
    m["reproducibility"]["digest"]["sha256"] = hashlib.sha256("".join(hexes).encode("utf-8")).hexdigest()
    m["statement"]["sha256"] = stmt_digest(stmt)
    run("N12 PV-7 digest computed over the text of the digests", stmt, m, files, False,
        ["reproducibility_unbound"])

    # P3: one undecidable row -> denominator = scored (14) < sampleSize (15). This is the
    # conformant case PV-2/ PV-5 are designed to allow; numerator 0 matches the real per-task n_true.
    e3 = json.loads(json.dumps(e2))
    e3["counts"]["scored"] = 14; e3["counts"]["cannotDetermine"] = 1
    e3["claims"][0]["denominator"] = 14
    e3["claims"][0]["numerator"] = 0
    e3["claims"][0]["metricKind"] = "capability_positive"
    e3["statement"]["sha256"] = stmt_digest(stmt)
    run("P3 one undecidable row, denominator=scored=14", stmt, e3, files, True)

    # P4 — a threshold at 18 significant digits, with the honest verdict `passed=false`. A
    # verifier comparing in binary floats reads 9 >= 9.0000000000000001 as true (both round to
    # 9) and rejects this; the decimal comparison PV-3 requires accepts it. This guards the
    # comparison base, not just the integer arithmetic every other vector uses.
    s4 = json.loads(json.dumps(stmt)); s4["predicate"]["claims"][0]["threshold"] = "9.0000000000000001"
    e4 = json.loads(json.dumps(e2))
    e4["claims"][0]["threshold"] = "9.0000000000000001"
    e4["claims"][0]["passed"] = False
    e4["statement"]["sha256"] = stmt_digest(s4)
    run("P4 threshold at 18 significant digits (decimal, not binary float)", s4, e4, files, True)

    # I1 / I2 — the third verdict. A verifier that was not given the artifact cannot
    # decide, and must say so instead of reporting acceptance it did not earn. Both use
    # the run-bound evidence so that nothing else in the case is wrong.
    f_nods = {k: v for k, v in files.items() if k != "dataset"}
    run("I1 dataset absent (PV-1c not checkable)", stmt, e2, f_nods, "indeterminate", ["cannot_verify:dataset"])
    f_nopt = {k: v for k, v in files.items() if k != "per-task scores"}
    run("I2 per-task artifact absent (PV-2b not checkable)", stmt, e2, f_nopt, "indeterminate", ["cannot_verify:per-task scores"])

    print(f"{'vector':52} {'expect':6} {'verdict':11} {'ok':4} codes")
    for r in rows:
        print(f"{r['name']:52} {str(r['expect_ok']):6} {r['verdict']:11} {'PASS' if r['pass'] else 'FAIL':4} {','.join(r['codes']) or '-'}")
    print(f"\n{sum(r['pass'] for r in rows)}/{len(rows)} vectors behaved as specified")
    with open(os.path.join(HERE, "profile_vector_results.json"), "w", newline="\n") as f:
        json.dump(rows, f, indent=2)
    return rows


if __name__ == "__main__":
    if "--build" in sys.argv:
        stmt, ev, files = build(write=True)
        print("built examples; counts:", ev["counts"], "numerator:", ev["claims"][0]["numerator"])
        print("verified codes:", verify(stmt, ev, files) or "none")
    else:
        sys.exit(0 if all(r["pass"] for r in vectors()) else 1)
