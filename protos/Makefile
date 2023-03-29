# Creates the go libs in the go/ directory.
go:
	protoc --go_out=../go --go_opt=module=github.com/in-toto/attestation/go $(shell find ./in_toto_attestation -name "*.proto")

.PHONY: go
