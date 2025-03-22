package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/config"
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/routes"
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
	log.Println(viper.GetString("mysql.host") + ":" + viper.GetString("mysql.port"))

	controller.InitGateway()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.MakeRoutes(r)

	go func() {
		controller.ServiceAliveChecker()
	}()
	// 启动服务器
	log.Fatal(r.Run(":" + viper.GetString("server.port")))
}
