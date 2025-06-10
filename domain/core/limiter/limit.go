package limiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/domain/i"
	"gorm.io/gorm"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type IpRateRecorder struct {
	//lastAccess time.Time
	//tokens     float64 //当前令牌数
	//rate       float64 //每秒产生令牌数
	//capacity   float64 //桶的容量
	cnt        atomic.Int64
	lastAccess atomic.Int64 // 最后访问时间戳，单位为纳秒
}

type AccessLimiter struct {
	redisClient    *redis.Client
	ipRateRecorder map[string]*IpRateRecorder //IP对应的限流器
	limiterLock    sync.RWMutex               //限流器读写锁
	//globalRate     int                     //全局默认请求速率
	//globalCapacity int                     //全局默认容量

	blackListKey   string          //黑名单键名
	blackListCache map[string]bool //黑名单缓存
	blackListLock  sync.RWMutex    //缓存读写锁
	//cacheTTL       time.Duration   //缓存有效期

	//白名单直接读数据库
	whiteListKey   string          //白名单键名
	whiteListCache map[string]bool //白名单缓存
	whiteListLock  sync.RWMutex    //白名单缓存读写锁
	db             *gorm.DB
}

func (this *AccessLimiter) IsWhiteListed(ip string) bool {
	this.whiteListLock.RLock()
	_, ok := this.whiteListCache[ip]
	this.whiteListLock.RUnlock()
	if !ok {
		return false
	} else {
		return true
	}
}

// 实现接口

// 全局一个
func NewAccessLimiter(db *gorm.DB) *AccessLimiter {
	// 创建Redis客户端
	rdb := domain.InitRedisClient()

	return &AccessLimiter{
		redisClient:    rdb,
		ipRateRecorder: make(map[string]*IpRateRecorder),
		blackListCache: make(map[string]bool),
		blackListKey:   "gateway:blacklist",
		whiteListCache: make(map[string]bool),
		db:             db,
		whiteListKey:   "gateway:whitelist",
	}
}

// 检查IP是否在黑名单中
func (this *AccessLimiter) IsBlackListed(ip string) bool {
	//检查本地缓存
	this.blackListLock.RLock()
	_, ok := this.blackListCache[ip]
	this.blackListLock.RUnlock()
	if !ok {
		return false
	} else {
		return true
	}
	//this.blackListLock.RUnlock()
	//如果本地缓存不存在，查询Redis
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//member, err := this.redisClient.SIsMember(ctx, this.blackListKey, ip).Result()
	//if err != nil {
	//	//检测失败，默认放行
	//	log.Printf("Error checking blacklist for IP %s: %v\n", ip, ok)
	//	return false //考虑中
	//}
	////更新本地缓存
	//this.blackListLock.Lock()
	//this.blackListCache[ip] = member
	//this.blackListLock.Unlock()
	//return member
}

// 添加IP到黑名单
func (this *AccessLimiter) AddToBlackList(ip string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 添加到Redis黑名单
	if err := this.redisClient.SAdd(ctx, this.blackListKey, ip).Err(); err != nil {
		log.Printf("Error adding IP %s to blacklist: %v\n", ip, err)
		return false
	}
	// 更新本地缓存
	this.blackListLock.Lock()
	this.blackListCache[ip] = true
	this.blackListLock.Unlock()
	return true
}

// 添加ip到白名单
func (this *AccessLimiter) AddToWhiteList(ip string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 添加到Redis白名单
	if err := this.redisClient.SAdd(ctx, this.whiteListKey, ip).Err(); err != nil {
		log.Printf("Error adding IP %s to whitelist: %v\n", ip, err)
		return false
	}

	// 更新本地缓存
	this.whiteListLock.Lock()
	this.whiteListCache[ip] = true
	this.whiteListLock.Unlock()
	return true
}

// 同步黑名单到本地缓存
func (this *AccessLimiter) SyncBlackList() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取Redis黑名单成员
	members, err := this.redisClient.SMembers(ctx, this.blackListKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get blacklist from Redis: %v", err)
	}

	newCache := make(map[string]bool)
	for _, member := range members {
		newCache[member] = true
	}

	this.blackListLock.Lock()
	this.blackListCache = newCache
	this.blackListLock.Unlock()
	return nil
}

func (this *AccessLimiter) SyncWhiteList() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取Redis白名单成员
	members, err := this.redisClient.SMembers(ctx, this.whiteListKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get whitelist from Redis: %v", err)
	}

	newCache := make(map[string]bool)
	for _, member := range members {
		newCache[member] = true
	}

	this.whiteListLock.Lock()
	this.whiteListCache = newCache
	this.whiteListLock.Unlock()
	return nil
}

func (this *AccessLimiter) StartCacheRefresher(duration time.Duration) {
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := this.SyncBlackList(); err != nil {
					log.Printf("从redis同步黑名单失败: %v\n", err)
				}
				// 同步白名单
				if err := this.SyncWhiteList(); err != nil {
					log.Printf("从数据库同步白名单失败: %v\n", err)
				} else {
					log.Println("同步白名单成功")
				}
				log.Println("同步黑名单白名单成功")
			}
		}
	}()
}

func (this *AccessLimiter) AllowIp(ip string) bool {
	this.blackListLock.RLock()
	if _, ok := this.blackListCache[ip]; ok {
		this.blackListLock.RUnlock()
		log.Println("IP is blacklisted:", ip)
		return false
	}
	this.blackListLock.RUnlock()
	this.limiterLock.Lock()
	recorder, ok := this.ipRateRecorder[ip]
	if !ok {
		// 如果不存在，创建一个新的限流器
		recorder = &IpRateRecorder{
			cnt: atomic.Int64{},
		}
		this.ipRateRecorder[ip] = recorder
	}
	this.limiterLock.Unlock()
	// 打印调试
	//this.limiterLock.Lock()
	//log.Println("aaa", this.ipRateRecorder[ip].cnt.Load())
	//defer this.limiterLock.Unlock()
	// 更新访问时间
	recorder.cnt.Add(1)                              // 增加访问计数
	recorder.lastAccess.Store(time.Now().UnixNano()) // 更新最后访问时间戳
	// 检查访问频率是否超过限制
	if recorder.cnt.Load() > int64(domain.QueryPerMinLimitGlobal) {
		// 如果超过限制，加入黑名单
		if this.AddToBlackList(ip) {
			this.limiterLock.Lock()
			delete(this.ipRateRecorder, ip) // 删除限流器
			this.limiterLock.Unlock()
			log.Println("IP added to blacklist due to rate limit exceeded:", ip)
			return false
		}
	}
	log.Println("IP 访问成功:", ip)
	return true
}
func (this *AccessLimiter) StartIpRateRecorderFlusher() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				this.limiterLock.Lock()
				for ip, recorder := range this.ipRateRecorder {
					// 将cnt置0。若超过十分钟没有访问，则认为该IP不再活跃，删除
					if time.Now().UnixNano()-recorder.lastAccess.Load() > 10*time.Minute.Nanoseconds() {
						delete(this.ipRateRecorder, ip)
						log.Println("IP removed from rate limiter due to inactivity:", ip)
					} else {
						recorder.cnt.Store(0) // 重置访问计数
						log.Println("IP rate recorder flushed:", ip, "Count:", recorder.cnt.Load())
					}
				}
				this.limiterLock.Unlock()
			}
		}
	}()
}

var _ i.ILimiter = (*AccessLimiter)(nil)
