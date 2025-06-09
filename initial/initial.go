package initial

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/domain/core/limiter"
	"github.com/trancecho/mundo-gateway/global"
	"log"
)

// 配置
func InitVarFromConfigGlobal() {
	err := check()
	if err != nil {
		log.Fatalln(err)
	}
	// 初始化配置
	global.QueryPerMinLimitGlobal = viper.GetInt("limit.qpm")
}

func check() error {
	var checks []string
	// 手动配置检查项
	checks = []string{
		"limit.qpm",
	}
	for _, x := range checks {
		if !viper.IsSet(x) {
			return errors.New("配置项 " + x + " 未设置")
		}
	}
	return nil
}
func InitLimiterGlobal() {
	if global.LimiterGlobal == nil {
		// 初始化限流器
		global.LimiterGlobal = limiter.NewAccessLimiter(
			viper.GetInt("limit.global_rate"),
			viper.GetInt("limit.global_capacity"),
		)
		if global.LimiterGlobal == nil {
			log.Fatal("限流器初始化失败")
		}
	} else {
		log.Println("限流器已存在，跳过初始化")
	}
}
