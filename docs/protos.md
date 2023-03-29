# Spec protobuf definitions

Protobuf definitions for the in-toto spec and some predicates are
provided in the [protos/] directory. Pre-generated Go bindings of those
protos are available in the go/ directory.

**DISCLAIMER**: The protobuf definitions and Golang bindings will not be
considered stable until the v1.1 tagged release. Use at your own risk.

## Pre-requisites

On an Ubuntu-based system, install the following dependencies.

```shell
sudo apt install build-essential protobuf-compiler golang
```

## Regenerating Go proto libraries

[It's typical to keep generated Go code in the repository itself](https://go.dev/doc/articles/go_command#:~:text=and%20then%20check%20those%20generated%20source%20files%20into%20your%20repository)
since it makes users' lives much easier.

Proto libraries should be regenerated & committed after any change to the
proto files:

```shell
$ make protos
$ git commit -asm "update protos"
...
```

## Run the Go example

examples/go/example/main.go provides an example of how these protos can be used.

To try it:

```shell
$ make go_run
...
Read statement with predicateType https://example.com/unknownPred2
Predicate fields:{key:"foo"  value:{struct_value:{fields:{key:"bar"  value:{string_value:"baz"}}}}}
```

[protos/]: https://github.com/in-toto/attestation/tree/main/protos
