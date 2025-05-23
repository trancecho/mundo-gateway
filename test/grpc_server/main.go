package main

import (
	"context"
	"github.com/trancecho/mundo-gateway/test/ping/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

type serverB struct {
	grpcpingv1.UnimplementedPingServiceServer
}

func (s *serverB) Ping(ctx context.Context, req *grpcpingv1.PingRequest) (*grpcpingv1.PingResponse, error) {
	return &grpcpingv1.PingResponse{
		Message: "pong",
	}, nil
}

func main() {
	server := grpc.NewServer()
	grpcpingv1.RegisterPingServiceServer(server, &serverB{})
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Println("failed to listen:", err)
		return
	}
	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Println("failed to serve:", err)
		return
	}
}
