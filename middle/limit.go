package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain"
)

func LimitRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的IP地址
		ip := c.ClientIP()

		if domain.Limiter.IsBlackListed(ip) {
			c.JSON(500, gin.H{
				"error": "黑名单用户",
			})
			c.Abort()
			return
		}
		// 检查限流器
		if !domain.Limiter.AllowRequest(ip) {
			domain.Limiter.AddToBlackList(ip)
			c.JSON(500, gin.H{
				"error": "请求限流，已进入黑名单",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
