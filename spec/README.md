# Attestation Spec

An in-toto **attestation** is authenticated metadata about one or more software
artifacts, as per the [SLSA Attestation Model]. It has three layers that are
independent but designed to work together:

<!-- BEGIN: When updating below, also update ../README.md#specification -->

-   [Envelope]: Handles authentication and serialization.
-   [Statement]: Binds the attestation to a particular subject and unambiguously
    identifies the types of the predicate.
-   [Predicate]: Contains arbitrary metadata about the subject, with a
    type-specific schema.
-   [Bundle]: Defines a method of grouping multiple attestations together.

The [processing model] provides pseudocode showing how these layers fit
together.

<!-- END -->

See the [top-level README](../README.md) for background and examples.

## Parsing rules

The following rules apply to [Statement] and predicates that opt-in to this
model.

-   **Unrecognized fields:** Consumers MUST ignore unrecognized fields. This is
    to allow minor version upgrades and extension fields. Ignoring fields is
    safe due to the monotonic principle.

-   **Versioning:** Each type has a [SemVer2](https://semver.org) version number
    and the [TypeURI] reflects the major version number. A message is always
    semantically correct, but possibly incomplete, when parsed as any other
    version with the same major version number and thus the same [TypeURI].
    Minor version changes always follow the monotonic principle. NOTE: 0.X
    versions are considered major versions.

-   **Extension fields:** Producers MAY add extension fields to any JSON object
    by using a property name that is a [TypeURI]. The use of URI is to protect
    against name collisions. Consumers MAY parse and use these extensions if
    desired. The presence or absence of the extension field MUST NOT influence
    the meaning of any other field, and the field MUST follow the monotonic
    princple.

-   **Monotonic:** A policy is considered monotonic if ignoring an attestation,
    or a field within an attestation, will never turn a DENY decision into an
    ALLOW. A predicate or field follows the monotonic principle if the expected
    policy that consumes it is monotonic. Consumers SHOULD design policies to be
    monotonic. Example: instead of "deny if a 'has vulnerabilities' attestation
    exists", prefer "deny unless a 'no vulnerabilities' attestation exists".

See [versioning rules](versioning.md) for details and examples.

## Envelope

```jsonc
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "<Base64(Statement)>",
  "signatures": [{"sig": "<Base64(Signature)>"}]
}
```

The Envelope is the outermost layer of the attestation, handling authentication
and serialization. The format and protocol are defined in [DSSE] and adopted by
in-toto in [ITE-5]. It is a [JSON] object with the following fields:

`payloadType` _string, required_

> Identifier for the encoding of the payload. Always
> `application/vnd.in-toto+json`, which indicates that it is a JSON object with
> a `_type` field indicating its schema.

`payload` _string, required_

> Base64-encoded JSON [Statement].

`signatures` _array of objects, required_

> One or more signatures over `payloadType` and `payload`, as defined in [DSSE].

## Statement

```jsonc
{
  "_type": "https://in-toto.io/Statement/v0.1",
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

Version:
[0.1.0](https://github.com/in-toto/attestation/blob/v0.1.0/spec/README.md) (see
[parsing rules])

The Statement is the middle layer of the attestation, binding it to a particular
subject and unambiguously identifying the types of the [predicate]. It is a
[JSON] object with the following fields:

`_type` _string ([TypeURI]), required_

> Identifier for the schema of the Statement. Always
> `https://in-toto.io/Statement/v0.1` for this version of the spec.

`subject` _array of objects, required_

> Set of software artifacts that the attestation applies to. Each element
> represents a single software artifact.
>
> IMPORTANT: Subject artifacts are matched purely by digest, regardless of
> content type. If this matters to you, please comment on
> [GitHub Issue #28](https://github.com/in-toto/attestation/issues/28)

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
> the name is not meaningful, use "\_". For example, a [SLSA Provenance]
> attestation might use the name to specify output filename, expecting the
> consumer to only considers entries with a particular name. Alternatively, a
> vulnerability scan attestation might use the name "\_" because the results
> apply regardless of what the artifact is named.
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

A predicate has a required `predicateType` ([TypeURI]) identifying what the
predicate means, plus an optional `predicate` (object) containing additional,
type-dependent parameters.

Users are expected to choose a predicate type that fits their needs, or invent a
new one if no existing one satisfies. Predicate types are not registered.

The following popular predicate types may be of general interest:

-   [SLSA Provenance]: To describe the origins of a software artifact.
-   [Link]: For migration from [in-toto 0.9].
-   [SPDX]: A Software Package Data Exchange document.

### Predicate conventions

We recommend the following conventions for predicates:

-   Predicates SHOULD follow and opt-in to the [parsing rules], particularly the
    monotonic principle, and SHOULD explain what the parsing rules are.

-   Field names SHOULD use lowerCamelCase.

-   Timestamps SHOULD use [RFC 3339] syntax with timezone "Z" and SHOULD clarify
    the meaning of the timestamp. For example, a field named `timestamp` is too
    ambiguous; a better name would be `builtAt` or `allowedAt` or `scannedAt`.

-   References to other artifacts SHOULD be an object that includes a `digest`
    field of type [DigestSet]. Consider using the same type as [SLSA Provenance]
    `materials` if it is a good fit.

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

-   `artifactToVerify`: blob of data
-   `attestation`: JSON-encoded [Envelope]
-   `recognizedAttesters`: collection of (`name`, `publicKey`) pairs
-   `acceptableDigestAlgorithms`: collection of acceptable cryptographic hash
    algorithms (usually just `sha256`)

Steps:

-   Envelope layer:
    -   `envelope` := decode `attestation` as a JSON-encoded [Envelope]; reject
        if decoding fails
    -   `attesterNames` := empty set of names
    -   For each `signature` in `envelope.signatures`:
        -   For each (`name`, `publicKey`) in `recognizedAttesters`:
            -   Optional: skip if `signature.keyid` does not match `publicKey`
            -   If `signature.sig` matches `publicKey`:
                -   Add `name` to `attesterNames`
    -   Reject if `attesterNames` is empty
-   Intermediate state: `envelope.payloadType`, `envelope.payload`,
    `attesterNames`
-   Statement layer:
    -   Reject if `envelope.payloadType` != `application/vnd.in-toto+json`
    -   `statement` := decode `envelope.payload` as a JSON-encoded [Statement];
        reject if decoding fails
    -   Reject if `statement.type` != `https://in-toto.io/Statement/v0.1`
    -   `artifactNames` := empty set of names
    -   For each `s` in `statement.subject`:
        -   For each digest (`alg`, `value`) in `s.digest`:
            -   If `alg` is in `acceptableDigestAlgorithms`:
                -   If `hash(alg, artifactToVerify)` == `hexDecode(value)`:
                    -   Add `s.name` to `artifactNames`
    -   Reject if `artifactNames` is empty

Output (to be fed into policy engine):

-   `statement.predicateType`
-   `statement.predicate`
-   `artifactNames`
-   `attesterNames`

[Bundle]: bundle.md
[DSSE]: https://github.com/secure-systems-lab/dsse
[DigestSet]: field_types.md#DigestSet
[Envelope]: #envelope
[ITE-5]: https://github.com/in-toto/ITE/blob/master/ITE/5/README.adoc
[JSON]: https://www.json.org
[Link]: predicates/link.md
[Predicate]: #predicate
[RFC 3339]: https://tools.ietf.org/html/rfc3339
[SLSA Attestation Model]: https://slsa.dev/attestation-model
[SLSA Provenance]: https://slsa.dev/provenance
[SPDX]: predicates/spdx.md
[Statement]: #statement
[TypeURI]: field_types.md#TypeURI
[in-toto 0.9]: https://github.com/in-toto/docs/blob/v0.9/in-toto-spec.md
[parsing rules]: #parsing-rules
[processing model]: #processing-model
