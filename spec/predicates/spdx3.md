# Predicate type: SPDX

Type URI: https://spdx.dev/Document/v3

Version: SPDX 3 document version identified by the embedded document

## Purpose

A Bill of Materials type following version 3 of the [SPDX Specification].

This can represent software artifacts, software supply chains, AI models and
more. For a complete list, see the [SPDX 3 Scope].

## Non-goal

This proposal does not introduce new safety semantics. It only clarifies
how predicates identify and validate embedded CycloneDX and SPDX documents.

## Prerequisites

The in-toto [attestation] framework and a [SPDX generation tool].

## Model

This is a predicate type that fits within the larger [Attestation] framework.

## Schema

The schema of this predicate type is documented in the
[SPDX Specification].

### Parsing Rules

The parsing rules for the embedded document are documented in the
[SPDX Specification].

The predicate type `https://spdx.dev/Document/v3` identifies the SPDX 3
document family. Consumers must determine the embedded document's version
and applicable profiles using the version identification and conformance
rules defined by the corresponding SPDX specification. A compatible SPDX 3
minor-version update does not by itself require a new predicate type.

The predicate URI alone does not establish support for every SPDX 3 minor
version or profile. Consumers must validate against the applicable upstream
model, serialization rules, and profile constraints. A JSON-LD context
alone is not sufficient evidence of profile conformance.

When a policy requires a profile, consumers must verify both the profile's
required content and their ability to interpret it. An unsupported version
or profile must not be treated as satisfying that requirement. Missing
required information must not be interpreted as a successful assessment.

Development versions and release candidates must be identified as such.
This clarification does not declare an unreleased SPDX version or profile
to be supported.

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
