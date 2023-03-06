# Statement layer specification

Version: v1.0-draft

The Statement is the middle layer of the attestation, binding it to a
particular subject and unambiguously identifying the types of the
[Predicate].

## Schema

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

## Fields

The Statement is represented as a [JSON] object with the following fields.
Additional [parsing rules] apply.

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

`subject[*].name` _string, optional_

> Identifier to distinguish this artifact from others within the `subject`.
>
> The semantics are up to the producer and consumer and they MAY use it when
> evaluating policy. If the name is not meaningful, leave the field unset or
> use "\_". For example, a [SLSA Provenance] attestation might use the name
> to specify output filename, expecting the consumer to only consider
> entries with a particular name. Alternatively, a vulnerability scan
> attestation might leave name unset because the results apply regardless of
> what the artifact is named.
>
> If set, `name` SHOULD be unique within subject.

`predicateType` _string ([TypeURI]), required_

> URI identifying the type of the [Predicate].

`predicate` _object, optional_

> Additional parameters of the [Predicate]. Unset is treated the same as
> set-but-empty. MAY be omitted if `predicateType` fully describes the
> predicate.

[DigestSet]: digest_set.md
[Predicate]: predicate.md
[SLSA Provenance]: https://slsa.dev/provenance
[TypeURI]: scalar_field_types.md#TypeURI
[parsing rules]: README.md#parsing-rules
