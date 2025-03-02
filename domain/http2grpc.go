package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/trancecho/mundo-gateway/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"io"
	"sync"
	"time"
)

//type RegisterServiceBO struct {
//	Name     string `json:"name"`
//	Prefix   string `json:"prefix"`
//	Protocol string `json:"protocol"`
//	Address  string `json:"address"`
//}
//
//func RegisterServiceService(c *gin.Context, bo RegisterServiceBO) bool {
//	//var svc po.ServiceBO
//	//GatewayGlobal.DB.Find(&svc
//	GatewayGlobal.DB.Create(&po.ServiceBO{
//		Name:     bo.Name,
//		Prefix:   bo.Prefix,
//		Protocol: bo.Protocol,
//	})
//	return true
//}

type MethodCache struct {
	sync.RWMutex
	entries map[string]*methodMeta // key: grpcServer/service/method
}

type methodMeta struct {
	methodDesc  protoreflect.MethodDescriptor
	lastUpdated time.Time
}

var (
	methodCache = &MethodCache{
		entries: make(map[string]*methodMeta),
	}
	grpcTimeout = 5 * time.Second
	connPool    = struct {
		sync.RWMutex
		conns map[string]*grpc.ClientConn // key: grpcServer
	}{
		conns: make(map[string]*grpc.ClientConn),
	}
)

func GRPCProxyHandler(c *gin.Context, grpcMethod string, grpcService string, grpcServer string) {
	cacheKey := fmt.Sprintf("%s/%s/%s", grpcServer, grpcService, grpcMethod)

	// 尝试从缓存获取方法元数据
	if meta, ok := getMethodFromCache(cacheKey); ok {
		conn, err := getOrCreateConn(grpcServer)
		if err != nil {
			util.ServerError(c, 500, "连接失败")
			return
		}

		if handleCachedRequest(c, conn, meta.methodDesc, grpcService, grpcMethod) {
			return
		}
		// 调用失败则清除缓存条目
		removeFromCache(cacheKey)
	}

	// 缓存未命中或调用失败，重新反射获取
	conn, err := getOrCreateConn(grpcServer)
	if err != nil {
		util.ServerError(c, 500, "连接失败")
		return
	}

	methodDesc, err := reflectMethodDescriptor(conn, grpcService, grpcMethod)
	if err != nil {
		util.ServerError(c, 404, "方法不存在")
		return
	}

	// 更新缓存
	updateCache(cacheKey, methodDesc)

	// 处理请求
	handleCachedRequest(c, conn, methodDesc, grpcService, grpcMethod)
}

// 缓存操作函数
func getMethodFromCache(key string) (*methodMeta, bool) {
	methodCache.RLock()
	defer methodCache.RUnlock()
	entry, exists := methodCache.entries[key]
	return entry, exists
}

func updateCache(key string, desc protoreflect.MethodDescriptor) {
	methodCache.Lock()
	defer methodCache.Unlock()
	methodCache.entries[key] = &methodMeta{
		methodDesc:  desc,
		lastUpdated: time.Now(),
	}
}

func removeFromCache(key string) {
	methodCache.Lock()
	defer methodCache.Unlock()
	delete(methodCache.entries, key)
}

// 连接池管理
func getOrCreateConn(grpcServer string) (*grpc.ClientConn, error) {
	connPool.RLock()
	conn, exists := connPool.conns[grpcServer]
	connPool.RUnlock()

	if exists {
		state := conn.GetState()
		if state == connectivity.Ready || state == connectivity.Idle {
			return conn, nil
		}
		// 清理无效连接
		conn.Close()
		connPool.Lock()
		delete(connPool.conns, grpcServer)
		connPool.Unlock()
	}

	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()

	newConn, err := grpc.DialContext(ctx, grpcServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %v", err)
	}

	connPool.Lock()
	connPool.conns[grpcServer] = newConn
	connPool.Unlock()

	return newConn, nil
}

// 反射获取方法描述符
func reflectMethodDescriptor(conn *grpc.ClientConn, serviceName string, methodName string) (protoreflect.MethodDescriptor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()

	refClient := grpc_reflection_v1.NewServerReflectionClient(conn)
	stream, err := refClient.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer stream.CloseSend()

	if err := stream.Send(&grpc_reflection_v1.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: serviceName,
		},
	}); err != nil {
		return nil, err
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	return parseMethodDescriptor(resp, methodName)
}

// 处理缓存命中的请求
func handleCachedRequest(
	c *gin.Context,
	conn *grpc.ClientConn,
	methodDesc protoreflect.MethodDescriptor,
	serviceName string,
	methodName string,
) bool {
	reqMsg := dynamicpb.NewMessage(methodDesc.Input())
	if err := bindGRPCRequest(c, reqMsg); err != nil {
		util.ClientErr(c, 400, "请求参数错误")
		return false
	}

	reply := dynamicpb.NewMessage(methodDesc.Output())
	var header, trailer metadata.MD
	fullMethod := fmt.Sprintf("/%s/%s", serviceName, methodName)

	err := conn.Invoke(c.Request.Context(), fullMethod, reqMsg, reply,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	if err != nil {
		util.HandleGRPCError(c, err)
		return false
	}

	marshaler := protojson.MarshalOptions{UseProtoNames: true}
	data, err := marshaler.Marshal(reply)
	if err != nil {
		util.ServerError(c, 500, "响应序列化失败")
		return false
	}

	c.Data(200, "application/json", data)
	return true
}

// ... parseMethodDescriptor 和 bindGRPCRequest 保持不变 ...

func parseMethodDescriptor(resp *grpc_reflection_v1.ServerReflectionResponse, grpcMethod string) (protoreflect.MethodDescriptor, error) {
	fdResp := resp.GetFileDescriptorResponse()
	if fdResp == nil {
		return nil, errors.New("未获取到文件描述")
	}

	files := protoregistry.Files{}
	for _, fdBytes := range fdResp.FileDescriptorProto {
		fdProto := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(fdBytes, fdProto); err != nil {
			return nil, err
		}
		fd, err := protodesc.NewFile(fdProto, &files)
		if err != nil {
			return nil, err
		}
		files.RegisterFile(fd)
	}

	var methodDesc protoreflect.MethodDescriptor
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			for j := 0; j < svc.Methods().Len(); j++ {
				mth := svc.Methods().Get(j)
				if string(mth.Name()) == grpcMethod {
					methodDesc = mth
					return false // 找到直接退出
				}
			}
		}
		return true
	})

	if methodDesc == nil {
		return nil, errors.New("方法未找到")
	}
	return methodDesc, nil
}

func bindGRPCRequest(c *gin.Context, msg *dynamicpb.Message) error {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	return protojson.Unmarshal(data, msg)
}
