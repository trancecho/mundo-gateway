package domain

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/po"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

// GatewayGlobal 第一个功能，代理路由
type Gateway struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Prefixes []Prefix
	Services []ServiceBO
	globalKV sync.Map //可以先忽略
	//读写锁
	sync.RWMutex
	HTTPClient *http.Client //http客户端
}

// NewGateway 创建一个全局网关shili
func NewGateway() *Gateway {
	//todo全局网关注册后，还需要时刻更新属性列表，如prefixes，services
	var err error
	// 会自己注册一个地址的。
	var db *gorm.DB
	pwd := viper.GetString("mysql.pwd")
	dsn := "root:" + pwd + "@tcp(" + viper.GetString("mysql.host") + ":" + viper.GetString("mysql.port") + ")/md_gateway?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}

	err = db.AutoMigrate(&po.Service{}, &po.API{}, &po.Address{}, &po.GrpcMethodMeta{})
	if err != nil {
		log.Fatalln("failed to migrate database", err)
	}
	log.Println("connect database success!!!!!!!!")

	// c初始化service列表
	var services []po.Service
	db.Preload("Addresses").Preload("APIs.GrpcMethodMeta").Where("available=?", true).
		Find(&services)
	log.Println("services:", services)
	var prefixes []Prefix
	//在services中找到所有的prefix
	for _, service := range services {
		prefixes = append(prefixes, Prefix{Name: service.Prefix, ServiceId: service.ID})
	}

	// 初始化全局services列表
	var serviceBOs []ServiceBO
	for _, service := range services {
		var addresses []*Address
		for _, address := range service.Addresses {
			addresses = append(addresses, &Address{address.Address, time.Now(), true})
		}

		// 转换 APIs
		var apis []APIBO
		for _, api := range service.APIs {
			apis = append(apis, APIBO{
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

		serviceBOs = append(serviceBOs, ServiceBO{
			ServicePOId: service.ID,
			Prefix:      service.Prefix,
			Name:        service.Name,
			Addresses:   addresses,
			Protocol:    service.Protocol,
			curAddress:  0,
			APIs:        apis, // 添加 APIs
			Available:   service.Available,
		})
	}

	// 初始化 HTTP 客户端
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	return &Gateway{
		DB:         db,
		Services:   serviceBOs,
		Prefixes:   prefixes,
		HTTPClient: httpClient,
	}
}

// FlushGateway 重新获取service列表
func (this *Gateway) FlushGateway() {
	// todo 可以优化
	// 重新获取service列表
	var servicesPO []po.Service
	this.DB.Preload("Addresses").Preload("APIs").Where("available=?", true).
		Find(&servicesPO)
	// 初始化全局services列表
	var serviceBOs []ServiceBO
	for _, service := range servicesPO {
		var addresses []*Address
		for _, address := range service.Addresses {
			addresses = append(addresses, &Address{address.Address, time.Now(), true})
		}
		log.Println("addresses:", addresses)

		// 转换 APIs
		var apis []APIBO
		for _, api := range service.APIs {
			apis = append(apis, APIBO{
				Id:         api.ID,
				HttpPath:   api.HttpPath,
				HttpMethod: api.HttpMethod,
			})
		}

		serviceBOs = append(serviceBOs, ServiceBO{
			ServicePOId: service.ID,
			Prefix:      service.Prefix,
			Name:        service.Name,
			Addresses:   addresses,
			Protocol:    service.Protocol,
			curAddress:  0,
			APIs:        apis,
			Available:   service.Available,
		})
	}

	// 重新获取prefix列表
	var prefixes []Prefix
	for _, service := range servicesPO {
		prefixes = append(prefixes, Prefix{Name: service.Prefix, ServiceId: service.ID})
	}

	// 更新全局网关
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.Services = serviceBOs
	this.Prefixes = prefixes
}

//增加服务健康检查

//func RegisterMux(mux *runtime.ServeMux, service string) error {
//	serve, ok := GetServiceByName(service)
//	if !ok {
//		return fmt.Errorf("服务不存在")
//	}
//
//}
