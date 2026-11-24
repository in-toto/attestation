# Predicate type: AMEP Verdict

Type URI: https://github.com/stafordtituss/amep/verdict/v0.3

Version: 0.3

Predicate Name: AMEP Verdict

Authors: Staford Titus Selvaraj (@stafordtituss)

## Purpose

An AMEP Verdict records what a verifier concluded when it checked an AI
agent's actions against the authority the agent was given.

[AMEP] (Agentic Mandate Evidence Protocol) is a draft format for the signed
evidence of an agent session: the mandate a principal granted, any delegations
and approvals, each tool invocation, and the receipts tool servers signed for
what they executed. A verifier evaluates that evidence packet and returns, for
each invocation, one of three outcomes:

-   `preserved`: the action stayed within its authority;
-   `violated`: the evidence shows it went outside;
-   `unresolved`: the evidence is not enough to decide either way.

Each result names the deciding invariant, the party at fault when something
was violated (`blame`), and the missing evidence when nothing could be decided
(`gaps`).

The predicate is reproducible. It names the digest of the exact packet bytes
and of the key set that were evaluated, and two conforming verifiers produce
the same predicate for the same inputs. A consumer that holds the packet can
re-run any conforming verifier instead of trusting the signer; one that does
not can rely on the signer as with any attestation.

## Use Cases

### An agent changed a repository

A coding agent pushes a commit or opens a pull request using tools over the
Model Context Protocol. Its session evidence is an AMEP packet. A CI job or an
independent auditor evaluates the packet and publishes an AMEP Verdict about
it, alongside the commit:

-   `preserved`: the change was made by actions within the agent's mandate;
-   `violated`, invariant `I2`, blame `agent`: it wrote outside the resources
    the mandate allows.

An admission policy accepts the change only with a `preserved` verdict from a
trusted verifier. This asks of an agent the question supply-chain security
asks of a build: was this change made within authority?

### Payments and other irreversible actions

An agent pays an invoice. The verdict states whether the payment was within
the mandate's per-call limit and, if the mandate required an approval, whether
one was present and bound to this exact payment. The signed verdict is
archived with the payment record.

### Disputes and audits

Months later, a party asks whether an incident was the agent's fault, the
delegating principal's, the host's, the transport's or the tool server's. The
archived verdict gives the answer reached at the time and names the evidence
it was reached from; re-running a verifier on the archived packet reproduces
it.

### Detection

A security team verifies agent sessions after the fact. An `unresolved`
verdict with its `gaps` (for example, a missing receipt for invocation `i-7`)
tells them what evidence to collect, rather than reporting a failure.

### Why existing predicates do not cover this

-   [SLSA Verification Summary] and [Simple Verification Result] record that an
    artifact passed verification: a pass/fail decision and the passing
    properties. They do not express:
    -   a third outcome, `unresolved`. Missing evidence must read neither as a
        pass nor as a failure: a missing receipt is a gap, not a violation;
    -   results per action, since one session holds many actions with
        different outcomes;
    -   who is at fault and what evidence is missing, each from a fixed
        vocabulary;
    -   the inputs needed to reproduce the result. SVR explicitly does not
        include them; an AMEP Verdict is defined so that it can be reproduced.
-   [Test Result] records pass, warn or fail for a test run, with the same
    limits.
-   [SCAI Report] carries attribute assertions with evidence. An AMEP Verdict
    could be phrased as SCAI attributes, but the outcome, blame and gap
    vocabularies, and the determinism rule that makes the result reproducible,
    would live outside the schema.

The predicates compose. Both AMEP reference verifiers also publish the same
verdict as an SVR about the same subject, for consumers that only evaluate
SVRs:

-   `properties` is `["AMEP_VERDICT_PRESERVED"]` when the packet is
    `preserved`, and empty otherwise;
-   `verifier.policies` names the key set, by the same digest as
    `keyset.sha256`, and the AMEP version whose rules were applied.

The AMEP Verdict carries the detail needed when the outcome is anything else.

## Prerequisites

-   The in-toto Attestation Framework and its [Statement] layer.
-   AMEP v0.3.0-draft.1, Sections 10.5, 10.5.1 and 10.8 of the [AMEP]
    specification: the evaluation procedure, the outcome, invariant, blame and
    gap vocabularies, and verdict records.

## Model

-   **Subject:** the evidence packet, identified by the SHA-256 of the exact
    bytes the verifier evaluated. The name is conventionally `amep-packet`.
    Digesting the bytes, rather than a parsed form, means a verdict also covers
    a packet that does not parse, which evaluates to `unresolved`.
-   **Producer:** an AMEP verifier, such as a CI job, an auditor, an agent
    gateway or a monitoring service, signing with its own key.
