package dto

type APICreateReq struct {
	ServiceName string `json:"service_name"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	GrpcService string `json:"grpc_service"`
	GrpcMethod  string `json:"grpc_method"  `
}

type APIUpdateReq struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Method         string `json:"method"`
	GrpcMethodMeta struct {
		ServiceName string `json:"service_name"`
		MethodName  string `json:"method_name"`
	} `json:"grpc_method_meta"`
}
type APIDeleteReq struct {
	Id int64 `json:"id"`
}
