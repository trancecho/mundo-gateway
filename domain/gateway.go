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

	// 初始化全局 services 列表
	serviceBOs := maplist.NewMapList[ServiceBO]()
	for _, service := range services {
		addresses := maplist.NewMapList[Address]()
		for _, addr := range service.Addresses {
			addresses.Add(addr.ID, &Address{
				Address:   addr.Address,
				LastBeat:  time.Now(),
				IsHealthy: true,
			})
		}
		apis := maplist.NewMapList[APIBO]()
		for _, api := range service.APIs {
			apis.Add(api.ID, &APIBO{
				Id:         api.ID,
				HttpPath:   api.HttpPath,
				HttpMethod: api.HttpMethod,
				GrpcMethodMeta: GrpcMethodMetaBO{
					ApiId:       api.ID,
					ServiceName: api.GrpcMethodMeta.ServiceName,
					MethodName:  api.GrpcMethodMeta.MethodName,
				},
			})
		}

		serviceBOs.Add(service.ID, &ServiceBO{
			ServicePOId: service.ID,
			Prefix:      service.Prefix,
			Name:        service.Name,
			Protocol:    service.Protocol,
			Addresses:   addresses,
			APIs:        apis,
			curAddress:  0,
			Available:   service.Available,
		})

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

// FlushGateway 重新获取service列表
func (g *Gateway) FlushGateway() {
	var servicesPO []po.Service
	g.DB.Preload("Addresses").Preload("APIs").
		Where("available=?", true).
		Find(&servicesPO)

	newServices := maplist.NewMapList[ServiceBO]()
	newPrefixes := maplist.NewMapStringList[Prefix]()

	for _, service := range servicesPO {
		addresses := maplist.NewMapList[Address]()
		for _, addr := range service.Addresses {
			addresses.Add(addr.ID, &Address{
				Address:   addr.Address,
				LastBeat:  time.Now(),
				IsHealthy: true,
			})
		}
		log.Println("刷新服务地址列表", addresses)

		apis := maplist.NewMapList[APIBO]()
		for _, api := range service.APIs {
			apis.Add(api.ID, &APIBO{
				Id:         api.ID,
				HttpPath:   api.HttpPath,
				HttpMethod: api.HttpMethod,
				GrpcMethodMeta: GrpcMethodMetaBO{
					ApiId:       api.ID,
					ServiceName: api.GrpcMethodMeta.ServiceName,
					MethodName:  api.GrpcMethodMeta.MethodName,
				},
			})
		}
		log.Println("apis:", apis)
		newServices.Add(service.ID, &ServiceBO{
			ServicePOId: service.ID,
			Prefix:      service.Prefix,
			Name:        service.Name,
			Protocol:    service.Protocol,
			Addresses:   addresses,
			APIs:        apis,
			curAddress:  0,
			Available:   service.Available,
		})

		newPrefixes.Add(service.Prefix, &Prefix{
			Name:      service.Prefix,
			ServiceId: service.ID,
		})
	}
	log.Println("newServices:", newServices)
	log.Println("newPrefixes:", newPrefixes)

	g.RWMutex.Lock()
	g.Services = newServices
	g.Prefixes = newPrefixes
	val, ok := g.Services.Get(3)
	if ok {
		log.Println("flush gateway success!!!!!!!!", val)
	} else {
		log.Println("flush gateway success!!!!!!!!, but service 3 not found")
	}

	g.RWMutex.Unlock()

}

//增加服务健康检查

//func RegisterMux(mux *runtime.ServeMux, service string) error {
//	serve, ok := GetServiceByName(service)
//	if !ok {
//		return fmt.Errorf("服务不存在")
//	}
//
//}
