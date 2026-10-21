# Verifying an attestation against conformance vectors

[Validation model](validation.md) states, as pseudocode, what a consumer checks when
it reads one attestation. This page states how a verifier proves it does that, and
how two verifiers of one predicate find out whether they agree.

The difference matters because the pseudocode has no executable form. Two
implementations can each satisfy every line of it and still accept different bytes,
and neither one can discover that from this repository.

## The contract

A verifier is conforming for a predicate when, for every vector in that predicate's
manifest under [`conformance/`](../conformance/README.md), the verdict it reaches
equals the verdict the manifest declares, compared on the surface the manifest names.

Three properties of that sentence are the whole design.

**The expectation travels with the vector.** A manifest entry declares
`expected.verdict`, so a verifier needs no out-of-band knowledge of what a case is
for. A test fixture whose expected result lives in a test file is a fixture for one
implementation; a fixture that declares its own expectation is a fixture for all of
them.

**A refusal declares its codes.** Verdict-only comparison scores a verifier that
refuses the right statement for the wrong reason as agreeing with one that refuses it
for the right one. Those two are not interoperable, and the difference is invisible
until somebody debugs a rejection in production. So an entry of kind `reject` carries
the codes its refusal must name, and a verifier that refuses for a different reason
is reported as diverging rather than as passing.

**The comparison surface is declared, not assumed.** A manifest names which fields of
a verdict are compared, so one verifier's private diagnostics do not count against it
and two verifiers are not held to each other's internal structure.

## Running one

A verifier reads the manifest, runs every statement it names, and reports per vector.
Nothing about the invocation is predicate-specific:

```shell
<verifier> --manifest conformance/<predicate>/manifest.json
```

The exit status is the campaign's: zero when every vector reached its declared
verdict, non-zero otherwise, with the diverging vector ids named. A verifier that
cannot parse a vector reports that vector as an error rather than as a refusal,
because "this statement is invalid" and "I could not read this statement" are
different findings and only the first one is a verdict.

## A worked example

The reference verifier for the adversarial execution evidence predicate reads the
manifest form above and reports per vector against the declared verdict and, on a
refusal, against the declared codes. Its corpus for that predicate is 272 vectors at
its own suite revision 28, of which the four in
[`conformance/adversarial-execution-evidence/`](../conformance/adversarial-execution-evidence/)
are a sample chosen to exercise both verdicts and two different refusal codes.

Two independent implementations of that predicate, neither written by its author, have
run against that corpus and published their results: a Rust verifier with its own
I-JSON parser, RFC 8785 serializer, RFC 6962 Merkle root and Ed25519 tier, and a
Python one written to a pre-registered evidence order. Both reach 272 of 272. That
number is the argument for this page: it exists because there was a corpus to run,
and no other predicate in this repository has one.

## What conformance vectors do not establish

A verifier that matches every declared verdict is interoperable on the cases the
corpus covers. It is not correct, and the corpus is not a specification.

**A reading no vector exercises is untested rather than confirmed.** Where a
document underdetermines a reading and no vector forces it, two implementations
agreeing on that reading is evidence the text is determinate and not proof of it: a
third implementation can take the other reading in silence. The honest response is to
add the vector, not to record the agreement.

**Coverage is a property of the corpus, not of the verifier.** A corpus that
exercises none of a predicate's refusal paths certifies nothing about refusals, and a
verifier can pass it while refusing nothing at all. A manifest's `conditions` exist
so that gap is readable: a rule no vector's conditions name is a rule nothing checks.

**A vector proves a verifier reads bytes, never that a producer told the truth.** The
statements here are inputs. Whether the party that signed one observed what it
claims is a question about signatures and vantage, and no conformance corpus can
reach it.
