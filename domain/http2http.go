package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/util"
	"io"
	"log"
	"net/http"
)

func HTTPProxyHandler(c *gin.Context, err error, address string, serviceName string) {
	// 手动实现http代理请求
	// 构造代理请求，输入参数为：方法、地址、请求体
	proxyReq, err := http.NewRequest(c.Request.Method, address+c.Request.URL.Path, c.Request.Body)
	if err != nil {
		util.ServerError(c, 2, "创建请求失败")
		return
	}
	// 请求的query参数
	proxyReq.URL.RawQuery = c.Request.URL.RawQuery

	// 将原始请求头复制到代理请求
	for key, values := range c.Request.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// 创建 HTTP 客户端来发送代理请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		util.ServerError(c, 3, "代理请求失败。访问服务名称："+serviceName)
		return
	}
	defer resp.Body.Close()

	// 将目标服务的响应头转发到客户端
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 将目标服务的状态码转发到客户端
	c.Status(resp.StatusCode)

	// 将响应体转发到客户端
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		util.ServerError(c, 4, "转发响应体失败")
		return
	}

	// todo 日志记录：成功完成代理请求
	log.Println("成功代理请求", c.Request.URL, "到", address+c.Request.URL.Path)
}
