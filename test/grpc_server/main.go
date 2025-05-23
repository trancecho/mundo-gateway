package main

import (
	"context"
	pb "github.com/trancecho/mundo-gateway/test/grpc_b/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

type serverB struct {
	pb.UnimplementedPingServiceServer
}

func (s *serverB) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "pong",
	}, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterPingServiceServer(server, &serverB{})
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Println("failed to listen:", err)
		return
	}
	if err := server.Serve(listener); err != nil {
		log.Println("failed to serve:", err)
		return
	}
}
