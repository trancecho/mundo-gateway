package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/po"
	"github.com/trancecho/mundo-gateway/util"
	"log"
	"net/http/httputil"
	"net/url"
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
	//比如请求的是 /api/v1/user/1
	path := c.Request.URL.Path
	method := c.Request.Method
	// 获得地址和path对
	var servicePO po.Service
	var prefix string
	// 用Prefix列表匹配找到对应的服务
	// todo 可以进行o1优化，用map存储
	for _, curPrefix := range domain.GatewayGlobal.Prefixes {
		// 去除前缀
		if path[:len(curPrefix.Name)] == curPrefix.Name {
			path = path[len(curPrefix.Name):]
			prefix = curPrefix.Name
			break
		}
	}

	fmt.Println("path", path, "method", method, "prefix", prefix)

	// 找到可用服务地址
	affected := domain.GatewayGlobal.DB.Preload("APIs").Preload("Addresses").First(&servicePO, "prefix = ?", prefix).RowsAffected
	if affected == 0 {
		log.Println("未找到注册服务记录")
		util.ServerError(c, 1, "未找到注册服务记录")
		return
	}
	log.Println(servicePO)
	log.Println("servicePO", servicePO)
	var apiPO *po.API

	// 寻找服务方法和路由都匹配的API
	for _, api := range servicePO.APIs {
		if api.Path == path && api.Method == method {
			apiPO = &api
			break
		}
	}
	log.Println("apiPO", apiPO)
	if apiPO == nil {
		log.Println("未找到API记录", err)
		util.ClientErr(c, 3, "未找到API记录")
		return
	}
	var serviceBO domain.ServiceBO
	// 根据id获取bo
	for _, bo := range domain.GatewayGlobal.Services {
		if bo.ServicePOId == servicePO.ID {
			serviceBO = bo
		}
	}
	// 获得下一个地址
	address := serviceBO.GetNextAddress()

	log.Println("address", address)
	switch servicePO.Protocol {
	case "http":
		domain.HTTPProxyHandler(c, err, address)
	case "grpc":
		//domain.GRPCProxyHandler(c, err, address)
	default:
		log.Println("未知协议")
		util.ServerError(c, 4, "未知协议")
	}
}
