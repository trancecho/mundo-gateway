package domain

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/po"
	"github.com/trancecho/ragnarok/maplist"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GatewayGlobal 第一个功能，代理路由
type Gateway struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Prefixes maplist.MapStringList[Prefix]
	Services maplist.MapList[ServiceBO]

	globalKV sync.Map
	//读写锁
	sync.RWMutex
	HTTPClient *http.Client //

}

// NewGateway,创建一个全局网关实例
func NewGateway() *Gateway {
	// todo 全局网关注册后，还需要时刻更新属性列表，如prefixes，services
	var err error
	// 会自己注册一个地址的。
	var db *gorm.DB
	pwd := viper.GetString("mysql.pwd")
	dsn := "root:" + pwd +
		"@tcp(" + viper.GetString("mysql.host") + ":" + viper.GetString("mysql.port") + ")/" +
		viper.GetString("mysql.db") +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}

	err = db.AutoMigrate(
		&po.Service{},
		&po.API{},
		&po.Address{},
		&po.GrpcMethodMeta{},
	)
	if err != nil {
		log.Fatalln("failed to migrate database", err)
	}
	log.Println("connect database success!!!!!!!!")

	// 初始化 service 列表
	var services []po.Service
	db.Preload("Addresses").
		Preload("APIs.GrpcMethodMeta").
		Where("available=?", true).
		Find(&services)

	log.Println("services:", services)

	// 初始化 prefixes
	prefixes := maplist.NewMapStringList[Prefix]()

	// 在 services 中找到所有的 prefix
	for _, service := range services {
		prefixes.Add(service.Prefix, &Prefix{
			Name:      service.Prefix,
			ServiceId: service.ID,
		})
	}

	serviceBOs := maplist.NewMapList[ServiceBO]()
	for _, service := range services {
		serviceBOs.Add(service.ID, buildServiceBO(service, nil))
	}

	// 初始化 HTTP 客户端
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	gateway := &Gateway{
		DB:         db,
		Services:   serviceBOs,
		Prefixes:   prefixes,
		HTTPClient: httpClient,
	}

	return gateway
}

// buildServiceBO 从 PO 构建运行态 ServiceBO，可保留既有心跳时间。
func buildServiceBO(service po.Service, oldBeats map[string]time.Time) *ServiceBO {
	addresses := maplist.NewMapList[Address]()
	for _, addr := range service.Addresses {
		addresses.Add(addr.ID, &Address{
			Address:   addr.Address,
			LastBeat:  lastBeatFor(oldBeats, addr.Address),
			IsHealthy: true,
		})
	}
	apis := maplist.NewMapList[APIBO]()
	for _, api := range service.APIs {
		a := api
		apis.Add(a.ID, APIBOFromPO(a, service.Name))
	}
	return &ServiceBO{
		ServicePOId: service.ID,
		Prefix:      service.Prefix,
		Name:        service.Name,
		Protocol:    service.Protocol,
		Addresses:   addresses,
		APIs:        apis,
		curAddress:  0,
		Available:   service.Available,
	}
}

// FlushGateway 从 DB 全量刷新服务与 API 缓存（保留地址心跳时间）。
func (g *Gateway) FlushGateway() {
	var servicesPO []po.Service
	g.DB.Preload("Addresses").
		Preload("APIs.GrpcMethodMeta").
		Where("available=?", true).
		Find(&servicesPO)

	beats := snapshotAddressBeats()
	newServices := maplist.NewMapList[ServiceBO]()
	newPrefixes := maplist.NewMapStringList[Prefix]()

	for _, service := range servicesPO {
		newServices.Add(service.ID, buildServiceBO(service, beats[service.ID]))
		newPrefixes.Add(service.Prefix, &Prefix{
			Name:      service.Prefix,
			ServiceId: service.ID,
		})
	}

	g.RWMutex.Lock()
	g.Services = newServices
	g.Prefixes = newPrefixes
	g.RWMutex.Unlock()
	log.Printf("[gateway] FlushGateway 完成: %d 个服务", newServices.Size())
}

// SyncAllAPIsFromDB 仅刷新各服务内存中的 API 列表（轻量定时同步用）。
func (g *Gateway) SyncAllAPIsFromDB() {
	if g == nil {
		return
	}
	g.RWMutex.RLock()
	services := make([]*ServiceBO, 0, len(g.Services.List))
	for _, s := range g.Services.List {
		if s != nil {
			services = append(services, s)
		}
	}
	g.RWMutex.RUnlock()

	for _, s := range services {
		n := SyncServiceAPIsFromDB(s.ServicePOId, s.Name)
		if n > 0 {
			log.Printf("[gateway] 同步 API 缓存: %s (%d 条)", s.Name, n)
		}
	}
}

//增加服务健康检查

//func RegisterMux(mux *runtime.ServeMux, service string) error {
//	serve, ok := GetServiceByName(service)
//	if !ok {
//		return fmt.Errorf("服务不存在")
//	}
//
//}
