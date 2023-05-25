# in-toto attestation protobuf definitions

To ensure all in-toto libraries use the same common data format as defined by
the in-toto Attestation Framework spec, we provide protobuf definitions.
These enable us to pre-generate bindings for different languages that use the
same underlying spec format.

**DISCLAIMER**: The protobuf definitions and language bindings will not be
considered stable before the v1.1 tagged release. Use at your own risk.

## Predicates with protobuf definitions

In addition to the core in-toto attestation spec, the following attestation
predicates have protobuf definitions:

-   [SLSA Verification Summary]: SLSA verification decision about a software
    artifact.
-   [in-toto Link]: Generic predicate that records a software supply chain step.
-   [Test Result]: Expresses the result of a test run in the software supply
    chain.

## Supported language bindings

We currently support bindings for the following languages:

-   [go]
-   [python]
-   [java]

## Usage

We outline the package names to import the protobufs or language bindings in
your project.

To use any `.proto` definitions in this repo in your protobufs, import the
following packages as needed:

-   in-toto attestation layers: `in_toto_attestation/v1`
-   attestation predicates: `in_toto_attestation/predicates`

## Building the language bindings

Please read our protos [documentation] for instructions on building and
testing the supported language bindings.

[SLSA Verification Summary]: in_toto_attestation/predicates/vsa/
[in-toto Link]: in_toto_attestation/predicates/link/
[Test Result]: in_toto_attestation/predicates/test_result/
[documentation]: ../docs/protos.md
[go]: ../go/
[python]: ../python/
[java]: ../java/
