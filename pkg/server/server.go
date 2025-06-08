package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/noahkw/gohealthi/pkg/healthstats"
	"github.com/noahkw/gohealthi/pkg/ringbuffer"
	health "github.com/noahkw/gohealthi/proto"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 8337, "Server port")
var systemUsageQueue = ringbuffer.NewQueue[health.HealthResponse](15)

type server struct {
	health.UnimplementedHealthServer
}

func (s *server) GetHealth(_ context.Context, in *health.HealthRequest) (*health.HealthResponse, error) {
	log.Printf("Received %s", in.String())

	lastN := systemUsageQueue.GetLastN(int(in.Minutes))

	avgLastN, err := healthstats.SystemUsageMean(lastN)

	if err != nil {
		log.Printf("Error calculating average system usage: %v", err)
		return nil, err

	}

	return avgLastN, nil
}

func runPeriodically(duration time.Duration, callback func()) {
	ticker := time.NewTicker(duration)

	callback()

	go func() {
		for range ticker.C {
			callback()
		}
	}()
}

func Serve() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))

	if err != nil {
		log.Fatalf("failed to listen on port %v, %v", port, err)
	}

	runPeriodically(time.Minute, func() {
		currentStats, err := healthstats.CurrentSystemUsage()

		if err != nil {
			log.Printf("Error getting current system usage: %v", err)
			return
		}

		systemUsageQueue.Add(currentStats)
		log.Printf("current length health responses: %v", systemUsageQueue.Len())
	})

	serv := grpc.NewServer()

	health.RegisterHealthServer(serv, &server{})
	log.Printf("listening at %v", listener.Addr())

	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
