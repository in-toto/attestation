# Statement protobuf definitions

Protobuf definitions for the in-toto statement type and some predicates are provided in the
spec/ directory.  Pre-generated Go implementation of those protos are available in the go dir.

go/example/main.go provides an example of how these protos can be used.

To try it:

```shell
$ make run
...
Read statement with predicateType https://example.com/unknownPred2
Predicate fields:{key:"foo"  value:{struct_value:{fields:{key:"bar"  value:{string_value:"baz"}}}}}
```

Please consider providing a proto version of any new predicates proposed.

## Regenerating Go proto libraries

[It's typical to keep generated Go code in the repository itself](https://go.dev/doc/articles/go_command#:~:text=and%20then%20check%20those%20generated%20source%20files%20into%20your%20repository)
since it makes users lives much easier.

Proto libraries should be regenerated & commited after any change to the proto files:

```shell
$ make go_protos
go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
protoc --go_out=go --go_opt=paths=source_relative ./spec/predicates/vsa.proto ./spec/v1.0-draft/statement.proto
$ git commit -asm "update protos"
[statement_proto 5edb2c6] Update protos
...
```
