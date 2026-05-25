package domain

import (
	"log"
	"strings"

	"github.com/trancecho/mundo-gateway/po"
	"github.com/trancecho/ragnarok/maplist"
)

// APIBOFromPO 将持久化 API 转为内存路由对象。
func APIBOFromPO(api po.API, serviceName string) *APIBO {
	grpcSvc := ""
	grpcMethod := ""
	if api.GrpcMethodMeta.ServiceName != "" {
		grpcSvc = api.GrpcMethodMeta.ServiceName
		grpcMethod = api.GrpcMethodMeta.MethodName
	} else {
		grpcSvc = serviceName
		grpcMethod = api.HttpMethod
	}
	return &APIBO{
		Id:         api.ID,
		HttpPath:   api.HttpPath,
		HttpMethod: api.HttpMethod,
		GrpcMethodMeta: GrpcMethodMetaBO{
			ApiId:       api.ID,
			ServiceName: grpcSvc,
			MethodName:  grpcMethod,
		},
	}
}

// UpsertAPIBOToCache 将单条 API 写入（或覆盖）服务内存缓存。
func UpsertAPIBOToCache(apiPO *po.API, serviceName string) {
	if apiPO == nil || GatewayGlobal == nil {
		return
	}
	GatewayGlobal.RWMutex.Lock()
	defer GatewayGlobal.RWMutex.Unlock()

	for _, serviceBO := range GatewayGlobal.Services.List {
		if serviceBO == nil || serviceBO.ServicePOId != apiPO.ServiceId {
			continue
		}
		if serviceBO.APIs.IsEmpty() {
			serviceBO.APIs = maplist.NewMapList[APIBO]()
		}
		serviceBO.APIs.Add(apiPO.ID, APIBOFromPO(*apiPO, serviceName))
		return
	}
}

// SyncServiceAPIsFromDB 用数据库中的 API 列表刷新指定服务的内存路由表。
func SyncServiceAPIsFromDB(serviceID int64, serviceName string) int {
	if GatewayGlobal == nil {
		return 0
	}
	var apis []po.API
	if err := GatewayGlobal.DB.Where("service_id = ?", serviceID).Find(&apis).Error; err != nil {
		log.Println("SyncServiceAPIsFromDB 查询失败:", err)
		return 0
	}

	GatewayGlobal.RWMutex.Lock()
	defer GatewayGlobal.RWMutex.Unlock()

	for _, serviceBO := range GatewayGlobal.Services.List {
		if serviceBO == nil || serviceBO.ServicePOId != serviceID {
			continue
		}
		newAPIs := maplist.NewMapList[APIBO]()
		for _, api := range apis {
			a := api
			newAPIs.Add(a.ID, APIBOFromPO(a, serviceName))
		}
		serviceBO.APIs = newAPIs
		return len(apis)
	}
	return 0
}

// MatchAPIInMemory 在内存中匹配 HTTP 方法与路径。
func MatchAPIInMemory(serviceBO *ServiceBO, method, path string) *APIBO {
	if serviceBO == nil {
		return nil
	}
	for _, api := range serviceBO.APIs.List {
		if api == nil || api.HttpMethod != method {
			continue
		}
		if api.HttpPath == path {
			return api
		}
		if isPathMatch(api.HttpPath, path) {
			return api
		}
	}
	return nil
}

// ResolveAPI 先查内存，未命中则从 DB 加载并回填缓存（避免「库里有、转发没有」）。
func ResolveAPI(serviceBO *ServiceBO, method, path string) *APIBO {
	if serviceBO == nil {
		return nil
	}
	if api := MatchAPIInMemory(serviceBO, method, path); api != nil {
		return api
	}
	apiPO, err := GetAPIByPathAndMethod(path, method, serviceBO.Name)
	if err != nil || apiPO == nil {
		return nil
	}
	log.Printf("[gateway] 内存未命中，从 DB 回填 API: %s %s %s", serviceBO.Name, method, path)
	UpsertAPIBOToCache(apiPO, serviceBO.Name)
	return MatchAPIInMemory(serviceBO, method, path)
}

// GetServiceByPrefix 按 prefix 查找已加载服务（调用方需已持有读锁或仅在无并发写时使用）。
func GetServiceByPrefix(prefix string) *ServiceBO {
	if GatewayGlobal == nil {
		return nil
	}
	for _, bo := range GatewayGlobal.Services.List {
		if bo != nil && bo.Prefix == prefix {
			return bo
		}
	}
	return nil
}

// ReloadServiceByName 按服务名从 DB 载入到内存（心跳恢复用）。
func ReloadServiceByName(name string) bool {
	if GatewayGlobal == nil {
		return false
	}
	var servicePO po.Service
	if err := GatewayGlobal.DB.Where("name = ? AND available = ?", name, true).
		First(&servicePO).Error; err != nil {
		return false
	}
	return ReloadServiceIntoGateway(servicePO.ID)
}

// ReloadServiceByPrefix 按 prefix 从 DB 载入服务到内存。
func ReloadServiceByPrefix(prefix string) bool {
	if GatewayGlobal == nil {
		return false
	}
	var servicePO po.Service
	if err := GatewayGlobal.DB.Where("prefix = ? AND available = ?", prefix, true).
		First(&servicePO).Error; err != nil {
		return false
	}
	return ReloadServiceIntoGateway(servicePO.ID)
}

// PrefixRegistered 检查 prefix 是否已注册。
func PrefixRegistered(prefix string) bool {
	if GatewayGlobal == nil {
		return false
	}
	for _, p := range GatewayGlobal.Prefixes.List {
		if p != nil && p.Name == prefix {
			return true
		}
	}
	return false
}

func isPathMatch(routePath, actualPath string) bool {
	routeParts := strings.Split(routePath, "/")
	pathParts := strings.Split(actualPath, "/")
	if len(routeParts) != len(pathParts) {
		return false
	}
	for i := 0; i < len(routeParts); i++ {
		if strings.HasPrefix(routeParts[i], ":") {
			continue
		}
		if routeParts[i] != pathParts[i] {
			return false
		}
	}
	return true
}
