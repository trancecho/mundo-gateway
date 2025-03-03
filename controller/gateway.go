package controller

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
)

// 创建反向代理
func createReverseProxy(target string) (*httputil.ReverseProxy, error) {
	urlx, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(urlx), nil
}

func HandleRequestController(c *gin.Context) {
	var err error
	//比如请求的是 /api/v1/user/1，prefix是 /api/v1
	path := c.Request.URL.Path
	method := c.Request.Method
	// 获得地址和path对
	var prefix string
	// 用Prefix列表匹配找到对应的服务
	// todo 可以进行o1优化，用map存储
	for _, curPrefix := range domain.GatewayGlobal.Prefixes {
		if len(path) < len(curPrefix.Name) {
			continue
		}
		// 去除前缀
		if path[:len(curPrefix.Name)] == curPrefix.Name {
			path = path[len(curPrefix.Name):]
			prefix = curPrefix.Name
			break
		}
	}
	c.Request.URL.Path = path
	fmt.Println("访问路由：", "c.Request.URL.Path:", c.Request.URL.Path, "method:", method, "prefix:", prefix)
	// 尝试直接从缓存拿服务
	//var servicePO po.Service
	// 找到可用服务地址
	//affected := domain.GatewayGlobal.DB.Preload("APIs").Preload("Addresses").First(&servicePO, "prefix = ?", prefix).RowsAffected
	//if affected == 0 {
	//	log.Println("未找到注册服务记录")
	//	util.ServerError(c, 1, "未找到注册服务记录")
	//	return
	//}
	//log.Println("servicePO", servicePO)

	var serviceBO domain.ServiceBO
	log.Println("services", domain.GatewayGlobal.Services)
	for _, bo := range domain.GatewayGlobal.Services {
		// 一个prefix只存在一个服务
		if bo.Prefix == prefix {
			serviceBO = bo
			break
		}
	}
	var apiBO *domain.APIBO
	// 寻找服务方法和路由都匹配的API，如果没有就拦截
	for _, api := range serviceBO.APIs {
		log.Println("api", api, "method", method, "path", path)
		if api.Method == method && api.Path == path {
			apiBO = &api
			break
		}
	}
	//log.Println(serviceBO)
	if apiBO == nil {
		util.ClientErr(c, 3, "未找到API记录")
		return
	}

	// 获得下一个地址
	address := serviceBO.GetNextAddress()

	log.Println("address", address)
	switch serviceBO.Protocol {
	case "http":
		domain.HTTPProxyHandler(c, err, address, serviceBO.Name)
	case "grpc":
		//domain.GRPCProxyHandler(c, err, address)
	default:
		util.ServerError(c, 4, "未知协议")
	}
}

func InitGateway() {
	domain.GatewayGlobal = domain.NewGateway()
	//log.Println("初始化网关", domain.GatewayGlobal)
}
