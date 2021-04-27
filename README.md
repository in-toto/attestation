# In-toto Attestations

This repository defines the in-toto **attestation** format, which represents
authenticated metadata about a set of software artifacts. Attestations are
intended for consumption by automated policy engines, such as [in-toto] and
[Binary Authorization].

## IMPORTANT

This format is still in development and is not yet used by in-toto (see
[ITE-6]). Versions 0.x are unstable and are subject to change until 1.0. Use git
tags to switch between versions.

Latest stable version:
[0.1.0](https://github.com/in-toto/attestation/tree/v0.1.0)

## Introduction

An in-toto **attestation** is authenticated metadata about one or more software
artifacts, as per the [SLSA Attestation Model].

Examples of hypothetical attestations:

*   [Provenance]: GitHub Actions attests to the fact that it built a container
    image with digest "sha256:87f7fe…" from git commit "f0c93d…" in the "master"
    branch of "https://github.com/example/foo".
*   Code review: GitHub attests to the fact that Alice uploaded and Bob approved
    git commit "f0c93d…" in the "master" branch of
    "https://github.com/example/foo".
*   Test result: GitHub Actions attests to the fact that the npm tests passed on
    git commit "f0c93d…".
*   Vulnerability scan: Google Container Analysis attests to the fact that no
    vulnerabilities were found in container image "sha256:87f7fe…" at a
    particular time.
*   Policy decision: Binary Authorization attests to the fact that container
    image "sha256:87f7fe…" is allowed to run under GKE project "example-project"
    within the next 4 hours, and that it used the four attestations above and as
    well as the policy with sha256 hash "79e572" to make its decision.

Goals:

*   Standardize artifact metadata without being specific to the producer or
    consumer. This way CI/CD pipelines, vulnerability scanners, and other
    systems can generate a single set of attestations that can be consumed by
    anyone, such as [in-toto] or [Binary Authorization].
*   Make it possible to write automated policies that take advantage of
    structured information.
*   Fit within the [SLSA Framework][SLSA]. The [provenance] format defined here
    is the official SLSA recommendation.

## Tooling / how to use

None yet!

## Specification

See [spec/README.md](spec/README.md). Summary:

<!-- NOTE: When updating below, also update spec/README.md -->

*   [Envelope]: Handles authentication and serialization.
*   [Statement]: Binds the attestation to a particular subject and unambiguously
    identifies the types of the predicate.
*   [Predicate]: Contains arbitrary metadata about the subject, with a
    type-specific schema. This repo defines the following predicate types,
    though custom types are allowed:
    *   [Provenance]: To describe the origins of a software artifact.
    *   [Link]: For migration from [in-toto 0.9].
    *   [SPDX]: A Software Package Data Exchange document.

The [processing model] provides pseudocode showing how these layers fit
together.

## Examples

### Provenance example

A [Provenance]-type attestation describing how the
[curl 7.72.0 source tarballs](https://curl.se/download.html) were built,
pretending they were built on
[GitHub Actions](https://github.com/features/actions).

```json
{
  "payloadType": "application/vnd.in-toto+json",
  "payload": "ewogICJzdWJqZWN0IjogWwogICAg...",
  "signatures": [{"sig": "MeQyap6MyFyc9Y..."}]
}
```

where `payload` base64-decodes as the following [Statement]:

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [
    { "name": "curl-7.72.0.tar.bz2",
      "digest": { "sha256": "ad91970864102a59765e20ce16216efc9d6ad381471f7accceceab7d905703ef" }},
    { "name": "curl-7.72.0.tar.gz",
      "digest": { "sha256": "d4d5899a3868fbb6ae1856c3e55a32ce35913de3956d1973caccd37bd0174fa2" }},
    { "name": "curl-7.72.0.tar.xz",
      "digest": { "sha256": "0ded0808c4d85f2ee0db86980ae610cc9d165e9ca9da466196cc73c346513713" }},
    { "name": "curl-7.72.0.zip",
      "digest": { "sha256": "e363cc5b4e500bfc727106434a2578b38440aa18e105d57576f3d8f2abebf888" }}
  ],
  "predicateType": "https://in-toto.io/Provenance/v0.1",
  "predicate": {
    "builder": { "id": "https://github.com/Attestations/GitHubHostedActions@v1" },
    "recipe": {
      "type": "https://github.com/Attestations/GitHubActionsWorkflow@v1",
      "definedInMaterial": 0,
      "entryPoint": "build.yaml:maketgz"
    },
    "metadata": {
      "buildStartedOn": "2020-08-19T08:38:00Z"
    },
    "materials": [
      {
        "uri": "git+https://github.com/curl/curl-docker@master",
        "digest": { "sha1": "d6525c840a62b398424a78d792f457477135d0cf" }
      }, {
        "uri": "github_hosted_vm:ubuntu-18.04:20210123.1"
      }, {
        "uri": "git+https://github.com/actions/checkout@v2",
        "digest": {"sha1": "5a4ac9002d0be2fb38bd78e4b4dbde5606d7042f"}
      }, {
        "uri": "git+https://github.com/actions/upload-artifact@v2",
        "digest": { "sha1": "e448a9b857ee2131e752b06002bf0e093c65e571" }
      }, {
        "uri": "pkg:deb/debian/stunnel4@5.50-3?arch=amd64",
        "digest": { "sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa" }
      }, {
        "uri": "pkg:deb/debian/python-impacket@0.9.15-5?arch=all",
        "digest": { "sha256": "71fa2e67376c8bc03429e154628ddd7b196ccf9e79dec7319f9c3a312fd76469" }
      }, {
        "uri": "pkg:deb/debian/libzstd-dev@1.3.8+dfsg-3?arch=amd64",
        "digest": { "sha256": "91442b0ae04afc25ab96426761bbdf04b0e3eb286fdfbddb1e704444cb12a625" }
      }, {
        "uri": "pkg:deb/debian/libbrotli-dev@1.0.7-2+deb10u1?arch=amd64",
        "digest": { "sha256": "05b6e467173c451b6211945de47ac0eda2a3dccb3cc7203e800c633f74de8b4f" }
      }
    ]
  }
}
```

### Custom-type examples

In many cases, custom-type predicates would be a more natural fit, as shown
below. Such custom attestations are not yet supported by in-toto because the
layout format has no way to reference such attestations. Still, we show the
examples to explain the benefits for the new link format.

The initial step is often to write code. This has no materials and no real
command. The existing [Link] schema has little benefit. Instead, a custom
`predicateType` would avoid all of the meaningless boilerplate fields.

(Only the base64-decoded `payload` [Statement] is shown, since the outer
[Envelope] looks the same in all cases.)

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [
    { "digest": { "sha1":  "859b387b985ea0f414e4e8099c9f874acb217b94" } }
  ],
  "predicateType": "https://example.com/CodeReview/v1",
  "predicate": {
    "repo": {
      "type": "git",
      "uri": "https://github.com/example/my-project",
      "branch": "main"
    },
    "author": "mailto:alice@example.com",
    "reviewers": ["mailto:bob@example.com"]
  }
}
```

Test results are also an awkward fit for the Link schema, since the subject is
really the materials, not the products. Again, a custom `predicateType` is a
better fit:

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [
    { "digest": { "sha1": "859b387b985ea0f414e4e8099c9f874acb217b94" } }
  ],
  "predicateType": "https://example.com/TestResult/v1",
  "predicate": {
    "passed": true
  }
}
```

## Motivation

This project has two main goals:

1.  Support [use cases](#motivating-use-case) where the existing link schema is
    a poor fit. For example, test steps and vulnerability scans are not about
    "producing" a new artifact so they are awkward to represent in the current
    format.
2.  Support interoperability with
    [Binary Authorization](https://cloud.google.com/binary-authorization), which
    will support the agreed-upon format once finalized. This way we have a
    single ecosystem of software supply chain security.

Functional requirements:

*   Must support user-defined types and schemas, for two reasons:
    *   To allow in-toto users to more naturally express attestations, as
        explained above.
    *   Because Binary Authorization does not want to require its users to use
        the existing in-toto link schema, which is overly specific.
*   Should allow indexing of attestations by artifact ID, without having to
    understand the user-defined schema.
    *   Primary reason: To support generic attestation indexing/storage/fetching
        without requiring user configuration for every type of attestation.
    *   Secondary reason: To simplify the programming model of policies. The
        binding between artifact and attestation can be done in a framework
        without requiring type-dependent configuration.
    *   Implication: the association between attestation and primary artifact ID
        must be standardized.
*   Should allow identification of related artifacts IDs given an attestation,
    without having to understand the user-defined schema.
    *   Reason: To support "inline attestations," where the client fetches and
        sends all required attestations to the server for verification. The
        client does not know the policy ahead of time or understand all
        attestation types.
    *   Example: Given a provenance attestation for a docker image, it should be
        possible to identify all the materials generically.
    *   Implication: the association between attestation and related artifact
        IDs must be standardized.

Nonfunctional requirements:

*   Must support backwards compatible links that can be consumed by existing
    layout files.
*   Must differentiate between different types of related artifacts (only if
    related artifacts are standardized.) Examples: materials vs products,
    sources vs build tools.
    *   Should be type-dependent, rather than mandating "materials" and
        "products."

## Reasoning

### Reason for separate Statement and Predicate layers

The [Statement] layer has a fixed schema while the [Predicate] layer has an
arbitrary schema. Furthermore, the fixed Statement schema has a `subject` and
`predicateType`. There are two main reasons for this.

First, doing so allows policy engines to make decisions without requiring
Predicate-specific logic or configuration. Binary Authorization policies today
are purely about "does an attestation exist that is signed by X with subject Y",
and similarly in-toto layouts are about "does an attestation exist that is
signed by X with materials/products Z?"[1] These relatively simple policies are
quite powerful. With this proposal, such policies become more expressive without
any additional configuration: "does an attestation exist that is signed by X
having predicate type T, with subject Y/Z?"

Second, it enables lookup of attestations by `subject`, again without
Predicate-specific logic or configuration. Consider the policy described in the
[motivating use case](#motivating-use-case). There, the instruction is "fetch
attestations for artifact X". The lookup could be from a set of attestations
provided by the caller, or it could be from an external database keyed by
subject.[2] Without a standardized `subject` field, this would be significantly
harder.

The alternative is to not have a fixed Statement schema and instead have
`subject` be part of the Predicate. Doing so would require users to configure
the system for every possible Predicate type they wanted to support, in order to
instruct the system how to find the subject. Furthermore, because there would be
no standardization, concepts and models may not necessarily translate between
predicate types. For example, one predicate type might require an "or" between
artifact IDs, while another requires an "and." This difference would add
complexity and confusion.

## Footnotes

\[1]: The `expected_command` is only a warning, and `inspections` require
running external commands which is infeasible in many situations.

\[2]: That said, we strongly recommend against keying a database purely by
content hash. The reason is that such databases quickly run into scaling issues,
as explained in
[Building Secure and Reliable Systems](https://static.googleusercontent.com/media/landing.google.com/en//sre/static/pdf/Building_Secure_and_Reliable_Systems.pdf#page=364),
Chapter 14, page 328, "Ensure Unambiguous Provenance." Instead, we recommend
keying primarily by resource name, in addition to content hash.

[Binary Authorization]: https://cloud.google.com/binary-authorization
[Envelope]: spec/README.md#envelope
[ITE-6]: https://github.com/in-toto/ITE/pull/15
[Link]: spec/predicates/link.md
[Predicate]: spec/README.md#predicate
[Provenance]: spec/predicates/provenance.md
[SLSA Attestation Model]: https://github.com/slsa-framework/slsa-controls/blob/main/attestations.md
[SLSA]: https://github.com/slsa-framework/slsa
[SPDX]: spec/predicates/spdx.md
[Statement]: spec/README.md#statement
[in-toto 0.9]: https://github.com/in-toto/docs/blob/v0.9/in-toto-spec.md
[in-toto]: https://in-toto.io
[processing model]: spec/processing_model.md
