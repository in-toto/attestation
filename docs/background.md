# Background

We capture motivational context for the creation of the in-toto Attestation
Framework here.

## Motivation

This project was initiated with two main motivations:

1.  Support use cases where the prior in-toto link schema was a poor fit. For
    example, test steps and vulnerability scans are not about "producing" a new
    artifact so they are awkward to represent in the original format.
2.  Support interoperability with [Binary Authorization], which will support
    the agreed-upon format once finalized. This way we have a single ecosystem
    of software supply chain security.

## Requirements

Functional requirements:

-   Must support user-defined types and schemas, for two reasons:
    -   To allow in-toto users to more naturally express attestations, as
        explained above.
    -   Because Binary Authorization does not want to require its users to use
        the existing in-toto link schema, which is overly specific.
-   Should allow indexing of attestations by their artifact identifier, i.e., a
    digest, without having to understand the user-defined schema.
    -   Primary reason: To support generic attestation indexing/storage/fetching
        without requiring user configuration for every type of attestation.
    -   Secondary reason: To simplify the programming model of policies. The
        binding between artifact and attestation can be done in a framework
        without requiring type-dependent configuration.
    -   Implication: the association between attestation and primary artifact
        identifier must be standardized.
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

-   Must support backwards compatible Links that can be consumed by existing
    Layout files.
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
signed by X with materials/products Z?"[^1] These relatively simple policies are
quite powerful. With this proposal, such policies become more expressive without
any additional configuration: "does an attestation exist that is signed by X
having predicate type T, with subject Y/Z?"

Second, it enables lookup of attestations by `subject`, again without
Predicate-specific logic or configuration. Consider a validation policy as "fetch
attestations for artifact X". The lookup could be from a set of attestations
provided by the caller, or it could be from an external database keyed by
subject.[^2] Without a standardized `subject` field, this would be significantly
harder.

The alternative is to not have a fixed Statement schema and instead have
`subject` be part of the Predicate. Doing so would require users to configure
the system for every possible Predicate type they wanted to support, in order to
instruct the system how to find the subject. Furthermore, because there would be
no standardization, concepts and models may not necessarily translate between
predicate types. For example, one predicate type might require an "or" between
artifact IDs, while another requires an "and." This difference would add
complexity and confusion.

[^1]: The `expected_command` is only a warning, and `inspections` require running external commands which is infeasible in many situations.

[^2]: That said, we strongly recommend against keying a database purely by content hash. The reason is that such databases quickly run into scaling issues, as explained in [Building Secure and Reliable Systems](https://static.googleusercontent.com/media/landing.google.com/en//sre/static/pdf/Building_Secure_and_Reliable_Systems.pdf#page=364), Chapter 14, page 328, "Ensure Unambiguous Provenance." Instead, we recommend keying primarily by resource name, in addition to content hash.

[Binary Authorization]: https://cloud.google.com/binary-authorization
[Predicate]: spec/README.md#predicate
[Statement]: spec/README.md#statement
