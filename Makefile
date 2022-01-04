.PHONY: install protobuf

default: protobuf

install:
	go mod vendor
	go install -mod=mod github.com/golang/protobuf/protoc-gen-go@v1.5.2

protobuf:
	protoc --proto_path=./ --go_out=paths=source_relative:. ./utils/pb/*.proto
	protoc --proto_path=./ --go_out=paths=source_relative:. --rpc_out=paths=source_relative:. ./mem/mem.proto