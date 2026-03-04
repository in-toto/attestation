# in-toto Attestation Predicates

<!-- markdownlint-disable MD059 -->

Whenever possible, in-toto Attestation Framework users are encouraged
to choose existing attestation predicates that best fit their needs,
although the framework easily supports new predicate types for new use cases.

Anyone is welcome to contribute their new predicate back to the in-toto
community! Please see our [New Predicate Guidelines].

## Vetted Predicates

This directory contains predicate specification types that have gone through
our [vetting process], and may be of general interest:

| Predicate | Description |
| --- | --- |
| [CycloneDX](cyclonedx.md) | [CycloneDX] BOM for software artifacts. |
| [Link] | For migration from [in-toto 0.9] links. |
| [Promotion Record] | Records a container image move from registry to registry. |
| [Reference] | References documents that are relevant to some resource. |
| [Release] | Details an artifact that is part of a given release version. |
| [Runtime Traces] | Captures runtime traces of software supply chain operations. |
| [SCAI Report] | Evidence-based assertions about software artifact and supply chain attributes or behavior. |
| [Simple Verification Result] | Captures the result of an artifact's evaluation against policies. |
| [SLSA Provenance] | Describes how an artifact or set of artifacts was produced according to the [SLSA spec]. |
| [SLSA Verification Summary] | SLSA verification decision about a software artifact. |
| [SPDX 2.x] | [SPDX](https://spdx.dev)-formatted software bill of materials. |
| [SPDX 3.x] | [SPDX](https://spdx.dev) 3.X documents |
| [Test Result] | A generic schema to express results of any type of tests. |
| [VULNS] | Defines the metadata to share the results of vulnerability scanning on software artifacts. |

[CycloneDX]: https://cyclonedx.org/
[Link]: link.md
[New Predicate Guidelines]: ../../docs/new_predicate_guidelines.md
[Release]: release.md
[Runtime Traces]: runtime-trace.md
[SCAI Report]: scai.md
[VULNS]: vulns_02.md
[Promotion Record]: promotion_record.md
[Simple Verification Result]: svr.md
[SLSA Provenance]: provenance.md
[SLSA spec]: https://slsa.dev/provenance
[SLSA Verification Summary]: vsa.md
[SPDX 2.x]: spdx2.md
[SPDX 3.x]: spdx3.md
[Test Result]: test-result.md
[in-toto 0.9]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md#44-file-formats-namekeyid-prefixlink
[vetting process]: ../../docs/new_predicate_guidelines.md#vetting-process
[Reference]: reference.md
