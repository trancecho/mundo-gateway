package main

import (
	"context"
	"log"
	"net"

	pb "github.com/trancecho/mundo-gateway/test/grpcping/v1"
	"google.golang.org/grpc"
)

// PingServer 实现了 pb.PingServiceServer 接口
type PingServer struct {
	pb.UnimplementedPingServiceServer
}

// Ping 实现了 Ping 方法
func (s *PingServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "Pong: " + in.Message}, nil
}

func main() {
	// 监听本地端口 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建 gRPC 服务器
	s := grpc.NewServer()
	// 注册 PingService 服务
	pb.RegisterPingServiceServer(s, &PingServer{})
	log.Printf("server listening at %v", lis.Addr())
	// 启动服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
