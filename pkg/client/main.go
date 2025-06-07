package main

import (
	"context"
	"flag"
	health "github.com/noahkw/gohealthi/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var addr = flag.String("server", "localhost:8000", "server to connect to")

func main() {
	connection, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer connection.Close()

	client := health.NewHealthClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.GetHealth(ctx, &health.HealthRequest{})
	if err != nil {
		log.Fatalf("could not request health %v", err)
	}
	log.Printf("server health: %s", resp.String())
}
