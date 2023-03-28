# in-toto Attestation Predicates

Whenever possible, in-toto Attestation Framework users are encouraged
to choose existing attestation predicates that best fit their needs,
although the framework easily supports new predicate types for new use cases.

Anyone is welcome to contribute their new predicate back to the in-toto
community! Please see our [New Predicate Guidelines].

## Vetted Predicates

This directory contains predicate specification types that have gone through
our [vetting process], and may be of general interest:

-   [SLSA Provenance]: Describes how an artifact or set of artifacts was
    produced.
-   [Link]: For migration from [in-toto 0.9].
-   [SCAI Report]: Evidence-based assertions about software artifact and
    supply chain attributes or behavior.
-   [Runtime Traces]: Captures runtime traces of software supply chain
    operations.
-   [SLSA Verification Summary]: SLSA verification decision about a software
    artifact.
-   [CycloneDX]: CycloneDX BOM for software artifacts.

[Link]: link.md
[New Predicate Guidelines]: ../../docs/new_predicate_guidelines.md
[SCAI Report]: scai.md
[SLSA Provenance]: https://slsa.dev/provenance
[SLSA Verification Summary]: https://github.com/in-toto/attestation/blob/main/spec/predicates/vsa/vsa.md
[in-toto 0.9]: https://github.com/in-toto/docs/blob/master/in-toto-spec.md#44-file-formats-namekeyid-prefixlink
[vetting process]: ../../docs/new_predicate_guidelines.md#vetting-process
[Runtime Traces]: runtime-trace.md
[CycloneDX]: https://cyclonedx.org/
