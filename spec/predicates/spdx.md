# Predicate type: SPDX

Type URI: https://spdx.dev/Document

Version: 2.3, 3.0

## Purpose

A Bill of Materials type following the [SPDX Specification].

Version 2.3 can represent an "exportable" or "published" software artifact. It
can also be used as an entry point for other types of in-toto attestation when
performing policy decisions.

Version 3.0 can represent software artifacts, software supply chains, AI models
and more. For a complete list, see the [SPDX 3.0 Scope].

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
The `subject` contains whatever artifacts are to be associated with this SPDX
document.

## Example (Version 2.3)

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

## Example (Version 3.0)

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://spdx.dev/Document/v3.0",
  "predicate": {
    "@context": "https://spdx.org/rdf/3.0/spdx-context.jsonld",
    "@graph": [
        ...
    ]
  }
}
```

## Changelog and Migrations

### Version 2.3

-   Added version to predicateType

[Attestation]: ../README.md
[SPDX 3.0 scope]: https://spdx.github.io/spdx-spec/v3.0/scope/
[SPDX specification]: https://spdx.dev/use/specifications/
[SPDX generation tool]: https://spdx.dev/resources/tools/
