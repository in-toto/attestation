# Predicate type: SPDX

Type URI: https://spdx.dev/Document/v3

Version: 3.0

## Purpose

A Bill of Materials type following version 3 of the [SPDX Specification].

This can represent software artifacts, software supply chains, AI models and
more. For a complete list, see the [SPDX 3 Scope].

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
[SPDX Specification].

The `predicate` contains a JSON-encoded SPDX document.
The `subject` contains whatever artifacts are to be associated with this SPDX
document.

## Example

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://spdx.dev/Document/v3",
  "predicate": {
    "@context": "https://spdx.org/rdf/3.0/spdx-context.jsonld",
    "@graph": [
        ...
    ]
  }
}
```

[Attestation]: ../README.md
[SPDX 3 scope]: https://spdx.github.io/spdx-spec/v3.0/scope/
[SPDX specification]: https://spdx.dev/use/specifications/
[SPDX generation tool]: https://spdx.dev/resources/tools/
