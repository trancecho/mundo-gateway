package domain

import (
	"github.com/trancecho/mundo-gateway/controller"
	"github.com/trancecho/mundo-gateway/po"
	"log"
)

// ServiceBO BO 业务对象
type ServiceBO struct {
	ServicePOId int64
	Prefix      string
	Name        string
	Addresses   []string
	Protocol    string
	curAddress  int64
}

// NewService 构造函数
func NewService(servicePOId int64) *ServiceBO {
	return &ServiceBO{ServicePOId: servicePOId}
}

// GetNextAddress 使用轮询算法获取下一个后端服务地址
func (s *ServiceBO) GetNextAddress() string {
	// 不使用取模，可读性更强
	if s.curAddress >= int64(len(s.Addresses)) {
		s.curAddress = 0
	}
	address := s.Addresses[s.curAddress]
	s.curAddress++
	return address
}

// 创建服务
func CreateServiceService(dto *controller.ServiceCreateReq) (*po.Service, bool) {
	var err error
	var servicePO po.Service
	servicePO.Name = dto.Name
	// 根据name查找service
	affected := GatewayGlobal.DB.Find(&servicePO).RowsAffected
	if affected > 0 {
		// 找其地址列表
		var addresses []po.Address
		GatewayGlobal.DB.Where("service_id = ?", servicePO.ID).Find(&addresses)
		addresses = append(addresses, po.Address{ServiceId: servicePO.ID, Address: dto.Address})
		// 更新地址列表
		err = GatewayGlobal.DB.Save(&addresses).Error
		if err != nil {
			log.Println("服务新地址保存失败", err)
			return nil, false
		}
	} else {
		// 说明没有service，需要新建
		servicePO.Prefix = dto.Prefix
		servicePO.Protocol = dto.Protocol
		// 创建service
		err = GatewayGlobal.DB.Create(&servicePO).Error
		if err != nil {
			log.Println("服务创建失败", err)
			return nil, false
		}
	}
	return &servicePO, true
}

// 更新服务
func UpdateServiceService(dto *controller.ServiceUpdateReq) (*po.Service, bool) {
	var err error
	var servicePO po.Service
	servicePO.ID = dto.Id
	// 根据id查找service
	affected := GatewayGlobal.DB.Find(&servicePO).RowsAffected
	if affected == 0 {
		log.Println("服务不存在")
		return nil, false
	}
	servicePO.Name = dto.Name
	servicePO.Prefix = dto.Prefix
	servicePO.Protocol = dto.Protocol
	// 更新service
	err = GatewayGlobal.DB.Save(&servicePO).Error
	if err != nil {
		log.Println("服务更新失败", err)
		return nil, false
	}
	// 更新地址
	var addresses []po.Address
	GatewayGlobal.DB.Where("service_id = ?", servicePO.ID).Find(&addresses)
	addresses = append(addresses, po.Address{ServiceId: servicePO.ID, Address: dto.Address})
	err = GatewayGlobal.DB.Save(&addresses).Error
	if err != nil {
		log.Println("服务新地址保存失败", err)
		return nil, false
	}
	return &servicePO, true
}

// 删除服务
func DeleteServiceService(id int64) bool {
	var err error
	var servicePO po.Service
	servicePO.ID = id
	// 根据id查找service
	affected := GatewayGlobal.DB.Find(&servicePO).RowsAffected
	if affected == 0 {
		log.Println("服务不存在")
		return false
	}
	// 删除service
	err = GatewayGlobal.DB.Delete(&servicePO).Error
	if err != nil {
		log.Println("服务删除失败", err)
		return false
	}
	// 删除地址
	var addresses []po.Address
	GatewayGlobal.DB.Where("service_id = ?", servicePO.ID).Find(&addresses)
	err = GatewayGlobal.DB.Delete(&addresses).Error
	if err != nil {
		log.Println("服务地址删除失败", err)
		return false
	}
	return true
}

// 查询服务
func GetServiceService(id int64) (*po.Service, bool) {
	var servicePO po.Service
	servicePO.ID = id
	// 根据id查找service
	affected := GatewayGlobal.DB.Find(&servicePO).RowsAffected
	if affected == 0 {
		log.Println("服务不存在")
		return nil, false
	}
	// 找其地址列表
	var addresses []po.Address
	GatewayGlobal.DB.Where("service_id = ?", servicePO.ID).Find(&addresses)
	servicePO.Addresses = addresses
	return &servicePO, true
}

// 查询服务列表
func ListServicesService() ([]po.Service, bool) {
	var services []po.Service
	// 查询所有service
	GatewayGlobal.DB.Find(&services)
	// 查询所有service的地址
	for i := range services {
		var addresses []po.Address
		GatewayGlobal.DB.Where("service_id = ?", services[i].ID).Find(&addresses)
		services[i].Addresses = addresses
	}
	return services, true
}

type Prefix struct {
	Id        int64
	Name      string
	ServiceId int64
}

type Address struct {
	Id        int64
	ServiceId int64
	Address   string
}
