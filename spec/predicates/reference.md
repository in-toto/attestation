# Predicate type: Reference

Type URI: https://in-toto.io/attestation/reference/v0.1

Version: v0.1.0

## Purpose

The reference attestation references a set of documents that are relevant to
some resource. It is used to point to and sign metadata that is not
provided in-band.

## Use Cases

One use case is to provide an SBOM of an artifact without needing to sign
the entire SBOM and transmit it in-band. A consumer, such as
[GUAC](https://guac.sh/), retrieves the SBOM, verifies the digest, and uses it
for vulnerability management, compliance, analysis, or any other purpose.

Another use case is to provide evidence for a policy such as "deny unless a
reference to an SBOM is provided".

There is encodability overlap with the SCAI predicate, but it is not the best
fit because it doesn't support standardizing an attribute (which would be
desired here), the use case is distinct, and it is a larger prequisite.
A more opinionated predicate will increase usability and encourage adoption.

## Prerequisites

The
[in-toto Attestation Framework](https://github.com/in-toto/attestation/blob/main/spec/README.md)
and an SBOM specification such as [SPDX](https://spdx.dev/).

## Model

This predicate is intended to be generated and consumed throughout the software
supply chain. In addition, it is intended to be used in the analysis of it as a
whole.

## Data definition

The predicate grammar is provided in CDDL.
Undefined directives are imported from the base v1.2 specification.

```cddl
reference-predicate = (
  predicateType-label => "https://in-toto.io/attestation/reference/v0.1",
  predicate-label => reference-predicate-map
)
reference-predicate-map = {
  reference-attester-label => reference-attester-map
  reference-references-label => [ * ResourceDescriptor ]
}
reference-attester-label   = JC<"attester",   0>
reference-references-label = JC<"references", 1>

reference-attester-map = {
  reference-attester-id-label => uri-type
}
reference-attester-id-label = JC<"id", 0>
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).

The `references` field is repeated as an optimization: multiple references
SHOULD be treated as if each reference was represented in its own attestation
with all other content unchanged.

The reference(s) are associated with the set of subjects as a whole, instead of
being independently associated with each subject. See the
[example](#reference-to-an-sbom-for-multiple-artifacts) with two subjects.

### Fields

`attester.id`: string ([TypeUri](../v1/field_types.md#typeuri)), *required*

An identifier for the system that provides the document.

`references` array of [ResourceDescriptor](../v1/field_types.md#resource_descriptor) objects, *required*

The referred documents. The `downloadLocation` and `mediaType` MUST be provided
for each. If the file type is unknown, `application/octet-stream` SHOULD be
used.

## Examples

### Reference to an SBOM for an image

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{
    "name": "registry.example.com/my-project/my-image",
    "digest": {
      "sha256": "886a77e9b9c993e221d1843e3af5b1fbeea26d8962995e8562174a6aba0c7cc9"
    }
  }],
  "predicateType": "https://in-toto.io/attestation/reference/v0.1",
  "predicate": {
    "attester": {
      "id": "https://my-builder.com/v1"
    },
    "references": [{
      "downloadLocation": "https://cloud.com/my-sboms/sbom.spdx.json",
      "digest": {
        "sha256": "3b0d19b348f1e46a571c6e9df32897637c20dd9803261c2fc9cbe38b0c8422c4"
      },
      "mediaType": "application/spdx+json"
    }]
  }
}
```

### Reference to an SBOM for multiple artifacts

In this example, a single SBOM was generated for a set of build outputs by
scanning the source file system. Per the [parsing rules](#parsing-rules), this
attestation SHOULD be interpreted to mean that the SBOM was generated for both
subjects -- it will list dependencies of both foo and bar.

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "foo.deb",
      "digest": {
        "sha256": "f1847f7e67aa18d0063fda3cbde8b61f84384fbf35c3dd4cfa8c4822400b5a64"
      }
    },
    {
      "name": "bar.deb",
      "digest": {
        "sha256": "680177210742d53016820ec118b6d7dd0be62758ad6d7f503e96927ced4809b4"
      }
    }
  ],
  "predicateType": "https://in-toto.io/attestation/reference/v0.1",
  "predicate": {
    "attester": {
      "id": "https://my-builder.com/v1"
    },
    "references": [{
      "downloadLocation": "https://cloud.com/my-sboms/sbom.spdx.json",
      "digest": {
        "sha256": "3b0d19b348f1e46a571c6e9df32897637c20dd9803261c2fc9cbe38b0c8422c4"
      },
      "mediaType": "application/spdx+json"
    }]
  }
}
```

## Changelog and Migrations

Updated to include CDDL description.
