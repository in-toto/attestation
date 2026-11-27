#!/usr/bin/env node
// Second implementation of the eval-evidence-profile v0.1 validator, in Node.
// Purpose: independent re-implementation of the same checks on the same bytes, so that
// agreement between two languages is evidence a reader can recompute (not a claim).
// Scope: verify() + the committed example + the 15 mutation vectors. It does not rebuild
// the example (that needs the out-of-band corpora and is the first implementation's job).
//
// Usage: node validate_profile_js.mjs
// Exit: 0 all vectors behaved as specified; 1 mismatch; 2 out-of-band bundle missing; 3 published pair inconsistent.

import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";
import { fileURLToPath } from "node:url";

const HERE = path.dirname(fileURLToPath(import.meta.url));
const OOB = path.join(HERE, "out-of-band");
const EX = path.join(HERE, "examples");

const sha = (p) => crypto.createHash("sha256").update(fs.readFileSync(p)).digest("hex");
const shaText = (s) => crypto.createHash("sha256").update(s, "utf8").digest("hex");

function jcs(o) {
  if (o === null) return "null";
  const t = typeof o;
  if (t === "boolean") return o ? "true" : "false";
  if (t === "number") {
    if (!Number.isInteger(o)) throw new Error("jcs: non-integer number " + o);
    return String(o);
  }
  if (t === "string") return JSON.stringify(o);
  if (Array.isArray(o)) return "[" + o.map(jcs).join(",") + "]";
  if (t === "object") {
    const ks = Object.keys(o).sort();
    return "{" + ks.map((k) => JSON.stringify(k) + ":" + jcs(o[k])).join(",") + "}";
  }
  throw new Error("jcs: unsupported type " + t);
}

function stmtDigest(stmt) {
  const norm = JSON.parse(JSON.stringify(stmt));
  delete norm.signatures;
  for (const e of ((norm.predicate || {}).evidence || [])) {
    if (e.name === "eval-evidence-profile") e.digest.sha256 = "";
  }
  return shaText(jcs(norm));
}

// Exact decimal comparison. Using binary doubles here makes 9 >= 9.0000000000000001
// evaluate true (both round to 9), which silently accepts a false claim. Decimal strings
// are compared as scaled integers instead, so the second implementation cannot "pass"
// a threshold it never met.
function decParts(s) {
  const m = /^([+-]?)(\d*)(?:\.(\d*))?$/.exec(String(s).trim());
  if (!m) throw new Error("non-decimal: " + s);
  return { sign: m[1] === "-" ? -1 : 1, ip: m[2] || "0", fp: m[3] || "" };
}

function cmpExact(a, b) {
  const A = decParts(a), B = decParts(b);
  const f = Math.max(A.fp.length, B.fp.length);
  const ia = BigInt(A.ip + A.fp.padEnd(f, "0")) * BigInt(A.sign);
  const ib = BigInt(B.ip + B.fp.padEnd(f, "0")) * BigInt(B.sign);
  return ia < ib ? -1 : (ia > ib ? 1 : 0);
}

function cmpDecimal(a, op, b) {
  const c = cmpExact(a, b);
  if (op === ">=") return c >= 0;
  if (op === ">") return c > 0;
  if (op === "<=") return c <= 0;
  if (op === "<") return c < 0;
  throw new Error("bad comparator " + op);
}

const HEX64 = /^[0-9a-f]{64}$/;

function policyRuleResolves(ruleId, ev) {
  if (typeof ruleId !== "string" || !ruleId) return false;
  for (const section of ["exclusionPolicy", "aggregation"]) {
    const block = ev[section] || {};
    if (typeof block !== "object" || block.ruleId !== ruleId) continue;
    let dig = block.ruleDigest;
    if (dig && typeof dig === "object") dig = dig.sha256;
    if (typeof dig === "string" && HEX64.test(dig)) return true;
  }
  return false;
}

