# Creates the go libs in the go/ directory.
go:
	protoc --go_out=../go --go_opt=module=github.com/in-toto/attestation/go $(shell find ./in_toto_attestation -name "*.proto")

# Creates the python libs in the python/ directory.
python: 
	mkdir -p ../python
	protoc -I ./in_toto_attestation --python_out=../python $(shell find ./in_toto_attestation -name "*.proto")

.PHONY: go
