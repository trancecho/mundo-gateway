package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller"
	"log"
)

func init() {
	log.Println("gateway启动！！！！！")
}

func main() {
	controller.InitGateway()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	MakeRoutes(r)

	// 启动服务器
	log.Fatal(r.Run(":12388"))
}
