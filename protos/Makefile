# Creates the go libs in the go/ directory.
go:
	# All other Predicates
	protoc --go_out=../go \
		--go_opt=module=github.com/in-toto/attestation/go \
		$(shell find ./in_toto_attestation -name "*.proto")
	# CycloneDX 1.5
	# https://protobuf.dev/reference/go/go-generated/#package
	protoc --go_out=../go \
		--go_opt=Mcyclonedx/schema/bom-1.5.proto=predicates/cyclonedx/v1 \
		./cyclonedx/schema/bom-1.5.proto
	mv ../go/predicates/cyclonedx/v1/bom-1.5.pb.go ../go/predicates/cyclonedx/v1/bom.pb.go

# Creates the python libs in the python/ directory.
python:
	protoc --python_out=../python $(shell find ./in_toto_attestation -name "*.proto")

# Creates the java libs in the java/ directory.
java:
	protoc --java_out=../java $(shell find ./in_toto_attestation -name "*.proto")

.PHONY: go python java
