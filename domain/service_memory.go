package domain

import (
	"log"
	"time"

	"github.com/trancecho/mundo-gateway/po"
	"github.com/trancecho/ragnarok/maplist"
)

// snapshotAddressBeats 在刷新前保存各地址最近心跳时间。
func snapshotAddressBeats() map[int64]map[string]time.Time {
	out := make(map[int64]map[string]time.Time)
	if GatewayGlobal == nil {
		return out
	}
	for _, svc := range GatewayGlobal.Services.List {
		if svc == nil {
			continue
		}
		m := make(map[string]time.Time)
		for _, addr := range svc.Addresses.List {
			if addr != nil {
				m[addr.Address] = addr.LastBeat
			}
		}
		out[svc.ServicePOId] = m
	}
	return out
}

func lastBeatFor(beats map[string]time.Time, address string) time.Time {
	if t, ok := beats[address]; ok && !t.IsZero() {
		return t
	}
	return time.Now()
}

// RemoveAddressFromMemory 从内存摘掉单个后端地址（不触发全量 FlushGateway）。
func RemoveAddressFromMemory(serviceName, address string) {
	if GatewayGlobal == nil {
		return
	}
	GatewayGlobal.RWMutex.Lock()
	defer GatewayGlobal.RWMutex.Unlock()

	svc := GetServiceBO(serviceName)
	if svc == nil {
		return
	}

	newAddrs := maplist.NewMapList[Address]()
	for key, addr := range svc.Addresses.GetMap() {
		if addr != nil && addr.Address != address {
			newAddrs.Add(key, addr)
		}
	}
	svc.Addresses = newAddrs

	if !svc.Addresses.IsEmpty() {
		return
	}

	newServices := maplist.NewMapList[ServiceBO]()
	for _, s := range GatewayGlobal.Services.List {
		if s != nil && s.ServicePOId != svc.ServicePOId {
			newServices.Add(s.ServicePOId, s)
		}
	}
	GatewayGlobal.Services = newServices

	newPrefixes := maplist.NewMapStringList[Prefix]()
	for _, p := range GatewayGlobal.Prefixes.List {
		if p != nil && p.ServiceId != svc.ServicePOId {
			newPrefixes.Add(p.Name, p)
		}
	}
	GatewayGlobal.Prefixes = newPrefixes
	log.Println("[gateway] 服务无可用地址，已从内存卸载:", serviceName)
}

// RecoverOrAddAddress 心跳到达时，若地址不在内存中（可能因之前下线被软删除），
// 从 DB 恢复或新建地址记录，并加入内存 ServiceBO。
func RecoverOrAddAddress(bo *ServiceBO, address string) bool {
	if bo == nil || GatewayGlobal == nil {
		return false
	}
	// 尝试查找已存在但被软删除的地址记录
	var addrPO po.Address
	err := GatewayGlobal.DB.Unscoped().
		Where("service_id = ? AND address = ?", bo.ServicePOId, address).
		First(&addrPO).Error
	if err == nil {
		// 记录存在（可能被软删除），恢复或确保未删除
		if addrPO.DeletedAt.Valid {
			GatewayGlobal.DB.Unscoped().
				Model(&addrPO).
				Update("deleted_at", nil)
		}
	} else {
		// 记录不存在，新建
		addrPO = po.Address{
			ServiceId: bo.ServicePOId,
			Address:   address,
		}
		if err := GatewayGlobal.DB.Create(&addrPO).Error; err != nil {
			log.Println("[gateway] 创建心跳地址失败:", err)
			return false
		}
	}

	// 加入内存
	GatewayGlobal.RWMutex.Lock()
	bo.Addresses.Add(addrPO.ID, &Address{
		Address:   address,
		LastBeat:  time.Now(),
		IsHealthy: true,
	})
	GatewayGlobal.RWMutex.Unlock()
	log.Printf("[gateway] 心跳恢复地址: %s -> %s", bo.Name, address)
	return true
}

// ReloadServiceIntoGateway 将单个服务从 DB 载入或更新到内存。
func ReloadServiceIntoGateway(serviceID int64) bool {
	if GatewayGlobal == nil {
		return false
	}
	var servicePO po.Service
	err := GatewayGlobal.DB.Preload("Addresses").
		Preload("APIs.GrpcMethodMeta").
		Where("id = ?", serviceID).
		First(&servicePO).Error
	if err != nil {
		return false
	}

	beats := snapshotAddressBeats()
	oldBeats := beats[serviceID]

	addresses := maplist.NewMapList[Address]()
	for _, addr := range servicePO.Addresses {
		addresses.Add(addr.ID, &Address{
			Address:   addr.Address,
			LastBeat:  lastBeatFor(oldBeats, addr.Address),
			IsHealthy: true,
		})
	}
	apis := maplist.NewMapList[APIBO]()
	for _, api := range servicePO.APIs {
		a := api
		apis.Add(a.ID, APIBOFromPO(a, servicePO.Name))
	}

	bo := &ServiceBO{
		ServicePOId: servicePO.ID,
		Prefix:      servicePO.Prefix,
		Name:        servicePO.Name,
		Protocol:    servicePO.Protocol,
		Addresses:   addresses,
		APIs:        apis,
		curAddress:  0,
		Available:   servicePO.Available,
	}

	GatewayGlobal.RWMutex.Lock()
	defer GatewayGlobal.RWMutex.Unlock()
	GatewayGlobal.Services.Add(servicePO.ID, bo)
	GatewayGlobal.Prefixes.Add(servicePO.Prefix, &Prefix{
		Name:      servicePO.Prefix,
		ServiceId: servicePO.ID,
	})
	return true
}
