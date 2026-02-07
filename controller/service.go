package controller

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-gateway/controller/dto"
	"github.com/trancecho/mundo-gateway/domain"
	"github.com/trancecho/mundo-gateway/util"
	"github.com/trancecho/ragnarok/maplist"
)

//type ServiceDTO struct {
//	Name      string `json:"name"`
//	Prefix    string `json:"prefix"`
//	Protocol  string `json:"protocol"`
//	Addresses []Address
//}

//type Address struct {
//	Id        int64
//	ApiId int64
//	Address   string
//}

func CreateServiceController(c *gin.Context) {
	var req dto.ServiceCreateReq
	c.ShouldBindJSON(&req)
	if req.Name == "" {
		util.ClientError(c, 1, "name不能为空")
		return
	}
	if req.Prefix == "" {
		util.ClientError(c, 2, "prefix不能为空")
		return
	}
	if "/"+req.Name != req.Prefix {
		util.ClientError(c, util.QueryParamError, "prefix必须为/{name}")
		return
	}
	if req.Prefix == "gateway" {
		util.ClientError(c, 300, "prefix不能为gateway ")
	}
	if req.Protocol == "" {
		util.ClientError(c, 310, "protocol不能为空")
		return
	}
	if domain.GatewayGlobal == nil || domain.GatewayGlobal.Redis == nil {
		util.ServerError(c, util.DefaultError, "Redis 未初始化")
		return
	}

	// ✅ Redis 密码校验
	redisPassword, err := domain.GatewayGlobal.Redis.Get(c, "gateway:register:password").Result()
	if err != nil {
		util.ServerError(c, util.DefaultError, "无法读取注册密码，请联系管理员")
		return
	}
	if req.Password != redisPassword {
		util.ClientError(c, util.QueryParamError, "注册密码错误")
		return
	}
	//todo 更新
	ok := domain.LimiterGlobal.AddToWhiteList(c.ClientIP())
	if !ok {
		log.Println("IP白名单添加失败:", c.ClientIP())
	} else {
		log.Println("IP白名单添加成功:", c.ClientIP())
	}

	// 根据协议判断地址是否合规 目前只有http和grpc
	if req.Protocol == "http" {
		// 检查地址是否以 http:// 或 https:// 开头
		if !strings.HasPrefix(req.Address, "http://") && !strings.HasPrefix(req.Address, "https://") {
			util.ServerError(c, 500, "http协议地址不能为空")
			return
		}
	} else if req.Protocol == "grpc" {
		// 修复gRPC地址验证逻辑
		if !strings.HasPrefix(req.Address, "grpc://") && !strings.HasPrefix(req.Address, "grpcs://") {
			util.ServerError(c, 600, "grpc协议地址格式不正确，需以grpc://或grpcs://开头")
			return
		}
	} else {
		util.ServerError(c, 700, "协议不合规")
		return
	}

	servicePO, ok, err := domain.CreateServiceService(&req, c.ClientIP())
	if !ok && err != nil {
		util.ServerError(c, util.ResourceAlreadyExistsWarning, "服务创建失败:"+err.Error())
		return
	}
	//构造 Addresses
	addresses := maplist.NewMapList[domain.Address]()
	for _, addrPO := range servicePO.Addresses {
		addresses.Add(addrPO.ID, &domain.Address{
			Address:   addrPO.Address,
			LastBeat:  time.Now(),
			IsHealthy: true,
		})
	}

	//构造 ServiceBO（完整初始化）
	serviceBO := &domain.ServiceBO{
		ServicePOId: servicePO.ID,
		Prefix:      servicePO.Prefix,
		Name:        servicePO.Name,
		Protocol:    servicePO.Protocol,
		Available:   servicePO.Available,
		Addresses:   addresses,
		APIs:        maplist.NewMapList[domain.APIBO](),
	}

	// 3️⃣ 放入 Gateway
	domain.GatewayGlobal.Services.Add(servicePO.ID, serviceBO)

	// 4️⃣ Prefix 同样一次性加

	domain.GatewayGlobal.Prefixes.Add(serviceBO.Prefix, &domain.Prefix{
		Name:      servicePO.Prefix,
		ServiceId: servicePO.ID,
	})

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

	// 1️⃣ 从 Gateway 取已有 ServiceBO
	serviceBO, ok := domain.GatewayGlobal.Services.Get(servicePO.ID)
	if !ok || serviceBO == nil {
		util.ServerError(c, 900, "Gateway 中不存在该 Service")
		return
	}

	// 2️⃣ 原地更新字段（不破坏运行态）
	serviceBO.Name = servicePO.Name
	serviceBO.Prefix = servicePO.Prefix
	serviceBO.Protocol = servicePO.Protocol
	serviceBO.Available = servicePO.Available

	// 3️⃣ 更新 Prefix 映射（注意：不是 Add serviceID）
	domain.GatewayGlobal.Prefixes.Add(serviceBO.Prefix, &domain.Prefix{
		Name:      servicePO.Prefix,
		ServiceId: servicePO.ID,
	})

	util.Ok(c, "服务更新成功", gin.H{
		"service": servicePO,
	})
}

