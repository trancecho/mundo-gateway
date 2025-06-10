package main

import (
	gatewaySdk "github.com/trancecho/mundo-gateway-sdk"
	"log"
	"net/http"
	"time"
)

func main() {
	// 目标：我怎么请求gateway去拿到这个target
	target, err := gatewaySdk.NewClient("http://localhost:12388").GetTarget("ping")
	if err != nil {
		log.Println("get target error:", err)
		return
	}
	// http请求target
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", "http://"+target+"/ping", nil)
	if err != nil {
		log.Fatalln("create request error:", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("do request error:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("request failed with status:", resp.Status)
		return
	}
	log.Println("request successful, status:", resp.Status)
}
