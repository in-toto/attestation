# Introduction

An in-toto **attestation** is authenticated metadata about one or more software
artifacts, as per the [SLSA Attestation Model].

Examples of hypothetical attestations:

-   [Provenance][SLSA Provenance]: GitHub Actions attests to the fact that it
    built a container image with digest "sha256:87f7fe…" from git commit
    "f0c93d…" in the "main" branch of "https://github.com/example/foo".
-   Test result: GitHub Actions attests to the fact that the npm tests passed on
    git commit "f0c93d…".
-   Vulnerability scan: Google Container Analysis attests to the fact that no
    vulnerabilities were found in container image "sha256:87f7fe…" at a
    particular time.
-   Policy decision: Binary Authorization attests to the fact that container
    image "sha256:87f7fe…" is allowed to run under GKE project "example-project"
    within the next 4 hours, and that it used the four attestations above and as
    well as the policy with sha256 hash "79e572" to make its decision.

Goals:

-   Standardize artifact metadata without being specific to the producer or
    consumer. This way CI/CD pipelines, vulnerability scanners, and other
    systems can generate a single set of attestations that can be consumed by
    anyone, such as [in-toto] or [Binary Authorization].
-   Make it possible to write automated policies that take advantage of
    structured information.
-   Fit within the [SLSA Framework][SLSA].

## Goals

This project has two main goals:

1.  Support use cases where the existing link schema is
    a poor fit. For example, test steps and vulnerability scans are not about
    "producing" a new artifact so they are awkward to represent in the current
    format.
2.  Support interoperability with
    [Binary Authorization](https://cloud.google.com/binary-authorization), which
    will support the agreed-upon format once finalized. This way we have a
    single ecosystem of software supply chain security.

Functional requirements:

-   Must support user-defined types and schemas, for two reasons:
    -   To allow in-toto users to more naturally express attestations, as
        explained above.
    -   Because Binary Authorization does not want to require its users to use
        the existing in-toto link schema, which is overly specific.
-   Should allow indexing of attestations by artifact ID, without having to
    understand the user-defined schema.
    -   Primary reason: To support generic attestation indexing/storage/fetching
        without requiring user configuration for every type of attestation.
    -   Secondary reason: To simplify the programming model of policies. The
        binding between artifact and attestation can be done in a framework
        without requiring type-dependent configuration.
    -   Implication: the association between attestation and primary artifact ID
        must be standardized.
-   Should allow identification of related artifacts IDs given an attestation,
    without having to understand the user-defined schema.
    -   Reason: To support "inline attestations," where the client fetches and
        sends all required attestations to the server for verification. The
        client does not know the policy ahead of time or understand all
        attestation types.
    -   Example: Given a provenance attestation for a docker image, it should be
        possible to identify all the materials generically.
    -   Implication: the association between attestation and related artifact
        IDs must be standardized.

Nonfunctional requirements:

-   Must support backwards compatible links that can be consumed by existing
    layout files.
-   Must differentiate between different types of related artifacts (only if
    related artifacts are standardized.) Examples: materials vs products,
    sources vs build tools.
    -   Should be type-dependent, rather than mandating "materials" and
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
Predicate-specific logic or configuration. Consider a validation policy as "fetch
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
[Predicate]: spec/README.md#predicate
[SLSA Attestation Model]: https://slsa.dev/attestation-model
[SLSA Provenance]: https://slsa.dev/provenance
[SLSA]: https://github.com/slsa-framework/slsa
[Statement]: spec/README.md#statement
[in-toto]: https://in-toto.io
