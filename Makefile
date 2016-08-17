deps:
	glide install

proto:
	cd protos && protoc --go_out=plugins=grpc:. event.proto

build: deps proto
	cd app && go build .
