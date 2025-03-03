package po

type RequestLog struct {
	BaseModel
	RouteId      int64
	RequestType  string
	RequestBody  string
	ResponseBody string
}
type Service struct {
	BaseModel
	Name      string    `json:"name"`
	Prefix    string    `json:"prefix"`
	Protocol  string    `json:"protocol"`
	APIs      []API     `gorm:"foreignKey:ServiceId" json:"APIs"`
	GrpcAPIs  []GrpcAPI `gorm:"foreignKey:ServiceId" json:"grpc_apis"`
	Addresses []Address `gorm:"foreignKey:ServiceId" json:"addresses"`
}

//type Prefix struct {
//	BaseModel
//	Name      string
//	ServiceId int64
//}

type Address struct {
	BaseModel
	ServiceId int64  `json:"service_id"`
	Address   string `json:"address"`
}

type API struct {
	BaseModel
	ServiceId int64
	Path      string
	Method    string
}
type GrpcAPI struct {
	BaseModel
	ServiceId    int64  `json:"service_id"`
	RequestType  string `json:"request_type"`  // 请求类型 example.GetUserRequest
	ResponseType string `json:"response_type"` // 响应类型 example.UserResponse
	GrpcMethod   string `json:"grpc_method"`   // gRPC方法名 GetUser
	GrpcService  string `json:"grpc_service"`  // gRPC服务名 user.UserService
}
