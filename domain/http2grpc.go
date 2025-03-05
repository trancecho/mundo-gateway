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
	"log"
	"sync"
	"time"
)

type GrpcConnPool struct {
	sync.RWMutex
	conns map[string]*grpc.ClientConn // key: grpcServer
}

type MethodDescCache struct {
	sync.RWMutex
	cache map[string]protoreflect.MethodDescriptor
}

var (
	grpcTimeout = 5 * time.Second
	connPool    = &GrpcConnPool{
		conns: make(map[string]*grpc.ClientConn),
	}
	methodDescCache = &MethodDescCache{
		cache: make(map[string]protoreflect.MethodDescriptor),
	}
)

func GRPCProxyHandler(c *gin.Context, address string, apibo *APIBO) {
	//address为grpc://ip:port的格式，需要去掉前缀
	address = address[7:]
	grpcService := apibo.GrpcMethodMeta.ServiceName
	grpcMethod := apibo.GrpcMethodMeta.MethodName

	// 生成缓存key
	cacheKey := fmt.Sprintf("%s/%s/%s", address, grpcService, grpcMethod)
	log.Println("cacheKey:", cacheKey)

	// 尝试从缓存获取
	methodDescCache.RLock()
	methodDesc, exists := methodDescCache.cache[cacheKey]
	methodDescCache.RUnlock()
	// 获取连接
	conn, err := getOrCreateConn(address)
	if err != nil {
		util.ServerError(c, 500, "连接失败")
		return
	}

	if !exists {
		// 通过反射获取方法描述符
		methodDesc, err = reflectMethodDescriptor(conn, grpcService, grpcMethod)
		if err != nil {
			log.Println("反射获取方法描述符失败:", err)
			return
		}

		// 存入缓存
		methodDescCache.Lock()
		methodDescCache.cache[cacheKey] = methodDesc
		methodDescCache.Unlock()
	}

	// 处理请求
	handleCachedRequest(c, conn, methodDesc, grpcService, grpcMethod)
}

// 缓存操作函数
//func getMethodFromCache(apiId string) (*methodMeta, bool) {
//	methodCache.RLock()
//	defer methodCache.RUnlock()
//	entry, exists := methodCache.entries[key]
//	return entry, exists
//}

//func updateCache(ApiId int64, desc protoreflect.MethodDescriptor) {
//	GatewayGlobal.RWMutex.Lock()
//	defer GatewayGlobal.RWMutex.Unlock()
//	//根据 apiId 更新
//	for i := range GatewayGlobal.Services {
//		for j := range GatewayGlobal.Services[i].APIs {
//			if GatewayGlobal.Services[i].APIs[j].Id == ApiId {
//				GatewayGlobal.Services[i].APIs[j].GrpcMethodMeta.MethodDesc = desc
//				GatewayGlobal.Services[i].APIs[j].GrpcMethodMeta.LastUpdated = time.Now()
//			}
//		}
//	}
//}

//func removeFromCache(key string) {
//	methodCache.Lock()
//	defer methodCache.Unlock()
//	delete(methodCache.entries, key)
//}

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
		util.ClientError(c, 400, "请求参数错误")
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
