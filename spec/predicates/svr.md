# Predicate type: Simple Verification

Type URI: https://in-toto.io/attestation/svr/v0.1

Version: 0.1

Authors: Tom Hennen (@TomHennen), Andrew McNamara (@arewm)

## Purpose

A Simple Verification Result (SVR) communicates that an artifact has been
evaluated against one or more policies, and records the properties that were
verified, at the point-in-time of the evaluation. This enables consumers to
establish trust in a concise summary of verification results coming from a
trusted verifier, without requiring access to all underlying attestations or
policies.

## Use Cases

-   Allowing software consumers to delegate complex policy decisions to a
    trusted verifier, and rely on a simple summary of which properties were
    verified.
-   Supporting evolving verification needs (e.g., SLSA, vulnerability scanning,
    organizational policies) without coupling the SVR format to any specific
    framework.
-   Enabling verifiers to add additional information (e.g., policy references,
    expiration, property-level timing) as needed, while maintaining a simple
    core schema.

## Prerequisites

This predicate depends on the [in-toto Attestation Framework]. Familiarity with
the concept of attestation statements and the use of URIs for identifying
verifiers and resources is assumed.

## Model

An SVR is an attestation that a verifier (`verifier.id`) has checked the
subject artifact(s) at a specific time (`timeCreated`) and asserts a set of
verified properties (`properties`). Multiple SVRs can be created for the same
subject, and a single SVR can reference multiple subjects in the attestation
statement. The SVR does not include information for reproducing the
verification result. Instead, verifiers MAY generate additional provenance
attestations for the simple verification attestation as needed.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/svr/v0.1",
  "predicate": {
    "verifier": {
      "id": "<URI>"
    },
    "timeCreated": "<TIMESTAMP>",
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

> URI indicating the verifier's identity. It is RECOMMENDED to version the
> identifiers when their verification logic changes materially (e.g., different
> properties are verified, different validation rules are applied). This allows
> consumers to distinguish between verifier versions and make trust decisions
> accordingly.
>
> Note: Verifiers that use configurable policies MAY include additional fields
> (following the [in-toto extension field conventions]) to communicate policy
> information to consumers familiar with that verifier's implementation. Such
> fields are optional and verifier-specific.

**`timeCreated`, required** string (Timestamp)

> Timestamp indicating what time the verification occurred.

**`properties`, required** array of string

> Indicates the properties verified for the artifact. These SHOULD be scoped
> according to the framework being verified or the verifier's policy rules. For
> example, this could be a policy engine prefix like `AMPEL_` or `CONFORMA_`.

## Examples

### Basic Example

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "static",
      "digest": {"sha256": "3A244fd47e07d10..."}
    }
  ],
  "predicateType": "https://in-toto.io/attestation/svr/v0.1",
  "predicate": {
    "verifier": {
      "id": "https://example.com/publication_verifier/v2"
    },
    "timeCreated": "1985-04-12T23:20:50.52Z",
    "properties": [
      "SLSA_BUILD_LEVEL_3",
      "SLSA_SOURCE_LEVEL_3",
      "ORG_SOURCE_STATIC_ANALYSIS",
      "CONFORMA_HERMETIC_BUILD_TASK",
      "COMPANY_RELEASE_APPROVED"
    ]
  }
}
```

### Time-Bounded Properties Example

This example demonstrates how time-bounded properties can be expressed through
naming conventions. The property `VERIFIER_VULN_SCANNED_1` represents a contract
between the attestation creator (the policy engine) and consumers: it indicates
that a vulnerability scan was performed within 1 day of the attestation's 
`timeCreated` timestamp.

Consumers verify this by checking that the attestation's `timeCreated` is
within their acceptable time window (e.g., no more than 1 day old). This
pattern allows verifiers to communicate temporal guarantees about their checks
without adding explicit expiration fields to the predicate.

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "myapp",
      "digest": {"sha256": "5d5b09f6dcb2d53..."}
    }
  ],
  "predicateType": "https://in-toto.io/attestation/svr/v0.1",
  "predicate": {
    "verifier": {
      "id": "https://example.com/security_scanner/v3"
    },
    "timeCreated": "2024-03-15T10:30:00Z",
    "properties": [
      "VERIFIER_VULN_SCANNED_1",
      "VERIFIER_NO_CRITICAL_VULNERABILITIES",
      "VERIFIER_NO_HIGH_VULNERABILITIES"
    ]
  }
}
```

In this example, `VERIFIER_VULN_SCANNED_1` means the scan was performed within
one day. A consumer checking this attestation on 2024-03-15 would accept it,
but would reject it on 2024-03-17 (more than 1 day later).

### Optional Policy Information Example

This example shows how a verifier might include optional policy information
using extension fields on the verifier object. This is useful when the verifier
uses configurable policies and wants to communicate which policy was applied.
Consumers who understand this verifier's implementation can use this
information; others can safely ignore it per the in-toto parsing rules.

```jsonc
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "webapp",
      "digest": {"sha256": "abc123def456..."}
    }
  ],
  "predicateType": "https://in-toto.io/attestation/svr/v0.1",
  "predicate": {
    "verifier": {
      "id": "https://example.com/policy_engine/v1",
      // Optional extension field specific to this verifier
      "policy": {
        "uri": "https://example.com/policies/production-deployment",
        "digest": {"sha256": "789abc012def..."}
      }
    },
    "timeCreated": "2024-03-15T14:22:00Z",
    "properties": [
      "SLSA_BUILD_LEVEL_3",
      "EXAMPLE_DEPLOYMENT_APPROVED"
    ]
  }
}
```

The `policy` field is an extension on the verifier object that consumers
familiar with this verifier can use to identify exactly which policy was
applied, enabling reproducibility or auditing of the verification decision.

## Changelog and Migrations

-   0.1: Initial version.

[in-toto Attestation Framework]: ../v1/README.md
[in-toto extension field conventions]: ../v1/README.md#parsing-rules
