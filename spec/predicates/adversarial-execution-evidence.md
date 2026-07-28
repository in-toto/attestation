# Predicate type: Adversarial Execution Evidence

Type URI: https://in-toto.io/attestation/adversarial-execution-evidence/v0.6

Version: 0.6.0

Predicate Name: Adversarial Execution Evidence

> Status: DRAFT submission for vetting. The schema below matches a production
> implementation; field-level feedback is the point of this PR, and naming or
> structural changes from review are expected before vetting completes.

## Purpose

Records the evidence produced by deliberately executing an untrusted or
under-trusted software artifact (an agent tool, an MCP server, a plugin, a
build step) inside an instrumented containment substrate while a known corpus
of adversarial inputs is thrown at it. The predicate carries what was thrown,
what the substrate was configured to catch, what was actually observed (each
observation as an independently signed record: every interception, the armed
vantage the run was observed under, and the seal that the vantage stayed armed
to run-end), and an explicit statement of the observation's coverage bounds.

The design goal is that a consumer can recompute the outcome from the
attestation alone, with no call back to the producer's infrastructure and no
dependency on a document that does not travel with the statement. That goal
is met on one half and not on the other, and the split is worth stating at
the top. The reduction of the carried rows to a `result` is deterministic and
recomputable: the coverage denominator is committed by digest, so the
producer cannot assert it unilaterally, and each observation record verifies
on its own before it is read. The construction of those rows is not
recomputable. Mapping an observation to attack semantics is an
assembly-plane assertion no verifier can check, because the substrate sees a
dropped packet or a changed inode and has no notion of where one attack
begins and another ends. The substrate cryptographically proves what was
observed; the assembly plane asserts what it means.

A producer therefore cannot claim more than the carried evidence supports. A
producer claiming less is not detectable from the statement. `batchRoot` is
recomputed over the records the statement carries, so a producer that drops
an inconvenient interception and recomputes the root emits a self-consistent
statement that passes every check in this document. What the root binds is
the carried set against a party who cannot re-sign the enclosing envelope,
which is a network attacker rather than the assembly plane the substrate key
separation is written against. Completeness of that set against the run is
nowhere proven.

## Use Cases

-   An admission controller (e.g. a Kubernetes policy engine) gating a
    third-party MCP server or agent tool image on evidence that it was
    executed against a named attack corpus under an enforcing catch policy,
    with the policy digest and network posture pinned in the evidence.
-   An auditor re-verifying, offline and without trusting the producer's
    infrastructure, that a specific interception happened: the signed
    observation record binds the destination, the payload commitment, and the
    substrate context.
-   A security team comparing two runs of the same artifact: because the
    corpus manifest is digest-committed at attack granularity, a consumer
    can check "both runs assessed the same attacks" rather than take it on
    the producer's word.

Existing predicates cover adjacent but different ground. [Runtime Traces]
carries raw observed activity from a monitor, with no corpus binding, no
coverage denominator, and no per-event signature. [SCAI] carries
evidence-backed attribute assertions but does not model an adversarial corpus
or recomputable outcomes. [VSA] and [SVR] carry policy verdicts computed at
verification time, downstream of evidence like this. [Test Result] carries
test outcomes without cryptographic binding of the inputs or the
interceptions. This predicate is the evidence layer those verdict predicates
can consume. This predicate makes no cross-predicate claim: composing it with
a sibling execution predicate (for example a runtime trace of a different
execution) does not yield end-to-end coverage, and a consumer MUST NOT infer a
composite guarantee unless its policy binds both attestations to the same
execution (for example through a shared subject digest and run identifier).

## Prerequisites

