package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/noahkw/gohealthi/pkg/models"
	"log"
	"net"
	"time"

	"github.com/noahkw/gohealthi/pkg/healthstats"
	"github.com/noahkw/gohealthi/pkg/ringbuffer"
	health "github.com/noahkw/gohealthi/proto"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 8337, "Server port")

type Server struct {
	health.UnimplementedHealthServer
	systemUsageQueue *ringbuffer.Queue[*models.SystemUsage]
}

func (healthServer *Server) GetHealth(_ context.Context, in *health.HealthRequest) (*health.HealthResponse, error) {
	log.Printf("Received %v", in.String())

	lastN := healthServer.systemUsageQueue.GetLastN(int(in.Minutes))

	avgLastN, err := healthstats.SystemUsageMean(lastN)

	if err != nil {
		log.Printf("Error calculating average system usage: %v", err)
		return nil, err

	}

	return models.ToHealthResponse(avgLastN), nil
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

	systemUsageQueue := ringbuffer.NewQueue[*models.SystemUsage](15)

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

	healthServer := &Server{
		systemUsageQueue: systemUsageQueue,
	}

	health.RegisterHealthServer(serv, healthServer)
	log.Printf("listening at %v", listener.Addr())

	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
