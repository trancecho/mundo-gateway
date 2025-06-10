package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
)

func LimitRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的IP地址
		ip := c.ClientIP()

		if domain.LimiterGlobal.IsBlackListed(ip) {
			util.ServerError(c, util.RateLimitExceeded, "请联系管理员解封")
			c.Abort()
			return
		}
		// 检查限流器
		if !domain.LimiterGlobal.AllowIp(ip) {
			util.ServerError(c, util.RateLimitExceeded, "请求过于频繁")
			c.Abort()
			return
		}

		c.Next()
	}
}
