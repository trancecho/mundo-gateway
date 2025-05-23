package main

import (
	"context"
	pb "github.com/trancecho/mundo-gateway/test/ping/v1"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	dial, err := grpc.NewClient("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("no", err)
	}
	defer dial.Close()
	client := pb.NewPingServiceClient(dial)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println("Ping response: ", r.GetMessage())
}
