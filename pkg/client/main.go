package main

import (
	"context"
	"flag"
	"log"
	"time"

	health "github.com/noahkw/gohealthi/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("server", "localhost:8000", "server to connect to")
var N = flag.Int("n", 5, "number of minutes to average health stats over")

func main() {
	flag.Parse()

	log.Printf("Querying health stats from %s for the last %d minutes", *addr, *N)

	connection, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}

	defer connection.Close()

	client := health.NewHealthClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.GetHealth(ctx, &health.HealthRequest{Minutes: int32(*N)})
	if err != nil {
		log.Fatalf("could not request health %v", err)
	}
	log.Printf("server health: %s", resp.String())
}
