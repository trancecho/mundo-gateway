package dto

type APICreateReq struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type APIUpdateReq struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}
type APIDeleteReq struct {
	Id int64 `json:"id"`
}
