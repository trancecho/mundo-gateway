package limiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/trancecho/mundo-gateway/domain"
	"log"
	"sync"
	"time"
)

type RateLimiter struct {
	lastAccess time.Time
	tokens     float64 //当前令牌数
	rate       float64 //每秒产生令牌数
	capacity   float64 //桶的容量
}

type AccessLimiter struct {
	redisClient    *redis.Client
	rateLimiters   map[string]*RateLimiter //IP对应的限流器
	limiterMutex   sync.RWMutex            //限流器读写锁
	globalRate     int                     //全局默认请求速率
	globalCapacity int                     //全局默认容量
	blackListCache map[string]bool         //黑名单缓存
	blackListKey   string                  //黑名单键名
	cacheMutex     sync.RWMutex            //缓存读写锁
	cacheTTL       time.Duration           //缓存有效期
	lastSync       time.Time               //上次同步时间
}

func NewAccessLimiter(globalRate, globalCapacity int) *AccessLimiter {
	// 创建Redis客户端
	rdb := domain.InitRedisClient()

	return &AccessLimiter{
		redisClient:    rdb,
		rateLimiters:   make(map[string]*RateLimiter),
		blackListCache: make(map[string]bool),
		globalRate:     globalRate,
		globalCapacity: globalCapacity,
		blackListKey:   "gateway:blacklist",
		lastSync:       time.Now(),
	}
}

// IsBlackListed 检查IP是否在黑名单中
func (A *AccessLimiter) IsBlackListed(ip string) bool {
	//检查本地缓存
	A.cacheMutex.RLock()
	cached, exists := A.blackListCache[ip]
	if exists {
		A.cacheMutex.RUnlock()
		return cached
	}
	A.cacheMutex.RUnlock()
	//如果本地缓存不存在，查询Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	member, ok := A.redisClient.SIsMember(ctx, A.blackListKey, ip).Result()
	if ok != nil {
		//检测失败，默认放行
		log.Printf("Error checking blacklist for IP %s: %v\n", ip, ok)
		return false //考虑中
	}
	//更新本地缓存
	A.cacheMutex.Lock()
	A.blackListCache[ip] = member
	A.cacheMutex.Unlock()
	return member
}

// AddToBlackList 添加IP到黑名单
func (A *AccessLimiter) AddToBlackList(ip string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 添加到Redis黑名单
	if err := A.redisClient.SAdd(ctx, A.blackListKey, ip).Err(); err != nil {
		log.Printf("Error adding IP %s to blacklist: %v\n", ip, err)
		return false
	}
	// 更新本地缓存
	A.cacheMutex.Lock()
	A.blackListCache[ip] = true
	A.cacheMutex.Unlock()
	return true
}

// SyncBlackList 同步黑名单到本地缓存
func (A *AccessLimiter) SyncBlackList() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取Redis黑名单成员
	members, err := A.redisClient.SMembers(ctx, A.blackListKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get blacklist from Redis: %v", err)
	}

	newCache := make(map[string]bool)
	for _, member := range members {
		newCache[member] = true
	}

	A.cacheMutex.Lock()
	A.blackListCache = newCache
	// 更新上次同步时间
	A.lastSync = time.Now()
	A.cacheMutex.Unlock()

	return nil
}

// AllowRequest 限流检查
func (A *AccessLimiter) AllowRequest(ip string) bool {
	A.limiterMutex.Lock()
	defer A.limiterMutex.Unlock()

	limiter, ok := A.rateLimiters[ip]
	if !ok {
		// 创建新的限流器
		limiter = &RateLimiter{
			lastAccess: time.Now(),
			tokens:     float64(A.globalCapacity),
			rate:       float64(A.globalRate),
			capacity:   float64(A.globalCapacity),
		}
		A.rateLimiters[ip] = limiter
	}

	now := time.Now()
	elapsed := now.Sub(limiter.lastAccess).Seconds()
	limiter.tokens += elapsed * limiter.rate
	if limiter.tokens > limiter.capacity {
		limiter.tokens = limiter.capacity // 限制令牌数不超过容量
	}
	limiter.lastAccess = now

	//检查令牌
	if limiter.tokens > 1 {
		limiter.tokens -= 1 // 消耗一个令牌
		return true
	}
	return false
}

func (A *AccessLimiter) StartCacheRefresher(duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for range ticker.C {
			if err := A.SyncBlackList(); err != nil {
				log.Printf("Error syncing blacklist: %v\n", err)
			} else {
				log.Println("BlackList cache refreshed successfully")
			}
		}
	}()
}
