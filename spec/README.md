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
    *   [Rekord](predicate/rekord.md): Recording signed metadata

The [processing model] provides pseudocode showing how these layers fit
together.

<!-- END -->

See the [top-level README](../README.md) for background and examples.

## Envelope

```jsonc
{
  "payloadType": "https://in-toto.io/Statement/v1-json",
  "payload": "<Base64(Statement)>",
  "signatures": [{"sig": "<Base64(Signature)>"}]
}
```

The Envelope is the outermost layer of the attestation, handling authentication
and serialization. The format and protocol are defined in [signing-spec] and
adopted by in-toto in [ITE-5]. It is a [JSON] object with the following fields:

`payloadType` *string, required*

>   Always `https://in-toto.io/Statement/v1-json` (for the [Statement] defined
>   below).

`payload` *string, required*

>   Base64-encoded JSON [Statement].

`signatures` *array of objects, required*

>   One or more signatures over `payloadType` and `payload`, as defined in
>   [signing-spec].

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
  "predicateType": "<URI>",
  "predicate": { ... }
}
```

The Statement is the middle layer of the attestation, binding it to a particular
subject and unambiguously identifying the types of the [predicate]. It is a
[JSON] object with the following fields:

`subject` _array of objects, required_

> Set of software artifacts that the attestation applies to. Each element
> represents a single software artifact.
>
> IMPORTANT: Subject artifacts are matched purely by digest, regardless of
> content type. If this matters to you, please open a GitHub issue to discuss.

`subject[*].digest` _object ([DigestSet]), required_

> Collection of cryptographic digests for the contents of this artifact.
>
> Two DigestSets are considered matching if ANY of the fields match. The
> producer and consumer must agree on acceptable algorithms. If there are no
> overlapping algorithms, the subject is considered not matching.

`subject[*].name` _string, required_

> Identifier to distinguish this artifact from others within the `subject`.
>
> The semantics are up to the producer and consumer. Because consumers evaluate
> the name against a policy, the name SHOULD be stable between attestations. If
> the name is not meaningful, use "\_". For example, a [Provenance] attestation
> might use the name to specify output filename, expecting the consumer to only
> considers entries with a particular name. Alternatively, a vulnerability scan
> attestation might use the name "\_" because the results apply regardless of
> what the artifact is named.
>
> MUST be non-empty and unique within `subject`.

`predicateType` _string ([TypeURI]), required_

> URI identifying the type of the [Predicate].

`predicate` _object, optional_

> Additional parameters of the [Predicate]. Unset is treated the same as
> set-but-empty. MAY be omitted if `predicateType` fully describes the
> predicate.

## Predicate

```jsonc
"predicateType": "<URI>",
"predicate": {
    // arbitrary object
}
```

The Predicate is the innermost layer of the attestation, containing arbitrary
metadata about the [Statement]'s `subject`.

A predicate has a requried `predicateType` ([TypeURI]) identifying what the
predicate means, plus an optional `predicate` (object) containing additional,
type-dependent parameters.

Users are expected to choose a predicate type that fits their needs, or invent a
new one if no existing one satisfies. Predicate types are not registered.

### Pre-defined predicates

This repo defines the following predicate types:

*   [Provenance]: To describe the origins of a software artifact.
*   [Link]: For migration from [in-toto 0.9].
*   [SPDX]: A Software Package Data Exchange document.

### Predicate conventions

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
    *   `envelope` := decode `attestation` as a JSON-encoded [Envelope]; reject
        if decoding fails
    *   `attesterNames` := empty set of names
    *   For each `signature` in `envelope.signatures`:
        *   For each (`name`, `publicKey`) in `recognizedAttesters`:
            *   Optional: skip if `signature.keyid` does not match `publicKey`
            *   If `signature.sig` matches `publicKey`:
                *   Add `name` to `attesterNames`
    *   Reject if `attesterNames` is empty
*   Intermediate state: `envelope.payloadType`, `envelope.payload`,
    `attesterNames`
*   Statement layer:
    *   Reject if `envelope.payloadType` !=
        `https://in-toto.io/Attestation/v1-json`
    *   `statement` := decode `envelope.payload` as a JSON-encoded [Statement];
        reject if decoding fails
    *   `artifactNames` := empty set of names
    *   For each `s` in `statement.subject`:
        *   For each digest (`alg`, `value`) in `s.digest`:
            *   If `alg` is in `acceptableDigestAlgorithms`:
                *   If `hash(alg, artifactToVerify)` == `hexDecode(value)`:
                    *   Add `s.name` to `artifactNames`
    *   Reject if `artifactNames` is empty

Output (to be fed into policy engine):

*   `statement.predicateType`
*   `statement.predicate`
*   `artifactNames`
*   `attesterNames`

[DigestSet]: field_types.md#DigestSet
[Envelope]: #envelope
[ITE-5]: https://github.com/in-toto/ITE/pull/13
[JSON]: https://www.json.org
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
[signing-spec]: https://github.com/secure-systems-lab/signing-spec