The in-toto Attestation Framework, plus an understanding of
[DSSE](https://github.com/secure-systems-lab/dsse) (each observation record is
a DSSE-shaped envelope) and [RFC 8785 (JCS)](https://www.rfc-editor.org/rfc/rfc8785)
canonical JSON, which the digest bindings are defined over. Producers MUST
enforce the RFC 7493 (I-JSON) safe-integer profile on canonicalized content:
integers with magnitude at or above 2^53 MUST be rejected, so every rail
(producer and verifier, in any language) derives identical bytes.

The whole statement JSON is parsed as strict I-JSON: a duplicate member
anywhere in the statement, at any depth and not only inside a covering record
payload, makes the statement malformed. A lenient parser that silently keeps
the last of a repeated member would let two rails disagree on identical bytes,
so a verifier MUST reject a duplicate member statement-wide, fail-closed.

Strict I-JSON also constrains the bytes of every string, and for the same
reason. A verifier MUST reject, statement-wide and fail-closed, any statement
in which a string literal is not a well-formed sequence of Unicode scalar
values. That means: the statement MUST be valid UTF-8, with no overlong form
and no surrogate encoded directly in UTF-8 (CESU-8); a `\u` escape naming a
high surrogate MUST be immediately followed by a `\u` escape naming a low
surrogate, and an unpaired surrogate escape of either half is malformed; and a
string MUST NOT contain a raw unescaped character below U+0020. A `\u` escape
MUST consist of exactly four hexadecimal digits, with no sign, no whitespace
and no radix prefix, so that a reader built on a permissive integer parser does
not accept `\u+041` where a strict one rejects it. The profile also excludes the
Unicode noncharacters -- the code points U+FDD0 through U+FDEF, and U+nFFFE and
U+nFFFF in every plane -- which RFC 7493 section 2.1 forbids in the same sentence
as surrogates. A noncharacter is a valid scalar value that nothing substitutes
for, so unlike an ill-formed sequence it is not a cross-rail decoding split; it
is excluded so that a verifier implementing the RFC 7493 label does not reject a
record another verifier accepts, and it is rejected wherever a string literal
appears, at any depth and in both member-name and value position.

This rule exists because a lenient decoder does not fail on ill-formed bytes,
it substitutes U+FFFD for them, and every check downstream of the decode then
reads a string the producer never wrote. Where a digest is recomputed from
decoded strings rather than compared against carried bytes -- which is how the
`observationVocabulary` digest is defined below -- a producer could otherwise
emit ill-formed bytes, derive the digest over the substituted form, and obtain
a statement that one conforming verifier calls valid and another calls
malformed. A verifier MUST therefore apply this check to the raw bytes, before
any decoded string is read.

A verifier MUST reject, fail-closed, a statement whose JSON nesting depth
exceeds 128. Nesting depth is the number of arrays and objects that are open
at a given point, counting the outermost `{` of the statement as depth 1;
scalar values do not increase it. The bound is normative because it is not a
resource limit alone: with no bound stated, implementations pick their own, and
two conforming verifiers then disagree about whether identical bytes are
evidence at all over the entire range between their choices. The counting rule
is stated because implementations that increment per parsed value rather than
per open container arrive one level apart from an identical constant. Record
payloads are parsed under the same bound.

The identical-bytes requirement has a string half. On every signed canonical
surface (object member names in covering record payloads and the
`observationVocabulary.labels`/`caught` arrays), strings MUST be BMP-only:
no code point above U+FFFF, no surrogate pair. RFC 8785 sorts object members
by UTF-16 code unit; a verifier that instead compares Unicode code points
orders a supplementary-plane name differently from one in U+E000 through
U+FFFF, so two otherwise-conforming verifiers could disagree on whether
identical bytes are canonical, which under the coverage validity gate is
attestation-valid versus attestation-invalid on the same bytes. Restricting
the sorted strings to the BMP makes UTF-16 code-unit order and code-point
order coincide, so that divergence is unconstructible. A verifier MUST treat
a violation exactly as it treats non-canonical bytes: a supplementary-plane
member name makes the covering payload cover nothing, and a
supplementary-plane vocabulary entry makes the statement malformed.

These bounds close the divergences the text can foresee: a stated depth, a fixed
sort order, a pinned encoding. They do not close the ones it cannot. Where the
text underdetermines a reading and no conformance vector exercises it, two
implementations agreeing on that reading is evidence the text is determinate, not
proof of it -- the reading is untested rather than confirmed, and a third
implementation could differ there in silence. Conformance is established by
vectors; an agreement no vector has exercised is a candidate for the next vector,
not a settled rule.

**Run binding.** For any statement carrying at least one `basis: substrate`
row, the run binding digest is the lowercase 64-hex SHA-256 of the RFC 8785
canonicalization of the object `{"aeeBindingVersion": "1", "catchPolicy":
"<catchPolicy.digest.sha256>", "corpus": "<corpus.digest.sha256>",
"networkPosture": "<networkPosture.digest.sha256>", "runEntropy":
"<runEntropy.digest.sha256>", "subject": "<subject[0].digest.sha256>",
"substrate": "<substrate.digest.sha256>"}`. `runEntropy` is a run-start
value the substrate emits and commits inside the arming record's signature;
its pre-image is the substrate's run-start checkpoint, so two executions
sharing every other input still derive distinct bindings. The pre-image
SHOULD additionally fold in a publicly datable value that was unpredictable
before its round (a drand round output, or an epoch identifier in the
RFC 9334 Section 10.3 sense), in addition to, never in place of, the
substrate-unique run-start component, with the round reference recoverable
by the consumer (carrying it in the arming payload as producer vocabulary
suffices, since the digest binds it). A signature over a value that did not
exist before its round cannot predate the round, so the arming record gains
a proven earliest-possible signing time, a floor. `issuedAt` remains the
asserted ceiling; the pair is deliberately not a two-sided proof, by the
asserted-versus-attested rule this predicate applies everywhere. The floor
bounds recency only where consumer policy couples the folded round to its
freshness window (the producer selects the round, so an uncoupled round
proves age, never freshness), and a beacon inside the producer's own trust
domain yields no floor against that producer. The value is fetched at
arming time, never cached: a stale round silently folded as current would
defeat the same coupling. Public rounds also make two consumers'
`runEntropy`-reuse observations comparable against a shared public time
axis rather than against the producer's clock.
For this predicate `subject` MUST contain exactly one entry on a statement
of any basis; a statement carrying zero or more than one subject is
malformed, regardless of whether any row is `basis: substrate`. Separately,
each of the six digest inputs MUST carry a `sha256` digest whose value is
already lowercase 64-hex; a substrate-row-carrying statement violating this
digest-input requirement is malformed. Values are taken verbatim (no case-folding, no
null fill). A statement whose rows are all `basis: artifact` derives no
binding and need not carry `runEntropy`. A verifier derives the digest from
the statement alone; no field carries it. Every substrate-signed
observation record commits to the run by carrying this digest inside its
signed payload (see the reserved members under `observationRecords`). The
binding is anti-splice: a record signed under a different subject, corpus,
catch policy, network posture, substrate, or run-start entropy value cannot
be spliced in. It is not anti-forge and not a freshness challenge: it
carries no verifier nonce, and identical-configuration re-runs are
distinguished only by the substrate-emitted `runEntropy` value, so a
consumer that must exclude replay of a genuine record into a later
identical-configuration run does so by rejecting reuse of a `runEntropy`
value it has already seen. `aeeBindingVersion` names this construction;
a future version that changes the construction (another hash algorithm,
additional inputs, multiple subjects or substrates) names a new binding
version, and a verifier MUST reject, fail-closed, a binding version it does
not implement rather than attempt more than one construction. An arming
record's payload MAY carry an explicit `aeeBindingVersion` member declaring
its construction; a verifier reads it before deriving and rejects it
fail-closed (the arming record covers nothing) when the value is a version it
does not implement, distinguishably from a run-binding digest mismatch. An
absent member defaults to version `1`; the carried value never drives the
derivation (a verifier derives only under the version it implements, so a
record declaring `1` but constructed otherwise still fails on the digest). A future
minor version admitting multiple subjects or multiple substrates binds all
of them in canonical name-then-digest order.

## Model

The producer is a containment substrate operator: a functionary that runs the
subject artifact inside an isolated, instrumented environment (a microVM, a
sandbox, an eBPF-supervised process), injects the corpus, and signs what the
substrate observed: each interception, the armed vantage it occurred under,
and the seal that the vantage stayed armed. The subject is the executed
artifact, by digest. Every
attestation references a substrate by subject (the `substrate` field is
required); that substrate SHOULD in turn carry its own attestation, e.g. build
provenance for the substrate image, so the evidence can inherit a substrate
trust chain rather than a bare name. Verdicts (pass/fail against an
organization's policy) are deliberately out of scope; they belong in a
downstream summary predicate such as [VSA], computed over this evidence.

## Schema

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    { "name": "<executed-artifact-name>", "digest": { "sha256": "<64-hex>" } }
  ],
  "predicateType": "https://in-toto.io/attestation/adversarial-execution-evidence/v0.6",
  "predicate": {
    "result": "fail",
    "observationEnvironment": {
      "substrate": { "name": "<substrate-attestation-subject>", "digest": { "sha256": "<64-hex>" } },
      "corpus": {
        "name": "<corpus-name>",
        "uri": "pkg:<producer>/<corpus>@<version>",
        "digest": { "sha256": "<64-hex-JCS-digest-of-manifest>" },
        "manifest": { "classes": { "CO": ["CO-EXFIL-1"] } }
      },
      "catchPolicy": { "digest": { "sha256": "<64-hex-JCS-digest-of-catch-policy>" } },
      "networkPosture": { "posture": "sinkhole", "digest": { "sha256": "<64-hex>" } },
      "observationVocabulary": {
        "digest": { "sha256": "<64-hex-JCS-digest-of-vocabulary>" },
        "labels": ["egress_captured", "no_egress"],
        "caught": ["egress_captured"]
      },
      "runEntropy": { "digest": { "sha256": "<64-hex run-start value>" } }
    },
    "coverage": {
      "assessedClasses": ["CO"],
      "outOfScope": {},
      "routedElsewhere": {}
    },
    "attackResults": [
      {
        "attackId": "CO-EXFIL-1",
        "containmentObserved": "egress_captured",
        "basis": "substrate",
        "method": "intercepted",
        "actualLayer": "policy.egress_sinkhole",
        "observationRefs": [0]
      }
    ],
    "observationRecords": [
      {
        "payload": "<base64(canonical +json bytes the substrate signed)>",
        "payloadType": "<producer-defined media type ending in +json>",
        "signatures": [ { "keyid": "<hex>", "sig": "<base64>" } ]
      }
    ],
    "batchRoot": "<64-hex RFC 6962 root over observation records>",
    "doesNotAssert": [ "<explicit negative-scope statements>" ],
    "issuedAt": "2026-06-23T16:08:07Z"
  }
}
```

### Parsing Rules

The predicate opts in to the framework's standard parsing rules, including the
monotonic principle, with one deliberate strengthening: `result` is not an
independent claim. A consumer MUST be able to recompute it from the rest of
the predicate (rules under `result` below), and a `result` the recompute does
not reproduce makes the attestation invalid. Observation record `payload`s
follow the same verify-then-read discipline: the fields inside a payload mean
nothing until its signature verifies against a key the consumer trusts. The
`result` recompute is a function of the carried predicate alone: it never
reads `observationRecords`, signature-verification outcomes, or any consumer
trust decision. A `result` that varied with the consumer's trust anchors would
not be recomputable. The coverage validity requirements below (which read
record payloads but not signatures or consumer policy) and the evidence tier
(which reads signatures against consumer policy) are separate gates from the
recompute: the validity gate can invalidate an attestation, and the tier ranks
a row, but neither alters `result`.

A verifier proceeds in two stages. Stage one is byte-pure: four validity
steps, each a function of the carried statement alone, and all four are
consumption preconditions: (1) statement well-formedness, including the
vocabulary rules and, for substrate-carrying statements, run-binding
derivability; (2) the coverage validity requirements; (3) the `result`
recompute; (4) manifest and vocabulary digest integrity. Stage two is
trust-relative: the envelope signature and the per-row evidence tier
against consumer key policy, then the strength orderings and the rest of
consumer policy, including the anchor comparison under Consumer policy
obligations. Only the consumption preconditions stated under Coverage
validity and the evidence tier are normative in this ordering; the
sequencing itself is informative.

A design invariant follows from the recompute: any per-observation property
that the recompute or the documented consumer gating reads travels on the
row itself, as a required member, other than the record-borne binding
members defined under `observationRecords`, with a closed vocabulary,
fail-closed on missing or unknown values. Run-level pins in
`observationEnvironment` never substitute for a row-level property, because
the recompute reads rows. The instruments that corroborate a property may be
run-scoped: a run-level `arming` or `sealed` record backs a clean row's claim
through the row's own `observationRefs`, and never substitutes for a row
member. The arming and sealed records attest that a vantage was armed and
stayed armed run-wide, not that the specific channel for this row's attack
class was armed; per-channel arming completeness stays producer vocabulary
bounded by the pinned `networkPosture` digest.

### Fields

`result` _string, required_

One of `pass`, `degraded`, `fail` (lowercase). Defined as a total,
deterministic, severity-independent function of the predicate: `fail` iff any
`attackResults` row carries a containment-observed label from the carried
caught set (`observationVocabulary.caught`), a label outside the carried
`observationVocabulary.labels` (fail-closed), or a missing or
out-of-vocabulary `basis` or `method` (fail-closed, same rule); otherwise `degraded`
iff `coverage.outOfScope` or `coverage.routedElsewhere` is non-empty;
otherwise `pass`. A `pass` is coverage-bounded-observed: it states what was
assessed and makes no general safety claim. There is intentionally no
severity threshold, policy ruleset, or free-text reason here; policy belongs
downstream.

Two further bounds sit beside that one. The substrate is a passive sensor,
not an orchestrator: it observes what crosses the vantages it was armed at,
and it has no view of the runner that injects the corpus. It therefore cannot
verify that the runner executed every attack the manifest names. The set of
attacks actually executed is a producer assertion. Coverage integrity
compares the carried rows against the carried manifest, which catches an
omission the producer failed to declare in `coverage`, but neither comparison
reaches the run, so a consumer reading a `pass` is relying on producer
integrity that nothing was silently skipped.

The second bound is that four fields sit deliberately outside the run
binding, because the substrate operates on content digests and has no view of
presentation metadata: `corpus.name`, `corpus.uri` (RECOMMENDED as a purl),
`substrate.name`, and `subject[0].name`. No digest, no record payload, and no
gate reads them, so they are attacker-modifiable on a statement that remains
fully valid. Calling them unauthenticated understates that for a reader who
assumes the enclosing envelope signature protects everything inside it: a
statement can be relabeled as evidence about a differently named corpus,
substrate, or artifact and still satisfy every requirement in this document.
Consumers MUST anchor identity on the digests rather than on these names, as
Consumer policy obligations below requires for the corpus and the substrate.

**Coverage validity (derived from carried bytes; a violation is malformed).**
For every `basis: substrate` row, the following MUST hold or the attestation
is invalid, exactly as a missing `actualLayer` is invalid. These read record
payloads but never signatures or consumer policy, so they are a pure function
of the carried statement:

-   its `observationRefs` is non-empty and every index is in range for
    `observationRecords`;
-   the referenced records match the class the row requires: a caught row
    with `method: intercepted` references at least one `interception`
    record; a `method: reconstructed` row references at least one
    `examination` record; a clean row with `method: intercepted`
    references at least one `arming` record and at least one covering
    `sealed` record (a `sealed` record covers under the conditions stated
    at its class definition);
-   every referenced payload parses as a canonical `+json` object (see
    `observationRecords`), carries the reserved members, and its
    `aeeRunBinding` equals the run binding digest derived from this
    statement;
-   the row's `method` is no stronger than the weakest `aeeMethod` across
    its covering records (`reconstructed` is weaker than `intercepted`);
-   `batchRoot` recomputes over `observationRecords` (see `batchRoot`).

These requirements are consumption preconditions, not optional lints: a
consumer that consumes `result`, credits any row, or applies either
strength ordering MUST first evaluate them, and on failure the attestation
is invalid and its `result` MUST NOT be consumed, the same handling as
any malformed statement.

What these requirements are is worth stating as plainly as what they
require, because their name invites a reader to take them for a security
gate on their own, and on their own they are not one. Every check above
reads record content, and record content means nothing until the record's
signature verifies, which is the verify-then-read discipline stated under
Parsing Rules. Against a party able to author a record, none of these checks
costs a secret: `aeeRunBinding` is derived from six digests the statement
already carries in plaintext, `aeePostureDigest` is the pinned
`networkPosture` digest carried beside them, and an `arming` or `sealed`
payload describing a vantage that never existed satisfies every constraint
here. These requirements are therefore structural well-formedness
constraints. They become security properties only in combination with
observation-record signature verification against a substrate key the
consumer trusts, which is the evidence tier below. A verifier that evaluates
coverage validity and skips that verification has checked that the producer
filled the form in correctly.

A caught row is one whose `containmentObserved` label is in the carried
caught set (`observationVocabulary.caught`); a clean row is one whose
label is in the carried `observationVocabulary.labels` and not in the
caught set. Both sets travel in the attestation, so the distinction is a
function of carried bytes. A `basis: substrate` row whose
`containmentObserved`, `basis`, or `method` is fail-closed (outside the
carried vocabulary) cannot satisfy the class-match requirement and is
therefore invalid; a `basis: artifact` fail-closed row sits at the bottom
of both orderings as before.

**Evidence tier (derived, never carried).** Given a valid attestation, a
consumer MUST, before crediting any `basis: substrate` row or applying
either strength ordering, derive a per-row evidence tier: a `basis:
artifact` row is `declared`; a `basis: substrate` row is `attested` when
every covering record's signature verifies against a key the consumer's
policy names as a substrate observation key, and `unattested` otherwise. A
consumer with no policy-pinned substrate root MUST treat every `basis:
substrate` row as `unattested` and MUST NOT infer the substrate root from
the predicate. The tier is total and deterministic given the consumer's
key policy; it never alters `result`. Consumer policy MAY subdivide
`attested` into stricter refinements (for example requiring a
hardware-attested observation key, or agreement of multiple keys); a
refinement refines, never reorders, the three tiers, and tier names
beginning with `aee` are reserved. A carried predicate member named
`evidenceTier`, or any predicate-level member beginning with the reserved
prefix `aee`, MUST be ignored and MUST NOT alter the derivation.

`observationEnvironment` _object, required_

The digest-pinned context the evidence was earned under. Five required
members: `substrate` (the subject reference of the substrate's own
attestation), `corpus` (name, uri, RECOMMENDED as a purl; the JCS digest of
the embedded `manifest`, and the `manifest` itself: a map from assessment
class to the complete array of attack identifiers it defines; an attackId
MUST NOT appear under more than one class), `catchPolicy` (JCS digest of the
parsed catch-policy document, so an empty or permissive policy is
distinguishable from an enforcing one), `networkPosture` (the
substrate-authoritative egress posture, e.g. `no_network`, `allowlist`,
`sinkhole`, with its configuration digest), and `observationVocabulary` (the
producer's versioned observation label set carried in the attestation:
`labels`, the complete array of `containmentObserved` values the producer can
emit, and `caught`, the subset whose observation constitutes a caught
containment event; both arrays sorted ascending by UTF-16 code unit (RFC 8785
Section 3.2.3) with no duplicates and every entry BMP-only (see
Prerequisites), `caught` a subset of `labels`, and `digest` the JCS digest of
the object `{"caught": [...], "labels": [...]}`; a statement violating any
of these is malformed).
The recompute and the coverage validity requirements read only this carried
set; the producer's published documentation is commentary on the same
vocabulary, never a normative input, so archived attestations remain
verifiable after the producer's documentation moves or disappears. A sixth
member, `runEntropy` (the substrate-emitted run-start value defined under Run
binding), is required exactly when any row carries `basis: substrate`. The
manifest pre-image travels in the attestation, so a verifier re-derives
`corpus.digest` offline and any edit to the assessed set (a dropped attack, a
renamed class) fails that check.

`coverage` _object, required_

The coverage bound: `assessedClasses` (array of class codes actually
assessed), `outOfScope` and `routedElsewhere` (maps from class code to a
reason string; empty objects when complete). Disclosing a gap moves the class
into one of these maps and forces `result` to `degraded`, which is the honest
alternative to leaving it out and quietly reporting a narrower run as a full
one. The three sets are a disjoint partition of the manifest's classes: a
class appears in exactly one of `assessedClasses`, `outOfScope`,
`routedElsewhere` (a move, not a copy). A class in more than one of the three,
or a manifest class in none, is malformed - a class both assessed and
disclosed as a gap is contradictory.

`attackResults` _array of objects, required_

One row per executed attack: `attackId` (must appear in the manifest),
`containmentObserved` (a label from the carried
`observationVocabulary.labels`; consumers treat labels outside the carried
set as fail-closed), `basis` (required; the observation's vantage, see
below), `method` (required; the observation's directness, see below),
`actualLayer` (required; which enforcement layer acted, see below), and
`observationRefs`
(indexes into `observationRecords` binding this row to the observation
records that cover it). An `interception` index MAY be referenced by more
than one row. A producer MUST NOT reference a record from a row whose
attack the record's committed payload does not evidence; that is a
producer obligation outside every gate: no validity requirement,
recompute input, or tier evaluation reads it, so a conforming verifier
neither can nor may invent an evidencing heuristic for shared references.
No two `attackResults` rows may carry the same `attackId`: one row per
executed attack is a well-formedness invariant, and a statement with a
duplicate `attackId` across rows is malformed. Coverage integrity
set-compares row `attackId`s against the manifest, so a duplicate would
silently collapse under set semantics; uniqueness is enforced separately,
before that comparison, not left to it.
Wherever `observationRefs` is present - on any row, regardless of `basis`,
and including rows on which nothing normative reads it - every index MUST be
in range for `observationRecords`; an out-of-range index is a structural
integrity fault that makes the statement malformed, fail-closed and
independent of any gate, so a reference that does not resolve is never
silently ignored.
A row MAY carry `observationSelectors`, an array of
producer-defined string tokens positionally parallel to `observationRefs`,
each naming the sub-observation within the referenced record's committed
payload that this row rests on; token content is producer vocabulary,
nothing normative reads it, and selector presence or absence changes no
gate outcome. `arming`, `sealed`, and
`examination` indexes MAY likewise be shared: one run-level record covers
every row earned under it. The single normative reading of this value is its
membership in the carried caught set; see `method` for the
exhaustion-of-meaning rule.
Coverage integrity is checked at attack granularity: the union of attackIds
for the assessed classes must exactly equal the manifest's. That granularity
is what stops a failing attack from being quietly omitted inside a class the
producer still reports as assessed.

`basis` states each observation's vantage, with a closed two-value
vocabulary:

-   `substrate`: every input the row's claim depends on was obtained at a
    vantage the executed artifact could neither forge nor suppress (a
    network boundary, syscall supervision, a hypervisor's read of guest
    state). This names a class of vantage, cooperation-independence, not
    the identity of the enforcing substrate: a passive tap, an inline gate,
    and an adversarial corpus endpoint logging the connection it received
    all sit at this basis and differ only in enforcement role, which is
    `actualLayer`'s question.
-   `artifact`: at least one input the claim depends on derives from output
    the executed artifact itself produced (its stdout/stderr, exit status,
    or self-emitted logs).

An input is artifact-sourced when the claim relies on it as a channel
whose content the artifact can populate arbitrarily without performing
the claimed containment event: testimony about an event rather than the
event itself. An egress capture is not artifact-sourced even though the
packet bytes were artifact-authored, because the artifact cannot cause
the boundary to record an egress without performing one; its stdout is
artifact-sourced because the artifact can write anything there at no
cost.

`basis` is the vantage of the claim's weakest input. A derived observation
inherits `artifact` from any artifact-sourced input it consumed: a state
diff computed by substrate machinery over the artifact's own logs is
`artifact`, however trusted the machinery, because the artifact could have
populated what the machinery read without performing any containment
event.

A producer MUST NOT declare `basis: substrate` on a row it cannot cover
under the coverage validity requirements above: such a row is not merely
mislabeled, it makes the attestation invalid.

`method` states how the row's claim was established, with a closed
two-value vocabulary:

-   `intercepted`: the claim rests on events captured as they occurred. On
    a clean row, a live capture vantage was armed for the attack and no
    capture was attributed to it.
-   `reconstructed`: the claim derives from state examined after the fact,
    e.g. a snapshot-to-snapshot diff. A reconstruction can miss a transient
    raised and undone between the states it compares; on a clean row that
    tolerance is part of the claim.

Like `basis`, `method` composes by weakest input: a claim inherits
`reconstructed` from any state-derived input it depends on. Post-hoc
decode of an event stream captured as it occurred (a packet capture
parsed later, a hardware trace decoded after the run) does not demote a
row, provided the capture channel was armed for the claimed event class
before the event; a row that fuses a live capture with after-the-fact
state examination is `reconstructed`.

`method` describes how the observation was made, not how the row was
attributed to its attack. The normative content of a `containmentObserved`
value is exhausted by its membership in the carried caught set: nothing
normative (neither the `result` recompute, the coverage validity gate,
the evidence tier, nor either strength ordering) reads anything else from
the label. An axis earns its own required member exactly when a normative
reader consumes it, which is why `basis` and `method` are members (the
recompute and the gate read them) and attribution nuance is not: no
normative reader consumes attribution strength (a hash-pinned payload
versus a time-window pairing) or attribution tolerance on clean rows,
including window bleed between same-layer siblings, so they remain
non-normative producer vocabulary documented beside the vocabulary's
definition, and a consumer MUST NOT move a `result` or an evidence tier on
them.

Both fields are REQUIRED on every row and both vocabularies are closed: a
missing value, or any value outside them, is fail-closed exactly as an
out-of-vocabulary `containmentObserved` label is: the row forces
`result` to `fail` and can support nothing stronger. The 0.4 values
`substrate_observed`, `artifact_reported`, and `inferred` are
out-of-vocabulary in 0.5, with no alias, for the same
single-canonicalization reason the old `does_not_assert` spelling is
rejected. `inferred` has no successor value because it conflated the two
axes: a hypervisor snapshot diff and a diff parsed from the artifact's own
logs were both "derived indirectly" while carrying opposite vantages.
Under 0.5 the first is (`substrate`, `reconstructed`), the second is
(`artifact`, `reconstructed`), and no row can read as more independent
than its weakest input.

The two axes bind a consumer's confidence on opposite sides. A `fail`'s
supporting rows are the rows that force `result` to `fail`: rows whose
`containmentObserved` label is in the carried caught set, plus rows
fail-closed on a missing or out-of-vocabulary member. `basis` bounds a
`fail` over that supporting set: it answers whether the artifact could
have manufactured the observation. Only rows deriving `attested` occupy the
`substrate` rank: an `unattested` substrate row ranks, in both orderings,
exactly where an `artifact` row ranks. Ignoring a record, or failing to
verify one, can only weaken the claim, never strengthen it. A `fail` whose
supporting rows are all `artifact` SHOULD be treated as a weaker claim than
one whose supporting set carries a `substrate` row, and a consumer MAY
reject it; a `fail`
from a (`substrate`, `reconstructed`) row is still an observation the
artifact could not manufacture without performing a containment-relevant
event, weaker than an interception only in that it was derived after the
fact. A row fail-closed on `basis` or `method` sits at the bottom of both
orderings: its vantage is unknown, so it can strengthen nothing. `method`
bounds a `pass`: it answers whether a real event could have slipped past
the observation. A `pass` whose clean rows are all `attested` at
(`substrate`, `intercepted`), each covered by both the run's arming record
and a covering sealed record, the seal reporting under its own signature
that the vantage stayed armed to run-end, that the run-wide count of dropped
observations is zero or within a bound the same signed payload declares,
and that the posture it re-attests equals the pinned one, makes the
strongest absence claim this predicate can carry: a live vantage was armed
and stayed armed for this run and the committed record set attributes no
interception to this attack. It is not a claim that no interception
occurred; because the arming and sealed records are run-wide it is not a
claim that this row's specific channel was armed; and it is a claim about
the run this attestation carries, never about a run population. Nothing
in this predicate proves that other runs of the same configuration did not
occur or were not withheld, and run-population completeness (for example
via a run ledger or monotonic counter) is a consumer or policy concern.
The optional run-sequence members defined under `observationRecords` move
cherry-picking from invisible to gap-evident across whatever set a
producer does publish (a gap, a duplicated sequence number, a shared
predecessor, or a duplicated genesis is detectable by any consumer holding
both attestations), without changing this non-claim; their definition
states the ordering-only scope and the registration-receipt completion. An
`pass` resting on any `reconstructed` clean row SHOULD be read as
tolerating transients between the observed states; a `pass` resting on any
`artifact` clean row, or on an `unattested` substrate clean row, is
self-reported absence, the weakest. The
clean-row ordering applies equally to the clean rows of a `degraded`
result's assessed classes, with `degraded` additionally bounded by its
disclosed coverage gap. A consumer MAY reject, never downgrade: an
attestation containing `unattested` substrate rows MAY be rejected
outright, and the signed bytes are never rewritten. Declining to credit a
rank is not a downgrade of signed bytes, exactly as declining to trust a
carried `result` is not. Both orderings are consumer guidance; the evidence
tier is an input to that guidance, not a verdict, and the predicate still
carries none.

A `basis: substrate` row's coverage is a validity requirement, not weak
evidence: a caught (`substrate`, `intercepted`) row with empty or
out-of-range `observationRefs`, or one whose referenced records do not
class-match, run-bind, cap `method`, or recompute `batchRoot`, makes the
attestation invalid (see Coverage validity). These are facts about carried
bytes alone; whether the referenced records verify, and against whose key,
is the evidence tier's separate question, because an answer that varied
with the consumer's trust anchors cannot live in a validity rule. A clean
row's `intercepted` still has no per-event record by definition: an armed
vantage that captured nothing produces no interception to sign. Its
covering instruments are the run-level `arming` record (a substrate-signed
statement that a live capture vantage was armed before corpus injection)
and the `sealed` record (a substrate-signed statement that the vantage
stayed armed to run-end with no dropped observation), which is the shape
absence evidence takes elsewhere: an attested launch measurement, a
hermetic-build flag, a monitor configuration attested as first-class
evidence. Where the producer emits a checkpoint chain, "armed before the
first observed event" is a chain-order fact (the arming record is the
chain head and each interception carries a higher sequence); where it does
not, `armedAt` ordering is producer-asserted and only the arming instant
is attested. One or more arming or sealed records MAY cover a run; each
referenced record must independently satisfy its class constraints.

Fields divide by the identity whose signature backs them.
Substrate-covered, through the coverage validity gate and evidence tier:
`basis` and `method` on rows deriving `attested`, and the content of every
verified observation record. Producer-asserted, backed only by the
enclosing envelope: `containmentObserved` labels and their attribution
nuance, `basis` and `method` on `artifact` rows, `actualLayer`,
`coverage`, `doesNotAssert`, and the assembly of the predicate itself.
Which keys count as substrate observation keys is consumer key policy,
resolved where signer identity is always resolved. The substrate
observation key MUST NOT be accessible to the subject artifact, and SHOULD
be held apart from the producer's assembly plane; a consumer's policy MAY
additionally require that the key signing any covering observation record
differ from the key signing the enclosing Statement. The tier's value
against a dishonest producer is exactly that separation: where the
observation key is held apart from the assembly plane, the tier defeats a
pipeline with no substrate in the loop and cross-configuration splices,
because neither can be produced without a signature under the observation
key. It does not defeat record drops, and it does not defeat method
inflation. Both limits are structural rather than gaps in the checks above.
`batchRoot` recomputes over the carried records, so a party holding the
envelope key removes a record and recomputes a root that is self-consistent
over what remains. The `method` cap binds a row to the weakest `aeeMethod`
across the records that row references, and no record names its attack, since
a substrate signs at observation time and before attribution; the cap is
therefore per record and never per attack, so re-pointing a row's
`observationRefs` at an `intercepted` record signed for a different attack
raises that row's `method` with every substrate signature still verifying.
What holding the observation key apart defeats is the manufacture of
substrate evidence, not the assembly plane's selection and arrangement of
substrate evidence that genuinely exists. Where one party holds both keys
(the common single-root deployment), the tier instead defeats only a party
that holds neither key: the envelope signature closes the carried set to
such a party, and the run binding and the signed `aeeMethod` close record
substitution and forgery, so a tamperer with no key cannot splice, drop, or
inflate an already-signed set. A substrate operator who signs false
evidence, or who runs no substrate at all and signs an arming record anyway,
remains outside this predicate's threat model, as for every self-asserted
field. Coverage is
therefore only as trustworthy as the named key's un-compromised lifetime;
the single trust root is a single point of total failure, and a consumer's
policy MAY bound a named key with a validity window checked against
`issuedAt`. Consumers MAY additionally coherence-check row claims against
the pinned `observationEnvironment`: a `substrate` row claiming a
network-boundary observation under a `networkPosture` that provides no
interception path at that boundary is incoherent, and a consumer MAY
reject on that ground.

`actualLayer` names the enforcement layer that acted on the row's
containment event. It is required on every row; a row missing the member
is malformed under the framework's standard parsing rules and the
attestation is invalid, rather than the row forcing `fail`. That altitude
is deliberate and follows from the design invariant under Parsing Rules:
fail-closed-row semantics are reserved for members the recompute or the
documented consumer gating reads (`containmentObserved`, `basis`,
`method`); `actualLayer` is read by neither, so its absence is a
malformed statement, not weak evidence. On a row whose
`containmentObserved` label is from the carried
`observationVocabulary.labels` but not in the caught set (a clean row:
nothing acted), the producer MUST emit the literal string `none`. `none` is
explicit rather than the field being
omitted so that "no layer needed to act" is distinguishable from an
accidental omission. `none` is also valid on a caught row, and there
states that the containment event was observed but no enforcement layer
acted: the observing vantage was positioned to see, not to act (a passive
tap, a monitor-only deployment). This is deliberate: enforcement role
travels here and only here, so `basis` never has to encode who could act.
Whether anything was positioned to see is answered by `basis` and
`method`: a clean row carrying (`basis: substrate`, `method:
intercepted`) states a live substrate vantage was armed and no capture
was attributed, which is the strongest claim a `pass` can rest on, a claim
the run-level `arming` and `sealed` records under `observationRecords` now
carry under the substrate's own signature, bounded to the vantage's
existence and continuity rather than to the absence of any event; a
`reconstructed` clean row makes the bounded version of that claim
described under `method`.

`observationRecords` _array of objects, optional_

One DSSE envelope per observation: `payload` (base64 of the exact
canonical bytes the substrate signed at observation time), `payloadType`,
and `signatures`. A consumer verifies each record's signature, DSSE PAE
over `(payloadType, payload)`, before reading any field inside the
payload. Any record used to cover a `basis: substrate` row MUST carry a
JSON object payload that is canonical per RFC 8785 and valid I-JSON per
RFC 7493 (no duplicate members, integers within the safe range, member
names BMP-only per Prerequisites, every string a well-formed sequence of
Unicode scalar values per Prerequisites, and nesting within the bound stated
there), whose
media type ends in `+json`, and which carries these reserved members as
top-level fields; a record whose payload is not so parseable, or whose
media type is not `+json`, covers nothing:

-   `aeeRunBinding` _string_: the run binding digest defined under
    Prerequisites.
-   `aeeKind` _string_: `interception` (per-event capture, covers caught
    rows); `arming` (run-level: a live, cooperation-independent capture
    vantage was armed for the run before corpus injection; payload MUST
    carry `armedAt` in RFC 3339 with a zero UTC offset (`Z` or `+00:00`, never a
    non-zero offset such as `+05:00`) no later than `issuedAt` and
    `aeePostureDigest` equal to the pinned `networkPosture` digest, and
    its `aeeMethod` MUST be `intercepted`); `sealed` (run-level: the
    vantage stayed armed to run-end; payload MUST carry `aeeStillArmed`,
    a boolean; `aeeDropCount`, an integer counting run-wide dropped
    observations; and `aeePostureDigest`, the effective posture at
    run-end; it MAY carry `aeeDropBound`, a producer-declared integer
    bound; its `aeeMethod` MUST be `intercepted`); or `examination` (the
    substrate examined artifact-independent state after the fact; its
    `aeeMethod` MUST be `reconstructed`, and its payload SHOULD identify
    the states compared).
-   `aeeMethod` _string_: `intercepted` or `reconstructed`; how the
    substrate observed, stated inside the signature.

A record violating any constraint of its declared `aeeKind` (including a
missing `armedAt` on an `arming` record, an `armedAt` after `issuedAt`, an
`armedAt` with a non-zero UTC offset, or
an `examination` record signed `aeeMethod: intercepted`) covers nothing.
A `sealed` record covers no clean row unless its `aeeStillArmed` is
`true`, its `aeeDropCount` is zero or does not exceed an `aeeDropBound`
declared in the same signed payload, and its `aeePostureDigest` equals
both the arming record's and the pinned `networkPosture` digest, each a
check on signed carried bytes, so failing it is a coverage validity
failure, never a silent pass.

An `arming` record's payload MAY additionally carry three reserved members
that chain runs under the same substrate key: `aeeRunSeq` (a positive
safe-range integer), `aeePrevRunBinding` (the lowercase 64-hex run binding
digest of the predecessor run, absent exactly when `aeeRunSeq` is `1`),
and `aeeChainScope` (the population the sequence counts, declared as a
duplicate-free array of dimension tokens drawn from the closed vocabulary
registered below, sorted in the same canonical order as
`observationVocabulary.labels` (UTF-16 code-unit order, RFC 8785 section
3.2.3); REQUIRED whenever `aeeRunSeq` is present). The chain is always
structurally under one substrate key; each token names a further within-key
partition attribute already carried elsewhere in the attestation and fixes
where a consumer reads that attribute's value. The declared array is the
_dimension set_; the _evaluated tuple_ is the projection of the
substrate-key value and each declared token onto its registered attribute
value for this run (computed, never carried). The recommended minimum is
`["subject"]`; the empty array is the single global per-key counter that
makes every chain rule below vacuous and leaks the producer's total run
volume across its customers.

The `aeeChainScope` vocabulary is closed and each token pins a projection
to a value already carried on the wire: `subject` to
`subject[0].digest.sha256`, `corpus` to `observationEnvironment.corpus.digest`,
and `networkPosture` to `networkPosture.digest.sha256`. The substrate key is
the structural outer axis and is never a token. Values are not carried in
the member; a consumer projects each declared token onto its registered
field for this run. Minor versions MAY append tokens (each with a pinned
projection) and MUST NOT redefine an existing one; an unrecognized token
fails closed, as every closed vocabulary in this spec does.

Within one attestation these members are syntax-checked in the
reserved-member walk and nothing else normative reads them: the coverage
validity requirements, the `result` recompute, and the evidence tier are
unchanged. A violation of the syntax rules (a non-positive or non-integer
`aeeRunSeq`, a malformed `aeePrevRunBinding`, a missing `aeeChainScope`
when the sequence is present, a non-array `aeeChainScope`, an array carrying
a token outside the registered vocabulary, an array not in canonical order
(the same canonicality rule as `observationVocabulary.labels`: UTF-16
code-unit order, duplicate-free), or any of the three present without
`aeeRunSeq`) is handled as any reserved-member violation: the record
covers nothing. Their value is across attestations, as consumer policy over
whatever set a producer publishes. A consumer compares each attestation's
declared dimension set against the set its policy demands: an equal set is
admissible; a strictly finer set (a superset of dimensions) is
scope-narrowing, fragmenting every run into a singleton chain so no gap,
fork, or duplicate genesis can arise and the chain proves nothing; a
strictly coarser set (a subset of dimensions) pools distinct subjects, so a
withheld run of the demanded subject is deniable as a sibling's private run
and a sibling's run can occupy the withheld sequence position. A consumer
that has demanded a scope admits only the equal set, neither finer nor
coarser. Among admitted attestations the rules key on the evaluated tuple,
not the token set: a skipped `aeeRunSeq` under one tuple is a gap; two
attestations carrying the same `aeeRunSeq` under one tuple are a fork; two
carrying the same `aeePrevRunBinding` share a predecessor; and two genesis
records (absent `aeePrevRunBinding`) under one tuple are equivocation of the
same grade as a shared predecessor. Keying on the tuple is load-bearing:
genesis-per-subject-value is the normal case, and only a second genesis
under an identical tuple is a reset. A chain reset is not a fresh start. The
members claim ordering under the substrate key, nothing more: nothing on the
wire anchors when an arming record was signed relative to the run's outcome,
so commit-before-outcome holds only in combination with the publicly datable
run-entropy floor under Prerequisites or an external registration receipt. A
numeric gap is unexplained absence (crashed, private, and discarded runs all
produce gaps innocently), never fraud evidence in itself. Even a contiguous,
fork-free, correctly-scoped chain does not prove population completeness: a
producer may still mint a dense, gap-free set of passing runs after the
fact. Fork consistency among the published set is the ceiling of what any
self-contained attestation set can establish; the demand-disclosure yield is
that a consumer policy MAY require a contiguous, fork-free chain over the
runs offered to it. The external completion is a registration receipt:
committing each arming record to an append-only transparency log at run
start (for example SCITT, RFC 9943, with COSE receipts, RFC 9942) upgrades
gap-evidence to third-party-auditable non-omission, and that machinery is
deliberately outside this predicate. These members are semantic
(armed, stayed armed, nothing dropped, posture unchanged) rather than
mechanism-specific: how a substrate establishes them (a checkpoint chain,
a sequence counter, a hardware watchdog) stays producer territory. A
record whose `aeeKind` the consumer does not recognize covers nothing and
is otherwise ignored, while still contributing its leaf to `batchRoot`;
minor versions MAY add kinds and MUST NOT change the covering semantics of
an existing kind. An unrecognized kind can only weaken, never
strengthen, a row (candidate future kinds, informatively: a hardware-quote
kind binding the vantage to a measured platform, and a `registration` kind
carrying a transparency-service receipt over the arming record). The `aee` member prefix is reserved for future versions
(`aeeVersion` is reserved for a payload contract version); everything else
in the payload stays producer territory.

_Precedent (informative)._ Reserved members inside a producer-defined
signed payload follow an established lineage rather than a novel
mechanism: an RFC 7519 JWT claims set is a producer-defined object from
which verifiers read registered claim names (`exp`, `aud`, `iss`), with
collision resistance by registration and prefixing; EAT (RFC 9711) applies
the same registered-claims pattern inside an attestation token, in the
RATS family; OCI image annotations reserve the `org.opencontainers.*`
prefix inside an otherwise free-form map; and the `+json` requirement is
RFC 6839 Section 3.1's structured-syntax license to parse a media type not
otherwise known. None of the cited standards' semantics apply here by
reference; the citations locate the pattern, not the rules. Two deliberate
departures from that lineage: where RFC 7519 tells consumers to ignore
unrecognized claims, this predicate is fail-closed (a colliding or
unrecognized `aee*` member can only weaken coverage, the record covering
nothing, never create it), and the verify-then-read discipline is
normative here (a payload's fields mean nothing until its signature
verifies), which closes the parse-before-verify class of deployment
mistake the JWT lineage is known for. No record is required to name its
attack. Substrates sign at observation time, before attribution, and
attribution strength remains producer vocabulary. What an
`interception` record carries is a commitment to an intercepted payload
rather than the payload itself, keeping the attestation publishable rather
than a sensitive-data store.

`batchRoot` _string, required when `observationRecords` is non-empty_

An RFC 6962 Merkle root over the observation records, SHA-256, with
domain-separated hashing: each leaf is `H(0x00 || the record's DSSE PAE
bytes)`, each internal node is `H(0x01 || left || right)`, the tree built
by the RFC 6962 recursive split (never by duplicating a trailing node to
pad the leaf count), leaves in `observationRecords` array order, a
single-record tree's root its leaf hash, and an empty array with no root.
Two byte-identical entries in `observationRecords` make the attestation
invalid: a record's canonical identity is its leaf hash, a positional
reference is shorthand for the leaf hash at that position, and a future
minor version may admit detached records addressed by leaf hash with
`batchRoot` unchanged. Carried once at the predicate level. A `batchRoot`
that does not recompute over the carried records makes the attestation
invalid. Because every `basis: substrate` row requires covering records
under Coverage validity, any valid attestation carrying a substrate row
carries a `batchRoot`, and a clean run's committed set includes its
`arming` and `sealed` records. The root commits to the set the statement
carries and not to the set the run produced: it is recomputed from the
carried records, so a party who can re-sign the enclosing envelope drops a
record, recomputes over what remains, and emits a statement whose root is
self-consistent. What `batchRoot` detects is alteration of the carried set
by a party who cannot re-sign the envelope, and what it establishes for
every other party is the internal consistency of that set, never its
completeness against the run. `batchRoot` is omitted only when
`observationRecords` is absent, in which case every `basis: substrate` row
fails Coverage validity, so a valid recordless attestation carries only
`basis: artifact` rows.

`doesNotAssert` _array of strings, optional_

Explicit negative scope: statements the producer declares this evidence makes
no claim about (e.g. behavior outside the thrown corpus, host integrity
beyond the substrate's own attestation). Advisory: a verifier MUST NOT
require it, and nothing in it weakens the required checks. `doesNotAssert`
is the single canonical spelling: earlier internal versions used a
snake_case spelling (see the changelog), which is not accepted as an alias,
since two accepted spellings would mean two canonicalizations for the same
content. Migrating old producer output to the new name is a producer
concern that the wire format does not carry.

`issuedAt` _string (RFC 3339 timestamp), required_

When the producer signed the evidence bundle.

## Example

The schema block above is a complete statement whose `corpus.digest` is
re-derivable from the embedded manifest (`{"classes":{"CO":["CO-EXFIL-1"]}}`
canonicalized under RFC 8785 and hashed); the remaining digests are
placeholders.

### Consumer policy obligations

Two expectations are consumer policy, resolved outside the attestation and
never read from it: which keys count as substrate observation keys (the
evidence tier's input), and which corpus and substrate this consumer
expects. A consumer MUST pin, out of band, the corpus digest and the
substrate digest it expects for the deployment it is admitting into, and
at consumption MUST compare them against
`observationEnvironment.corpus.digest` and
`observationEnvironment.substrate.digest`; on mismatch the attestation is
not admitted, exactly as an attestation whose covering signatures do not
verify is not admitted. The comparison is deliberately not a validity
gate: validity is a function of carried bytes alone and holds identically
for every consumer, while the expected corpus and substrate differ per
consumer. An anchor-mismatched attestation is valid evidence about the
wrong context. Verification surfaces SHOULD expose one consumer-facing
admission result that conjoins validity, tier-policy satisfaction, and the
anchor comparison, so a result-only consumer cannot read a
valid-but-wrong-context attestation as admissible.

### Consumer policy example (non-normative)

Naming the substrate observation keys is policy, not wire format. A
minimal policy in the style of a witness layout, keyed on both the derived
tier and the row's `method` so a reconstructed clean row is not admitted
as a live one. The check that matters is signature verification against a
key pinned out of band; a record's `keyid` is an unauthenticated lookup
hint and never the check itself:

```rego
# Pinned out of band; never read from the predicate.
substrate_keys := {"substrate-2026": "<ed25519-public-key-bytes>"}

