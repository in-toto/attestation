# Predicate type: CycloneDX

Type URI: https://cyclonedx.org/bom

Version: 1.4

## Purpose

A Software Bill of Materials type following the [CycloneDX standard].

This allows representing "exportable", or "published" software artifacts,
services, vulnerability information, and more. For a complete list of
capabilities see [CycloneDX Capabilities].

## Prerequisites

The in-toto [attestation] framework and a [CycloneDX BOM generation tool].

## Model

This is a predicate type that fits within the larger [Attestation] framework.

## Data definition

The schema of this predicate type is documented in the
[CycloneDX Specification] for vnd.cyclonedx+json.

As of 2024, CycloneDX does not have an official vnd.cyclonedx+cbor content type for CBOR encoding.
To embed a JSON CycloneDX into the predicate, use the [`TN()` transformation] of application/json for wrapping the JSON encoding.

```cddl
cyclonedx-predicate = (
  predicateType-label => "https://cyclonedx.org/bom/v1.4",
  predicate-label => JC<cyclonedx-map, cyclonedx-json-tunnel>
)
cyclonedx-map = object
cyclonedx-json-tunnel = { &(json-embedding: -4478722) => 1668546867(bytes) }
```

### Parsing Rules

The parsing rules for this predicate type are documented in the
[CycloneDX Specification].

### Fields

The fields that make up this predicate type are documented in the
[CycloneDX Specification].

The `predicate` contains a JSON-encoded CycloneDX BOM.
The `subject` contains whatever software artifacts are to be associated with
this CycloneDX BOM document.

## Example

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://cyclonedx.org/bom/v1.4",
  "predicate": {
    "bomFormat": "CycloneDX",
    "specVersion": "1.4",
    "serialNumber": "urn:uuid:3e671687-395b-41f5-a30f-a58921a69b79",
    "version": 1,
    "components": [
        {
        "type": "library",
        "name": "acme-library",
        "version": "1.0.0"
        }
    ]
    ...
  }
}
```

## Changelog and Migrations

Not applicable for this initial version.

[Attestation]: ../README.md
[CycloneDX standard]: https://cyclonedx.org/specification/overview
[CycloneDX Capabilities]: https://cyclonedx.org/capabilities/
[CycloneDX Specification]: https://github.com/CycloneDX/specification/tree/1.4/schema
[CycloneDX BOM generation tool]: https://cyclonedx.org/tool-center
[`TN()` transformation]: https://www.rfc-editor.org/rfc/rfc9277.html#ct-tags
