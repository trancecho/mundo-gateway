package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain"
)

func CheckCondition() gin.HandlerFunc {
	return func(c *gin.Context) {
		service := c.GetHeader("Service")
		server, ok := domain.GetServiceByName(service)
		if !ok {
			c.JSON(404, gin.H{
				"error": "Service not found",
			})
			c.Abort()
			return
		}
		if server.Available == false {
			c.JSON(503, gin.H{
				"error": "Service unavailable",
			})
			c.Abort()
			return
		}
		c.Set("ServiceName", service)
		c.Set("Address", server.Addresses)
		c.Next()
	}
}
