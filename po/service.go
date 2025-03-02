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
	Name      string
	Prefix    string
	Protocol  string
	APIs      []API     `gorm:"foreignKey:ServiceId"`
	GrpcAPIs  []GrpcAPI `gorm:"foreignKey:ServiceId"`
	Addresses []Address `gorm:"foreignKey:ServiceId"`
}

//type Prefix struct {
//	BaseModel
//	Name      string
//	ServiceId int64
//}

type Address struct {
	BaseModel
	ServiceId int64
	Address   string
}

type API struct {
	BaseModel
	ServiceId int64
	Path      string
	Method    string
}
type GrpcAPI struct {
	BaseModel
	ServiceId    int64
	RequestType  string // 请求类型 example.GetUserRequest
	ResponseType string // 响应类型 example.UserResponse
	GrpcMethod   string // gRPC方法名 GetUser
	GrpcService  string // gRPC服务名 user.UserService
}