-   **Consumer:** an admission or deployment policy, an archive, or a party
    re-checking the verdict.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ "name": "amep-packet", "digest": { "sha256": "<HEX>" } }],

  // Predicate:
  "predicateType": "https://github.com/stafordtituss/amep/verdict/v0.3",
  "predicate": {
    "packet": { "sha256": "<HEX>" },
    "keyset": { "sha256": "<HEX>" },
    "verdict": {
      "outcome": "<OUTCOME>",
      "invariant": "<INVARIANT>",
      "blame": "<PARTY>",
      "gaps": [{ "need": "<NEED>", "ref": "<ID>", "cause": "<CAUSE>" }, ...]
    },
    "results": [
      {
        "id": "<INVOCATION_ID>",
        "outcome": "<OUTCOME>",
        "invariant": "<INVARIANT>",
        "blame": "<PARTY>",
        "gaps": [{ "need": "<NEED>", "ref": "<ID>", "cause": "<CAUSE>" }, ...]
      },
      ...
    ],
    "verifier": { "name": "<STRING>", "version": "<STRING>" }
  }
}
```

### Parsing Rules

This predicate follows the framework's [standard parsing rules], including the
equivalence of unset, null and empty values: an unset `blame` means no party is
blamed, and unset `gaps` means none.

In addition:

-   Producers MUST serialize the statement in the JSON Canonicalization Scheme
    ([RFC 8785]) form. Consumers MUST reject a payload with duplicate member
    names.
-   **Determinism:** with `verifier` removed, two conforming verifiers given
    the same packet bytes and key set MUST produce byte-identical predicates.
    Free-text reasons differ between implementations and are therefore not
    part of the predicate.
-   `blame` is set only when `outcome` is `violated`; `gaps` is non-empty only
    when `outcome` is `unresolved`.
-   Consumers MUST ignore unrecognized fields. Adding an outcome, party or need
    value, or changing the meaning of a field, requires a new predicate
    version.

### Fields

**`packet.sha256`, required** string (hex)

> SHA-256 of the exact packet bytes evaluated. MUST equal the subject's
> `sha256` digest.

**`keyset.sha256`, required** string (hex)

> SHA-256 of the JCS form of the key set the packet was verified against:
> each key id's public key and signer roles.

**`verdict`, required** object

> The packet verdict: the first `violated` result, else the first `unresolved`
> one, else `preserved`. A packet with no results is `unresolved`.

**`verdict.outcome`, `results[].outcome`, required** string

> One of `preserved`, `violated`, `unresolved`.

**`verdict.invariant`, `results[].invariant`, required** string

> The deciding AMEP invariant, `I1` to `I11`, or `-` when the outcome is
> `preserved`.

**`verdict.blame`, `results[].blame`, optional** string

> On a `violated` outcome, the party whose own signed record first departs
> from authority: `agent`, `delegator`, `host`, `transport` or `server`.

**`verdict.gaps`, `results[].gaps`, optional** array of objects

> On an `unresolved` outcome, the evidence that would let evaluation continue.
> `need` is from the fixed set in AMEP Section 10.5.1 (for example
> `receipt`, `revocation_status`, `well_formed_packet`); `ref` is the id
> of the object it concerns. A gap for a refused signature also carries
> `cause`: `missing`, `unknown_key`, `role` or `invalid`. Other gaps omit it.

**`results`, required** array of objects

> One entry per result, in evaluation order. `id` is the invocation id, or `-`
> for a result about the packet as a whole.

**`verifier`, optional** object

> The implementation that produced the predicate (`name`, `version`). It is
> the only field that conforming verifiers may fill differently.

## Example

The statement both reference verifiers produce for the AMEP conformance vector
`must-reject/payment-over-the-per-call-limit-without-approval`: an agent made a
payment above the mandate's per-call limit without the approval the mandate
requires.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "amep-packet",
      "digest": {
        "sha256": "b0b46388661ff3e4dd396d31c0172bc3c82f123e109592552d4dfff596c0d5d0"
      }
    }
  ],
  "predicateType": "https://github.com/stafordtituss/amep/verdict/v0.3",
  "predicate": {
    "packet": {
      "sha256": "b0b46388661ff3e4dd396d31c0172bc3c82f123e109592552d4dfff596c0d5d0"
    },
    "keyset": {
      "sha256": "150056c1465cdc2bdb8a14c3df7f8c272d0ef43df378481f0d574c396bbf9c54"
    },
    "verdict": {
      "outcome": "violated",
      "invariant": "I2",
      "blame": "agent",
      "gaps": []
    },
    "results": [
      {
        "id": "i-1",
        "outcome": "violated",
        "invariant": "I2",
        "blame": "agent",
        "gaps": []
      }
    ],
    "verifier": {
      "name": "amep-protocol",
      "version": "0.3.0a4"
    }
  }
}
```

Implementations: the Python package `amep-protocol` (module
`amep_protocol.verdicts`) and the npm package `@stafordtituss/amep` (export
`@stafordtituss/amep/verdict`). Both sign statements in [DSSE] envelopes, and
a differential test checks that they produce byte-identical statements for all
282 AMEP conformance vectors (three suites). Both are by the same author.

Outside use: the Verifiable Agent Trust Envelope project keeps signed verdict
envelopes for its lost-response and recovery case under `reference/amep/`. Its
maintainer re-ran the pinned verifier and reported that all five verdicts
matched.

## Changelog and Migrations

Not applicable: this is the first version.

[AMEP]: https://github.com/stafordtituss/amep
[DSSE]: https://github.com/secure-systems-lab/dsse
[RFC 8785]: https://www.rfc-editor.org/rfc/rfc8785
[SCAI Report]: scai.md
[SLSA Verification Summary]: vsa.md
[Simple Verification Result]: svr.md
[Statement]: ../v1/statement.md
[Test Result]: test-result.md
[standard parsing rules]: ../v1/README.md#parsing-rules
