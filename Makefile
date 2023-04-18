go_setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go

protos: go_setup
	make -C protos go

go_run:
	go run examples/go/main.go

.PHONY: protos go_setup go_run
