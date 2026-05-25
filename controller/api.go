package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller/dto"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
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

	existingAPI, err := domain.GetAPIByPathAndMethod(req.Path, req.Method, req.ServiceName)
	if err == nil && existingAPI != nil {
		domain.UpsertAPIBOToCache(existingAPI, req.ServiceName)
		c.JSON(200, gin.H{
			"err_code": 410100,
			"message":  "API已存在",
			"data":     gin.H{"api": *existingAPI},
		})
		return
	}

	apiPO, err := domain.CreateAPIService(&req)
	if err != nil {
		if err.Error() == "API已存在" {
			existingAPI, _ = domain.GetAPIByPathAndMethod(req.Path, req.Method, req.ServiceName)
			if existingAPI != nil {
				domain.UpsertAPIBOToCache(existingAPI, req.ServiceName)
				c.JSON(200, gin.H{
					"err_code": 410100,
					"message":  "API已存在",
					"data":     gin.H{"api": *existingAPI},
				})
				return
			}
		}
		util.ServerError(c, 200, "API创建失败:"+err.Error())
		return
	}

	domain.UpsertAPIBOToCache(apiPO, req.ServiceName)
	if domain.GetServiceBO(req.ServiceName) == nil {
		domain.ReloadServiceIntoGateway(apiPO.ServiceId)
	}

	util.Ok(c, "API创建成功", gin.H{
		"api": *apiPO,
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
	apiPO, err := domain.UpdateAPIService(&req)
	if err != nil {
		util.ServerError(c, 5, "API更新失败")
		return
	}

	domain.UpsertAPIBOToCache(apiPO, req.Name)
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
	apiPO, err := domain.DeleteAPIService(req.Id)
	if err != nil {
		util.ServerError(c, 2, "API删除失败")
		return
	}
	util.Ok(c, "API删除成功", nil)
	_ = apiPO
}

func GetAPIController(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		util.ServerError(c, 3, "api的id格式错误")
		return
	}
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
	serviceName := c.Query("serviceName")
	if serviceName == "" {
		serviceName = "all"
	}
	apis, err := domain.ListAPIService(serviceName)
	if err != nil {
		util.ServerError(c, 5, "API列表获取失败")
		return
	}
	util.Ok(c, "API列表获取成功", gin.H{
		"apis": apis,
	})
}
