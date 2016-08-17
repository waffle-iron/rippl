.PHONY: default
default: build ;
.DEFAULT_GOAL := build

deps:
	glide install

proto: lib/rpc/*.pb.go
	cd protos && protoc --go_out=plugins=grpc:. event.proto

test:
    true

build: deps
	cd app && go build .

travis: test
