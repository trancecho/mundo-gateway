package controller

import (
	"fmt"
	"github.com/trancecho/mundo-gateway/global"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
)

//// 创建反向代理
//func createReverseProxy(target string) (*httputil.ReverseProxy, error) {
//	urlx, err := url.Parse(target)
//	if err != nil {
//		return nil, err
//	}
//	return httputil.NewSingleHostReverseProxy(urlx), nil
//}

func HandleRequestController(c *gin.Context) {
	var err error
	//比如请求的是 /api/v1/user/1，prefix是 /api/v1
	path := c.Request.URL.Path
	method := c.Request.Method
	// 获得地址和path对
	var prefix string

	// 提取路径中第一个 `/` 之后的内容
	pathParts := strings.SplitN(path, "/", 3)
	if len(pathParts) < 3 {
		util.ClientError(c, 100, "路径不合法")
		return
	}
	// 这里的path是去掉了第一个 `/` 的部分
	path = "/" + pathParts[2]
	// 这里的prefix是第一个 `/` 之后的部分
	prefix = "/" + pathParts[1]
	var prefixOkFlag bool
	// 用Prefix列表匹配找到对应的服务
	// todo 可以进行o1优化，用map存储
	for _, curPrefix := range global.GatewayGlobal.Prefixes {
		log.Println("curPrefix", curPrefix)
		log.Println("prefix", prefix)
		if curPrefix.Name == prefix {
			prefixOkFlag = true
			break
		}
	}
	if !prefixOkFlag {
		util.ClientError(c, 200, "prefix不合法或服务挂了")
		return
	}
	c.Request.URL.Path = path
	fmt.Println("访问路由：", "c.Request.URL.Path:", c.Request.URL.Path, "method:", method, "prefix:", prefix)
	//尝试直接从缓存拿服务
	//var servicePO po.Service
	//找到可用服务地址
	//affected := domain.GatewayGlobal.DB.Preload("APIs").Preload("Addresses").First(&servicePO, "prefix = ?", prefix).RowsAffected
	//if affected == 0 {
	//	log.Println("未找到注册服务记录")
	//	util.ServerError(c, 1, "未找到注册服务记录")
	//	return
	//}
	//log.Println("servicePO", servicePO)

	var serviceBO domain.ServiceBO
	log.Println("services", global.GatewayGlobal.Services)
	for _, bo := range global.GatewayGlobal.Services {
		// 一个prefix只存在一个服务
		if bo.Prefix == prefix {
			serviceBO = bo
			break
		}
	}
	var apiBO *domain.APIBO
	// 寻找服务方法和路由都匹配的API，如果没有就拦截(默认grpc服务也有http路由，apibo是http2grpc的映射)
	//
	//但是对于"/:xx/"的路由，需要进行特殊处理
	for _, api := range serviceBO.APIs {
		log.Println("api", api, "method", method, "path", path)

		// 先判断方法
		if api.HttpMethod != method {
			continue
		}

		// 精确匹配优先
		if api.HttpPath == path {
			apiBO = &api
			break
		}

		// 尝试模糊匹配 /user/:id -> /user/123
		if isPathMatch(api.HttpPath, path) {
			apiBO = &api
			break
		}
	}

	//log.Println(serviceBO)
	if apiBO == nil {
		util.ClientError(c, 3, "未找到API记录")
		return
	}

	// 获得下一个地址
	address := serviceBO.GetNextAddress()

	log.Println("address", address)
	switch serviceBO.Protocol {
	case "http":
		domain.HTTPProxyHandler(c, err, address, serviceBO.Name)
	case "grpc":
		domain.GRPCProxyHandler(c, address, apiBO)
	default:
		util.ServerError(c, 4, "未知协议")
	}
}

func InitGateway() {
	global.GatewayGlobal = domain.NewGateway()
	//log.Println("初始化网关", domain.GatewayGlobal)
}

func FlushAPIController(c *gin.Context) {
	// 重新加载API
	global.GatewayGlobal.FlushGateway()
	util.Ok(c, "网关刷新成功", nil)
}

// isPathMatch 用于支持带参数的路径匹配，例如 /user/:id 匹配 /user/123
func isPathMatch(routePath, actualPath string) bool {
	routeParts := strings.Split(routePath, "/")
	pathParts := strings.Split(actualPath, "/")

	if len(routeParts) != len(pathParts) {
		return false
	}

	for i := 0; i < len(routeParts); i++ {
		// 如果是 : 开头，认为是参数位，跳过
		if strings.HasPrefix(routeParts[i], ":") {
			continue
		}
		// 如果实际路径不一致就不匹配
		if routeParts[i] != pathParts[i] {
			return false
		}
	}
	return true
}
