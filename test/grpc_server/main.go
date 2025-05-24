package main

import (
	"context"
	sdk "github.com/trancecho/mundo-gateway-sdk"
	"github.com/trancecho/mundo-gateway/test/ping/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	gatewayClient := sdk.NewGatewaySDK("ping", "grpc://localhost:50052", "grpc", "http://localhost:12388")
	gatewayClient.RegisterServiceAddress()
	gatewayClient.StartHeartbeat()
	server := grpc.NewServer()
	grpcpingv1.RegisterPingServiceServer(server, &serverB{})
	reflection.Register(server)
	err2 := gatewayClient.AutoRegisterGRPCRoutes(server, "ping")
	if err2 != nil {
		log.Println("failed to register grpc routes:", err2)
	}

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
