# Predicate type: Container Image Promotion Record

Type URI: https://k8s.io/promo-tools/promotion/v1

Version: 1

Authors: Kubernetes Release Engineering @kubernetes/release-engineering

## Purpose

This predicate type records the promotion of a container image from a staging
registry to a production registry mirror. It provides an auditable record of
when, where, and by what system an image was copied, enabling consumers to
verify the provenance chain of images served from production registries.

## Use Cases

### Kubernetes Image Promotion Auditing

The Kubernetes project uses the [Kubernetes Image Promoter] to copy release
images from staging registries (e.g., `gcr.io/k8s-staging-*`) to production
registry mirrors (e.g., `registry.k8s.io`). A promotion record attestation is
generated when the copy is done, creating a verifiable link between the staging
source and the production destination. This allows consumers to confirm that
images in production registries were legitimately promoted from their expected
staging origins.

### Supply Chain Continuity

Promotion records bridge the gap between build-time provenance (e.g.,
[SLSA Provenance]) and the images actually served to end users. A build
provenance attestation tells you how an image was built and pushed to staging.
A promotion record tells you that the same image (by digest) was subsequently
promoted to production. Together, they provide end-to-end supply chain
traceability from source code to released artifact.

### Detecting Unauthorized Image Modifications

By recording the exact digest that was promoted, consumers can detect if an
image in a production registry has been replaced or modified after promotion.
If the digest in the production registry does not match the digest recorded in
the promotion record, the image may have been tampered with.

## Model

This predicate captures a single image promotion event performed by an
automated promotion system. The subject of the attestation is the promoted
container image identified by its digest. The predicate records the source and
destination image references, the promotion timestamp, and the identity of the
promotion system.

The promotion is a server-side registry copy that preserves the image digest.
Both the source and destination references therefore refer to the same image
content, but in different registry locations.

## Schema

```jsonc
{
  // Standard attestation fields:
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [{ ... }],

  // Predicate:
  "predicateType": "https://in-toto.io/attestation/promotion-record/v1",
  "predicate": {
    "srcRef": "<IMAGE_REFERENCE>",
    "dstRef": "<IMAGE_REFERENCE>",
    "digest": "<DIGEST>",
    "timestamp": "<TIMESTAMP>",
    "builderId": "<URI>"
  }
}
```

### Parsing Rules

This predicate follows the
[in-toto Attestation Framework's parsing rules](../v1/README.md#parsing-rules).

`srcRef` _string_, _required_

The fully qualified source image reference in the staging registry
(e.g., `gcr.io/k8s-staging-foo/bar:v1.2.3`). This is the location the image
was copied from.

`dstRef` _string_, _required_

The fully qualified destination image reference in the production registry
(e.g., `registry.k8s.io/foo/bar:v1.2.3`). This is the location the image was
copied to.

`digest` _string_, _required_

The image digest in the standard `algorithm:hex` format
(e.g., `sha256:abc123...`). Because promotion is a server-side copy, this
digest is the same at both the source and destination.

`timestamp` _string ([Timestamp])_, _required_

An [RFC 3339] formatted timestamp recording when the promotion occurred.

`builderId` _string ([ResourceURI])_, _required_

A URI identifying the promotion system that performed the copy, including its
version (e.g., `https://k8s.io/promo-tools@v4.0.8`). This allows consumers
to evaluate trust in the promoter and track which version was used.

## Example

A single image promoted from the `k8s-staging-kas-network-proxy` staging
registry to the production `registry.k8s.io`:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "registry.k8s.io/kas-network-proxy/proxy-agent",
      "digest": {
        "sha256": "d5b3f4e8c7a1b2d9e6f0c8a4b7d1e3f5a9c2b6d0e4f8a1c5b9d3e7f0a2c6b8"
      }
    }
  ],
  "predicateType": "https://k8s.io/promo-tools/promotion/v1",
  "predicate": {
    "srcRef": "gcr.io/k8s-staging-kas-network-proxy/proxy-agent:v0.31.0",
    "dstRef": "registry.k8s.io/kas-network-proxy/proxy-agent:v0.31.0",
    "digest": "sha256:d5b3f4e8c7a1b2d9e6f0c8a4b7d1e3f5a9c2b6d0e4f8a1c5b9d3e7f0a2c6b8",
    "timestamp": "2025-07-15T14:30:00Z",
    "builderId": "https://k8s.io/promo-tools@v4.0.8"
  }
}
```

[Kubernetes Image Promoter]: https://github.com/kubernetes-sigs/promo-tools
[SLSA Provenance]: https://slsa.dev/provenance
[ResourceURI]: ../v1/field_types.md#resourceuri
[Timestamp]: ../v1/field_types.md#timestamp
[RFC 3339]: https://tools.ietf.org/html/rfc3339
