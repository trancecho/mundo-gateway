package initial

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/domain/core/limiter"
	"log"
	"time"
)

// 配置
func InitVarFromConfigGlobal() {
	err := check()
	if err != nil {
		log.Fatalln(err)
	}
	// 初始化配置
	domain.QueryPerMinLimitGlobal = viper.GetInt("limit.qpm")
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
	if domain.LimiterGlobal == nil {
		// 初始化限流器
		domain.LimiterGlobal = limiter.NewAccessLimiter()
		if domain.LimiterGlobal == nil {
			log.Fatal("限流器初始化失败")
		}
		// 每分钟（从其他节点同步。如果是本地的话，会立刻刷新）
		domain.LimiterGlobal.StartCacheRefresher(time.Minute * 1)
		domain.LimiterGlobal.StartIpRateRecorderFlusher()
	} else {
		log.Println("限流器已存在，跳过初始化")
	}
}
