# In-toto Attestations

This repository defines the in-toto **attestation** format, which represents
authenticated metadata about a set of software artifacts. Attestations are
intended for consumption by automated policy engines, such as [in-toto] and
[Binary Authorization].

## IMPORTANT

This specification is a work in progress:

*   This doc was spun off from [ITE-6] and some parts have not yet been updated.
*   Details are subject to change until the initial release.
*   There are not yet any producers or consumers of this format.

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

An attestation has three layers that are independent but designed to work
together:

*   [Envelope]: Handles authentication and serialization.
*   [Statement]: Binds the attestation to a particular subject and
    unambiguously identifies the types of the predicate.
*   [Predicate]: Contains arbitrary metadata about the subject, with a
    type-specific schema. This repo defines the following predicate types,
    though custom types are allowed:
    *   [Provenance]: To describe the origins of a software artifact.
    *   [Link]: For migration from [in-toto 0.9].
    *   [SPDX](spec/predicates/spdx.md): A Software Package Data Exchange document.

It may help to first look at [Examples](#examples) to get an idea.

### Envelope

```jsonc
{
  "payloadType": "https://in-toto.io/Statement/v1-json",
  "payload": "...",
  "signatures": [<SIGNATURES>]
}
```

The Envelope is defined in
[signing-spec](https://github.com/secure-systems-lab/signing-spec) and adopted
by in-toto in [ITE-5].

### Statement

```jsonc
{
  "subject": [
    {
      "name": "<NAME>",
      "digest": {"<ALGORITHM>": "<HEX_VALUE>"}
    },
    ...
  ],
  "predicateType": "...",
  "predicate": {
    <PREDICATE>
  }
}
```

*   Type URI: https://in-toto.io/Statement/v1-json (value of `payloadType` in
    [Envelope])
*   Encoding: [JSON](https://www.json.org)
*   Schema: [statement.proto](spec/statement.proto)

The `subject` describes the set of software artifacts that the attestation
applies to. Each entry has a `name` and at least one `digest`.

The subject `digest`, of type [DigestSet], is a collection of alternate content
hashes of a single artifact. Two DigestSets are considered matching if ANY of
the fields match. The producer and consumer must agree on acceptable algorithms.
If there are no overlapping algorithms, the subject is considered not matching.

The subject `name` differentiates between artifacts. The semantics are up to the
producer and consumer. Because consumers evaluate the name against a policy, it
should be stable between attestations. If the name is not meaningful, use "\_".
For example, a [Provenance] attestation might use the name to specify output
filename, expecting the consumer to only considers entries with a particular
name. Alternatively, a vulnerability scan attestation might use the name "\_"
because the results apply regardless of what the artifact is named.

IMPORTANT: Subject artifacts are matched purely by digest, regardless of content
type. If this matters to you, please open a GitHub issue to discuss.

The `predicateType` and `predicate` together form the [Predicate], describing
metadata about the artifacts referenced by`subject`.

See [processing model](#processing-model) for more details.

### Predicate

The required `predicateType` is a [URI][RFC 3986] describing the overall meaning
of the attestation as well as the schema of `predicate`. The optional
`predicate` contains additional details.

The predicate can contain arbitrary information. Users are expected to choose a
predicate type that fits their needs, or invent a new one if no existing one
satisfies. Type URIs are not registered; the natural namespacing of URIs is
sufficient to prevent collisions.

This repo defines the following predicate types:

*   [Provenance]: To describe the origins of a software artifact.
*   [Link]: For migration from [in-toto 0.9].
*   [SPDX](spec/predicates/spdx.md): A Software Package Data Exchange document.

We recommend the following conventions for predicates:

*   Field names SHOULD use lowerCamelCase.

*   Timestamps SHOULD use [RFC 3339] syntax with timezone "Z" and SHOULD clarify
    the meaning of the timestamp. For example, a field named `timestamp` is too
    ambiguous; a better name would be `builtAt` or `allowedAt` or `scannedAt`.

*   References to other artifacts SHOULD be an object that includes a `digest`
    field of type [DigestSet]. Consider using the same type as [Provenance]
    `materials` if it is a good fit.

*   Predicates SHOULD be designed to encourage policies to be "monotonic,"
    meaning that deleting an attestation will never turn a DENY decision into an
    ALLOW. One reason because verifiers MUST ignore unrecognized subject digest
    types; if no subject is recognized, the attestation is effectively deleted.
    Example: instead of "deny if a 'has vulnerabilities' attestation exists",
    prefer "deny unless a 'no vulnerabilities' attestation exists".

Predicate designers are free to limit what subject types are valid for a given
predicate type. For example, suppose a "Gerrit code review" predicate only
applies to git commit subjects. In that case, a producer of such attestations
should never use a subject other than a git commit.

## Processing model

The following pseudocode shows how to verify and extract metadata about a single
artifact from a single attestation. The expectation is that consumers will feed
the resulting metadata into a policy engine.

TODO: Explain how to process multiple artifacts and/or multiple attestations.

Inputs:

*   `artifactToVerify`: blob of data
*   `attestation`: JSON-encoded [Envelope]
*   `recognizedAttesters`: collection of (`name`, `publicKey`) pairs
*   `acceptableDigestAlgorithms`: collection of acceptable cryptographic hash
    algorithms (usually just `sha256`)

Steps:

*   Envelope layer:
    *   Decode `attestation` as a JSON-encoded [Envelope]; reject if decoding
        fails
    *   Initialize `attesterNames` as an empty set of names
    *   For each `signature` in the envelope:
        *   For each (`name`, `publicKey`) in `recognizedAttesters`:
            *   Optional: skip if `signature.keyid` does not match `publicKey`
            *   If `signature.sig` matches `publicKey`:
                *   Add `name` to `attesterNames`
    *   Reject if `attesterNames` is empty
*   Intermediate state: `payloadType`, `payload`, `attesterNames`
*   Statement layer:
    *   Reject if `payloadType` != `https://in-toto.io/Attestation/v1-json`
    *   Decode `payload` as a JSON-encoded [Statement], reject if decoding fails
    *   Initialize `artifactNames` as an empty set of names
    *   For each subject `s` in the statement:
        *   For each digest (`alg`, `value`) in `s.digest`
            *   If `alg` is in `acceptableDigestAlgorithms`:
                *   If `hash(alg, artifactToVerify)` == `hexDecode(value)`:
                    *   Add `s.name` to `artifactNames`
    *   Reject if `artifactNames` is empty

Output (to be fed into policy engine):

*   `predicateType`
*   `predicate`
*   `artifactNames`
*   `attesterNames`

## Examples

### Provenance example

A [Provenance]-type attestation describing how the
[curl 7.72.0 source tarballs](https://curl.se/download.html) were built,
pretending they were built on
[GitHub Actions](https://github.com/features/actions).

```json
{
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
  "predicateType": "https://in-toto.io/Provenance/v1",
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
        "digest": { "sha1": "d6525c840a62b398424a78d792f457477135d0cf" },
        "mediaType": "application/vnd.git.commit",
        "tags": ["source"]
      }, {
        "uri": "github_hosted_vm:ubuntu-18.04:20210123.1",
        "tags": ["base-image"]
      }, {
        "uri": "git+https://github.com/actions/checkout@v2",
        "digest": {"sha1": "5a4ac9002d0be2fb38bd78e4b4dbde5606d7042f"},
        "mediaType": "application/vnd.git.commit",
        "tags": ["dev-dependency"]
      }, {
        "uri": "git+https://github.com/actions/upload-artifact@v2",
        "digest": { "sha1": "e448a9b857ee2131e752b06002bf0e093c65e571" },
        "mediaType": "application/vnd.git.commit",
        "tags": ["dev-dependency"]
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

```json
{
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
    "reviewers": ["mailto:bob@example.com"],
  }
}
```

Test results are also an awkward fit for the Link schema, since the subject is
really the materials, not the products. Again, a custom `predicateType` is a
better fit:

```json
{
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
[DigestSet]: spec/field_types.md#DigestSet
[Envelope]: #envelope
[ITE-5]: https://github.com/MarkLodato/ITE/blob/ite-5/ITE/5/README.md
[ITE-6]: https://github.com/in-toto/ITE/pull/15
[Link]: spec/predicates/link.md
[Predicate]: #predicate
[Provenance]: spec/predicates/provenance.md
[RFC 3339]: https://tools.ietf.org/html/rfc3339
[RFC 3986]: https://tools.ietf.org/html/rfc3986
[SLSA Attestation Model]: https://github.com/slsa-framework/slsa-controls/blob/main/attestations.md
[SLSA]: https://github.com/slsa-framework/slsa
[Statement]: #statement
[in-toto 0.9]: https://github.com/in-toto/docs/blob/v0.9/in-toto-spec.md
[in-toto]: https://in-toto.io
