# Creates the go libs in the go/ directory.
go:
	protoc --go_out=../go --go_opt=module=github.com/in-toto/attestation/go $(shell find ./in_toto_attestation -name "*.proto")

# Creates the python libs in the python/ directory.
python:
	protoc --python_out=../python $(shell find ./in_toto_attestation -name "*.proto")

# Creates the java libs in the java/ directory.
java:
	protoc --java_out=../java $(shell find ./in_toto_attestation -name "*.proto")

.PHONY: go python java
