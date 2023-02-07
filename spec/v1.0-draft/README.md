# in-toto Attestation Framework Specification

An **in-toto attestation** is authenticated metadata about one or more
software artifacts[^1]. The intended consumers are automated policy engines,
such as [in-toto-verify] and [Binary Authorization].

It has three layers that are independent but designed to work together:

-   [Envelope]: Handles authentication and serialization.
-   [Statement]: Binds the attestation to a particular subject and
    unambiguously identifies the types of the predicate.
-   [Predicate]: Contains arbitrary metadata about the subject, with a
    type-specific schema.
-   [Bundle]: Defines a method of grouping multiple attestations together.

The [validation model] provides pseudocode showing how these layers fit
together. See the [documentation](../docs) for more background and examples.

## Parsing rules

The following rules apply to [Statement] and [Predicates] that opt-in to this
model.

-   **Unrecognized fields:** Consumers MUST ignore unrecognized fields. This
    is to allow minor version upgrades and extension fields. Ignoring fields
    is safe due to the monotonic principle.

-   **Versioning:** Each type has a [SemVer2](https://semver.org) version
    number and the [TypeURI] reflects the major version number. A message is
    always semantically correct, but possibly incomplete, when parsed as any
    other version with the same major version number and thus the same
    [TypeURI]. Minor version changes always follow the monotonic principle.
    NOTE: 0.X versions are considered major versions.

-   **Extension fields:** Producers MAY add extension fields to any JSON
    object by using a property name that is a [TypeURI]. The use of URI is
    to protect against name collisions. Consumers MAY parse and use these
    extensions if desired. The presence or absence of the extension field
    MUST NOT influence the meaning of any other field, and the field MUST
    follow the monotonic princple.

-   **Monotonic:** A policy is considered monotonic if ignoring an
    attestation, or a field within an attestation, will never turn a DENY
    decision into an ALLOW. A predicate or field follows the monotonic
    principle if the expected policy that consumes it is monotonic.
    Consumers SHOULD design policies to be monotonic. Example: instead of
    "deny if a 'has vulnerabilities' attestation exists", prefer "deny
    unless a 'no vulnerabilities' attestation exists".

See [versioning rules](../versioning.md) for details and examples.

## Envelope

The Envelope is the outermost layer of the attestation, handling
authentication and serialization. The format and protocol are defined in
[DSSE] and adopted by in-toto in [ITE-5].

### Schema

```jsonc
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "<Base64(Statement)>",
  "signatures": [{"sig": "<Base64(Signature)>"}]
}
```

### Fields

`payloadType` _string, required_

> Identifier for the encoding of the payload. Always
> `application/vnd.in-toto+json`, which indicates that it is a JSON object with
> a `_type` field indicating its schema.

`payload` _string, required_

> Base64-encoded JSON [Statement].

`signatures` _array of objects, required_

> One or more signatures over `payloadType` and `payload`, as defined in [DSSE].

## Statement

The Statement is the middle layer of the attestation, binding it to a
particular subject and unambiguously identifying the types of the
[Predicate].

### Schema

Version:
[1.0.0](https://github.com/in-toto/attestation/blob/v1.0/spec/README.md) (see
[parsing rules])

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1.0",
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

### Fields

`_type` _string ([TypeURI]), required_

> Identifier for the schema of the Statement. Always
> `https://in-toto.io/Statement/v1.0` for this version of the spec.

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

The Predicate is the innermost layer of the attestation, containing arbitrary
metadata about the [Statement]'s `subject`.

### Schema

```jsonc
"predicateType": "<URI>",
"predicate": {
    // arbitrary object
}
```

### Fields

A predicate has a required `predicateType` ([TypeURI]) identifying what the
predicate means, plus an optional `predicate` (object) containing additional,
type-dependent parameters.

Users are expected to choose an [existing predicate type] that
fits their needs, or develop a new one if no existing one satisfies.
New predicate types MAY be vetted by the in-toto attestation maintainers.

[^1]: This is compatible with the [SLSA Attestation Model].

[Binary Authorization]: https://cloud.google.com/binary-authorization
[Bundle]: bundle.md
[DSSE]: https://github.com/secure-systems-lab/dsse
[DigestSet]: field_types.md#DigestSet
[Envelope]: #envelope
[ITE-5]: https://github.com/in-toto/ITE/blob/master/ITE/5/README.adoc
[JSON]: https://www.json.org
[Predicate]: #predicate
[RFC 3339]: https://tools.ietf.org/html/rfc3339
[SLSA Attestation Model]: https://slsa.dev/attestation-model
[SLSA Provenance]: https://slsa.dev/provenance
[Statement]: #statement
[TypeURI]: field_types.md#TypeURI
[existing predicate type]: predicates/README.md
[in-toto-verify]: https://github.com/in-toto/in-toto#verification
[parsing rules]: #parsing-rules
[validation model]: ../../docs/validation.md