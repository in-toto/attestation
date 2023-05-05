# Predicate layer specification

Version: v1.0

The Predicate is the innermost layer of the attestation, containing arbitrary
metadata about the [Statement]'s `subject`.

## Schema

```jsonc
"predicateType": "<URI>",
"predicate": {
    // arbitrary object
}
```

## Fields

A predicate has a required `predicateType` ([TypeURI]) identifying what the
predicate means, plus an optional `predicate` [JSON] object containing
additional, type-dependent parameters.

Users are expected to choose an [existing predicate type] that
fits their needs, or develop a new one if no existing one satisfies.
New predicate types MAY be vetted by the in-toto attestation maintainers.

Additional [parsing rules] apply.

[JSON]: https://www.json.org
[Statement]: statement.md
[TypeURI]: field_types.md#TypeURI
[parsing rules]: README.md#parsing-rules
[existing predicate type]: ../predicates
