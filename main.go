package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/config"
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/domain/core/point"
	"github.com/trancecho/mundo-gateway/initial"
	"github.com/trancecho/mundo-gateway/job"
	"github.com/trancecho/mundo-gateway/middle"
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
	initial.InitVarFromConfigGlobal()

	log.Println(viper.GetString("mysql.host") + ":" + viper.GetString("mysql.port"))

	controller.InitGateway()
	point.InitGlobalPoints()

	// redis
	domain.GatewayGlobal.Redis = initial.InitRedisClient()
	if domain.GatewayGlobal.Redis == nil {
		log.Fatal("Redis 初始化失败，程序终止")
	}
	job.StartPasswordRefreshTask()

	initial.InitLimiterGlobal()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middle.Middleware())
	r.Use(middle.LimitRequest()) // 限流中间件

	routes.MakeRoutes(r)

	go func() {
		controller.ServiceAliveChecker()
	}()
	// 启动服务器
	log.Fatal(r.Run(":" + viper.GetString("server.port")))
}
