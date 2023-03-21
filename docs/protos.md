# Spec protobuf definitions

Protobuf definitions for the in-toto spec and some predicates are
provided in the spec/ directory.  Pre-generated Go implementation of those
protos are available in the go/ directory.

## Pre-requisites

On an Ubuntu-based system, install the following dependencies.

```shell
sudo apt install protobuf-compiler golang
```

## Regenerating Go proto libraries

[It's typical to keep generated Go code in the repository itself](https://go.dev/doc/articles/go_command#:~:text=and%20then%20check%20those%20generated%20source%20files%20into%20your%20repository)
since it makes users' lives much easier.

Proto libraries should be regenerated & commited after any change to the
proto files:

```shell
$ make go_protos
$ git commit -asm "update protos"
...
```

## Run the Go example

go/example/main.go provides an example of how these protos can be used.

To try it:

```shell
$ make run
...
Read statement with predicateType https://example.com/unknownPred2
Predicate fields:{key:"foo"  value:{struct_value:{fields:{key:"bar"  value:{string_value:"baz"}}}}}
```
