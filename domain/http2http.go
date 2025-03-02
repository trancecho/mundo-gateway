package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/util"
	"io"
	"log"
	"net/http"
)

func HTTPProxyHandler(c *gin.Context, err error, address string) {
	//现在有了地址
	// 手动实现http代理请求
	proxyReq, err := http.NewRequest(c.Request.Method, address+c.Request.URL.Path, c.Request.Body)
	if err != nil {
		log.Println("请求创建失败:", err)
		util.ServerError(c, 2, "创建请求失败")
		return
	}

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
		log.Println("发送代理请求失败:", err)
		util.ServerError(c, 2, "代理请求失败")
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
		log.Println("转发响应体失败:", err)
		util.ServerError(c, 2, "转发响应体失败")
		return
	}

	// todo 日志记录：成功完成代理请求
	//todo需要自动加前缀http://
	log.Println("成功代理请求", c.Request.URL.Path, "到", address)
}
