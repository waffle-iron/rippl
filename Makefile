deps:
	glide install

proto: lib/rpc/*.pb.go
	cd protos && protoc --go_out=plugins=grpc:. event.proto

build: deps proto
	cd app && go build .
