package main

import (
	"github.com/gin-gonic/gin"
	gatewaySdk "github.com/trancecho/mundo-gateway-sdk"
	"log"
	"sync"
)

func main() {
	var err error
	client := gatewaySdk.NewGatewayService("ping", "http://localhost:6666", "http", "http://localhost:12388",
		gatewaySdk.NewMyRedisTokenGetter("localhost:6379", ""))
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	client.HttpConn(r)
	//BlackListTest(client)
	//	启动服务
	err = r.Run(":6666")
	if err != nil {
		panic(err)
	}
}

func BlackListTest(client *gatewaySdk.GatewayService) {
	// 启动一个协程，并发100访问gateway
	wg := sync.WaitGroup{}
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func(i int) {
			pong := client.Ping()
			log.Println(pong, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
