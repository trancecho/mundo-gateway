package domain

import (
	"errors"
	"github.com/trancecho/mundo-gateway/controller/dto"
	"github.com/trancecho/mundo-gateway/po"
	"gorm.io/gorm"
	"log"
	"time"
)

// ServiceBO BO 业务对象
type ServiceBO struct {
	ServicePOId int64
	Prefix      string
	Name        string
	Addresses   []*Address
	APIs        []APIBO
	Protocol    string
	curAddress  int64
	Available   bool
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
	return address.Address
}

func GetServiceBO(name string) *ServiceBO {
	for _, service := range GatewayGlobal.Services {
		if service.Name == name {
			return &service
		}
	}
	return nil
}

func (s *ServiceBO) GetAddressBO(address string) *Address {
	for _, addr := range s.Addresses {
		if addr.Address == address {
			return addr
		}
	}
	return nil
}

func UnregisterServiceService(name string, address string) bool {
	// 开启事务
	tx := GatewayGlobal.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("事务回滚，发生错误:", r)
		}
	}()

	// 根据 name 查找 service
	var servicePO po.Service
	if err := tx.Where("name = ?", name).First(&servicePO).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("服务不存在:", name)
		} else {
			log.Println("查询服务失败:", err)
		}
		tx.Rollback()
		return false
	}

	// 根据 service_id 和 address 查找 addresses
	var addresses []po.Address
	if err := tx.Where("service_id = ? AND address = ?", servicePO.ID, address).Find(&addresses).Error; err != nil {
		log.Println("查询地址失败:", err)
		tx.Rollback()
		return false
	}

	if len(addresses) == 0 {
		log.Println("地址不存在:", address)
		tx.Rollback()
		return false
	}

	// 删除地址记录
	if err := tx.Where("service_id = ? AND address = ?", servicePO.ID, address).Delete(&po.Address{}).Error; err != nil {
		log.Println("删除地址失败:", err)
		tx.Rollback()
		return false
	}

	// 检查是否还有其它地址
	var remainingAddresses []po.Address
	if err := tx.Where("service_id = ?", servicePO.ID).Find(&remainingAddresses).Error; err != nil {
		log.Println("查询剩余地址失败:", err)
		tx.Rollback()
		return false
	}

	// 如果地址列表是最后一个地址，将 service 的 available 字段置为 false
	if len(remainingAddresses) == 0 {
		if err := tx.Model(&po.Service{}).
			Where("id = ?", servicePO.ID).
			Update("available", false).Error; err != nil {
			log.Println("更新服务状态失败:", err)
			tx.Rollback()
			return false
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Println("提交事务失败:", err)
		tx.Rollback()
		return false
	}

	// 更新内存中的数据
	GatewayGlobal.FlushGateway()

	// 日志记录
	log.Println("服务地址删除成功:", servicePO.Name, address)
	return true
}

