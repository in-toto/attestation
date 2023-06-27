# Spec protobuf definitions

Protobuf definitions for the in-toto spec and some predicates are
provided in the [protos/] directory. We provide a list of
[supported language bindings] for the spec and predicates.

**DISCLAIMER**: The protobuf definitions and language bindings will not be
considered stable before the v1.1 tagged release. Use at your own risk.

## Pre-requisites

On an Ubuntu-based system, install the following dependencies.

```shell
sudo apt install build-essential protobuf-compiler golang python3 python3-pip
```

## Protobuf programming practices

You should follow standard [protobuf programming practices] when developing
a protobuf definition.

**NOTE**: This means that while [specification documents] [by convention use]
lowerCamelCase for field names, the protobuf definitions use lower_snake_case
for field names per the standard protobuf convention.

We establish the following project specific practices, in addition to the
standard protobuf programming practices:

### Package versioning

To enable consumers to support multiple versions of the in-toto attestation
spec or predicates[^1], we maintain versioned sub-packages for the protos.

This means, the protos for a new major version should be placed under new
`vMAJ` sub-package under the respective protos package:

-   spec in `protos/in\_toto\_attestation/`
-   predicates in `protos/in\_toto\_attestation/<predicate>/`

### Updates

Minor version updates to the protobufs are expected to be fully backwards
compatible (per [semver guidelines]), and should be made directly to the
proto defintions in the corresponding `vMAJ` sub-package.

To ensure backwards compatibility of the updated definition, please review
the message type [update guidelines].

## Regenerating proto libraries

[It's typical](https://go.dev/doc/articles/go_command#:~:text=and%20then%20check%20those%20generated%20source%20files%20into%20your%20repository)
to keep code generated from protobuf definitions in the repository itself,
since it makes users' lives much easier. However, do NOT manually regenerate
and check in the libraries if your change modifies or adds protos.

To ensure libraries are generated using consistent tooling, we have
[automated their generation](/.github/workflows/make-protos.yml).
Therefore, they will be regenerated automatically, after your change is
merged.

[^1]: This is especially helpful during transitions between major versions of the spec or predicate.

[protobuf programming practices]: https://protobuf.dev/programming-guides/proto3
[update guidelines]: https://protobuf.dev/programming-guides/proto3/#updating
[protos/]: ../protos/
[semver guidelines]: https://semver.org/#summary
[by convention use]: ../docs/new_predicate_guidelines.md#predicate-conventions
[specification documents]: ../spec/
[supported language bindings]: ../protos/README.md#supported-language-bindings
