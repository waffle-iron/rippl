package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/ripplxyz/rippl/lib/rpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	srv := server{bus: make(chan rpc.Event)}
	rpc.RegisterRipplServer(grpcServer, &srv)
	grpcServer.Serve(lis)
}

type server struct {
	bus chan rpc.Event
}

func (s *server) GetEvents(plugin *rpc.Plugin, stream rpc.Rippl_GetEventsServer) error {
	fmt.Println(plugin.Name)
	for i := range s.bus {
		fmt.Println(i.EventType.String())
	}
	return nil
}