attested(row) {
  every i in row.observationRefs {
    rec := input.predicate.observationRecords[i]
    # verify DSSE PAE(payloadType, payload) against a pinned key;
    # rec.signatures[_].keyid selects WHICH pinned key to try, nothing more
    pae_verified(rec, substrate_keys)
  }
}

deny[msg] {
  row := input.predicate.attackResults[_]
  row.basis == "substrate"
  not attested(row)
  msg := sprintf("substrate row %s is unattested", [row.attackId])
}

deny[msg] {
  row := input.predicate.attackResults[_]
  admit_only_live
  row.method != "intercepted"
  msg := sprintf("row %s covered by reconstruction, not live interception", [row.attackId])
}
```

`attested` is coverage-of-existence, not temporal completeness; transient
tolerance travels on the `method` axis, so an admission rule that needs a
live observation keys on `method: intercepted` as well as the tier.

## Changelog and Migrations

A versioning discipline this predicate commits to: a member is born exactly
when a normative reader consumes it. In particular, if a future version
makes the shared-reference evidencing obligation checkable (for example by
committing per-attack expected artifacts in the corpus manifest),
attribution strength acquires a normative reader at that version and
becomes a required member then, not retroactively and not through a
verifier-invented heuristic in the meantime.

Versions 0.1–0.2 were internal producer iterations; 0.3 was the first shape
proposed for vetting. Relative to those internal versions, 0.3 removed all
verdict/policy semantics (moved downstream), moved intercepted-payload bytes
out in favor of commitments, moved `batchRoot` from per-record to
predicate-level, and adopted the I-JSON safe-integer profile on every rail.

0.4 incorporates review feedback on 0.3:

-   Added the required per-row `basis` field on `attackResults` (closed
    vocabulary `substrate_observed` / `artifact_reported` / `inferred`,
    fail-closed on unknown values), so each observation carries its own
    vantage and consumers can gate on it.
-   Pinned `actualLayer` clean-run behavior: rows with no containment event
    carry the literal `none` rather than omitting the field.
-   Renamed `does_not_assert` to `doesNotAssert` to match the lowerCamelCase
    convention. The rename is in place with no alias: the old spelling is
    rejected, keeping a single canonicalization per content.

0.5 revises 0.4 after review:

-   Split the per-row `basis` field into two orthogonal required fields:
    `basis` (closed vocabulary `substrate` / `artifact`) now names only the
    observation's vantage, defined by its weakest input with a stated
    artifact-sourcing criterion, and the new `method` (closed vocabulary
    `intercepted` / `reconstructed`) names its directness, with the same
    weakest-input composition rule. The 0.4 values `substrate_observed`,
    `artifact_reported`, and `inferred` are rejected, not aliased, under
    the same single-canonicalization rule as the `does_not_assert` rename;
    `inferred` has no successor because it conflated the two axes.
-   Made `actualLayer` required on every row (missing member: malformed
    statement, deliberately a different altitude than the fail-closed row
    members, per the stated design invariant) and extended its literal
    `none` to caught rows, where it states observed-but-not-enforced, so
    enforcement role never leaks into `basis`.
-   Added consumer strength orderings on both sides with a defined
    supporting set (`basis` bounds a `fail`, `method` bounds a `pass`,
    fail-closed rows at the lattice bottom, clean-row ordering extended to
    `degraded`), a coherence check of row claims against the pinned
    `observationEnvironment`, and a row-internal check that intercepted
    caught rows reference verifiable intercept records, all consumer-side.
-   Stated the row-travel design invariant under Parsing Rules and the
    producer-claim trust boundary for `basis`/`method`.

0.6 folds in the review of 0.5:

-   `basis: substrate` is now backed by substrate-signed coverage at two
    gates. Byte-checkable coverage (references resolve in range and
    class-match; every covering payload is canonical `+json` carrying the
    reserved members with `aeeRunBinding` equal to the derived run
    binding; `method` capped by the weakest signed `aeeMethod`;
    `batchRoot` recomputes) is a VALIDITY requirement and a consumption
    precondition. A consumer that consumes `result` or credits any row
    MUST evaluate it first, and a violation makes the attestation invalid,
    independent of any consumer. The one trust-relative step (the
    covering signatures verify against a consumer-named substrate key) is
    a per-row evidence tier (`attested` / `unattested` / `declared`); a
    consumer with no pinned substrate root treats every substrate row as
    `unattested`. An `unattested` substrate row ranks with `artifact` in
    both orderings: rank, never relabel; the MAY-reject-never-downgrade
    rule is retained verbatim.
-   Caught intercepted rows are covered by `interception` records,
    reconstructed rows by `examination` records, and clean intercepted
    rows by BOTH a run-level `arming` record and a `sealed` record whose
    signed payload reports the vantage stayed armed with a zero or
    self-bounded run-wide drop count and an unchanged posture digest. The
    strongest absence claim is bounded to the vantage's existence and
    continuity for the carried run, not to the absence of any event and
    not to a run population.
-   The observation vocabulary now travels in the attestation
    (`observationVocabulary`: labels, caught subset, JCS digest), so the
    recompute and the validity gate are pure functions of carried bytes
    and archived attestations stay verifiable without the producer's
    documentation.
-   Renamed `interceptRecords` to `observationRecords` and `interceptRefs`
    to `observationRefs` (old spellings rejected, no alias). Record
    signatures are DSSE PAE over `(payloadType, payload)`; coverage
    payloads MUST be canonical `+json`. `batchRoot` is pinned to RFC 6962
    with domain separation, duplicate records rejected, and is required
    whenever records exist. A new `runEntropy` digest in
    `observationEnvironment` folds a substrate-emitted run-start value
    into a versioned run binding (`aeeBindingVersion: 1`) so
    identical-configuration re-runs derive distinct bindings; the binding
    is anti-splice, not a freshness challenge, and identical-config replay
    is bounded by a stateful consumer rejecting `runEntropy` reuse.
-   Replaced the producer-claim trust-boundary paragraph with a field
    partition and an honest key model: the tier defeats substrate-free
    minting only where the observation key is held apart from the
    assembly plane; under a single trust root it defeats only a keyless
    downstream tamperer, and a key-holding operator who signs fiction
    stays outside the threat model. Coverage is only as trustworthy as the
    named key's un-compromised lifetime.
-   Stated the single normative read of `containmentObserved` (carried
    caught-set membership) and the criterion that an axis earns its own
    member only when a normative reader consumes it; attribution strength
    and tolerance remain non-normative. Pinned the recompute's
    independence from the validity gate and the tier. Stated the
    composition and run-population non-claims. Unknown `aeeKind` covers
    nothing and is otherwise ignored (fail-closed forward compatibility).
-   Completed the canonical-bytes profile with its string half: the
    vocabulary arrays are sorted ascending by UTF-16 code unit explicitly,
    and every signed canonical surface (covering payload member names,
    both vocabulary arrays) is BMP-only, rejected as malformed: within
    the BMP, code-unit and code-point order coincide, so conforming
    verifiers cannot split on sort order.
-   Numbered the byte-pure validity steps and separated them from the
    trust-relative stage; stated the consumer anchor obligations (pinned
    expected corpus and substrate digests, compared at consumption, with
    a single conjoined admission result recommended for verification
    surfaces) as consumer policy rather than a validity gate.
-   Recommended a publicly datable, round-unpredictable component in the
    run-entropy pre-image (proven signing-time floor; asserted ceiling
    unchanged), with the qualifications that make the floor real.
-   Added the optional `aeeRunSeq` / `aeePrevRunBinding` /
    `aeeChainScope` arming-payload members: cross-run gap evidence under
    a declared scope, ordering-only, with equivocation semantics for
    forks and duplicated geneses and a stated registration-receipt
    completion path.
-   Reclassified the shared-reference evidencing rule as a producer
    obligation outside every gate, and recorded the member-birth
    versioning discipline in this changelog. Documented the
    registered-claims lineage for reserved payload members as an
    informative note.
-   Retracted three claims that building the reference verifier and
    executing attacks against it disproved. `batchRoot` is recomputed from
    the carried records, so it establishes the internal consistency of the
    carried set and never its completeness: neither a dropped interception
    nor a dropped `arming` or `sealed` record changes a root recomputed
    over what remains, and the `method` cap is per record rather than per
    attack, so re-pointing a row's `observationRefs` at another attack's
    record inflates the row's `method` with every substrate signature still
    verifying. Each retracted sentence is replaced by a statement of the
    party the mechanism does bound, namely one who cannot re-sign the
    enclosing envelope. Coverage validity is now stated as a set of
    structural well-formedness constraints that become security properties
    only in combination with observation-record signature verification.
    Recorded two further limitations beside the coverage-bounded-observed
    one (the executed attack set is a producer assertion, and the four
    names outside the run binding are attacker-modifiable on a valid
    statement), and split the recompute goal into its recomputable
    reduction and its asserted construction. No normative requirement
    changed.

[Runtime Traces]: runtime-trace.md
[SCAI]: scai.md
[SVR]: svr.md
[Test Result]: test-result.md
[VSA]: vsa.md
