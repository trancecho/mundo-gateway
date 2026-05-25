package job

import (
	"log"
	"time"

	"github.com/trancecho/mundo-gateway/domain"
)

// StartAPISyncJob 定时将 DB 中的 API 同步到内存，避免依赖重启网关。
func StartAPISyncJob(interval time.Duration) {
	if interval <= 0 {
		interval = 2 * time.Minute
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if domain.GatewayGlobal == nil {
				continue
			}
			domain.GatewayGlobal.SyncAllAPIsFromDB()
		}
	}()
	log.Printf("[gateway] API 定时同步已启动，间隔 %s", interval)
}
