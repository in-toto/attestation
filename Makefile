go_setup:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

protos: go_setup
	make -C protos go

go_run:
	go run examples/go/main.go

.PHONY: protos go_setup go_run
