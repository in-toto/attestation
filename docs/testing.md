# in-toto attestation implementation tests

We provide a set of basic tests for the different language
bindings for the in-toto attestation layers.

## Testing the Go bindings

The go packages `go/v1` and `go/predicates` provide a number of tests
for the statement and predicate layers.

### Running the Go tests

To run all tests:

```shell
make go_test
```

### Writing new Go tests

Please use the standard [Golang testing package] to write tests
for new predicates. For example tests, take a look at the `*_test.go`
files in the `go/` directory tree.

At a minimum, we suggest testing JSON marshalling and unmarshalling
of the Go language bindings.
