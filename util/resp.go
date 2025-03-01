package util

import (
	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, message string, data gin.H) {
	c.JSON(200, gin.H{
		"errCode": 20000,
		"message": message,
		"data":    data,
	})
}

// 规范：错误码提供三位，从000开始
func ClientErr(c *gin.Context, errCode int, message string) {
	c.JSON(400, gin.H{
		"errCode": 400000 + errCode,
		"message": message,
	})
}

func ServerError(c *gin.Context, errCode int, message string) {
	c.JSON(500, gin.H{
		"errCode": 500000 + errCode,
		"message": message,
	})
}