function verify(stmt, ev, files) {
  const codes = [];
  const p = stmt.predicate || {};
  const evAt = p.evaluatedAt || "";

  if (!["receipt", "public_model", "release_gate"].includes(p.subjectProfile)) codes.push("invalid_enum:subjectProfile");
  if (!["self_attested", "third_party", "reproduced", "enclave_attested"].includes(p.assuranceLevel)) codes.push("invalid_enum:assuranceLevel");

  // PV-1
  const c = ev.counts || {};
  const { total, evaluated, excluded, scored, cannotDetermine: cd } = c;
  if ([total, evaluated, excluded, scored, cd].some((x) => x === undefined || x === null)
      || total !== evaluated + excluded || evaluated !== scored + cd) {
    codes.push("count_arithmetic");
  }
  if (p.sampleSize !== evaluated) codes.push("sample_size_mismatch");

  // PV-1c
  const dsPath = files["dataset"];
  if (dsPath && fs.existsSync(dsPath)) {
    const ds = JSON.parse(fs.readFileSync(dsPath, "utf8"));
    let rows = null;
    if (ds && typeof ds === "object" && !Array.isArray(ds)) {
      for (const k of ["tasks", "suite_tasks", "rows", "items"]) {
        if (Array.isArray(ds[k])) { rows = ds[k].length; break; }
      }
    } else if (Array.isArray(ds)) rows = ds.length;
    if (rows === null || total !== rows) codes.push("total_unbound");
  } else {
    codes.push("cannot_verify:dataset");
  }

  // PV-2 / PV-3 / PV-5
  const sClaims = p.claims || [];
  const eClaims = ev.claims || [];
  if (sClaims.length !== 1 || eClaims.length !== 1) {
    codes.push("claim_binding");
  } else {
    for (const key of ["metric", "comparator", "threshold", "passed"]) {
      if (!(key in sClaims[0]) || sClaims[0][key] !== eClaims[0][key]) codes.push("claim_binding");
    }
  }
  for (const cl of eClaims) {
    const den = cl.denominator, num = cl.numerator;
    if (den !== scored) codes.push("denominator_mismatch");
    if (!Number.isInteger(num) || num < 0 || (Number.isInteger(den) && num > den)) codes.push("numerator_out_of_range");
    if (![">=", ">", "<=", "<"].includes(cl.comparator)) codes.push("invalid_enum:comparator");
    try {
      if (cl.passed !== cmpDecimal(num, cl.comparator, cl.threshold)) codes.push("verdict_not_recomputable");
    } catch (e) { codes.push("verdict_not_recomputable"); }
    if (cl.metricKind === "failure_absence" && cd && cl.passed === true) {
      const policy = cl.cannotDeterminePolicy;
      if (!policy) codes.push("cannot_determine_folded");
      else if (!policyRuleResolves(policy, ev)) codes.push("cannot_determine_policy_unbound");
    }
  }

  // PV-2b
  const ptPath = files["per-task scores"];
  if (ptPath && fs.existsSync(ptPath) && eClaims.length) {
    const tag = (p.model || {}).name;
    let rec = null;
    for (const line of fs.readFileSync(ptPath, "utf8").split("\n")) {
      if (!line.trim()) continue;
      const r = JSON.parse(line);
      if (r.tag === tag) { rec = r; break; }
    }
    if (rec === null) codes.push("numerator_unbound");
    else {
      const nTrue = Object.values(rec.per_task || {}).filter((v) => v).length;
      if (eClaims[0].numerator !== nTrue) codes.push("numerator_unbound");
    }
  } else if (eClaims.length) {
    codes.push("cannot_verify:per-task scores");
  }

  // PV-4
  const xp = ev.exclusionPolicy || {};
  if ((excluded || 0) > 0 && !(xp.ruleDigest && xp.excludedRowsDigest && xp.description)) codes.push("exclusion_opaque");
  // presence is not resolution: with rows dropped and no row list supplied, the honest
  // verdict is indeterminate (found by N13, where a well-formed fake digest was accepted)
  if ((excluded || 0) > 0 && !Object.keys(files).some((k) => k === "excluded rows" || k === "excluded")) codes.push("cannot_verify:excluded rows");

  // PV-6
  const g = ev.generation || {};
  const ct = g.chatTemplate || {};
  const rt = g.runtime || {};
  if (!((ct.digest) && g.decode && rt.version && rt.digest)) {
    codes.push("generation_undisclosed");
  } else {
    try {
      const served = Date.parse(String(ct.servedAt).replace("Z", "+00:00"));
      const evalat = Date.parse(evAt.replace("Z", "+00:00"));
      const lag = (evalat - served) / 1000;
      if (Number.isNaN(served) || Number.isNaN(evalat) || lag < 0 || lag > 7 * 86400) codes.push("template_not_run_bound");
    } catch (e) { codes.push("template_not_run_bound"); }
  }

  // PV-7
  const r = ev.reproducibility || {};
  if ((r.grade === undefined ? "none" : r.grade) !== "none") {
    const inputs = r.inputs || [];
    const hexes = [];
    if (Array.isArray(inputs) && inputs.length === 4) {
      for (const it of inputs) {
        if (typeof it === "string" && it.startsWith("sha256:") && it.length === 71 && HEX64.test(it.slice(7))) hexes.push(it.slice(7));
      }
    }
    if (hexes.length !== 4 || r.algorithm !== "sha256-concat-v1") {
      codes.push("reproducibility_unbound");
    } else {
      // sha256 over the raw 32-byte digest values (no separators, no re-hashing of the text form)
      const want = crypto.createHash("sha256").update(Buffer.from(hexes.join(""), "hex")).digest("hex");
      if (want !== ((r.digest || {}).sha256)) codes.push("reproducibility_unbound");
    }
    if (["third_party", "reproduced"].includes(r.grade)) {
      const repro = (p.evidence || []).filter((e) => String(e.name || "").startsWith("reproduction"));
      const rb = r.reproducer || {};
      if (!repro.length || !(rb.key_id && rb.signature)) codes.push("grade_unbound");
    }
  }

  // evidence digests
  for (const e of (p.evidence || [])) {
    const name = e.name || "";
    const ref = files[name];
    if (ref === undefined) codes.push("cannot_verify:" + name);
    else {
      const d = HEX64.test(ref) ? ref : sha(ref);
      if (d !== ((e.digest || {}).sha256)) codes.push("evidence_digest_mismatch:" + name);
    }
  }

  // statement binding
  const declared = (ev.statement || {}).sha256;
  if (declared && declared !== stmtDigest(stmt)) codes.push("statement_binding_mismatch");

  return [...new Set(codes)].sort();
}

