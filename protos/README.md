# in-toto attestation protobuf definitions

To ensure all in-toto libraries use the same common data format as defined by
the in-toto Attestation Framework spec, we provide protobuf definitions.
These enable us to pre-generate bindings for different languages that use the
same underlying spec format.

**DISCLAIMER**: The protobuf definitions and Golang bindings will not be
considered stable until the v1.1 tagged release. Use at your own risk.

## Predicates with protobuf definitions

In addition to the core in-toto attestation spec, the following attestation
predicates have protobuf definitions:

-   [SLSA Verification Summary]: SLSA verification decision about a software
    artifact.

## Supported language bindings

We currently support bindings for the following languages:

-   [go]

## Usage

We outline the package names to import the protobufs or language bindings in
your project.

### Protos

To use any `.proto` definitions in this repo in your protobufs, import the
following packages as needed:

-   in-toto attestation layers: `in_toto_attestation/v1`
-   attestation predicates: `in_toto_attestation/predicates`

### Go

The Go bindings for the attestations layers and predicates are provided in
the `github.com/in-toto/attestation/go/v1` and
`github.com/in-toto/attestation/go/predicates` packages, respectively.

## Building the language bindings

Please read our protos [documentation] for instructions on building and
testing the supported language bindings.

[SLSA Verification Summary]: in_toto_attestation/predicates/vsa/
[documentation]: ../docs/protos.md
[go]: ../go/
