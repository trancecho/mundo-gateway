package controller

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
)

func HandleRequestController(c *gin.Context) {
	path := c.Request.URL.Path
	method := c.Request.Method

	// 例如 /timerme/api/foo -> prefix=/timerme, path=/api/foo
	pathParts := strings.SplitN(path, "/", 3)
	if len(pathParts) < 3 || pathParts[1] == "" || pathParts[2] == "" {
		util.ClientError(c, 100, "路径不合法")
		return
	}
	path = "/" + pathParts[2]
	prefix := "/" + pathParts[1]

	domain.GatewayGlobal.RWMutex.RLock()
	if !domain.PrefixRegistered(prefix) {
		domain.GatewayGlobal.RWMutex.RUnlock()
		util.ClientError(c, 200, "prefix不合法或服务挂了")
		return
	}
	serviceBO := domain.GetServiceByPrefix(prefix)
	domain.GatewayGlobal.RWMutex.RUnlock()

	if serviceBO == nil {
		domain.ReloadServiceByPrefix(prefix)
		domain.GatewayGlobal.RWMutex.RLock()
		serviceBO = domain.GetServiceByPrefix(prefix)
		domain.GatewayGlobal.RWMutex.RUnlock()
	}
	if serviceBO == nil {
		util.ClientError(c, 201, "服务未加载，请稍后重试或检查服务注册")
		return
	}

	c.Request.URL.Path = path
	fmt.Println("访问路由：", "c.Request.URL.Path:", c.Request.URL.Path, "method:", method, "prefix:", prefix)

	apiBO := domain.ResolveAPI(serviceBO, method, path)
	if apiBO == nil {
		log.Printf("[gateway] 未找到API: %s %s %s (内存API数=%d)",
			serviceBO.Name, method, path, serviceBO.APIs.Size())
		util.ClientError(c, 3, "未找到API记录")
		return
	}

	address := serviceBO.GetNextAddress()
	if address == "" {
		util.ServerError(c, 202, "服务无可用地址，请检查后端是否在线并已发送心跳")
		return
	}

	switch serviceBO.Protocol {
	case "http":
		domain.HTTPProxyHandler(c, nil, address, serviceBO.Name)
	case "grpc":
		domain.GRPCProxyHandler(c, address, apiBO)
	default:
		util.ServerError(c, 4, "未知协议")
	}
}

func InitGateway() {
	domain.GatewayGlobal = domain.NewGateway()
	domain.GatewayGlobal.FlushGateway()
	log.Println("初始化网关", domain.GatewayGlobal)
}

func FlushAPIController(c *gin.Context) {
	domain.GatewayGlobal.FlushGateway()
	util.Ok(c, "网关刷新成功", nil)
}
