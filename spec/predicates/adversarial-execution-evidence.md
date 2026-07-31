# Predicate type: Adversarial Execution Evidence

Type URI: https://in-toto.io/attestation/adversarial-execution-evidence/v0.7

Version: 0.7.0

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
producer claiming less is not detectable from the statement, and that
sentence is exact rather than cautious: a statement that withdraws a claim is
a statement an honest producer with weaker instruments emits from the same
configuration, so no function of the carried bytes refuses the one without
refusing the other, and no quantity of additional signed material changes
that, because additional material is material a withdrawing producer also
declines to carry.

Between claiming more and claiming less, one rule decides where a commitment
can help and where it cannot: a commitment carried on a record binds
precisely the attacks that need that record, and none of the attacks that can
delete it. `batchRoot` is on the wrong side of that rule and it is worth
keeping the reason in sight. It is recomputed over the records the statement
carries, so it can never detect a missing member; what it binds is the
carried set against a party who cannot re-sign the enclosing envelope, which
is a network attacker rather than the assembly plane the substrate key
separation is written against. The run-end `sealed` record is on the right
side of it, because a party deleting an interception cannot also delete the
seal and still present a statement carrying a `basis: substrate` row. The
run-start `arming` record is on the right side of it for coverage inflation,
because inflation fabricates rows that must point at run-level records the
inflated statement therefore has to keep. Completeness of the record set
against the run is still nowhere proven: what the run-end commitment below
establishes is that the carried set is the emitted set, never that the
emitted set is everything that happened.

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
canonicalization of the object `{"aeeBindingVersion": "2", "catchPolicy":
"<catchPolicy.digest.sha256>", "corpus": "<corpus.digest.sha256>",
"networkPosture": "<the lowercase 64-hex SHA-256 of the RFC 8785
canonicalization of the carried networkPosture object>",
"observationVocabulary": "<observationVocabulary.digest.sha256>",
"runEntropy": "<runEntropy.digest.sha256>", "subject":
"<subject[0].digest.sha256>", "substrate": "<substrate.digest.sha256>"}`.
Every input is a property of the run's configuration and is fixed before
corpus injection. That is not incidental and it is the test any proposed
input must pass: the arming record carries this digest inside its own
signature and is signed before injection, so a value the producer could not
know at that moment would make the arming record unsignable, and an outcome
of the run can therefore never be an input here. `runEntropy` is a run-start
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
`catchPolicy`, `corpus`, `runEntropy`, `substrate` and `subject[0]` MUST each
carry a `sha256` digest whose value is already lowercase 64-hex, and so MUST
`networkPosture`, whose pinned digest this construction no longer reads
verbatim but which is still compared byte for byte against the
`aeePostureDigest` a record carries; a substrate-row-carrying statement
violating this digest requirement is malformed. The
`observationVocabulary` digest carries no rule of its own here, because the
digest-integrity step recomputes it from the arrays beside it and a value
that is not lowercase 64-hex cannot equal that recompute; restating the
requirement would add a check that could never be the one to fail. Values
are taken verbatim (no case-folding, no
null fill). A statement whose rows are all `basis: artifact` derives no
binding and need not carry `runEntropy`. A verifier derives the digest from
the statement alone; no field carries it. Every substrate-signed
observation record commits to the run by carrying this digest inside its
signed payload (see the reserved members under `observationRecords`). The
binding is anti-splice: a record signed under a different subject, corpus,
catch policy, network posture, observation vocabulary, substrate, or
run-start entropy value cannot
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
absent member defaults to the version this document defines; the carried
value never drives the
derivation (a verifier derives only under the version it implements, so a
record declaring the implemented version but constructed otherwise still
fails on the digest). Defaulting the absent member to the implemented
version rather than to any fixed number is what keeps the member optional:
a default pinned to a superseded version would reject every statement that
simply declines to carry it. A future
minor version admitting multiple subjects or multiple substrates binds all
of them in canonical name-then-digest order.

Two inputs distinguish version 2 from version 1, both of them configuration
the statement already carries, so neither costs a byte on the wire and
neither needs a new comparison: each closes through the equality every
record's `aeeRunBinding` is already put to.

Version 1's `networkPosture` input was the value of that member's own
`digest.sha256`. That left the `posture` string beside it outside every
signature. The posture configuration this predicate digests is not carried
anywhere in the statement, so nothing can check the string against the
digest, and a party holding only the envelope key could replace one posture
value with another, change no digest, and break no signature. Version 2
takes the canonical digest of the carried `networkPosture` object instead,
so the string, its pinned digest, and any further member a producer carries
there all sit inside the binding. The object the binding covers is the
carried one: a producer that adds, removes or edits a `networkPosture`
member after the arming record is signed derives a binding its own records
do not carry, and its statement is invalid on that ground.

`observationVocabulary` was not an input at all. Its `caught` array decides
which labels are caught, and the coverage validity requirements and the
`result` recompute both read it, so a producer that narrows the caught set
after the run turns a caught row into a clean one. Nothing resisted that:
the vocabulary's own digest is verified only against the arrays beside it,
so it re-derives for free, and no record's binding moved. Binding the
carried digest closes it, since a narrowed vocabulary derives a different
run binding and every record then fails the comparison.

