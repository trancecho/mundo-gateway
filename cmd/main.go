package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/config"
	"github.com/trancecho/mundo-gateway/controller"
	"log"
)

func init() {
	log.Println("gateway启动！！！！！")
}

func main() {
	// 先加载配置文件
	config.GlobalConfig = config.NewConfig()
	cfg := config.GlobalConfig
	err := cfg.Init()
	if err != nil {
		log.Fatal("配置文件加载失败", err)
	}
	log.Println(cfg.Mysql)

	controller.InitGateway()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	MakeRoutes(r)

	// 启动服务器
	log.Fatal(r.Run(":12388"))
}
