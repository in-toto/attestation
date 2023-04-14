# Spec protobuf definitions

Protobuf definitions for the in-toto spec and some predicates are
provided in the [protos/] directory. We provide a list of
[supported language bindings] for the spec and predicates.

**DISCLAIMER**: The protobuf definitions and Golang bindings will not be
considered stable until the v1.1 tagged release. Use at your own risk.

## Pre-requisites

On an Ubuntu-based system, install the following dependencies.

```shell
sudo apt install build-essential protobuf-compiler golang
```

## Protobuf programming practices

You should follow standard [protobuf programming practices] when developing
a protobuf definition. In addition, we establish the following practices.

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

## Regenerating Go proto libraries

[It's typical to keep generated Go code in the repository itself](https://go.dev/doc/articles/go_command#:~:text=and%20then%20check%20those%20generated%20source%20files%20into%20your%20repository)
since it makes users' lives much easier. However, to ensure libraries are
generated using consistent tooling, we have
[automated their generation](/.github/workflows/make-protos.yml). Therefore, if
your change modifies or adds protos, do NOT regenerate and check in the
libraries. After your change is merged, they will be regenerated automatically.

## Run the Go example

examples/go/main.go provides an example of how these protos can be used.

To try it:

```shell
$ make go_run
...
Read statement with predicateType https://example.com/unknownPred2
Predicate fields:{key:"foo"  value:{struct_value:{fields:{key:"bar"  value:{string_value:"baz"}}}}}
```

[^1]: This is especially helpful during transitions between major versions of the spec or predicate.

[protobuf programming practices]: https://protobuf.dev/programming-guides/proto3
[update guidelines]: https://protobuf.dev/programming-guides/proto3/#updating
[protos/]: ../protos/
[semver guidelines]: https://semver.org/#summary
[supported language bindings]: ../protos/README.md#supported-language-bindings
