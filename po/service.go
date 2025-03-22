package po

type RequestLog struct {
	BaseModel
	RouteId      int64
	RequestType  string
	RequestBody  string
	ResponseBody string
}

// todo基于角色的权限访问控制暂时先不做，交给服务本身
type Service struct {
	BaseModel
	Name     string `json:"name"`   //grpc http
	Prefix   string `json:"prefix"` //grpc  比如/v1/test/ping 就会代理grpc服务的
	Protocol string `json:"protocol"`
	CertPath string `json:"cert_path"` //证书路径
	APIs     []API  `gorm:"foreignKey:ServiceId" json:"APIs"`
	//GrpcAPIs  []GrpcAPI `gorm:"foreignKey:ApiId" json:"grpc_apis"`
	Addresses []Address `gorm:"foreignKey:ServiceId" json:"addresses"`
	Available bool      `json:"available"` // 是否可用
}

//type Prefix struct {
//	BaseModel
//	Name      string
//	ApiId int64
//}

type Address struct {
	BaseModel
	ServiceId int64  `json:"service_id"`
	Address   string `json:"address"`
}

type API struct {
	BaseModel
	ServiceId  int64
	HttpPath   string
	HttpMethod string

	//有过期时间的
	GrpcMethodMeta GrpcMethodMeta `gorm:"foreignKey:ApiId" json:"-"`
}

type GrpcMethodMeta struct {
	BaseModel
	ApiId       int64
	ServiceName string // 服务名
	MethodName  string // 方法名
}

//type GrpcAPI struct {
//	BaseModel
//	ApiId    int64  `json:"service_id"`
//	RequestType  string `json:"request_type"`  // 请求类型 example.GetUserRequest
//	ResponseType string `json:"response_type"` // 响应类型 example.UserResponse
//	GrpcMethod   string `json:"grpc_method"`   // gRPC方法名 GetUser
//	GrpcService  string `json:"grpc_service"`  // gRPC服务名 user.UserService
//}
