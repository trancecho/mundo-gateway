package initial

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"time"
)

// 初始化 Redis 并自动 ping 测试
func InitRedisClient() *redis.Client {
	if !viper.IsSet("redis.addr") ||
		!viper.IsSet("redis.password") ||
		!viper.IsSet("redis.db") {
		log.Fatalln("❌ Redis 配置未设置，请检查配置文件")
	}
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	rdb := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          db,
		DialTimeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if pong, err := rdb.Ping(ctx).Result(); err != nil {
		log.Printf("❌ Redis 连接失败: %v", err)
		return nil
	} else {
		log.Printf("✅ Redis 连接成功: %s\n", pong)
	}

	return rdb
}
