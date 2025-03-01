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
	APIs      []API
	Addresses []Address
}
type Prefix struct {
	BaseModel
	Name      string
	ServiceId int64
}

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
