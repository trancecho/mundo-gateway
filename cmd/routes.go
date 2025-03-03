package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/util"
	"strings"
)

func MakeRoutes(r *gin.Engine) {

	// Service
	r.POST("/gateway/service", controller.CreateServiceController)
	r.PUT("/gateway/service", controller.UpdateServiceController)
	r.DELETE("/gateway/service", controller.DeleteServiceController)
	r.GET("/gateway/service", controller.GetServiceController)
	r.GET("/gateway/service/list", controller.ListServiceController)

	// API
	r.POST("/gateway/api", controller.CreateAPIController)
	r.PUT("/gateway/api", controller.UpdateAPIController)
	r.DELETE("/gateway/api", controller.DeleteAPIController)
	r.GET("/gateway/api", controller.GetAPIController)
	r.GET("/gateway/api/list", controller.ListAPIController)

	// 通配符路由，排除 /gateway 前缀
	r.NoRoute(func(c *gin.Context) {
		// 检查路径是否以 /gateway 开头
		if strings.HasPrefix(c.Request.URL.Path, "/gateway") {
			// 如果是 /gateway 前缀，返回 404
			util.ClientErr(c, 404, "Not Found")
			return
		}

		// 否则，执行通配符逻辑
		controller.HandleRequestController(c)
	})
}
