package main

import (
	"context"
	gatewaySdk "github.com/trancecho/mundo-gateway-sdk"
	pb "github.com/trancecho/mundo-gateway/test/ping/v1"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	// 目标：我怎么请求gateway去拿到这个target
	target, err := gatewaySdk.NewClient("http://localhost:12388").GetTarget("ping")
	if err != nil {
		log.Println("get target error:", err)
		return
	}
	dial, err := grpc.NewClient(target, grpc.WithInsecure())
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