// ---------------------------------------------------------------- harness
const MISSING = ["smoke_suite_v1.json", "qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl", "census_rescore.jsonl", "census_digests.json"]
  .filter((n) => !fs.existsSync(path.join(OOB, n)));
if (MISSING.length) {
  console.error("out-of-band bundle missing: " + MISSING.join(", ") + "\nexit 2 means 'not verified', not 'specification violated'.");
  process.exit(2);
}

const evPath = path.join(EX, "qwen2.5-coder-3b-q2K.evidence.json");
const stmtPath = path.join(EX, "qwen2.5-coder-3b-q2K.statement.json");
let stmt = JSON.parse(fs.readFileSync(stmtPath, "utf8"));
let ev = JSON.parse(fs.readFileSync(evPath, "utf8"));

const files = {
  "eval-evidence-profile": evPath,
  "generations": path.join(OOB, "qwen2.5-coder_3b-instruct-q2_K.smoke.jsonl"),
  "per-task scores": path.join(OOB, "census_rescore.jsonl"),
  "run record": path.join(OOB, "census_digests.json"),
  "dataset": path.join(OOB, "smoke_suite_v1.json"),
};

const base = verify(stmt, ev, files).filter((c) => !c.startsWith("cannot_verify:"));
if (process.argv.includes("--digest")) {
  // cross-language canonicalization check: the statement digest must be identical in both
  // implementations, otherwise a whole class of binding checks would silently diverge.
  console.log("JS_stmt_digest " + stmtDigest(stmt));
  process.exit(0);
}
if (base.join(",") !== "template_not_run_bound") {
  console.error("published example is self-inconsistent: " + base.join(","));
  process.exit(3);
}