// 创建服务
func CreateServiceService(dto *dto.ServiceCreateReq) (*po.Service, bool, error) {
	var err error
	var servicePO po.Service
	// 根据name查找service
	affected := GatewayGlobal.DB.Where("prefix=?", dto.Prefix).
		Find(&servicePO).RowsAffected
	log.Println(dto)
	log.Println("affected", affected, servicePO)
	// 如果 name 已经存在，说明是更新地址。
	if affected > 0 {
		// 如果protocl不同，则报错
		if servicePO.Protocol != dto.Protocol {
			log.Println("服务协议不同")
			return nil, false, errors.New("匹配到对应服务。但是服务协议不同")
		}
		// 激活已存在的
		if servicePO.Available == false {
			GatewayGlobal.DB.Model(&po.Service{}).Where("id = ?", servicePO.ID).Update("available", true)
		}
		// 找其地址列表
		var addresses []po.Address
		GatewayGlobal.DB.Where("service_id = ?", servicePO.ID).Find(&addresses)
		// 查看是否有该地址
		for _, address := range addresses {
			if address.Address == dto.Address {
				log.Println("地址已存在")
				return nil, false, errors.New("地址已存在")
			}
		}
		addresses = append(addresses, po.Address{ServiceId: servicePO.ID, Address: dto.Address})
		// 更新地址列表
		err = GatewayGlobal.DB.Save(&addresses).Error
		if err != nil {
			log.Println("服务新地址保存失败", err)
			return nil, false, err
		}

		// ✅ 加载地址（这是关键补充步骤）
		err = GatewayGlobal.DB.Preload("Addresses").
			Where("id = ?", servicePO.ID).First(&servicePO).Error
		if err != nil {
			log.Println("刷新地址失败", err)
			return nil, false, err
		}

		//更新po的地址列表
		servicePO.Addresses = addresses
	} else {
		// 说明没有service，需要新建
		servicePO.Name = dto.Name
		servicePO.Prefix = dto.Prefix
		servicePO.Protocol = dto.Protocol
		servicePO.Available = true
		servicePO.Addresses = []po.Address{{Address: dto.Address}}
		// 创建service
		err = GatewayGlobal.DB.Create(&servicePO).Error
		if err != nil {
			log.Println("服务创建失败", err)
			return nil, false, err
		}
	}
	return &servicePO, true, nil
}

// 更新服务
func UpdateServiceService(dto *dto.ServiceUpdateReq) (*po.Service, bool) {
	var err error
	var servicePO po.Service
	servicePO.ID = dto.Id
	// 根据id查找service
	affected := GatewayGlobal.DB.Where("id=?", servicePO.ID).First(&servicePO).RowsAffected
	if affected == 0 {
		log.Println("service不存在")
		return nil, false
	}
	if dto.Name != "" {
		servicePO.Name = dto.Name
	}
	if dto.Prefix != "" {
		servicePO.Prefix = dto.Prefix
	}
	if dto.Protocol != "" {
		servicePO.Protocol = dto.Protocol
	}
	// 更新service
	err = GatewayGlobal.DB.Save(&servicePO).Error
	if err != nil {
		log.Println("服务更新失败", err)
		return nil, false
	}
	// 更新地址
	//var addresses []po.Address
	//GatewayGlobal.DB.Where("service_id = ?", servicePO.ID).Find(&addresses)
	//addresses = append(addresses, po.Address{ApiId: servicePO.ID, Address: dto.Address})
	//err = GatewayGlobal.DB.Save(&addresses).Error
	//if err != nil {
	//	log.Println("服务新地址保存失败", err)
	//	return nil, false
	//}
	return &servicePO, true
}

func DeleteAddressService(id int64) bool {
	var err error
	var addressPO po.Address
	addressPO.ID = id
	// 根据id查找address
	affected := GatewayGlobal.DB.Where("id=?", addressPO.ID).
		Find(&addressPO).RowsAffected
	if affected == 0 {
		log.Println("地址不存在")
		return false
	}
	// 删除address
	err = GatewayGlobal.DB.Delete(&addressPO).Error
	if err != nil {
		log.Println("地址删除失败", err)
		return false
	}
	return true
}

// 删除服务
func DeleteServiceService(id int64) bool {
	var err error
	var servicePO po.Service
	servicePO.ID = id
	// 根据id查找service
	affected := GatewayGlobal.DB.Where("id=?", servicePO.ID).
		Find(&servicePO).RowsAffected
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

// 查询服务
func GetServiceByName(name string) (*po.Service, bool) {
	var servicePO po.Service
	servicePO.Name = name
	// 根据id查找service
	affected := GatewayGlobal.DB.Where("name = ?", name).First(&servicePO).RowsAffected
	if affected == 0 {
		log.Println("服务不存在:", name)
		return nil, false
	}
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
	Address   string
	LastBeat  time.Time
	IsHealthy bool
}
