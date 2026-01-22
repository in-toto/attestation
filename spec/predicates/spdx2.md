# Predicate type: SPDX

Type URI: https://spdx.dev/Document
Type URI: https://spdx.dev/Document/v2.3

Version: 2.3

## Purpose

A Software Bill of Materials type following version 2 of the
[SPDX Specification].

This allows to represent an "exportable" or "published" software artifact. It
can also be used as an entry point for other types of in-toto attestation when
performing policy decisions.

## Prerequisites

The in-toto [attestation] framework and a [SPDX generation tool].

## Model

This is a predicate type that fits within the larger [Attestation] framework.

## Schema

The schema of this predicate type is documented in the
[SPDX Specification].

### Parsing Rules

The parsing rules for this predicate type are documented in the
[SPDX Specification].

### Fields

The fields that make up this predicate type are documented in the
[SPDX specification].

The `predicate` contains a JSON-encoded SPDX document.
The `subject` contains whatever software artifacts are to be associated with
this SPDX document.

## Example

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://spdx.dev/Document/v2.3",
  "predicate": {
    "SPDXID" : "SPDXRef-DOCUMENT",
    "spdxVersion" : "SPDX-2.3",
    ...
  }
}
```

## Changelog and Migrations

### Version 2.3

-   Added version to predicateType

[Attestation]: ../README.md
[SPDX specification]: https://spdx.github.io/spdx-spec/v2.3
[SPDX generation tool]: https://spdx.dev/resources/tools/
