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
		util.ClientErr(c, 100, "name不能为空")
		return
	}
	if req.Path == "" {
		util.ClientErr(c, 2, "path不能为空")
		return
	}
	if req.Method == "" {
		util.ClientErr(c, 3, "method不能为空")
		return
	}
	// 创建API

	apiPO, err := domain.CreateAPIService(&req)
	if err != nil {
		util.ServerError(c, 4, "API创建失败")
		return
	}
	util.Ok(c, "API创建成功", gin.H{
		"api": apiPO,
	})
}

func UpdateAPIController(c *gin.Context) {
	var req dto.APIUpdateReq
	c.BindJSON(&req)
	if req.Id == 0 {
		util.ClientErr(c, 1, "id不能为空")
		return
	}
	if req.Name == "" {
		util.ClientErr(c, 2, "name不能为空")
		return
	}
	if req.Path == "" {
		util.ClientErr(c, 3, "path不能为空")
		return
	}
	if req.Method == "" {
		util.ClientErr(c, 4, "method不能为空")
		return
	}
	// 更新API
	apiPO, err := domain.UpdateAPIService(&req)
	if err != nil {
		util.ServerError(c, 5, "API更新失败")
		return
	}
	util.Ok(c, "API更新成功", gin.H{
		"api": apiPO,
	})
}

func DeleteAPIController(c *gin.Context) {
	var req dto.APIDeleteReq
	c.BindJSON(&req)
	if req.Id == 0 {
		util.ClientErr(c, 1, "id不能为空")
		return
	}
	// 删除API
	err := domain.DeleteAPIService(req.Id)
	if err != nil {
		util.ServerError(c, 2, "API删除失败")
		return
	}
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
