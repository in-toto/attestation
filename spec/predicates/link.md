# Predicate type: Link v1

Type URI: https://in-toto.io/Link/v1

## Purpose

A generic attestation type with the same schema as [in-toto 0.9]. This allows
existing in-toto users to make minimal changes to upgrade to the new attestation
format.

Most users should migrate to a more specific attestation type, such as
[Provenance](provenance.md).

## Schema

```jsonc
{
  "subject": { ... }
  "predicateType": "https://in-toto.io/Link/v1",
  "predicate": {
    "_type": "link",
    "name": "...",
    "command": "...",
    "materials": { ... },
    "products": { ... },
    "byproducts": { ... },
    "environment": { ... }
  }
}
```

The `predicate` has the same schema as the link's `signed` field in
[in-toto 0.9]. See that document for details.

The `subject` MUST contain whatever elements from `products` or `materials` make
sense. For example, a traditional "build" step would list the `products` in the
`subject`, whereas a "test" or "vulnerability scan" would like the relevant
`materials`.

## TODO

*   [ ] Bump up the in-toto version from 0.9 to 1.0 once
    [in-toto/docs issue #46](https://github.com/in-toto/docs/issues/46) is
    resolved.

[in-toto 0.9]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md#44-file-formats-namekeyid-prefixlink
