# in-toto Attestation Predicates

Whenever possible, in-toto Attestation Framework users are encouraged
to choose existing attestation predicates that best fit their needs,
although the framework easily supports new predicate types for new use cases.

Anyone is welcome to contribute their new predicate back to the in-toto
community! Please see our [New Predicate Guidelines].

## Vetted Predicates

This directory contains predicate specification types that have gone through
our [vetting process], and may be of general interest:

-   [CycloneDX]: CycloneDX BOM for software artifacts.
<!-- markdownlint-disable-next-line MD059 -->
-   [Link]: For migration from [in-toto 0.9].
-   [Reference]: References documents that are relevant to some resource.
-   [Release]: Details an artifact that is part of a given release version.
-   [Runtime Traces]: Captures runtime traces of software supply chain
    operations.
-   [SCAI Report]: Evidence-based assertions about software artifact and
    supply chain attributes or behavior.
-   [SLSA Provenance]: Describes how an artifact or set of artifacts was
    produced.
-   [SLSA Verification Summary]: SLSA verification decision about a software
    artifact.
-   [SPDX2] and [SPDX3]: SPDX-formatted BOM for software artifacts.
-   [Simple Verification Result]: Evidence that an artifact has been evaluated
    against one or more policies.
-   [Test Result]: A generic schema to express results of any type of tests.
-   [VULNS]: Defines the metadata to share the results of vulnerability scanning
    on software artifacts.

## Proposed Predicates

-   [Decision Receipt]: Captures access control decisions from AI agent tool
    calls and physical sensor attestations, with policy evidence and hash-chained
    integrity.

[CycloneDX]: https://cyclonedx.org/
[Decision Receipt]: decision-receipt.md
[Link]: link.md
[New Predicate Guidelines]: ../../docs/new_predicate_guidelines.md
[Release]: release.md
[Runtime Traces]: runtime-trace.md
[SCAI Report]: scai.md
[VULNS]: vulns_02.md
[SLSA Provenance]: https://slsa.dev/provenance
[SLSA Verification Summary]: vsa.md
[SPDX2]: spdx2.md
[SPDX3]: spdx3.md
[Test Result]: test-result.md
[in-toto 0.9]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md#44-file-formats-namekeyid-prefixlink
[vetting process]: ../../docs/new_predicate_guidelines.md#vetting-process
[Reference]: reference.md
[Simple Verification Result]: svr.md
