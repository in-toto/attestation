# Predicate type: Link

Type URI: https://in-toto.io/Link/v0.2

Version: 0.2.0

## Purpose

A generic attestation type with a schema isomorphic to [in-toto 0.9]. This
allows existing in-toto users to make minimal changes to upgrade to the new
attestation format.

Most users should migrate to a more specific attestation type, such as
[Provenance](provenance.md).

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/Link/v0.2",
  "predicate": {
    "name": "...",
    "command": "...",
    "materials": { ... },
    "byproducts": { ... },
    "environment": { ... }
  }
}
```

_(Note: This is a Predicate type that fits within the larger
[Attestation](../README.md) framework.)_

The `predicate` has the same schema as the link's `signed` field in
[in-toto 0.9] except:

-   `predicate._type` is omitted.  `predicateType` serves the same purpose.
-   `predicate.products` is omitted. `subject` serves the same purpose.

## Converting to old-style links

A Link predicate may be converted into an in-toto 0.9 link as follows:

-   Set `link` to be a copy of `predicate`.
-   Set `link.type` to `"link"`.
-   Set `link.products` to be a map from `subject[*].name` to
    `subject[*].digest`.

In Python:

```python
def convert(statement):
    assert statement.predicateType == 'https://in-toto.io/Link/v0.2'
    link = statement.predicate.copy()
    link['_type'] = 'link'
    link['products'] = {s['name'] : s['digest'] for s in statement.subject}
    return link
```

## TODO

-   [ ] Bump up the in-toto version from 0.9 to 1.0 once
    [in-toto/docs issue #46](https://github.com/in-toto/docs/issues/46) is
    resolved.

## Version History

-   0.2: Removed `_type` and `products`. Defined conversion rules.
-   0.1: Initial version.

[in-toto 0.9]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md#44-file-formats-namekeyid-prefixlink