Version 2 is unchanged in 0.7 and no version 3 is defined. Every commitment
0.7 adds travels either on a record payload, where a substrate signature
already covers it, or inside the corpus manifest, whose digest is already an
input here. The manifest gaining `expectedPayloads` therefore changes the
value of the `corpus` input on every statement carrying one and changes
nothing about how that input is built, so a producer editing
`expectedPayloads` after the arming record is signed derives a binding its
own records do not carry, by the same mechanism that already covers the class
map beside it.

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
  "predicateType": "https://in-toto.io/attestation/adversarial-execution-evidence/v0.7",
  "predicate": {
    "result": "fail",
    "observationEnvironment": {
      "substrate": { "name": "<substrate-attestation-subject>", "digest": { "sha256": "<64-hex>" } },
      "corpus": {
        "name": "<corpus-name>",
        "uri": "pkg:<producer>/<corpus>@<version>",
        "digest": { "sha256": "<64-hex-JCS-digest-of-manifest>" },
        "manifest": {
          "classes": { "CO": ["CO-EXFIL-1"] },
          "expectedPayloads": { "CO-EXFIL-1": ["<64-hex>"] }
        }
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
        "attribution": "pinned",
        "actualLayer": "policy.egress_sinkhole",
        "observationRefs": [0]
      }
    ],
    "observationRecords": [
      // aeeKind: interception. Covers the caught row above and carries
      // aeePayloadCommitment.
      {
        "payload": "<base64(canonical +json bytes the substrate signed)>",
        "payloadType": "<producer-defined media type ending in +json>",
        "signatures": [ { "keyid": "<hex>", "sig": "<base64>" } ]
      },
      // aeeKind: arming. Carries aeeAssessedAttacks.
      {
        "payload": "<base64(canonical +json bytes the substrate signed)>",
        "payloadType": "<producer-defined media type ending in +json>",
        "signatures": [ { "keyid": "<hex>", "sig": "<base64>" } ]
      },
      // aeeKind: sealed. Carries aeeObservedSet and aeeObservedAttacks, and is
      // required on any statement carrying a basis: substrate row whether or
      // not a row references it.
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

One of `fail`, `degraded`, `pass_indirect`, `pass` (lowercase), ordered
`fail` < `degraded` < `pass_indirect` < `pass`. Defined as a total,
deterministic, severity-independent function of the predicate, evaluated as
the minimum under that order of three independent conditions rather than as
a cascade, because worst-wins rather than evaluation order is the rule. The
first condition holds when any `attackResults` row carries a
containment-observed label from the carried caught set
(`observationVocabulary.caught`), a label outside the carried
`observationVocabulary.labels` (fail-closed), or a missing or
out-of-vocabulary `basis`, `method` or `attribution` (fail-closed, same
rule), and it contributes `fail`. The second holds when `coverage.outOfScope` or
`coverage.routedElsewhere` is non-empty, and it contributes `degraded`. The
third holds when any clean row, meaning a row whose `containmentObserved`
is in the carried labels and not in the carried caught set, carries a
`basis` other than `substrate` or a `method` other than `intercepted`, and
it contributes `pass_indirect`. A condition that does not hold contributes
`pass`. A `pass` is coverage-bounded-observed: it states what was
assessed and makes no general safety claim. There is intentionally no
severity threshold, policy ruleset, or free-text reason here; policy belongs
downstream.

`pass_indirect` is a coverage-complete result at least one of whose clean
rows rests on an observation that was indirect in vantage (`basis:
artifact`, the executed artifact's own account of itself) or indirect in
time (`method: reconstructed`, derived after the event rather than at it).
`pass` and `pass_indirect` make the same coverage claim and different
observation claims, and the ordering says only that the second is never the
stronger of the two. The distinction exists because the top result was
otherwise reachable by a statement disclaiming substrate observation
altogether. A party holding the enclosing envelope key but not the
substrate's observation key can move every row to `basis: artifact` and
then drop the observation records, the batch root and the run entropy that
those rows no longer require, and the statement it presents is well formed,
carries no substrate evidence at all, and reads at the top of the ordering.
That statement is byte-identical to one an honest producer with no
substrate vantage emits from the same configuration, so no function of the
carried bytes refuses the first without refusing the second, and refusing
both would remove the producer whose attack classes have no substrate
vantage to observe from. What the fourth value does instead is price both
below a live interception, which is the only distinction the carried bytes
support.

`pass_indirect` is not a revival of the retired 0.4 `inferred` value.
`inferred` was a row-level value whose conflation destroyed, at its only
carrier, which axis was weak; `pass_indirect` is a statement-level
reduction over `basis` and `method`, both of which remain required and
individually readable on every row, so nothing a consumer needs is
available only through the reduction. A consumer that needs to separate
indirectness of vantage from indirectness of time MUST read the two row
members and never the result token.

`pass_indirect` says nothing about signature verification, and the third
condition is deliberately not phrased over the evidence tier. A `pass` may
still rest on clean rows deriving `unattested`, which the clean-row
ordering ranks with `artifact`; that half of the weakness is key-relative,
belongs to the evidence tier rather than to the recompute, and a byte-pure
function cannot see it. Neither the token nor the tier substitutes for the
other, and a consumer crediting any `basis: substrate` row MUST derive the
tier whichever of the two top results the statement carries.

The default admission threshold is `result == "pass"`. A consumer MAY
accept `pass_indirect`, and a consumer relaxing its threshold below `pass`
MUST additionally key on each clean row's `basis` and `method` and on that
row's derived evidence tier, because below `pass` the ordinal stops
distinguishing them: a `degraded` reached through a disclosed coverage gap
and a `degraded` whose clean rows are all `artifact` carry the same token.

`attribution` enters the recompute through the fail-closed arm of the first
condition and nowhere else. A row declaring `paired` is not a weaker result,
it is a weaker binding between the row and the records that cover it, and the
two are different questions: the recompute asks what was observed, and
`attribution` asks how firmly this row is the row that observation belongs to.
Pricing `paired` in the result would charge an honest producer for a layer
whose committed value no corpus can predict, which the member's own definition
states is a permanent condition of some layers rather than a gap in this one.
A consumer that cares about the binding reads the member; the token never
carries it, exactly as it never carries the evidence tier.

Two further bounds sit beside that one. The substrate observes what crosses
the vantages it was armed at, and this document requires nothing of it beyond
that. Some deployments run a substrate that also dispatches the corpus, and
such a substrate holds, on its own clock, which probe it issued and which of
its own records fell inside that probe's window; `aeeObservedAttacks` exists
so that a substrate holding that correspondence can sign it instead of handing
it to the assembly plane unsigned. Neither shape reaches the bound stated
here. A substrate that never saw the runner cannot verify that the runner
executed every attack the manifest names, and a substrate that dispatched the
corpus itself still says nothing about an attack that produced no observation,
because a lower bound over what was observed is silent about what was not. The
set of attacks actually executed therefore remains a producer assertion under
both shapes. Coverage integrity compares the carried rows against the carried
manifest, which catches an omission the producer failed to declare in
`coverage`, but neither comparison reaches the run, so a consumer reading a
`pass` is relying on producer integrity that nothing was silently skipped.

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
of the carried statement, and that is what makes them runnable by a consumer
holding no keys. It is also their limit: a violation here is conclusive, since
no signature can rescue a statement that does not hang together, but the
absence of a violation concludes nothing on its own. These requirements
establish that the statement is well formed, never that it is true:

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

Five further coverage validity requirements hold on the statement, or on
every row rather than only on a `basis: substrate` row. Each is a function of
carried bytes on the same terms as the list above, and a violation of any of
them makes the attestation invalid:

-   a clean row resolves no `observationRefs` index to an `interception`
    record. A row stating that nothing was caught while pointing at a record
    in which the substrate signed that it intercepted traffic states both
    halves of a contradiction, and the check is a membership test that reads
    no signature and no key. It is stated over every row because the
    contradiction does not depend on the vantage the row declares;
-   every carried record whose payload `aeeKind` is `interception` is
    resolved by at least one `observationRefs` index on a caught row. An
    interception the statement carries and no caught row accounts for is an
    observation the substrate signed and the producer then reported nothing
    about. One record MAY be resolved by more than one row, so this costs
    none of the sharing this document already permits;
-   a statement carrying at least one `basis: substrate` row carries at least
    one `sealed` record that satisfies every constraint of its kind and whose
    `aeeRunBinding` equals the derived run binding, whether or not any row
    resolves an index to it. Before 0.7 a `sealed` record was required only
    to cover a clean intercepted row, so a statement whose rows were all
    caught carried none, and the run-end commitment defined under
    `observationRecords` had nowhere to live on exactly the statements a
    record deletion works against. A rule conditioned on the presence of the
    record it constrains is a rule a producer switches off by omission;
-   `aeeObservedSet` on every carried `sealed` record equals the value
    recomputed over the carried records, by the construction stated at that
    member. A seal committing to a record set the statement does not carry
    does not hang together, in the same sense as a `batchRoot` that does not
    recompute;
-   a row declaring `attribution: pinned` resolves at least one
    `observationRefs` index to an `interception` record, its `attackId`
    carries an entry in `corpus.manifest.expectedPayloads`, and every
    `interception` record it resolves carries in its `aeePayloadCommitment`
    at least one value from that entry. A row whose `attackId` carries no
    such entry MUST declare `paired`. The first part is not redundant beside
    the third: a requirement universally quantified over an empty set is
    vacuously true, so without it a producer deletes the interception
    records, relabels the row, resolves only run-level records, and still
    declares the stronger of the two values with nothing checking it.

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
costs a secret: `aeeRunBinding` is derived entirely from material the
statement already carries in plaintext, `aeePostureDigest` is the pinned
`networkPosture` digest carried beside them, and an `arming` or `sealed`
payload describing a vantage that never existed satisfies every constraint
here. These requirements are therefore structural well-formedness
constraints. They become security properties only in combination with
observation-record signature verification against a substrate key the
consumer trusts, which is the evidence tier below. A verifier that evaluates
coverage validity and skips that verification has checked that the producer
filled the form in correctly.

**What the 0.7 commitments do not close.** Each of the four members 0.7 adds
is either a bound on inflation or a conversion of an invisible edit into a
declaration a consumer can read. None is a detector, and the four stop in
different places, so the limits are stated one at a time rather than as a
single caveat.

`aeeObservedSet` binds the deletion of a record from the carried set, and only
while the statement still has to carry the seal. It does not reach a producer
that declines to run an attack, declines to report a row, or withdraws every
`basis: substrate` row: once no row declares that basis, `runEntropy`,
`observationRecords`, `batchRoot` and the `sealed` record become optional
together, and what remains is a statement an honest producer with no substrate
vantage emits from the same configuration. It inherits, without curing,
whatever the substrate did not record, because a commitment to the emitted set
is a claim about what the substrate emitted and never a claim that nothing else
happened; an evasion at a boundary that produces no record is faithfully absent
from a verified commitment. And it is a structural constraint like every other
requirement here: against a party that can sign records it is exactly as
forgeable as the rest of the payload.

`aeeObservedAttacks` binds the deletion and the relabelling in one member,
because deleting every interception record does not delete the seal's claim
that an attack was observed. Three things bound it in turn. It is a lower bound
by construction: a seal naming an attack obliges a caught row for that attack,
and a seal omitting one licenses nothing, so an observation the substrate could
not attribute subtracts from the claim rather than adding a false one. It is
available only to a substrate that holds the correspondence, and a substrate
that does not carries the empty array honestly, so a producer withdraws the
whole control by presenting a substrate that does not dispatch the corpus. What
the member changes there is visibility rather than reachability: the withdrawal
is an empty array on the wire instead of an absence nothing records, and a
consumer that demands a non-empty set refuses a statement rather than failing
to notice one. And it names no attack on the record that observed it, so it
cannot cap a `method` per attack and cannot tell two observations of one attack
apart.

`aeeAssessedAttacks` binds coverage inflation and nothing else. It cannot
bind coverage suppression, and that is a property of the target rather than a
weakness of the rule: a producer that withdraws its rows and discloses the
classes emits a statement byte-identical to the one an honest producer emits
after a run that could not assess them, so no rule over the carried statement
refuses the first without refusing the second. The rule that would close
suppression inside the format conditions the commitment's presence on the
disclosure itself, and its price is the artifact-only producer that has no
substrate and is honest about a class it did not cover, the fully-skipped run
this document deliberately keeps valid, and every run that loses coverage
part-way; that price is not paid here. The comparison is a subset rather than
an equality, which is what keeps the part-way loss expressible, and what the
subset bounds is the producer's own run-start declaration: a producer that
declares the whole manifest at run start has committed to the largest set the
manifest permits and is bounded there by coverage integrity alone. What the
member removes is the freedom to raise that declaration after the outcome is
known. It is also computed over the carried manifest, so it is blind to a
corpus whose manifest never declared the class a consumer wanted assessed.

`expectedPayloads`, `aeePayloadCommitment` and `attribution` bind the
permutation of the row-to-record assignment, and only on the layers where the
committed value is predictable from the corpus. They do not bind the deletion,
because their commitments ride the records a deletion removes; that ground
belongs to the run-end commitment and not to these. They are unavailable,
rather than merely absent, wherever a substrate commits to a value no corpus
author can compute in advance: a commitment taken under a per-run key, one over
bytes the substrate re-serialized rather than the bytes the artifact emitted,
one over bytes a scrubbing or capping policy altered, or one over a raw frame
prefix carrying values a kernel assigned. On those layers `paired` is the true
answer, and a universal `pinned` requirement would oblige a producer to state a
falsehood or oblige a substrate to give up the property that made its
commitment unpredictable. `paired` is therefore itself a free withdrawal: a
producer declaring it on every row satisfies every rule here, because a
statement claiming a weaker binding is a statement an honest producer with a
weaker binding emits. The closure against that producer is the consumer
obligation stated under Consumer policy obligations, and it is the closure
rather than a hook into one, in the same place and for the same stated reason
as the corpus and substrate pins.

One limit is common to all four and is stronger than any of them. A property
that compares this statement against another statement of the same run is
unreachable by any rule over one statement, whatever that statement carries: a
consumer holds one attestation, and a rule asking whether a label was different
before has no second statement to read. That is the run-population non-claim
this document already keeps loudest, wearing different clothes.

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
the embedded `manifest`, and the `manifest` itself, which carries a required
`classes`, a map from assessment class to the complete array of attack
identifiers it defines, and an optional `expectedPayloads` defined below; an
attackId MUST NOT appear under more than one class), `catchPolicy` (JCS digest of the
parsed catch-policy document, so an empty or permissive policy is
distinguishable from an enforcing one), `networkPosture` (the
substrate-authoritative egress posture, drawn from the closed vocabulary
registered below, with its configuration digest), and
`observationVocabulary` (the
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

`corpus.manifest.expectedPayloads` is an optional map from attack identifier
to the duplicate-free array of lowercase 64-hex commitment values a substrate
is expected to carry when it observes that attack. Every key MUST be an attack
identifier the same manifest's `classes` declares, every array MUST be
non-empty, sorted ascending by UTF-16 code unit and duplicate-free, and every
entry MUST be lowercase 64-hex; a manifest violating any of these is
malformed. The map is a sibling of `classes` rather than an enrichment of it,
because `classes` is the in-scope partition that the manifest floor and
coverage integrity read, and what an attack looks like on the wire can be
added, corrected or extended without repartitioning the corpus. It inherits
the manifest's digest pinning with no new mechanism: it sits inside the
pre-image `corpus.digest` is taken over, which is a run binding input, so it
is pinned out of band by the same consumer obligation that pins the classes
beside it.

The value an entry carries is the commitment as the substrate computes it, and
not a digest of the input the corpus author wrote down. Where a substrate
canonicalizes before committing, a path it normalizes or a host name it maps
to its A-label form, the corpus carries the canonical form; a corpus carrying
the raw form silently fails to match and every row it should have supported
falls back to `paired` while looking correct. This is a producer obligation
that no verifier can check, because a verifier holding the corpus and the
statement cannot tell an honest absence of an expectation from a mispredicted
one, and it is stated here because the failure is silent in the direction that
weakens the claim rather than in the direction that invalidates it.

`expectedPayloads` discharges the versioning commitment recorded in the
changelog: it is the per-attack expected artifact that example named, so
attribution strength acquires a normative reader at this version and becomes a
required row member here. It is not retroactive, and no verifier reading a
statement of an earlier version may invent the reader.

The `networkPosture.posture` vocabulary is closed. Four values are
registered: `allowlist`, egress permitted only to a declared destination
set; `no_network`, no egress path exists; `sinkhole`, egress is accepted and
diverted to a capture endpoint rather than reaching its destination; and
`unsafe_bypass_egress`, egress is unrestricted and uninstrumented. A
statement whose `posture` is absent, is not a string, or carries a value
outside that set is malformed, fail-closed, on the same terms as every other
closed vocabulary here; a minor version MAY append a value and MUST NOT
redefine a registered one.

The set is closed rather than illustrative for a reason that is not
housekeeping. A consumer is invited under `basis` to coherence-check a row
against the posture the run was contained under, and a substrate row
claiming a network-boundary observation under a posture that provides no
interception path at that boundary is the case that invites it. That check
cannot be written against a value whose meaning is undeclared: no verifier
can decide whether an unregistered posture provides an interception path, so
an open vocabulary leaves the check permanently unreachable while appearing
to offer it. Closing the set also settles a divergence this document was
carrying on its own: the proto beside it already described this vocabulary
as closed and fail-closed and already carried the fourth value, so an
implementer reading the two together had to pick one, and every shipped
implementation picked the closed reading.

A typing discipline this predicate commits to, stated here so that a later
reader inherits it rather than rediscovers the question. Each of these six
members is descriptor-shaped, and two of them are [ResourceDescriptor]s.
`substrate` and `catchPolicy` identify a resource and carry nothing beside
that identity, so they take the framework type; the `sha256` digest required
on each is a requirement the descriptor specification permits a context using
the type to impose, and a producer MAY additionally carry `uri`,
`downloadLocation` or `mediaType` there, none of which any rule in this
document reads. Reading a pinned `sha256` off a descriptor is already what
this predicate does in its most load-bearing place, since `subject` entries
are [ResourceDescriptor]s by the Statement specification and the run binding
reads `subject[0].digest.sha256`.

The other four members stay locally typed, and the reasons are stated here
rather than left to inference. Where a member carries the pre-image its digest
is taken over, that pre-image stays on the statement's own JSON surface:
`corpus` carries the `manifest` and `observationVocabulary` carries `labels`
and `caught`, and the only descriptor member that could hold either is
`content`, whose value is base64. Every byte-level rule Prerequisites imposes
is stated over the statement's JSON, namely the duplicate-member rule at any
depth, the well-formed-scalar-value requirement applied to the raw bytes
before any decoded string is read, the nesting bound of 128, and the BMP
restriction on canonical surfaces. Material inside a base64 member is outside
all four, so carrying a digest pre-image there would open a second
canonicalization boundary inside a signed statement, in a predicate whose
whole encoding profile exists so that two conforming verifiers cannot
disagree about identical bytes. Where a member instead carries further
normative material beside an identity, this predicate keeps the member
locally typed rather than extending a descriptor with members of its own,
which is the shape [Runtime Traces] uses for `monitor`. `runEntropy` is
offered as a reading rather than as a rule: its digest commits to a
substrate-emitted run-start value rather than describing a resource, so a
descriptor is the wrong vessel for it.

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
`attribution` (required; how firmly the row is bound to the records that
cover it, see below),
`actualLayer` (required; which enforcement layer acted, see below), and
`observationRefs`
(indexes into `observationRecords` binding this row to the observation
records that cover it). An `interception` index MAY be referenced by more
than one row. A producer MUST NOT reference a record from a row whose
attack the record's committed payload does not evidence. On a row declaring
`attribution: pinned` that obligation is checkable and is checked, by the
coverage validity requirement stated above: the corpus declares what the
attack's interception commits to and the verifier compares. On a row
declaring `paired` it remains an obligation outside every gate, because no
validity requirement, recompute input or tier evaluation reads it there, and
a conforming verifier neither can nor may invent an evidencing heuristic in
its place. The line between the two is exactly the line the corpus draws by
carrying an expectation or not.
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

That comparison is only as strong as the manifest it reads against, so the
manifest carries a floor of its own: it MUST declare at least one attack
identifier across all of its classes, and a manifest declaring zero attack
identifiers makes the statement malformed. The requirement is phrased over
attack identifiers rather than over classes because a manifest whose classes
object is empty and a manifest carrying a named class with an empty array
declare the same thing, nothing to execute, and only counting identifiers
closes both; the second is the more plausible of the two, since it reads as
a real assessment class. Without the floor a zero-attack manifest satisfies
coverage integrity vacuously, comparing an empty union against an empty
union, and the rest of the statement follows from there: zero rows means
zero `basis: substrate` rows, and with no substrate row this document
permits `runEntropy`, `observationRecords` and `batchRoot` to be absent, so
every structure that would have required a substrate signature drops out and
a valid `pass` about an arbitrary subject can be minted with no substrate
participation at all. A corpus declaring no adversarial inputs is not an
adversarial corpus, which is why this sits with well-formedness rather than
with `result`: scoring it would concede that a zero-attack run is a
legitimate statement that merely scores badly. The floor bounds the manifest
and nothing beyond it. A manifest that declares attack identifiers and
assesses none of them stays valid, and the honest fully-skipped run, which
discloses its classes under `outOfScope` and scores `degraded`, is untouched
by this requirement, because its manifest declares an attack identifier.

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
reader consumes it, which is why `basis` and `method` are members: the
recompute and the gate read them.

`attribution` is the axis that acquires such a reader at this version. It
states how firmly the row is bound to the records that cover it, over a closed
two-value vocabulary:

-   `pinned`: the row rests on at least one interception whose committed value
    the corpus declared in advance, so a consumer holding the pinned corpus
    can check that the record this row cites is a record of this attack.
-   `paired`: the row rests on a correspondence established some other way,
    typically by the window in which the observation fell. The binding is a
    producer assertion.

Until 0.7 attribution strength was non-normative producer vocabulary, for a
stated reason: no normative reader consumed it. `expectedPayloads` makes it
checkable, so the member is born here rather than inherited, and the two
values are exactly the two readings that reason named, a hash-pinned payload
and a time-window pairing. Attribution _tolerance_ on clean rows, including
window bleed between same-layer siblings, is untouched by this and remains
non-normative producer vocabulary documented beside the observation
vocabulary's definition; a consumer MUST NOT move a `result` or an evidence
tier on it.

`pinned` is a claim about the binding and never about the observation. It does
not make the row's `method` stronger, does not raise its evidence tier, and
does not enter the `result` recompute except through the fail-closed arm every
required row member with a closed vocabulary shares. `paired` is a claim a
producer may always truthfully make, including where the corpus offers an
expectation, so it is a floor rather than a confession: what a consumer learns
from it is that this row does not carry the stronger binding, never that the
producer had one and withheld it.

`basis`, `method` and `attribution` are REQUIRED on every row and all three
vocabularies are closed: a
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
states the ordering-only scope and the registration-receipt completion. A
result resting on any `reconstructed` clean row SHOULD be read as
tolerating transients between the observed states, and a result resting on
any `artifact` clean row is self-reported absence, the weakest; both of
those statements recompute to `pass_indirect` rather than to `pass`, so the
first two ranks of this ordering are the two ranks the result token already
separates from a live interception, and a consumer reading only the token
still cannot tell them apart from each other. A `pass` resting on an
`unattested` substrate clean row is self-reported absence too, and is the
one rank of this ordering the recompute cannot express, because whether a
covering signature verifies is key-relative and the recompute is not. The
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
key. It does not defeat method inflation, and it defeats a record drop only
where the run-end commitment reaches. Both limits are structural rather than
gaps in the checks above.
`batchRoot` recomputes over the carried records, so a party holding the
envelope key removes a record and recomputes a root that is self-consistent
over what remains; what refuses that party is `aeeObservedSet` on the seal,
which is signed by a party that does not control the carried set, and it
refuses only while the statement still has to carry the seal. The `method`
cap binds a row to the weakest `aeeMethod`
across the records that row references, and no per-event record names its
attack, since
a substrate signs at observation time and before attribution; the cap is
therefore per record and never per attack, so re-pointing a row's
`observationRefs` at an `intercepted` record signed for a different attack
raises that row's `method` with every substrate signature still verifying.
`attribution: pinned` narrows that re-pointing wherever the corpus declares
an expectation for the row's attack, because the borrowed record must then
also carry that attack's committed value, and narrows nothing where it does
not.
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
field. Coverage is therefore only as trustworthy as the named key's
un-compromised lifetime; the single trust root is a single point of total
failure, and a consumer's policy MAY bound a named key with a validity
window. Where it does, that window MUST be evaluated against a
substrate-signed instant, and the one this predicate mandates is the
`armedAt` carried inside an `arming` record whose signature verifies under
the key being bounded; it MUST NOT be evaluated against `issuedAt`.
`issuedAt` is producer-asserted, sits outside every substrate signature,
and is not among the run binding digest's inputs, so a party holding the
envelope key moves it at will, changing no digest and breaking no
signature. The only constraint this document places on it is that
`armedAt` is no later than `issuedAt`, which every instant at or after the
run satisfies, and `armedAt` itself precedes the compromise of any key that
signed that run. A window evaluated against `issuedAt` therefore
rehabilitates, by back-dating alone, every record the revoked key ever
signed. A statement carrying no `arming` record that verifies under the
bounded key carries no substrate-signed instant for that key, so a consumer
bounding a key's validity refuses that statement rather than falling back
to `issuedAt`. Consumers MAY additionally coherence-check row claims against
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
and `signatures`, which MUST carry at least one entry. A consumer verifies
each record's signature, DSSE PAE over `(payloadType, payload)`, before
relying on any field inside the payload. The order is deliberate and the
wording is exact: the byte-pure gates below do read payload fields without
verifying anything, because they are structural and a consumer must be able
to run them with no key material at all, and what they produce is not a
finding a consumer may act on until the covering signatures have verified.
Reading ahead of verification is a stage, never a conclusion. Any record used
to cover a `basis: substrate` row MUST carry a
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
    rows; payload MUST carry `aeePayloadCommitment`); `arming` (run-level: a
    live, cooperation-independent capture
    vantage was armed for the run before corpus injection; payload MUST
    carry `armedAt` under the timestamp profile `issuedAt` defines, no
    later than `issuedAt`, `aeePostureDigest` equal to the pinned
    `networkPosture` digest, and `aeeAssessedAttacks`, and its
    `aeeMethod` MUST be `intercepted`);
    `sealed` (run-level: the vantage stayed armed to run-end; payload MUST
    carry `aeeStillArmed`, a boolean; `aeeDropCount`, an integer counting
    run-wide dropped observations; `aeePostureDigest`, the effective
    posture at run-end; `aeeObservedSet`; and `aeeObservedAttacks`; it MAY
    carry `aeeDropBound`, a producer-declared integer bound; its `aeeMethod`
    MUST be `intercepted`); or
    `examination` (the substrate examined artifact-independent state after
    the fact; its `aeeMethod` MUST be `reconstructed`, and its payload
    SHOULD identify the states compared).
-   `aeeMethod` _string_: `intercepted` or `reconstructed`; how the
    substrate observed, stated inside the signature.

A record violating any constraint of its declared `aeeKind` (including a
missing `armedAt` on an `arming` record, an `armedAt` after `issuedAt`, an
`armedAt` outside the timestamp profile, or an `examination` record signed
`aeeMethod: intercepted`) covers nothing.
A `sealed` record covers no clean row unless its `aeeStillArmed` is
`true`, its `aeeDropCount` is zero or does not exceed an `aeeDropBound`
declared in the same signed payload, and its `aeePostureDigest` equals
both the arming record's and the pinned `networkPosture` digest, each a
check on signed carried bytes, so failing it is a coverage validity
failure, never a silent pass.

Four reserved members carry the commitments this version adds. Each is
required on exactly one kind, each is a value the substrate holds at the
moment it signs that kind, and a record missing or malforming the member its
kind requires covers nothing, on the same terms as a missing `armedAt`.

`aeePayloadCommitment` _array of strings_, required on an `interception`
record: the commitment values this interception carries, duplicate-free,
sorted ascending by UTF-16 code unit (RFC 8785 Section 3.2.3, the
canonicality rule the vocabulary arrays already carry), non-empty, every
entry lowercase 64-hex. The member names what an `interception` record has
always carried and had no reserved spelling for. What a commitment is taken
over is the substrate's choice and this document does not constrain it,
because the choice is what keeps an attestation publishable rather than a
store of the traffic it observed; what a corpus can declare in advance is
bounded by that same choice, which is the subject of `expectedPayloads`. It
is an array rather than a single value because one record may be resolved by
more than one row and may commit to more than one observed value.

`aeeAssessedAttacks` _array of strings_, required on an `arming` record: the
attack identifiers this run declared, before corpus injection, that it would
assess. Duplicate-free, sorted ascending by UTF-16 code unit, every entry an
attack identifier the carried `corpus.manifest.classes` declares; a record
violating any of these covers nothing. The union of the manifest's
identifiers for the carried `coverage.assessedClasses` MUST be a subset of
this array, and a statement violating that is invalid.

The comparison is a subset rather than an equality, and the choice is
load-bearing in both directions. A subset refuses the claim to have assessed
more than was declared before any outcome was known, which is the whole of
what a run-start commitment can bind. An equality would additionally refuse
the honest run that declared two classes, lost one part-way, and disclosed the
loss, which is the shape the coverage maps exist to reward; and it would buy,
against the withdrawal it appears to catch, only the version of that
withdrawal that leaves the arming record in place, since a producer with no
`basis: substrate` row left may drop the record and the commitment with it.
A rule that costs an honest producer a legitimate shape in exchange for
catching the incomplete form of an attack is not a bargain.

The array is carried inside the signature rather than as a digest over a
pre-image on the statement's own JSON surface, which is the shape `corpus`
and `observationVocabulary` take. The two cases differ in who authors the
value. Those pre-images are producer material a verifier must read, so the
surface is the right home and the digest is the binding. This value is the
substrate's own assertion about a moment the assembly plane cannot revisit,
and splitting it would let the party under scrutiny author the array the
substrate is committing to. Carrying identifiers rather than a digest also
admits no value the substrate was not already committing to, since every
admissible entry appears in the carried manifest and the manifest's digest is
already a run binding input inside the same signature.

`aeeObservedSet` _string_, required on a `sealed` record: the lowercase
64-hex SHA-256 of the RFC 8785 canonicalization of the duplicate-free array,
sorted ascending by UTF-16 code unit, of the lowercase 64-hex leaf hashes of
every `interception` and `examination` record the substrate emitted for this
run, where a leaf hash is `H(0x00 || the record's DSSE PAE bytes)`, the same
leaf construction `batchRoot` uses. A verifier recomputes the same value over
the carried records of those two kinds and requires equality; a `sealed`
record whose `aeeObservedSet` does not equal that recompute covers nothing,
and because a `sealed` record is required on every statement carrying a
`basis: substrate` row, such a statement is invalid.

The seal is the only substrate signature made after the run is over, and
until this member it carried three facts about the vantage and none about the
observations. What the member adds is a commitment, by a party that does not
control the carried set, to the set that was emitted. A dropped record removes
a leaf and the values diverge. A record fabricated to keep a count intact
cannot be signed without the substrate key, a record borrowed from another run
fails the run binding comparison, and a record duplicated to the same end is
already invalid, because two byte-identical entries in `observationRecords`
make the attestation invalid. It commits to nothing the substrate would have
to interpret: no attack identifier, no outcome, no label, only a hash over
bytes the substrate itself produced.

Two producer obligations travel with it and neither is checkable by a
verifier. The seal is signed after every record it commits to, so a substrate
that examines state after run end seals after that examination rather than
before it; a seal that predates a record the statement carries is a seal its
own record set contradicts, and the statement is invalid on the recompute
without the verifier ever learning why. And every site at which the substrate
drops an observation rather than emitting it MUST increment `aeeDropCount`.
The two members are siblings: one commits to what the substrate recorded and
the other to what it failed to record, and a substrate that silently discards
an observation without counting it emits a wire claim that is wrong while
looking right, since the seal then commits to a set that is complete by its
own account and short by the run's.

`aeeObservedAttacks` _array of strings_, required on a `sealed` record: the
attack identifiers to which this run attributed at least one of its own
observations. Duplicate-free, sorted ascending by UTF-16 code unit, every
entry an attack identifier the carried `corpus.manifest.classes` declares; a
record violating any of these covers nothing. For every identifier in the
array the statement MUST carry an `attackResults` row with that `attackId`
whose `containmentObserved` is in the carried caught set, and a statement
violating that is invalid.

The array is a lower bound and the rule reads in one direction only. A seal
naming an attack obliges a caught row for that attack; a seal omitting one
licenses nothing, and in particular does not oblige a clean row. That is what
makes the member sound without requiring the substrate to resolve every
ambiguous case: an observation it could not attribute is left out, which
subtracts from what the seal claims and can never add a claim that is false.
Over-inclusion is the direction that would be unsound, and a substrate that
attributes by disjoint dispatch windows cannot produce it.

The empty array is the honest value, and it is required rather than
omissible. A substrate that does not dispatch the corpus holds no
correspondence to sign and says so by carrying nothing in the array; a
substrate that dispatched the corpus and attributed no observation says the
same thing and means something different. Allowing the member to be absent
instead would make the whole control escapable by omission, which is the
defect the mandatory `sealed` record above exists to close and would be
reintroduced one level down.

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
mistake the JWT lineage is known for. No per-event record is required to
name its attack, and none may be: a substrate signs an `interception` or an
`examination` at observation time, before attribution, so a record naming an
attack would be signing the runner's account of what it was doing rather than
its own account of what it saw. The run-level `sealed` record names attacks
because it is signed at run end, by which time the substrate that dispatched
the corpus holds the correspondence between its own probes and its own
records; that is its assertion and not the runner's, and it is why the
attack-naming member sits on that kind alone. `aeePayloadCommitment` keeps
the same discipline on the per-event side: an `interception` record carries a
commitment to an intercepted payload rather than the payload itself, which
keeps the attestation publishable rather than a sensitive-data store, and the
comparison that turns the commitment into an attribution happens in the
verifier, against a corpus pinned out of band, rather than in the substrate.

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

`issuedAt` _Timestamp, required_

When the producer signed the evidence bundle. `Timestamp` is the framework's
field type, which requires RFC 3339 in the UTC timezone, and this document
pins the two choices that type leaves open. A statement is canonicalized
and digested as its bytes, so no verifier may normalize the field before
reading it and the admissible set has to be written down; left open, one
implementation is quietly stricter than another and the divergence surfaces
only when a statement crosses between them. The date-time separator and the
zone designator MUST be uppercase, never the lowercase `t` and `z` that
RFC 3339 also admits, and the zone designator MUST be `Z`, `+00:00` or
`-00:00`, never a non-zero offset such as `+05:00`. `-00:00` is admitted
rather than excluded because RFC 3339 Section 4.3 gives it the meaning that
the instant in UTC is known while the offset to local time is not, which
describes where the producer stood and not when it signed, and the instant
is the only thing this document reads from the field. A statement whose
`issuedAt` is absent, is not RFC 3339, or is RFC 3339 outside this profile
is malformed. `armedAt` carries the same profile, defined here and cited
from the arming record so that the two fields cannot drift apart.

## Example

The schema block above is a complete statement whose `corpus.digest` is
re-derivable from the embedded manifest, which is the whole `manifest` object
including `expectedPayloads` and not the `classes` map alone, canonicalized
under RFC 8785 and hashed; the remaining digests are placeholders.

### Consumer policy obligations

Four expectations are consumer policy, resolved outside the attestation and
never read from it: which keys count as substrate observation keys (the
evidence tier's input), which corpus and substrate this consumer
expects, which assessment classes it demands a run to have assessed, and
what it requires of the row-to-record binding. The first two are stated
first because they are the older pair; the second two are the closures for
two attacks no rule over the carried statement can reach, and they are stated
here rather than left as guidance for that reason.

A consumer MUST pin, out of band, the corpus digest and the
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

A consumer MUST pin, out of band, the set of assessment classes it requires a
run to have assessed, and at consumption MUST compare that set against
`coverage.assessedClasses`; a statement disclosing a demanded class under
`outOfScope` or `routedElsewhere` is not admitted, exactly as an
anchor-mismatched statement is not admitted. This is the only obligation in
this section whose value a consumer must derive from what it wants rather than
from what a producer published. The corpus anchor pins bytes and says nothing
about what those bytes must contain, so a consumer that pinned a digest it
copied out of a producer's bundle has pinned whatever that producer chose to
ship. The demand is stated over classes rather than over attack identifiers
because identifiers are corpus-version scoped: a consumer pinning them re-pins
on every corpus revision, and the list it re-pins to is one it read out of the
producer's own manifest. A pin whose value comes from the party being checked
is not a demand.

The obligation exists because withdrawn coverage is invisible in the carried
bytes and provably so, and because the consumer holds the only fact that
decides it, which is what it asked for. It therefore never has to tell an
honest skipped run from a suppressed one: it refuses both, on the ground that
it demanded the class and the class was not assessed, and the producer's
intent stops being the question. A consumer that demands no class MUST record
that decision explicitly rather than reach it by omission, and MUST NOT fold
the decision into the corpus and substrate pins, which decline something else:
declining those concedes evidence about any corpus while leaving every
self-consistency requirement in this document standing, and declining this one
concedes that a producer may withdraw any class it likes, against which
nothing in this document stands at all. Absence here is not a degraded
control; it is the absence of one.

Finally, a consumer decides what it requires of the row-to-record binding and
of the run-end attack set. A consumer MAY require `attribution: pinned` on
every row whose attack the pinned corpus carries an `expectedPayloads` entry
for, and MAY require a non-empty `aeeObservedAttacks` on the seal. Each is a
refusal of a declared weakness rather than the detection of a hidden one: a
producer that declares `paired` throughout, or that presents a substrate
holding no correspondence to sign, violates no requirement in this document,
and the statement it emits is one an honest producer in the same position
emits. What the format does is put the weakness on the wire where a policy can
read it. A consumer that declines either requirement is admitting a statement
whose row-to-record assignment rests on the producer's word, and the conjoined
admission result recommended above SHOULD say so rather than report a bare
pass.

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
live observation keys on `method: intercepted` as well as the tier. The
second rule above is the row-level half of the partition the `result`
recompute now also reduces to a token: a policy gating on `result ==
"pass"` already excludes every statement that rule would deny, and a policy
relaxed to admit `pass_indirect` MUST keep the rule, because the token
states that some clean row is indirect and never which one.

## Changelog and Migrations

A versioning discipline this predicate commits to: a member is born exactly
when a normative reader consumes it. In particular, if a future version
makes the shared-reference evidencing obligation checkable (for example by
committing per-attack expected artifacts in the corpus manifest),
attribution strength acquires a normative reader at that version and
becomes a required member then, not retroactively and not through a
verifier-invented heuristic in the meantime.

That commitment came due. 0.7 carries `expectedPayloads` in the corpus
manifest, which is the example the sentence named, so attribution strength is
a required row member from 0.7 and from no earlier version, spelled
`attribution`. The clause about retroactivity is what decided that 0.7 is a
new predicate type rather than a revision of 0.6: a rule that attaches to one
version and not to its predecessor needs the wire to say which version a
reader is holding, and `predicateType` is the only member that says so.

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
    into a versioned run binding so
    identical-configuration re-runs derive distinct bindings; the binding
    is anti-splice, not a freshness challenge, and identical-config replay
    is bounded by a stateful consumer rejecting `runEntropy` reuse.
-   Moved the run binding to `aeeBindingVersion: 2`. The pre-image gains
    `observationVocabulary`, the carried vocabulary digest, so that
    narrowing the caught set after the run breaks every record's binding
    rather than re-deriving for free; and its `networkPosture` input
    becomes the canonical digest of the carried `networkPosture` object
    rather than the value of that object's own `digest` member, which
    brings the posture string inside the signature it was sitting beside.
    Both inputs are configuration already on the wire, so the change costs
    no bytes and adds no comparison. Version 1 is retired with no alias and
    no dual-accept window: a statement built under it derives a digest no
    record carries, and a record declaring version 1 explicitly covers
    nothing. The absent-member default is now stated as the implemented
    version rather than as a fixed number, so that omitting the optional
    declaration stays legal across a version change.
-   Closed the `networkPosture.posture` vocabulary at four registered
    values and made an unregistered one malformed, resolving a divergence
    in which this document introduced the values as an example while the
    proto beside it and every shipped implementation treated them as a
    closed, fail-closed set.
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
-   Required the corpus manifest to declare at least one attack identifier
    across all of its classes; a manifest declaring none makes the statement
    malformed. Coverage integrity otherwise passes vacuously on an empty
    union, zero rows carry no `basis: substrate` row, and the statement then
    legally omits `runEntropy`, `observationRecords` and `batchRoot`, which
    admitted a valid `pass` about an arbitrary subject with no substrate
    participation at all. The
    requirement is stated over attack identifiers rather than over classes so
    that a named class with an empty array is closed alongside an empty
    classes object, and it leaves the honest fully-skipped run (attack
    identifiers declared, every class disclosed under `outOfScope`, scoring
    `degraded`) valid.
-   Typed `issuedAt` as the framework's `Timestamp` rather than as a
    lowercase RFC 3339 timestamp, which is what the protobuf schema already
    did, and stated on the field the timestamp profile the type leaves open:
    uppercase designators, and a zone designator of `Z`, `+00:00` or
    `-00:00`. The zone rule was previously written only on `armedAt`, so a
    statement whose `issuedAt` carried `+05:00` was conformant while being
    off-guideline, and the case rule was written nowhere, which had already
    split two independently written verifiers on the same bytes. `armedAt`
    now cites the profile instead of restating half of it.
-   Typed `observationEnvironment.substrate` and
    `observationEnvironment.catchPolicy` as the framework's
    `ResourceDescriptor`, the type the sibling predicates already import, and
    stated the rule the other four members of that object are held under. The
    JSON member names and the wire shape are unchanged, so no signed byte, no
    digest, no signature and no conformance vector moves; in the protobuf
    schema two locally declared messages become that import. The rule is that
    a member carrying the pre-image its digest is taken over keeps that
    pre-image on the statement's own JSON surface, because the only descriptor
    member that could hold it is base64 `content`, and material inside a
    base64 member is outside every byte-level rule Prerequisites states.
-   Corrected the key-validity window recommendation, which named `issuedAt`
    as the operand a consumer checks a named key's validity against. That
    field is producer-asserted, sits outside every substrate signature, and
    is not among the run binding digest's inputs, and the one rule this
    document states about it survives moving it later, so the window was
    recommended on the single temporal value the party that key separation
    exists to constrain writes at will: back-dating it rehabilitates every
    record a since-revoked key ever signed, at no signature and no digest.
    The operand is now the `armedAt` inside an `arming` record that verifies
    under the bounded key, and a statement carrying no such record is refused
    rather than falling back. The operand is stated normatively rather than
    left to the consumer because an admission policy written independently of
    this sentence had already bounded evidence age against the same field and
    believed it bounded; one implementer choosing the defeated operand is a
    mistake, and two choosing it separately is a property of how the field
    reads. No wire byte, digest, signature or conformance vector moves: the
    correction is to a consumer obligation in stage two, which the byte-pure
    validity gate does not reach.
-   Added a fourth `result` value, `pass_indirect`, ordered between
    `degraded` and `pass`, and restated the recompute as the minimum of
    three independent conditions rather than as a cascade. The added
    condition holds when a clean row carries a `basis` other than
    `substrate` or a `method` other than `intercepted`. The top result was
    otherwise reachable by a statement carrying no substrate evidence at
    all: a party holding the enclosing envelope key alone moves every row
    to `basis: artifact`, drops the records, the batch root and the run
    entropy that a substrate row would have required, and lands above the
    run it downgraded. That mutant is byte-identical to the statement an
    honest producer with no substrate vantage emits, measured over every
    finding-bearing vector in the conformance suite, so the two are not
    separable by any function of the carried bytes and the value prices
    both rather than refusing either. The condition reads only required row
    members with closed vocabularies, so the recompute stays byte-pure, and
    it is deliberately phrased over declared vantage and directness rather
    than over the evidence tier, so the tier's independence from `result`
    survives unchanged and the one weakness the tier owns, an `unattested`
    substrate clean row, keeps the top token as it always did. Four accept
    vectors and one reject vector move their expected result; no wire
    member, digest, signature or record moves.

0.7 makes checkable four things 0.6 could only describe. It is breaking on
the wire, on the corpus and on every published conformance record:

-   The predicate type is `.../v0.7` and 0.6 is retired with no alias and no
    dual-accept window, under the same single-canonicalization rule that
    retired the 0.4 basis values and the `does_not_assert` spelling. The bump
    is neither housekeeping nor free. Three members become REQUIRED and a
    record that was conditionally required becomes unconditional, so a
    statement valid under 0.6 can be malformed under 0.7. The versioning
    discipline above says a member becomes required at the version that gives
    it a normative reader and never retroactively, and that rule is
    unimplementable unless the wire says which version a reader is holding.
    `aeeBindingVersion` cannot say it, because it scopes the run binding
    construction alone and none of these members is a binding input. A minor
    version was not available either: this document licenses a minor version
    to append a posture value, a chain dimension or a record kind without
    changing the covering semantics of an existing kind, and every change
    below changes the covering semantics of `arming` and `sealed`. So the
    discipline decides the bump, and the review calendar does not.
-   A `sealed` record is now required on every statement carrying a
    `basis: substrate` row rather than only where a clean intercepted row
    needs covering, and it carries `aeeObservedSet`, a commitment to the leaf
    hashes of every `interception` and `examination` record the substrate
    emitted. Until this, the only substrate signature made after the run was
    over carried three facts about the vantage and none about the
    observations, and a party holding the enclosing envelope key could delete
    an inconvenient interception, recompute a root over what remained, and
    emit a statement that satisfied every requirement in the document. The
    requirement is unconditional because a rule conditioned on the presence of
    the record it constrains is a rule a producer switches off by omission,
    and the statements it was switched off on were exactly the statements the
    deletion works against.
-   Two structural requirements join coverage validity and need no signed
    data at all: a clean row may resolve no index to an `interception`
    record, and every carried `interception` record is resolved by at least
    one caught row. The first refuses a relabelled row that still points at
    the record its caught form cited; the second refuses the escalation of
    dropping the reference instead of the record. Neither reaches a producer
    that drops the record itself, which is what the run-end commitment is for,
    and the second is worth nothing against a producer that also withdraws
    every substrate row, which is stated where it is defined rather than
    discovered later.
-   `aeeObservedAttacks` on the seal names the attacks the run attributed at
    least one of its own observations to, and obliges a caught row for each.
    It is the one member here that survives the deletion of every interception
    record, because the seal's claim does not travel on the records the
    deletion removes. It is a lower bound in one direction by construction, it
    is available only to a substrate that dispatched the corpus, and the empty
    array is required rather than omissible so that a substrate holding no
    correspondence declares that on the wire instead of leaving an absence
    nothing records.
-   `aeeAssessedAttacks` on the arming record names the attacks the run
    declared, before injection, that it would assess, and the assessed set
    carried at run end must be a subset of it. It binds coverage inflation,
    which is the withdrawal's mirror image and the only half of that pair a
    commitment can reach: inflation must keep the run-level records its
    fabricated rows point at, and withdrawal need keep nothing at all. The
    comparison is deliberately a subset and not an equality, so that a run
    which loses coverage part-way can still disclose the loss.
-   `corpus.manifest.expectedPayloads`, `aeePayloadCommitment` on
    `interception` records, and the required row member `attribution` over the
    closed vocabulary `pinned` / `paired` together make the shared-reference
    evidencing obligation checkable, which is the change the versioning
    discipline above pre-authorised and which is what makes attribution
    strength a member at this version. A row declaring `pinned` must resolve
    at least one interception and every interception it resolves must carry a
    value the corpus declared for that attack; the quantifier is paired with
    an existence requirement because a universally quantified rule over an
    empty set is vacuously true, and without the pairing a producer could
    delete the records and keep the stronger label.
-   What none of the four closes is stated beside the coverage validity
    requirements, one member at a time, because the four stop in different
    places and a single caveat would flatten them. Two of the closures in this
    version are consumer obligations rather than validity rules, and they are
    written as obligations under Consumer policy obligations for the same
    stated reason the corpus and substrate pins are: the deciding fact is one
    the consumer holds and the statement cannot carry.
-   The run binding is unchanged and stays at version 2. Every commitment here
    travels either on a record payload or inside the corpus manifest, whose
    digest is already an input, so the construction does not move; the value
    of the `corpus` input moves on every statement carrying
    `expectedPayloads`, which is the binding working rather than changing.
-   Every published conformance record breaks, including any built by an
    independent implementer, and that cost is named rather than absorbed: an
    independent re-run is the strongest external evidence this document has
    that its text is determinate, and a breaking version spends it. The
    predicate is pre-adoption, so the price is payable once and never again.

[ResourceDescriptor]: ../v1/resource_descriptor.md
[Runtime Traces]: runtime-trace.md
[SCAI]: scai.md
[SVR]: svr.md
[Test Result]: test-result.md
[VSA]: vsa.md
