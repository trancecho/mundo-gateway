package dto

type ServiceCreateReq struct {
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Password string `json:"password"`
}

type ServiceBeatReq struct {
	ServiceName string `json:"service_name"`
	Address     string `json:"address"`
}
type ServiceUpdateReq struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	Protocol string `json:"protocol"`
}

type AddressDeleteReq struct {
	Id int64 `json:"id"`
}