const rows = [];
function run(name, s, e, f, expectOk, expectCodes = []) {
  const codes = verify(s, e, f);
  const hard = codes.filter((c) => !c.startsWith("cannot_verify:"));
  const verdict = hard.length ? "reject" : (codes.length ? "indeterminate" : "accept");
  const want = expectOk === true ? "accept" : (expectOk === false ? "reject" : "indeterminate");
  const ok = verdict === want && expectCodes.every((c) => codes.includes(c));
  rows.push({ name, expect_ok: expectOk, verdict, codes, pass: ok });
}
const clone = (o) => JSON.parse(JSON.stringify(o));

run("P1 published example as-is (template captured late)", stmt, ev, files, false, ["template_not_run_bound"]);

const e2 = clone(ev);
e2.generation.chatTemplate.servedAt = "2026-08-12T01:19:14Z";
e2.statement.sha256 = stmtDigest(stmt);
run("P2 same data, template digest run-bound", stmt, e2, files, true);

let m = clone(ev); m.counts.total = 160;
run("N1 arithmetic broken (total 160)", stmt, m, files, false, ["count_arithmetic"]);

m = clone(ev); m.claims[0].denominator = 164;
run("N2 denominator = total instead of evaluated", stmt, m, files, false, ["denominator_mismatch"]);

m = clone(ev); m.claims[0].passed = true;
run("N3 verdict flipped to pass with 0/15", stmt, m, files, false, ["verdict_not_recomputable"]);

m = clone(ev);
m.counts.total = 20; m.counts.excluded = 5;
m.exclusionPolicy = { ruleId: "stage-2", description: "synthetic: 5 rows held back, no rule digest" };
run("N4 excluded rows without rule digest", stmt, m, files, false, ["exclusion_opaque"]);

m = clone(ev);
m.counts.cannotDetermine = 2; m.counts.scored = 13;
m.claims[0].denominator = 13; m.claims[0].metricKind = "failure_absence"; m.claims[0].passed = true;
run("N5 cannot-determine folded into a passing verdict", stmt, m, files, false, ["cannot_determine_folded"]);

m = clone(ev); m.generation.chatTemplate.servedAt = "2026-08-01T00:00:00Z";
run("N6 template recorded eleven days before the run (outside the 7-day window)", stmt, m, files, false, ["template_not_run_bound"]);

m = clone(ev); delete m.generation.decode;
run("N7 decode settings omitted", stmt, m, files, false, ["generation_undisclosed"]);

m = clone(ev); m.reproducibility.digest.sha256 = "0".repeat(64);
run("N8 reproducibility digest not recomputable", stmt, m, files, false, ["reproducibility_unbound"]);

const s2 = clone(stmt); s2.predicate.evidence[2].digest.sha256 = "9".repeat(64);
run("N9 evidence digest tampered", s2, ev, files, false, ["evidence_digest_mismatch:per-task scores"]);

const s3 = clone(stmt); s3.predicate.evidence = s3.predicate.evidence.slice(0, 1);
run("N10 evidence list emptied except profile", s3, ev, { ...files }, false, ["statement_binding_mismatch"]);

const e3 = clone(e2);
e3.counts.scored = 14; e3.counts.cannotDetermine = 1;
e3.claims[0].denominator = 14; e3.claims[0].numerator = 0; e3.claims[0].metricKind = "capability_positive";
e3.statement.sha256 = stmtDigest(stmt);
run("P3 one undecidable row, denominator=scored=14", stmt, e3, files, true);

const fNoDs = Object.fromEntries(Object.entries(files).filter(([k]) => k !== "dataset"));
run("I1 dataset absent (PV-1c not checkable)", stmt, e2, fNoDs, "indeterminate", ["cannot_verify:dataset"]);

const fNoPt = Object.fromEntries(Object.entries(files).filter(([k]) => k !== "per-task scores"));
run("I2 per-task artifact absent (PV-2b not checkable)", stmt, e2, fNoPt, "indeterminate", ["cannot_verify:per-task scores"]);

