# Attestation type: Link v1

## Purpose

A generic attestation type with the same schema as [in-toto 0.9].

Most users should migrate to a more specific attestation type. This is defined
to allow existing customers to work with their existing in-toto configuration.

## Schema

```jsonc
{
  "attestation_type": "https://in-toto.io/Provenance/v1",
  "subject": {
    "<URI-OR-PATH>": {
      "<DIGEST_TYPE>": "<DIGEST_VALUE>"
    }
  },
  "materials": {
    "<URI-OR-PATH>": {
      "<DIGEST_TYPE>": "<DIGEST_VALUE>"
    }
  },
  "_name": "<STRING>",
  "command": "<STRING>",
  "byproducts": { /* object */ },
  "environment": { /* object */ }
}
```

[ArtifactCollection]: field_types.md#ArtifactCollection
[TypeURI]: field_types.md#TypeURI
[in-toto 0.9]: https://github.com/in-toto/docs/blob/v0.9/in-toto-spec.md
