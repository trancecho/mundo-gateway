package domain

import (
	"github.com/trancecho/mundo-gateway/po"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

// GatewayGlobal 第一个功能，代理路由
type Gateway struct {
	DB       *gorm.DB
	Prefixes []Prefix
	Services []ServiceBO
	globalKV sync.Map //可以先忽略
	//读写锁
	sync.RWMutex
}

// NewGateway 创建一个全局网关shili
func NewGateway() *Gateway {
	//todo全局网关注册后，还需要时刻更新属性列表，如prefixes，services
	var err error
	// 会自己注册一个地址的。
	var db *gorm.DB
	dsn := "root:123456@tcp(127.0.0.1:13306)/md_gateway?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}

	err = db.AutoMigrate(&po.Service{}, &po.API{}, &po.Address{})
	if err != nil {
		log.Fatalln("failed to migrate database", err)
	}
	log.Println("connect database success!!!!!!!!")

	// c初始化service列表
	var services []po.Service
	db.Find(&services)
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
		serviceBOs = append(serviceBOs, ServiceBO{
			ServicePOId: service.ID,
			Prefix:      service.Prefix,
			Name:        service.Name,
			Addresses:   addresses,
			Protocol:    service.Protocol,
			curAddress:  0,
		})
	}
	// 网关实例。repo先建立初始化，服务数据预发现
	return &Gateway{
		DB:       db,
		Services: serviceBOs,
		Prefixes: prefixes}
}

// FlushGateway 重新获取service列表
func (g *Gateway) FlushGateway() {

	// 重新获取service列表
	var servicesPO []po.Service
	g.DB.Find(&servicesPO)
	// 初始化全局services列表
	var serviceBOs []ServiceBO
	for _, service := range servicesPO {
		var addresses []string
		for _, address := range service.Addresses {
			addresses = append(addresses, address.Address)
		}
		serviceBOs = append(serviceBOs, ServiceBO{
			ServicePOId: service.ID,
			Prefix:      service.Prefix,
			Name:        service.Name,
			Addresses:   addresses,
			Protocol:    service.Protocol,
			curAddress:  0,
		})
	}

	// 重新获取prefix列表
	var prefixes []Prefix
	for _, service := range servicesPO {
		prefixes = append(prefixes, Prefix{Name: service.Prefix, ServiceId: service.ID})
	}
	// 更新全局网关
	// 写操作一起处理
	g.RWMutex.Lock()
	defer g.RWMutex.Unlock()
	g.Services = serviceBOs
	g.Prefixes = prefixes
}
