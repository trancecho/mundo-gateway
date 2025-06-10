package domain

import (
	"errors"
	"github.com/trancecho/mundo-gateway/controller/dto"
	"github.com/trancecho/mundo-gateway/po"
)

// APIBO 这是一种映射，而http默认直接转发，grpc需要映射
type APIBO struct {
	Id         int64
	HttpPath   string
	HttpMethod string

	GrpcMethodMeta GrpcMethodMetaBO
}

type GrpcMethodMetaBO struct {
	ApiId       int64
	ServiceName string // 服务名
	MethodName  string // 方法名
}

func NewAPIBO(id int64, serviceId int64, path string, method string) *APIBO {
	return &APIBO{
		Id:         id,
		HttpPath:   path,
		HttpMethod: method,
		GrpcMethodMeta: GrpcMethodMetaBO{
			ApiId:       id,
			ServiceName: "",
			MethodName:  "",
		},
	}
}

func CreateAPIService(req *dto.APICreateReq) (*po.API, error) {
	var err error
	var apiPO po.API
	var servicePO po.Service
	db := GatewayGlobal.DB
	servicePO.Name = req.ServiceName

	// 如果没有service直接返回
	affected := db.Where("name = ?", req.ServiceName).
		First(&servicePO).RowsAffected
	if affected == 0 {
		return nil, errors.New("服务不存在")
	}
	//如果有service，寻找api是否存在
	apiPO.HttpPath = req.Path
	apiPO.HttpMethod = req.Method
	apiPO.ServiceId = servicePO.ID
	affected = db.Where("service_id = ? and http_path = ? and method = ?", servicePO.ID, req.Path, req.Method).
		Find(&apiPO).RowsAffected
	if affected > 0 {
		return nil, errors.New("API已存在")
	}
	if servicePO.Protocol == "http" {
		// api不存在
		db.Create(&apiPO)
		for _, service := range GatewayGlobal.Services {
			if service.ServicePOId == servicePO.ID {
				service.APIs = append(service.APIs, APIBO{
					Id:         apiPO.ID,
					HttpPath:   apiPO.HttpPath,
					HttpMethod: apiPO.HttpMethod,
					GrpcMethodMeta: GrpcMethodMetaBO{
						ApiId:       apiPO.ID,
						ServiceName: "",
						MethodName:  "",
					},
				})
			}
		}
		return &apiPO, nil
	} else if servicePO.Protocol == "grpc" {
		////如果有service，寻找api是否存在
		apiPO.GrpcMethodMeta.ServiceName = req.GrpcService
		apiPO.GrpcMethodMeta.MethodName = req.GrpcMethod
		//affected = db.Where("service_id = ? and grpc_method_meta.service_name = ? and grpc_method_meta.method_name = ?",
		//	servicePO.ID, req.GrpcService, req.Method).
		//	Find(&apiPO).RowsAffected
		//if affected > 0 {
		//	return nil, errors.New("API已存在")
		//}
		err = db.Create(&apiPO).Error
		if err != nil {
			return nil, err
		}
		for _, service := range GatewayGlobal.Services {
			if service.ServicePOId == servicePO.ID {
				service.APIs = append(service.APIs, APIBO{
					Id:         apiPO.ID,
					HttpPath:   apiPO.HttpPath,
					HttpMethod: apiPO.HttpMethod,
					GrpcMethodMeta: GrpcMethodMetaBO{
						ApiId:       apiPO.ID,
						ServiceName: req.GrpcService,
						MethodName:  req.GrpcMethod,
					},
				})
			}
		}
		return &apiPO, nil
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
func UpdateAPIService(dto *dto.APIUpdateReq) (*po.API, error) {
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
	apiPO.HttpPath = dto.Path
	apiPO.HttpMethod = dto.Method
	apiPO.GrpcMethodMeta.ServiceName = dto.GrpcMethodMeta.ServiceName
	apiPO.GrpcMethodMeta.MethodName = dto.GrpcMethodMeta.MethodName

	err = db.Save(&apiPO).Error

	for _, service := range GatewayGlobal.Services {
		if service.ServicePOId == apiPO.ServiceId {
			for _, api := range service.APIs {
				if api.Id == apiPO.ID {
					api.HttpPath = apiPO.HttpPath
					api.HttpMethod = apiPO.HttpMethod
					api.GrpcMethodMeta.ServiceName = dto.GrpcMethodMeta.ServiceName
					api.GrpcMethodMeta.MethodName = dto.GrpcMethodMeta.MethodName
				}
			}
		}
	}

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
	for _, service := range GatewayGlobal.Services {
		if service.ServicePOId == apiPO.ServiceId {
			for i, api := range service.APIs {
				if api.Id == apiPO.ID {
					service.APIs = append(service.APIs[:i], service.APIs[i+1:]...)
				}
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

// ListAPIService 获取API列表
func ListAPIService(serviceName string) ([]*po.API, error) {
	var apiPOs []*po.API
	db := GatewayGlobal.DB
	if serviceName == "all" {
		err := db.Find(&apiPOs).Error
		if err != nil {
			return nil, err
		}
	} else {
		var servicePO po.Service
		// 查询serviceName对应的id
		affected := db.Where("name = ?", serviceName).
			First(&servicePO).RowsAffected
		if affected == 0 {
			return nil, errors.New("服务不存在")
		}
		err := db.Where("service_id = ?", servicePO.ID).Find(&apiPOs).Error
		if err != nil {
			return nil, err
		}
	}

	return apiPOs, nil
}

// GetAPIByPathAndMethod 检查数据库中是否已存在相同路径和方法的 API
func GetAPIByPathAndMethod(path string, method string, serviceName string) (*po.API, error) {
	db := GatewayGlobal.DB
	// 查询serviceName对应的id
	var servicePO po.Service
	affected := db.Where("name = ?", serviceName).
		First(&servicePO).RowsAffected
	if affected == 0 {
		// 如果没有找到对应的服务，返回 nil
		return nil, nil
	}
	var apiPO po.API

	// 查询数据库中是否存在相同路径和方法的 API
	affected = db.Where("http_path = ? AND http_method = ? and service_id=?", path, method, servicePO.ID).
		First(&apiPO).RowsAffected

	// 如果存在，更新service_id并返回该 API
	if affected > 0 {
		return &apiPO, nil
	}

	// 如果不存在，返回 nil
	return nil, nil
}

// DeleteAPIsByServiceID
func DeleteAPIsByServiceID(serviceID int64) error {
	// 假设使用 ORM 来删除所有与服务 ID 相关的 API
	var apis []po.API
	db := GatewayGlobal.DB

	// 查询所有与服务 ID 相关的 API
	if err := db.Where("service_id = ?", serviceID).Find(&apis).Error; err != nil {
		return err
	}

	// 删除找到的所有 API
	if err := db.Delete(&apis).Error; err != nil {
		return err
	}

	return nil
}
