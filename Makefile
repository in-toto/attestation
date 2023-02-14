go_setup:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

# Creates the go libs in the go/ directory.
go_protos: go_setup
	protoc --go_out=go --go_opt=paths=source_relative $(shell find ./spec -name "*.proto")

run:
	go run go/example/main.go
