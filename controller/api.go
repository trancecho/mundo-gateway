package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller/dto"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
	"strconv"
)

func CreateAPIController(c *gin.Context) {
	var req dto.APICreateReq
	c.BindJSON(&req)
	if req.ServiceName == "" {
		util.ClientError(c, 100, "name不能为空")
		return
	}
	if req.Path == "" {
		util.ClientError(c, 110, "path不能为空")
		return
	}
	if req.Method == "" {
		util.ClientError(c, 120, "method不能为空")
		return
	}

	// 检查是否已存在相同的路径和方法,这里开始写bug
	existingAPI, err := domain.GetAPIByPathAndMethod(req.Path, req.Method)
	if err == nil && existingAPI != nil {
		util.ClientError(c, 10100, "API路径和方法已存在")
		return
	}

	// 创建API
	apiPO, err := domain.CreateAPIService(&req)
	if err != nil {
		util.ServerError(c, 200, "API创建失败:"+err.Error())
		return
	}
	domain.GatewayGlobal.FlushGateway()

	util.Ok(c, "API创建成功", gin.H{
		"api": apiPO,
	})
}

func UpdateAPIController(c *gin.Context) {
	var req dto.APIUpdateReq
	c.BindJSON(&req)
	if req.Id == 0 {
		util.ClientError(c, 1, "id不能为空")
		return
	}
	if req.Name == "" {
		util.ClientError(c, 2, "name不能为空")
		return
	}
	if req.Path == "" {
		util.ClientError(c, 3, "path不能为空")
		return
	}
	if req.Method == "" {
		util.ClientError(c, 4, "method不能为空")
		return
	}
	// 更新API
	apiPO, err := domain.UpdateAPIService(&req)
	if err != nil {
		util.ServerError(c, 5, "API更新失败")
		return
	}
	domain.GatewayGlobal.FlushGateway()

	util.Ok(c, "API更新成功", gin.H{
		"api": apiPO,
	})
}

func DeleteAPIController(c *gin.Context) {
	var req dto.APIDeleteReq
	c.BindJSON(&req)
	if req.Id == 0 {
		util.ClientError(c, 1, "id不能为空")
		return
	}
	// 删除API
	err := domain.DeleteAPIService(req.Id)
	if err != nil {
		util.ServerError(c, 2, "API删除失败")
		return
	}
	domain.GatewayGlobal.FlushGateway()

	util.Ok(c, "API删除成功", nil)
}

func GetAPIController(c *gin.Context) {
	var err error
	var id int
	id, err = strconv.Atoi(c.Query("id"))
	if err != nil {
		util.ServerError(c, 3, "api的id格式错误")
		return
	}
	// 获取API
	apiPO, err := domain.GetAPIService(int64(id))
	if err != nil {
		util.ServerError(c, 4, "API获取失败")
		return
	}

	util.Ok(c, "API获取成功", gin.H{
		"api": apiPO,
	})
}

func ListAPIController(c *gin.Context) {
	// 获取API列表
	apis, err := domain.ListAPIService()
	if err != nil {
		util.ServerError(c, 5, "API列表获取失败")
		return
	}
	util.Ok(c, "API列表获取成功", gin.H{
		"apis": apis,
	})
}
