package job

import (
	"context"
	"github.com/trancecho/mundo-gateway/domain"
	"log"
	"time"

	"github.com/google/uuid"
)

func StartPasswordRefreshTask() {
	refreshPassword()
	ticker := time.NewTicker(1 * time.Hour) // 每小时刷新一次
	go func() {
		for {
			select {
			case <-ticker.C:
				refreshPassword()
			}
		}
	}()
}

func refreshPassword() {
	ctx := context.Background()
	newPassword := generatePassword()
	err := domain.GatewayGlobal.Redis.Set(ctx, "gateway:register:password", newPassword, 2*time.Hour).Err()
	if err != nil {
		log.Println("刷新注册密码失败:", err)
		return
	}
	log.Println("🔐 注册密码已刷新:", newPassword)
	// 可选：写入数据库、消息队列、通知平台等
}

func generatePassword() string {
	return uuid.NewString()[:8] // 生成 8 位随机密码
}
