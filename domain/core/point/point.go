package point

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/proto/point/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"strings"
	"sync"
	"time"
)

// InitGlobalPoints 提供一个初始化函数，返回接口实现
func InitGlobalPoints() {
	domain.PointsGlobal = NewPointClientServer()
}

type PointClientServer struct {
	PointClient     point.UserServiceClient
	PointConn       *grpc.ClientConn
	PointClientOnce sync.Once
	PointServerAddr string
	PointTimeout    time.Duration
}

func NewPointClientServer() *PointClientServer {
	//service, ok := domain.GetServiceByName("point")
	//if !ok {
	//	err = fmt.Errorf("未查找到积分服务")
	//	return err
	//}
	//if service.Available == false {
	//	err = fmt.Errorf("积分服务不可用")
	//	return err
	//}
	//// 检查地址列表是否为空
	//if len(service.Addresses) == 0 {
	//	return fmt.Errorf("服务 %s 没有可用地址", service.Name)
	//}
	var Client PointClientServer
	serverAddr := "nsnhqnldwwab.sealosbja.site:443"
	Client.PointTimeout = time.Second * 10
	Client.PointClientOnce.Do(func() {
		Client.PointServerAddr = serverAddr

		// 从地址中去除 grpcs:// 前缀（如果存在）
		address := Client.PointServerAddr
		address = strings.TrimPrefix(address, "grpcs://")

		// 确保地址包含端口号
		if !strings.Contains(address, ":") {
			address = address + ":443" // 默认 HTTPS 端口
		}

		// 使用 TLS 凭证连接
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false, // 生产环境应设为 false 以验证证书
		})

		pointConn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
		if err == nil {
			Client.PointConn = pointConn
			Client.PointClient = point.NewUserServiceClient(pointConn)
		}
	})
	return &Client
}

func (p *PointClientServer) DialToPointsSystem(token string, userID string, points int64, experience int64, reason string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.PointTimeout)
	defer cancel()
	//添加metadata
	md := metadata.New(map[string]string{
		"Authorization": token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	// 发送远程调用请求
	r, err := p.PointClient.UpdatePointsAndExperience(ctx, &point.UpdatePointsRequest{
		UserId:          userID,
		DeltaPoints:     points,
		DeltaExperience: experience,
		Reason:          reason,
	})
	if err != nil {
		return "", fmt.Errorf("远程调用积分系统失败: %v", err)
	}

	return r.Message, nil
}

func (p *PointClientServer) SignToPointsSystem(token string, userID string) (*point.CommonResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.PointTimeout)
	defer cancel()
	//添加metadata
	md := metadata.New(map[string]string{
		"Authorization": token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	// 发送远程调用请求
	r, err := p.PointClient.Sign(ctx, &point.SignRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("远程调用积分系统失败: %v", err)
	}
	return r, nil
}

func (p *PointClientServer) GetPointsInfo(token string, userID string) (*point.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.PointTimeout)
	defer cancel()
	//添加metadata
	md := metadata.New(map[string]string{
		"Authorization": token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	// 发送远程调用请求
	r, err := p.PointClient.GetUserInfo(ctx, &point.GetUserInfoRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("远程调用积分系统失败: %v", err)
	}
	return r, nil
}

func (p *PointClientServer) AdminStats(token string, userID string) (*point.AdminStats, error) {
	// 创建带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), p.PointTimeout)
	defer cancel()
	//添加metadata
	md := metadata.New(map[string]string{
		"Authorization": token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	// 发送远程调用请求
	r, err := p.PointClient.GetAdminStats(ctx, &point.GetUserInfoRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("远程调用积分系统失败: %v", err)
	}
	return r, nil
}

//// DialToPointsSystem 调用积分系统更新用户积分和经验
//func DialToPointsSystem(token string, userID string, points int64, experience int64, reason string) (string, error) {
//	// 创建带有超时的上下文
//	ctx, cancel := context.WithTimeout(context.Background(), pointTimeout)
//	defer cancel()
//	//添加metadata
//	md := metadata.New(map[string]string{
//		"Authorization": token,
//	})
//	ctx = metadata.NewOutgoingContext(ctx, md)
//	// 发送远程调用请求
//	r, err := pointClient.UpdatePointsAndExperience(ctx, &point.UpdatePointsRequest{
//		UserId:          userID,
//		DeltaPoints:     points,
//		DeltaExperience: experience,
//		Reason:          reason,
//	})
//	if err != nil {
//		return "", fmt.Errorf("远程调用积分系统失败: %v", err)
//	}
//
//	return r.Message, nil
//}

//func SignToPointsSystem(token string, userID string) (*point.CommonResponse, error) {
//	// 创建带有超时的上下文
//	ctx, cancel := context.WithTimeout(context.Background(), pointTimeout)
//	defer cancel()
//	//添加metadata
//	md := metadata.New(map[string]string{
//		"Authorization": token,
//	})
//	ctx = metadata.NewOutgoingContext(ctx, md)
//	// 发送远程调用请求
//	r, err := pointClient.Sign(ctx, &point.SignRequest{
//		UserId: userID,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("远程调用积分系统失败: %v", err)
//	}
//	return r, nil
//}

//func GetPointsInfo(token string, userID string) (*point.UserInfo, error) {
//	// 创建带有超时的上下文
//	ctx, cancel := context.WithTimeout(context.Background(), pointTimeout)
//	defer cancel()
//	//添加metadata
//	md := metadata.New(map[string]string{
//		"Authorization": token,
//	})
//	ctx = metadata.NewOutgoingContext(ctx, md)
//	// 发送远程调用请求
//	r, err := pointClient.GetUserInfo(ctx, &point.GetUserInfoRequest{
//		UserId: userID,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("远程调用积分系统失败: %v", err)
//	}
//	return r, nil
//}

//func AdminStats(token string, userID string) (*point.AdminStats, error) {
//	// 创建带有超时的上下文
//	ctx, cancel := context.WithTimeout(context.Background(), pointTimeout)
//	defer cancel()
//	//添加metadata
//	md := metadata.New(map[string]string{
//		"Authorization": token,
//	})
//	ctx = metadata.NewOutgoingContext(ctx, md)
//	// 发送远程调用请求
//	r, err := pointClient.GetAdminStats(ctx, &point.GetUserInfoRequest{
//		UserId: userID,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("远程调用积分系统失败: %v", err)
//	}
//	return r, nil
//}
