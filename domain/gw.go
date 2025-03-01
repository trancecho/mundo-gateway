package domain

import (
	"github.com/trancecho/mundo-gateway/po"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

// Gateway 第一个功能，代理路由
type Gateway struct {
	DB       *gorm.DB
	Prefixes []po.Prefix
	Services []po.Service
	globalKV sync.Map
}

// NewGateway 创建一个全局网关shili
func NewGateway() *Gateway {
	var err error
	// 会自己注册一个地址的。
	var db *gorm.DB
	dsn := "root:123456@tcp(127.0.0.1:13306)/md_gateway?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}

	err = db.AutoMigrate(&po.Service{})
	if err != nil {
		log.Fatalln("failed to migrate database", err)
	}
	log.Println("connect database success!!!!!!!!")

	// 接下来初始化路由表
	var services []po.Service
	db.Find(&services)

	// 网关实例。repo先建立初始化，服务数据预发现
	return &Gateway{DB: db, Services: services}
}
