package domain

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-gateway/po"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GatewayGlobal 第一个功能，代理路由
type Gateway struct {
	DB         *gorm.DB
	Redis      *redis.Client
	Prefixes   []Prefix
	Services   []ServiceBO
	ServiceMap map[int64]ServiceBO
	globalKV   sync.Map //可以先忽略
	//读写锁
	sync.RWMutex
	HTTPClient *http.Client //

}

// NewGateway 创建一个全局网关shili
func NewGateway() *Gateway {
	//todo全局网关注册后，还需要时刻更新属性列表，如prefixes，services
	var err error
	// 会自己注册一个地址的。
	var db *gorm.DB
	pwd := viper.GetString("mysql.pwd")
	dsn := "root:" + pwd + "@tcp(" + viper.GetString("mysql.host") + ":" + viper.GetString("mysql.port") + ")/" + viper.GetString("mysql.db") + "?charset=utf8mb4&parseTime=True&loc=Local"
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
			MaxIdleConns:        10000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	// 初始化 ServiceMap，作为内存中的主索引
	serviceMap := make(map[int64]ServiceBO, len(serviceBOs))
	for _, bo := range serviceBOs {
		serviceMap[bo.ServicePOId] = bo
	}

	gateway := &Gateway{
		DB:         db,
		Services:   serviceBOs,
		ServiceMap: serviceMap,
		Prefixes:   prefixes,
		HTTPClient: httpClient,
	}

	return gateway

}

// FlushGateway 重新获取service列表
func (this *Gateway) FlushGateway() {
	// todo 可以优化,把全量更新，改成增量更新
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

	// 根据最新的 serviceBOs 构建 ServiceMap
	serviceMap := make(map[int64]ServiceBO, len(serviceBOs))
	for _, bo := range serviceBOs {
		serviceMap[bo.ServicePOId] = bo
	}

	// 更新全局网关
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.Services = serviceBOs
	this.ServiceMap = serviceMap
	this.Prefixes = prefixes

}

// PartialFlushGateway
func (this *Gateway) PartialFlushGateway(serviceID int64) {
	// 先从 DB 查询当前 service 最新状态（避免长时间持有锁）
	var service po.Service
	err := this.DB.Preload("Addresses").
		Where("id = ?", serviceID).First(&service).Error

	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	//
	if this.ServiceMap == nil {
		this.ServiceMap = make(map[int64]ServiceBO, len(this.Services))
		for _, svc := range this.Services {
			this.ServiceMap[svc.ServicePOId] = svc
		}
	}
	if err == gorm.ErrRecordNotFound || (err == nil && !service.Available) {
		delete(this.ServiceMap, serviceID)
	} else if err != nil {
		log.Println("PartialFlushGateway 查询 service 失败:", err)
		return
	} else {
		var addresses []*Address
		for _, addr := range service.Addresses {
			addresses = append(addresses, &Address{
				Address:   addr.Address,
				LastBeat:  time.Now(),
				IsHealthy: true,
			})
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

		newBO := ServiceBO{
			ServicePOId: service.ID,
			Prefix:      service.Prefix,
			Name:        service.Name,
			Addresses:   addresses,
			Protocol:    service.Protocol,
			Available:   service.Available,
			APIs:        apis,
		}

		// ServiceMap 中覆盖 / 新增这一条
		this.ServiceMap[serviceID] = newBO
	}

	// 统一用 ServiceMap 重新构建 Services 和 Prefixes
	this.Services = this.Services[:0]
	this.Prefixes = this.Prefixes[:0]
	for id, bo := range this.ServiceMap {
		this.Services = append(this.Services, bo)
		this.Prefixes = append(this.Prefixes, Prefix{
			Name:      bo.Prefix,
			ServiceId: id,
		})
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
