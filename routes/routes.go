package routes

import (
	"github.com/trancecho/mundo-gateway/domain/core/point"
	"github.com/trancecho/mundo-gateway/middle"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/util"
)

func MakeRoutes(r *gin.Engine) {

	// Service
	r.POST("/gateway/service", controller.CreateServiceController)
	r.PUT("/gateway/service", controller.UpdateServiceController)
	r.DELETE("/gateway/service", controller.DeleteServiceController)
	r.DELETE("/gateway/service/address", controller.DeleteServiceAddressController)
	r.GET("/gateway/service", controller.GetServiceController)
	r.GET("/gateway/service/list", controller.ListServiceController)
	r.POST("/gateway/service/beat", controller.ServiceAliveSignalController)

	// API
	r.POST("/gateway/api", controller.CreateAPIController)
	r.PUT("/gateway/api", controller.UpdateAPIController)
	r.DELETE("/gateway/api", controller.DeleteAPIController)
	r.GET("/gateway/api", controller.GetAPIController)
	r.GET("/gateway/api/list", controller.ListAPIController)

	r.GET("/gateway/flush", controller.FlushAPIController)
	r.GET("/gateway/ping", func(c *gin.Context) {
		util.Ok(c, "pong", nil)
	})

	//健康检查接口
	r.GET("/gateway/service/health", controller.HealthStatusHandler)

	//积分服务接口
	s1 := r.Group("/gateway/point", middle.CheckCondition())
	{
		s1.POST("/change", point.ChangePointAndExperience)
		s1.POST("/sign", point.Sign)
		s1.GET("/info", point.PointsInfo)
		s1.GET("/stats", point.GetStats)
	}

	// 通配符路由，排除 /gateway 前缀
	r.NoRoute(func(c *gin.Context) {
		// 检查路径是否以 /gateway 开头
		if strings.HasPrefix(c.Request.URL.Path, "/gateway") {
			// 如果是 /gateway 前缀，返回 404
			util.ClientError(c, 404, "不支持gateway前缀")
			return
		}

		// 否则，执行通配符逻辑
		controller.HandleRequestController(c)
	})
}
