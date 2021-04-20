# Attestation Spec

An in-toto **attestation** is authenticated metadata about one or more software
artifacts, as per the [SLSA Attestation Model]. It has three layers that are
independent but designed to work together:

<!-- BEGIN: When updating below, also update ../README.md#specification -->

*   [Envelope]: Handles authentication and serialization.
*   [Statement]: Binds the attestation to a particular subject and unambiguously
    identifies the types of the predicate.
*   [Predicate]: Contains arbitrary metadata about the subject, with a
    type-specific schema. This repo defines the following predicate types,
    though custom types are allowed:
    *   [Provenance]: To describe the origins of a software artifact.
    *   [Link]: For migration from [in-toto 0.9].
    *   [SPDX]: A Software Package Data Exchange document.

The [processing model] provides pseudocode showing how these layers fit
together.

<!-- END -->

See the [top-level README](../README.md) for background and examples.

## Envelope

```jsonc
{
  "payloadType": "https://in-toto.io/Statement/v1-json",
  "payload": "...",
  "signatures": [<SIGNATURES>]
}
```

The Envelope is defined in
[signing-spec](https://github.com/secure-systems-lab/signing-spec) and adopted
by in-toto in [ITE-5].

## Statement

```jsonc
{
  "subject": [
    {
      "name": "<NAME>",
      "digest": {"<ALGORITHM>": "<HEX_VALUE>"}
    },
    ...
  ],
  "predicateType": "...",
  "predicate": {
    <PREDICATE>
  }
}
```

*   Type URI: https://in-toto.io/Statement/v1-json (value of `payloadType` in
    [Envelope])
*   Encoding: [JSON](https://www.json.org)
*   Schema: [statement.proto](spec/statement.proto)

The `subject` describes the set of software artifacts that the attestation
applies to. Each entry has a `name` and at least one `digest`.

The subject `digest`, of type [DigestSet], is a collection of alternate content
hashes of a single artifact. Two DigestSets are considered matching if ANY of
the fields match. The producer and consumer must agree on acceptable algorithms.
If there are no overlapping algorithms, the subject is considered not matching.

The subject `name` differentiates between artifacts. The semantics are up to the
producer and consumer. Because consumers evaluate the name against a policy, it
should be stable between attestations. If the name is not meaningful, use "\_".
For example, a [Provenance] attestation might use the name to specify output
filename, expecting the consumer to only considers entries with a particular
name. Alternatively, a vulnerability scan attestation might use the name "\_"
because the results apply regardless of what the artifact is named.

IMPORTANT: Subject artifacts are matched purely by digest, regardless of content
type. If this matters to you, please open a GitHub issue to discuss.

The `predicateType` and `predicate` together form the [Predicate], describing
metadata about the artifacts referenced by`subject`.

See [processing model](#processing-model) for more details.

## Predicate

The required `predicateType` is a [URI][RFC 3986] describing the overall meaning
of the attestation as well as the schema of `predicate`. The optional
`predicate` contains additional details.

The predicate can contain arbitrary information. Users are expected to choose a
predicate type that fits their needs, or invent a new one if no existing one
satisfies. Type URIs are not registered; the natural namespacing of URIs is
sufficient to prevent collisions.

This repo defines the following predicate types:

*   [Provenance]: To describe the origins of a software artifact.
*   [Link]: For migration from [in-toto 0.9].
*   [SPDX](spec/predicates/spdx.md): A Software Package Data Exchange document.

We recommend the following conventions for predicates:

*   Field names SHOULD use lowerCamelCase.

*   Timestamps SHOULD use [RFC 3339] syntax with timezone "Z" and SHOULD clarify
    the meaning of the timestamp. For example, a field named `timestamp` is too
    ambiguous; a better name would be `builtAt` or `allowedAt` or `scannedAt`.

*   References to other artifacts SHOULD be an object that includes a `digest`
    field of type [DigestSet]. Consider using the same type as [Provenance]
    `materials` if it is a good fit.

*   Predicates SHOULD be designed to encourage policies to be "monotonic,"
    meaning that deleting an attestation will never turn a DENY decision into an
    ALLOW. One reason because verifiers MUST ignore unrecognized subject digest
    types; if no subject is recognized, the attestation is effectively deleted.
    Example: instead of "deny if a 'has vulnerabilities' attestation exists",
    prefer "deny unless a 'no vulnerabilities' attestation exists".

Predicate designers are free to limit what subject types are valid for a given
predicate type. For example, suppose a "Gerrit code review" predicate only
applies to git commit subjects. In that case, a producer of such attestations
should never use a subject other than a git commit.

## Processing model

The following pseudocode shows how to verify and extract metadata about a single
artifact from a single attestation. The expectation is that consumers will feed
the resulting metadata into a policy engine.

TODO: Explain how to process multiple artifacts and/or multiple attestations.

Inputs:

*   `artifactToVerify`: blob of data
*   `attestation`: JSON-encoded [Envelope]
*   `recognizedAttesters`: collection of (`name`, `publicKey`) pairs
*   `acceptableDigestAlgorithms`: collection of acceptable cryptographic hash
    algorithms (usually just `sha256`)

Steps:

*   Envelope layer:
    *   Decode `attestation` as a JSON-encoded [Envelope]; reject if decoding
        fails
    *   Initialize `attesterNames` as an empty set of names
    *   For each `signature` in the envelope:
        *   For each (`name`, `publicKey`) in `recognizedAttesters`:
            *   Optional: skip if `signature.keyid` does not match `publicKey`
            *   If `signature.sig` matches `publicKey`:
                *   Add `name` to `attesterNames`
    *   Reject if `attesterNames` is empty
*   Intermediate state: `payloadType`, `payload`, `attesterNames`
*   Statement layer:
    *   Reject if `payloadType` != `https://in-toto.io/Attestation/v1-json`
    *   Decode `payload` as a JSON-encoded [Statement], reject if decoding fails
    *   Initialize `artifactNames` as an empty set of names
    *   For each subject `s` in the statement:
        *   For each digest (`alg`, `value`) in `s.digest`
            *   If `alg` is in `acceptableDigestAlgorithms`:
                *   If `hash(alg, artifactToVerify)` == `hexDecode(value)`:
                    *   Add `s.name` to `artifactNames`
    *   Reject if `artifactNames` is empty

Output (to be fed into policy engine):

*   `predicateType`
*   `predicate`
*   `artifactNames`
*   `attesterNames`

[DigestSet]: field_types.md#DigestSet
[Envelope]: #envelope
[ITE-5]: https://github.com/MarkLodato/ITE/blob/ite-5/ITE/5/README.md
[Link]: predicates/link.md
[Predicate]: #predicate
[Provenance]: predicates/provenance.md
[RFC 3339]: https://tools.ietf.org/html/rfc3339
[RFC 3986]: https://tools.ietf.org/html/rfc3986
[SLSA Attestation Model]: https://github.com/slsa-framework/slsa-controls/blob/main/attestations.md
[SPDX]: predicates/spdx.md
[Statement]: #statement
[TypeURI]: field_types.md#TypeURI
[in-toto 0.9]: https://github.com/in-toto/docs/blob/v0.9/in-toto-spec.md
[processing model]: #processing-model
