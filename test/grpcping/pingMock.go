package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/sdk"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/reflection"
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
	return &pb.PingResponse{Message: "Pong through GRPC: " + in.Message}, nil
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
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())

	var g errgroup.Group
	g.Go(func() error {
		log.Println("gRPC server started on :50051")
		return s.Serve(lis)
	})

	// HTTP（Gin）服务
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	g.Go(func() error {
		log.Println("HTTP server started on :12389")
		return r.Run(":12389")
	})

	// 初始化 SDK
	gatewaySDK := sdk.NewGatewaySDK("http://localhost:12388")

	// 自动注册 Gin 路由到网关
	if err := gatewaySDK.AutoRegisterGinRoutes(r, "ping_v2"); err != nil {
		fmt.Println("Failed to auto-register routes:", err)
	} else {
		fmt.Println("All routes auto-registered successfully")
	}

	err = gatewaySDK.AutoRegisterGRPCRoutes(s, "grpc_ping_v2")
	if err != nil {
		fmt.Println("Failed to auto-register gRPC routes:", err)
	} else {
		fmt.Println("All gRPC routes auto-registered successfully")
	}

	// 等待任意一个出错退出
	if err := g.Wait(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
