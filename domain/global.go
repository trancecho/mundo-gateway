package domain

import (
	"github.com/trancecho/mundo-gateway/domain/i"
)

// 全局变量

// 全局变量声明
var (
	GatewayGlobal *Gateway
	PointsGlobal  i.PointClientInterface
	LimiterGlobal i.ILimiter
)

// 全局事件
var ()

var QueryPerMinLimitGlobal int // 单个ip每分钟查询限制
