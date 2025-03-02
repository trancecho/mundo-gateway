package domain

import (
	"errors"
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/po"
)

type API struct {
	Id        int64
	ServiceId int64
	Path      string
	Method    string
}

func NewAPI(id int64, serviceId int64, path string, method string) *API {
	return &API{
		Id:        id,
		ServiceId: serviceId,
		Path:      path,
		Method:    method,
	}
}

func CreateAPIService(dto *controller.APICreateReq) (*po.API, error) {
	var err error
	var apiPO po.API
	var servicePO po.Service
	db := GatewayGlobal.DB
	servicePO.Name = dto.Name

	// 如果没有service直接返回
	affected := db.Find(&servicePO).RowsAffected
	if affected == 0 {
		return nil, errors.New("服务不存在")
	}

	//如果有service，寻找api是否存在
	apiPO.Path = dto.Path
	apiPO.Method = dto.Method
	affected = db.Find(&apiPO).RowsAffected
	if affected > 0 {
		return nil, errors.New("API已存在")
	}

	// 创建API
	apiPO.ServiceId = servicePO.ID
	err = db.Create(&apiPO).Error
	if err != nil {
		return nil, err
	}

	return &apiPO, nil
}

// GetAPIService 获取API
func GetAPIService(id int64) (*po.API, error) {
	var apiPO po.API
	db := GatewayGlobal.DB
	affected := db.First(&apiPO, id).RowsAffected
	if affected == 0 {
		return nil, errors.New("API不存在")
	}
	return &apiPO, nil
}

// UpdateAPIService 更新API
func UpdateAPIService(dto *controller.APIUpdateReq) (*po.API, error) {
	var err error
	var apiPO po.API
	db := GatewayGlobal.DB
	apiPO.ID = dto.Id

	// 查找api
	affected := db.Find(&apiPO).RowsAffected
	if affected == 0 {
		return nil, errors.New("API不存在")
	}

	// 更新api
	apiPO.Path = dto.Path
	apiPO.Method = dto.Method
	err = db.Save(&apiPO).Error
	if err != nil {
		return nil, err
	}

	return &apiPO, nil
}

// DeleteAPIService 删除API
func DeleteAPIService(id int64) error {
	var apiPO po.API
	db := GatewayGlobal.DB
	affected := db.First(&apiPO, id).RowsAffected
	if affected == 0 {
		return errors.New("API不存在")
	}
	err := db.Delete(&apiPO).Error
	if err != nil {
		return err
	}
	return nil
}

// ListAPIService 获取API列表
func ListAPIService() ([]*po.API, error) {
	var apiPOs []*po.API
	db := GatewayGlobal.DB
	err := db.Find(&apiPOs).Error
	if err != nil {
		return nil, err
	}
	return apiPOs, nil
}