// ---- N13 / P5 / N14: added after the competitor review (the reject face was thinnest) ----
// N13: PV-4 says a digest "MUST resolve to a row list". A well-formed 64-hex digest that
// points at nothing was ACCEPTED before this vector existed — the finding, not the test.
{
  const s7 = clone(stmt);
  s7.predicate.sampleSize = 10;
  const e7 = clone(e2);
  Object.assign(e7.counts, { total: 15, evaluated: 10, excluded: 5, scored: 9, cannotDetermine: 1 });
  e7.claims[0].denominator = 9;
  e7.exclusionPolicy = {
    ruleId: "stage-2",
    description: "synthetic: 5 rows held back",
    ruleDigest: { sha256: "1".repeat(64) },
    excludedRowsDigest: { sha256: crypto.createHash("sha256").update("excluded-rows-nobody-committed", "utf8").digest("hex") },
  };
  e7.statement.sha256 = stmtDigest(s7);
  run("N13 excluded-rows digest well-formed but resolves to nothing", s7, e7, files, "indeterminate", ["cannot_verify:excluded rows"]);
}

// P5 / N14: the 7-day window is inclusive at both ends; one second outside must reject.
{
  const e5 = clone(e2);
  e5.generation.chatTemplate.servedAt = "2026-08-05T01:19:14Z";
  e5.statement.sha256 = stmtDigest(stmt);
  run("P5 template served exactly 7 days before the run", stmt, e5, files, true);
}
{
  const e6 = clone(e2);
  e6.generation.chatTemplate.servedAt = "2026-08-05T01:19:13Z";
  e6.statement.sha256 = stmtDigest(stmt);
  run("N14 template one second outside the 7-day window", stmt, e6, files, false, ["template_not_run_bound"]);
}

// ---- N11 / N12 / P4: the three vectors added after the first external review ----------
// N11: PV-5 requires the cannot-determine policy to name a rule whose digest is ON the artifact.
// N5 covers a missing policy; nothing covered an unresolvable one until this vector existed.
{
  const m = clone(ev);
  m.counts.cannotDetermine = 2;
  m.counts.scored = 13;
  m.claims[0].denominator = 13;
  m.claims[0].metricKind = "failure_absence";
  m.claims[0].passed = true;
  m.claims[0].cannotDeterminePolicy = "stage-9";
  run("N11 cannot-determine policy names an unregistered rule", stmt, m, files, false, ["cannot_determine_policy_unbound"]);
}

// N12: PV-7 read the wrong way round — digest computed over the *text* of the four digests
// instead of over their 32-byte values. Both readings can claim "conformant", so the suite
// should reject this one by name.
{
  const e4 = clone(e2);
  const hexes = e4.reproducibility.inputs.map((x) => x.slice(7));
  e4.reproducibility.digest.sha256 = crypto.createHash("sha256").update(hexes.join(""), "utf8").digest("hex");
  e4.statement.sha256 = stmtDigest(stmt);
  run("N12 PV-7 digest computed over the text of the digests", stmt, e4, files, false, ["reproducibility_unbound"]);
}

// P4: threshold at 18 significant digits. A verifier comparing in binary floats sees
// 9 >= 9.0000000000000001 as true and rejects the honest verdict; a decimal verifier accepts.
{
  const s4 = clone(stmt);
  s4.predicate.claims[0].threshold = "9.0000000000000001";
  const e5 = clone(e2);
  e5.claims[0].threshold = "9.0000000000000001";
  e5.claims[0].passed = false;
  e5.statement.sha256 = stmtDigest(s4);
  run("P4 threshold at 18 significant digits (decimal, not binary float)", s4, e5, files, true);
}

const pad = (s, n) => (s.length >= n ? s : s + " ".repeat(n - s.length));
console.log(pad("vector", 62) + pad("expect", 7) + pad("verdict", 12) + pad("ok", 5) + "codes");
for (const r of rows) {
  console.log(pad(r.name, 62) + pad(String(r.expect_ok), 7) + pad(r.verdict, 12) + pad(r.pass ? "PASS" : "FAIL", 5) + (r.codes.join(",") || "-"));
}
const passed = rows.filter((r) => r.pass).length;
console.log("\n" + passed + "/" + rows.length + " vectors behaved as specified (Node implementation)");
fs.writeFileSync(path.join(HERE, "profile_vector_results_js.json"), JSON.stringify(rows, null, 2) + "\n");
process.exit(passed === rows.length ? 0 : 1);
