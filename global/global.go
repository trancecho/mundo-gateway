package global

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/domain/i"
	"log"
)

// 全局变量

// 全局变量声明
var (
	GatewayGlobal *domain.Gateway
	PointsGlobal  i.PointClientInterface
	LimiterGlobal i.ILimiter
)

// 全局事件
var ()

var QueryPerMinLimitGlobal int // 单个ip每分钟查询限制
