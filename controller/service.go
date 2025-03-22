package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller/dto"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
	"log"
	"strconv"
	"strings"
	"time"
)

//type ServiceDTO struct {
//	Name      string `json:"name"`
//	Prefix    string `json:"prefix"`
//	Protocol  string `json:"protocol"`
//	Addresses []Address
//}

//type Address struct {
//	Id        int64
//	ApiId int64
//	Address   string
//}

func CreateServiceController(c *gin.Context) {
	var req dto.ServiceCreateReq
	c.ShouldBindJSON(&req)
	if req.Name == "" {
		util.ClientError(c, 1, "name不能为空")
		return
	}
	if req.Prefix == "" {
		util.ClientError(c, 2, "prefix不能为空")
		return
	}
	if req.Prefix == "gateway" {
		util.ClientError(c, 300, "prefix不能为gateway ")
	}
	if req.Protocol == "" {
		util.ClientError(c, 310, "protocol不能为空")
		return
	}
	// 根据协议判断地址是否合规 目前只有http和grpc
	if req.Protocol == "http" {
		// 检查地址是否以 http:// 或 https:// 开头
		if !strings.HasPrefix(req.Address, "http://") && !strings.HasPrefix(req.Address, "https://") {
			util.ServerError(c, 500, "http协议地址不能为空")
			return
		}
	} else if req.Protocol == "grpc" {
		if req.Address[:7] != "grpc://" {
			log.Println(req.Address[:7])
			util.ServerError(c, 600, "grpc协议地址不能为空")
			return
		}
	} else {
		util.ServerError(c, 700, "协议不合规")
		return
	}

	servicePO, ok := domain.CreateServiceService(&req)
	if !ok {
		domain.GatewayGlobal.FlushGateway()
		util.ServerError(c, 4, "服务创建失败")
		return
	}
	domain.GatewayGlobal.FlushGateway()
	util.Ok(c, "服务创建成功", gin.H{
		"service": servicePO,
	})
}

func UpdateServiceController(c *gin.Context) {
	var req dto.ServiceUpdateReq
	c.ShouldBindJSON(&req)
	if req.Id == 0 {
		util.ServerError(c, 100, "id不能为空")
		return
	}
	if req.Name == "" && req.Prefix == "" && req.Protocol == "" {
		util.ServerError(c, 200, "name、prefix、protocol不能同时为空")
		return
	}

	servicePO, ok := domain.UpdateServiceService(&req)
	if !ok {
		util.ServerError(c, 800, "服务更新失败")
		return
	}
	domain.GatewayGlobal.FlushGateway()

	util.Ok(c, "服务更新成功", gin.H{
		"service": servicePO,
	})
}

func DeleteServiceController(c *gin.Context) {
	var err error
	var id int
	id, err = strconv.Atoi(c.Query("id"))
	if err != nil {
		util.ServerError(c, 3, "id格式错误")
		return
	}
	idInt64 := int64(id)
	if id == 0 {
		util.ServerError(c, 1, "id不能为空")
		return
	}

	// 删除与该服务相关的所有API
	err = domain.DeleteAPIsByServiceID(idInt64) // 调用删除API的函数
	if err != nil {
		util.ServerError(c, 2, "删除相关API失败")
		return
	}

	ok := domain.DeleteServiceService(idInt64)
	if !ok {
		util.ServerError(c, 2, "服务删除失败")
		return
	}
	domain.GatewayGlobal.FlushGateway()

	util.Ok(c, "服务删除和相关api删除成功", nil)
}

func DeleteServiceAddressController(c *gin.Context) {
	var err error
	var id int
	id, err = strconv.Atoi(c.Query("id"))
	if err != nil {
		util.ServerError(c, 3, "id格式错误")
		return
	}
	idInt64 := int64(id)
	if id == 0 {
		util.ServerError(c, 1, "id不能为空")
		return
	}
	ok := domain.DeleteAddressService(idInt64)
	if !ok {
		util.ServerError(c, 2, "服务地址删除失败")
		return
	}
	domain.GatewayGlobal.FlushGateway()

	util.Ok(c, "服务地址删除成功", nil)
}

func ListServiceController(c *gin.Context) {
	services, ok := domain.ListServicesService()
	if !ok {
		util.ServerError(c, 1, "服务列表获取失败")
		return
	}
	util.Ok(c, "服务列表", gin.H{
		"services": services,
	})

}

func GetServiceController(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		util.ServerError(c, 1, "id不能为空")
		return
	}
	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ServerError(c, 3, "id格式错误")
		return
	}

	service, ok := domain.GetServiceService(idInt64)
	if !ok {
		util.ServerError(c, 2, "服务获取失败")
		return
	}
	util.Ok(c, "服务获取成功", gin.H{
		"service": service,
	})
}

// 服务心跳
func ServiceAliveSignalController(c *gin.Context) {
	var req dto.ServiceBeatReq
	c.ShouldBindJSON(&req)
	if req.ServiceName == "" || req.Address == "" {
		util.ServerError(c, 100, "服务名和地址不能为空")
		return
	} else {
		boPtr := domain.GetServiceBO(req.ServiceName)
		if boPtr == nil {
			util.ServerError(c, 200, "服务心跳失败")
			return
		}
		if boPtr.GetAddressBO(req.Address) == nil {
			util.ServerError(c, 300, "服务地址不存在")
			return
		}
		// 更新服务的心跳时间
		boPtr.GetAddressBO(req.Address).LastBeat = time.Now()
		log.Println("服务心跳成功", req.Address, req.ServiceName)
	}
}

func ServiceAliveChecker() {
	// 定时检查服务的心跳
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, serviceBO := range domain.GatewayGlobal.Services {
				for _, address := range serviceBO.Addresses {
					// 如果服务超过30秒没有心跳，则认为服务不可用
					if time.Since(address.LastBeat) > 30*time.Second {
						log.Println("服务不可用", serviceBO.Name, address.Address)
						// 删除服务地址
						domain.UnregisterServiceService(serviceBO.Name, address.Address)
					}
				}
			}
		}
	}
}
