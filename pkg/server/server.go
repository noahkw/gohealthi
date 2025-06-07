package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/noahkw/gohealthi/pkg/healthstats"
	health "github.com/noahkw/gohealthi/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = flag.Int("port", 8337, "Server port")

type server struct {
	health.UnimplementedHealthServer
}

func (s *server) GetHealth(_ context.Context, in *health.HealthRequest) (*health.HealthResponse, error) {
	log.Printf("Received %s", in.String())

	currentSystemUsage, err := healthstats.CurrentSystemUsage()
	if err != nil {
		log.Printf("Error getting system usage: %v", err)
		return nil, err
	}

	return &currentSystemUsage, nil
}

func Serve() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))

	if err != nil {
		log.Fatalf("failed to listen on port %v, %v", port, err)
	}

	serv := grpc.NewServer()

	health.RegisterHealthServer(serv, &server{})
	log.Printf("listening at %v", listener.Addr())

	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
