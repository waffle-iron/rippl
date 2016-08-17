package main

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/ripplxyz/rippl/lib/rpc"
	"github.com/ventu-io/slog/basic"
	"github.com/ventu-io/slog"
	"github.com/ventu-io/slf"

	"fmt"
)

var pluginRegistry map[string]chan rpc.Event

func init() {
	bh := basic.New()
	lf := slog.New()

	lf.SetLevel(slf.LevelDebug)

	lf.AddEntryHandler(bh)
	slf.Set(lf)
}

func logger() slf.StructuredLogger {
	return slf.WithContext("rippl")
}

func main() {
	pluginRegistry = make(map[string]chan rpc.Event)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	srv := server{}
	rpc.RegisterRipplServer(grpcServer, &srv)
	logger().Infof("Server listening on %s", lis.Addr().String())
	grpcServer.Serve(lis)
}

type server struct {
}

func (s *server) GetEvents(plugin *rpc.Plugin, stream rpc.Rippl_GetEventsServer) error {
	bus := make(chan rpc.Event)
	pluginRegistry[plugin.Name] = bus

	logger().Debugf("%s connected", plugin.Name)
	defer cleanupPlugin(plugin.Name)
	for {
		select {
		case <- stream.Context().Done():
			return stream.Context().Err()
		case e := <- bus:
			stream.Send(&e)
		}

	}
	return nil
}

func cleanupPlugin(pluginName string) {
	logger().Infof("Cleaning up %s", pluginName)
	delete(pluginRegistry, pluginName)
}