func DeleteServiceController(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id == 0 {
		util.ServerError(c, 1, "id格式错误或为空")
		return
	}

	// 1️⃣ 先从 Gateway 内存中取出 Service
	serviceBO, ok := domain.GatewayGlobal.Services.Get(id)
	if !ok || serviceBO == nil {
		util.ServerError(c, 404, "服务不存在或已被删除")
		return
	}

	// 2️⃣ 加写锁，防止并发转发
	domain.GatewayGlobal.RWMutex.Lock()
	defer domain.GatewayGlobal.RWMutex.Unlock()

	// 3️⃣ 从内存下线（优先）
	domain.GatewayGlobal.Services.Remove(id)

	domain.GatewayGlobal.Prefixes.Remove(serviceBO.Prefix)

	// 4️⃣ 显式清空运行态结构（语义更干净）
	if serviceBO.Addresses.IsEmpty() == false {
		serviceBO.Addresses.Clear()
	}
	if serviceBO.APIs.IsEmpty() == false {
		serviceBO.APIs.Clear()
	}

	// 5️⃣ 再删 DB（即使 DB 失败，也不影响网关稳定性）
	if err := domain.DeleteAPIsByServiceID(id); err != nil {
		log.Println("删除API失败:", err)
	}

	if ok := domain.DeleteServiceService(id); !ok {
		log.Println("删除Service失败:", id)
	}

	util.Ok(c, "服务及相关API删除成功", nil)
}

