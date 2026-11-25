# Predicate type: CycloneDX

Type URI: https://cyclonedx.org/bom

Version: CycloneDX document version declared by `specVersion`

## Purpose

A Software Bill of Materials type following the [CycloneDX standard].

This allows representing "exportable", or "published" software artifacts,
services, vulnerability information, and more. For a complete list of
capabilities see [CycloneDX Capabilities].

## Non-goal

This proposal does not introduce new safety semantics. It only clarifies
how predicates identify and validate embedded CycloneDX and SPDX documents.

## Prerequisites

The in-toto [attestation] framework and a [CycloneDX BOM generation tool].

## Model

This is a predicate type that fits within the larger [Attestation] framework.

## Schema

The schema of this predicate type is documented in the
[CycloneDX Specification].

### Parsing Rules

The parsing rules for the embedded document are documented in the
[CycloneDX Specification].

The predicate type is `https://cyclonedx.org/bom`. The embedded document's
`specVersion` identifies its CycloneDX schema version; the BOM's `version`
field identifies a revision of that BOM, not a schema version. A change to
`specVersion` does not by itself require a different predicate type.

Consumers must select validation rules for the declared document version.
Recognizing the predicate type alone does not establish that a consumer
supports that version, every capability it contains, or an application
profile. A policy that requires a particular capability or profile must
check that requirement explicitly. Unsupported versions or semantics must
not be treated as satisfying such a policy.

Development schemas and release candidates must be identified as such;
this specification does not imply support for an unreleased schema.

The versioned URI used in the previous example,
`https://cyclonedx.org/bom/v1.4`, differs from the official CycloneDX
predicate type. Correcting the example does not define an alias or require
consumers to silently rewrite existing signed statements.

### Fields

The fields that make up this predicate type are documented in the
[CycloneDX Specification].

The `predicate` contains a JSON-encoded CycloneDX BOM.
The `subject` contains whatever software artifacts are to be associated with
this CycloneDX BOM document.

## Example

The digest below is illustrative and must be replaced with the artifact's
actual digest when producing an attestation.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "acme-library-1.0.0.bin",
      "digest": {
        "sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
      }
    }
  ],
  "predicateType": "https://cyclonedx.org/bom",
  "predicate": {
    "bomFormat": "CycloneDX",
    "specVersion": "1.4",
    "serialNumber": "urn:uuid:3e671687-395b-41f5-a30f-a58921a69b79",
    "version": 1,
    "metadata": {
      "component": {
        "type": "library",
        "name": "acme-library",
        "version": "1.0.0",
        "hashes": [
          {
            "alg": "SHA-256",
            "content": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
          }
        ]
      }
    },
    "components": []
  }
}
```

## Changelog and Migrations

Clarify document-version selection and profile validation, and align the
example predicate type with the official CycloneDX URI. The example retains
a released CycloneDX 1.4 document; it does not add native safety fields.

[Attestation]: ../README.md
[CycloneDX standard]: https://cyclonedx.org/specification/overview
[CycloneDX Capabilities]: https://cyclonedx.org/capabilities/
[CycloneDX Specification]: https://cyclonedx.org/specification/overview/
[CycloneDX BOM generation tool]: https://cyclonedx.org/tool-center
