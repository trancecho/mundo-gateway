package domain

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"time"
)

// 初始化 Redis 并自动 ping 测试
func InitRedisClient() *redis.Client {
	if viper.IsSet("redis.addr") == false {
		log.Println("❌ redis.addr未配置，请检查配置文件")
		return nil
	}
	if viper.IsSet("redis.password") == false {
		log.Println("❌ redis.password未配置，请检查配置文件")
		return nil
	}
	if viper.IsSet("redis.db") == false {
		log.Println("❌ redis.db未配置，请检查配置文件")
		return nil
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
