go_setup:
	go install google.golang.org/protobuf/cmd/protoc-gen-go

protos: go_setup
	make -C protos go
	make -C protos python
	make -C protos java

go_run:
	go run examples/go/main.go

go_test: go_setup
	go test ./...

py_test:
	python3 -m unittest tests/python/test_*.py

.PHONY: protos go_setup go_run go_test py_test
