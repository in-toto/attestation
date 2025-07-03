# Predicate type: Simple Verification

Type URI: https://in-toto.io/attestation/simple-verification/v0.1

Version: 0.1

Authors: Tom Hennen (@TomHennen), Andrew McNamara (@arewm)

## Purpose

A Simple Verification Attestation (SVA) communicates that an artifact has been evaluated against one or more policies, and records the properties that were verified, at the point-in-time of the evaluation. This enables consumers to trust a concise summary of verification results, without requiring access to all underlying attestations or policies.

## Use Cases

-   Allowing software consumers to delegate complex policy decisions to a trusted verifier, and rely on a simple summary of which properties were verified.
-   Supporting evolving verification needs (e.g., SLSA, vulnerability scanning, organizational policies) without coupling the SVA format to any specific framework.
-   Enabling verifiers to add additional information (e.g., policy references, expiration, property-level timing) as needed, while maintaining a simple core schema.

## Prerequisites

This predicate depends on the [in-toto Attestation Framework]. Familiarity with the concept of attestation statements and the use of URIs for identifying verifiers and resources is assumed.

## Model

An SVA is an attestation that a verifier (`verifier.id`) has checked an artifact (`resourceUri`) at a specific time (`timeCreated`) and asserts a set of verified properties (`properties`). The SVA does not include information for reproducing the verification result. Instead, verifiers MAY generate additional provenance attestations for the simple verification attestation as needed.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/simple-verification/v0.1",
  "predicate": {
    "verifier": {
      "id": "<URI>"
    },
    "timeCreated": "<TIMESTAMP>",
    "resourceUri": "<artifact-URI>",
    "properties": ["<VERIFIED_PROPERTY>", ...]
  }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).

### Fields

**`verifier`, required** object

> Identifies the entity that performed the verification.

**`verifier.id`, required** string (TypeURI)

> URI indicating the verifier's identity.

**`timeCreated`, required** string (Timestamp)

> Timestamp indicating what time the verification occurred.

**`resourceUri`, required** string (ResourceURI)

> URI that identifies the resource associated with the artifact being verified.

**`properties`, required** array of string

> Indicates the properties verified for the artifact. These SHOULD be scoped according
> to the framework being verified or the verifier's policy rules.

## Example

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "static",
      "digest": {"sha256": "3A244fd47e07d10..."}
    }
  ],
  "predicateType": "https://in-toto.io/attestation/simple-verification/v0.1",
  "predicate": {
    "verifier": {
      "id": "https://example.com/publication_verifier"
    },
    "timeCreated": "1985-04-12T23:20:50.52Z",
    "resourceUri": "pkg:oci/static@sha256%3A244fd47e07d10...",
    "properties": [
      "SLSA_BUILD_LEVEL_3",
      "SLSA_SOURCE_LEVEL_3",
      "ORG_SOURCE_STATIC_ANALYSIS",
      "CONFORMA_HERMETIC_BUILD_TASK",
      "COMPANY_RELEASE_APPROVED",
    ]
  }
}
```

## Changelog and Migrations

-   0.1: Initial version.

[in-toto Attestation Framework]: ../v1/README.md
