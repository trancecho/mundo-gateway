package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/domain"
	"log"
)

func init() {
	log.Println("gateway启动！！！！！")
}

func main() {
	domain.GatewayGlobal = domain.NewGateway()
	log.Println(domain.GatewayGlobal)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Any("/*path", controller.HandleRequestController)
	// 启动服务器
	log.Fatal(r.Run(":12388"))
}
