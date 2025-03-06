package domain

import (
	"github.com/trancecho/mundo-gateway/config"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/trancecho/mundo-gateway/po"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GatewayGlobal 第一个功能，代理路由
type Gateway struct {
	DB       *gorm.DB
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
	pwd := config.GlobalConfig.Mysql.Pwd
	dsn := "root:" + pwd + "@tcp(127.0.0.1:13306)/md_gateway?charset=utf8mb4&parseTime=True&loc=Local"
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
	db.Preload("Addresses").Preload("APIs.GrpcMethodMeta").
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
		var addresses []string
		for _, address := range service.Addresses {
			addresses = append(addresses, address.Address)
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
func (g *Gateway) FlushGateway() {
	// todo 可以优化
	// 重新获取service列表
	var servicesPO []po.Service
	g.DB.Preload("Addresses").Preload("APIs").
		Find(&servicesPO)
	// 初始化全局services列表
	var serviceBOs []ServiceBO
	for _, service := range servicesPO {
		var addresses []string
		for _, address := range service.Addresses {
			addresses = append(addresses, address.Address)
		}

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
		})
	}

	// 重新获取prefix列表
	var prefixes []Prefix
	for _, service := range servicesPO {
		prefixes = append(prefixes, Prefix{Name: service.Prefix, ServiceId: service.ID})
	}

	// 更新全局网关
	g.RWMutex.Lock()
	defer g.RWMutex.Unlock()
	g.Services = serviceBOs
	g.Prefixes = prefixes
}
