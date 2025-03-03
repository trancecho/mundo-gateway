package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller/dto"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
	"strconv"
)

//type ServiceDTO struct {
//	Name      string `json:"name"`
//	Prefix    string `json:"prefix"`
//	Protocol  string `json:"protocol"`
//	Addresses []Address
//}

//type Address struct {
//	Id        int64
//	ServiceId int64
//	Address   string
//}

func CreateServiceController(c *gin.Context) {
	var dto dto.ServiceCreateReq
	c.ShouldBindJSON(&dto)
	if dto.Name == "" {
		util.ClientErr(c, 1, "name不能为空")
		return
	}
	if dto.Prefix == "" {
		util.ClientErr(c, 2, "prefix不能为空")
		return
	}
	if dto.Protocol == "" {
		util.ClientErr(c, 3, "protocol不能为空")
		return
	}
	// 根据协议判断地址是否合规
	if dto.Protocol == "http" {
		if dto.Address[:7] != "http://" || dto.Address[:8] != "https://" {
			util.ServerError(c, 500, "http协议地址不能为空")
			return
		}
	} else if dto.Protocol == "grpc" {
		if dto.Address[:6] != "grpc://" {
			util.ServerError(c, 600, "grpc协议地址不能为空")
			return
		}
	} else {
		util.ServerError(c, 700, "协议不合规")
		return
	}

	servicePO, ok := domain.CreateServiceService(&dto)
	if !ok {
		util.ServerError(c, 4, "服务创建失败")
		return
	}
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
	ok := domain.DeleteServiceService(idInt64)
	if !ok {
		util.ServerError(c, 2, "服务删除失败")
		return
	}
	util.Ok(c, "服务删除成功", nil)
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
