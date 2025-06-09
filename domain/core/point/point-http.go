package point

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/domain/core/point/point_dto"
	"github.com/trancecho/mundo-gateway/global"
	"github.com/trancecho/mundo-gateway/util"
	"log"
)

func ChangePointAndExperience(c *gin.Context) {
	var req point_dto.ChangePointReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ServerError(c, 400, "请求参数错误")
		return
	}
	system, err := global.PointsGlobal.DialToPointsSystem(req.Token, req.UserID, req.Point, req.Experience, req.Reason)
	if err != nil {
		log.Println(err)
		util.ServerError(c, 400, err.Error())
		return
	}
	util.Ok(c, "使用积分成功", gin.H{"Message": system})
}

func Sign(c *gin.Context) {
	var req point_dto.SignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ServerError(c, 400, "请求参数错误")
		return
	}
	Message, err := global.PointsGlobal.SignToPointsSystem(req.Token, req.UserId)
	if err != nil {
		log.Println(err)
		util.ServerError(c, 400, err.Error())
		return
	}
	util.Ok(c, "签到成功", gin.H{"Message": Message})
}

func PointsInfo(c *gin.Context) {
	var req point_dto.PointsInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ServerError(c, 400, "请求参数错误")
		return
	}
	Message, err := global.PointsGlobal.GetPointsInfo(req.Token, req.UserId)
	if err != nil {
		log.Println(err)
		util.ServerError(c, 500, err.Error())
		return
	}
	util.Ok(c, "获取积分信息成功", gin.H{"Message": Message})
}

func GetStats(c *gin.Context) {
	var req point_dto.AdminStatsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ServerError(c, 400, "请求参数错误")
		return
	}
	Message, err := global.PointsGlobal.AdminStats(req.Token, req.UserId)
	if err != nil {
		log.Println(err)
		util.ServerError(c, 500, err.Error())
		return
	}
	util.Ok(c, "获取信息成功", gin.H{"Message": Message})
}