func DeleteServiceAddressController(c *gin.Context) {
	addressID, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || addressID == 0 {
		util.ServerError(c, 1, "address id格式错误")
		return
	}

	// 1️⃣ 先删 DB，拿到 serviceID
	serviceID, ok := domain.DeleteAddressService(addressID)
	if !ok {
		util.ServerError(c, 2, "服务地址删除失败")
		return
	}

	// 2️⃣ 从 Gateway 内存中移除（加锁）
	domain.GatewayGlobal.RWMutex.Lock()
	defer domain.GatewayGlobal.RWMutex.Unlock()

	serviceBO, ok := domain.GatewayGlobal.Services.Get(serviceID)
	if !ok || serviceBO == nil {
		util.ServerError(c, 404, "服务不存在")
		return
	}

	if serviceBO.Addresses.IsEmpty() == true {
		util.ServerError(c, 500, "服务地址未初始化")
		return
	}

	// ✅ 正确：用 addressID 删除
	serviceBO.Addresses.Remove(addressID)

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

// 服务心跳
func ServiceAliveSignalController(c *gin.Context) {
	var req dto.ServiceBeatReq
	c.ShouldBindJSON(&req)
	if req.ServiceName == "" || req.Address == "" {
		util.ServerError(c, 100, "服务名和地址不能为空")
		return
	} else {
		boPtr := domain.GetServiceBO(req.ServiceName)
		if boPtr == nil {
			util.ServerError(c, 200, "服务心跳失败")
			return
		}
		if boPtr.GetAddressBO(req.Address) == nil {
			util.ServerError(c, 300, "服务地址不存在")
			return
		}
		// 更新服务的心跳时间
		boPtr.GetAddressBO(req.Address).LastBeat = time.Now()
		log.Println("服务心跳成功", req.Address, req.ServiceName)
	}
}

func ServiceAliveChecker() {
	// 定时检查服务的心跳，避免在遍历时直接修改集合（先收集再注销）
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		threshold := 30 * time.Second
		// 收集需要注销的地址，避免遍历时直接修改集合
		var toUnregister []struct {
			serviceName string
			address     string
		}

		// 只读锁收集信息
		domain.GatewayGlobal.RWMutex.RLock()
		if domain.GatewayGlobal.Services.Size() != 0 {
			for _, serviceBO := range domain.GatewayGlobal.Services.List {
				if serviceBO == nil {
					continue
				}
				// Addresses 可能为自定义集合，保留原有遍历方式并加空检查
				for _, address := range serviceBO.Addresses.List {
					if address == nil {
						continue
					}
					// 如果 LastBeat 为空或超过阈值，认为不可用
					if address.LastBeat.IsZero() || time.Since(address.LastBeat) > threshold {
						toUnregister = append(toUnregister, struct {
							serviceName string
							address     string
						}{serviceBO.Name, address.Address})
					}
				}
			}
		}
		domain.GatewayGlobal.RWMutex.RUnlock()

		// 在读锁之外执行注销操作，避免并发修改导致的竞态/panic
		for _, u := range toUnregister {
			log.Println("服务不可用，开始注销：", u.serviceName, u.address)
			// 保持与原逻辑一致，不关心返回值
			domain.UnregisterServiceService(u.serviceName, u.address)
		}
	}
}

type AddressStatus struct {
	Address   string    `json:"address"`
	IsHealthy bool      `json:"is_healthy"`
	LastBeat  time.Time `json:"last_beat"`
}

type ServiceStatus struct {
	Name      string          `json:"name"`
	Available bool            `json:"available"`
	Addresses []AddressStatus `json:"addresses"`
}

type HealthStatusHandlerReq struct {
	ServiceName string `form:"service_name"`
}

func HealthStatusHandler(c *gin.Context) {
	var req HealthStatusHandlerReq
	c.ShouldBindQuery(&req)
	// 允许不传 service_name：返回所有服务的健康状态；如果传了则同时返回 target
	var statuses []ServiceStatus

	// 读锁保护遍历
	domain.GatewayGlobal.RWMutex.RLock()
	if domain.GatewayGlobal.Services.Size() == 0 {
		domain.GatewayGlobal.RWMutex.RUnlock()
		util.ServerError(c, util.DefaultError, "服务列表为空")
		return
	}
	for _, service := range domain.GatewayGlobal.Services.List {
		if service == nil {
			continue
		}
		var addrStatuses []AddressStatus
		for _, addr := range service.Addresses.List {
			if addr == nil {
				continue
			}
			addrStatuses = append(addrStatuses, AddressStatus{
				Address:   addr.Address,
				IsHealthy: addr.IsHealthy,
				LastBeat:  addr.LastBeat,
			})
		}
		statuses = append(statuses, ServiceStatus{
			Name:      service.Name,
			Available: service.Available,
			Addresses: addrStatuses,
		})
	}
	domain.GatewayGlobal.RWMutex.RUnlock()

	// 如果传了 service_name，则尝试获取该服务的下一个目标地址，否则 target 为 nil
	var target interface{}
	if req.ServiceName != "" {
		servicebo := domain.GetServiceBO(req.ServiceName)
		if servicebo == nil {
			util.ClientError(c, util.QueryParamError, "服务不存在")
			return
		}
		target = servicebo.GetNextAddress()
	} else {
		target = nil
	}

	util.Ok(c, "服务健康状态", gin.H{
		"services": statuses,
		"target":   target,
	})
}
