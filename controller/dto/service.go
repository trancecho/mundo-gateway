package dto

type ServiceCreateReq struct {
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
}

type ServiceUpdateReq struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
}
