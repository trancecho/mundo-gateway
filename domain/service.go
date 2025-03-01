package domain

import "github.com/trancecho/mundo-gateway/po"

// Service BO 业务对象
type Service struct {
	Id     int64
	Name   string
	Prefix string
	//APIs       []API
	//Addresses  []Address
	curAddress int64
}

// NewService 构造函数
func NewService(id int64, name string, prefix string) *Service {
	return &Service{
		Id:     id,
		Name:   name,
		Prefix: prefix,
		//APIs:       apis,
		//Addresses:  addresses,
		curAddress: 0,
	}
}

// GetNextAddress 使用轮询算法获取下一个后端服务地址
func (s *Service) GetNextAddress(po *po.Service) string {
	// 不使用取模，可读性更强
	if s.curAddress >= int64(len(po.Addresses)) {
		s.curAddress = 0
	}
	address := po.Addresses[s.curAddress].Address
	s.curAddress++
	return address
}

type Prefix struct {
	Id        int64
	Name      string
	ServiceId int64
}

type Address struct {
	Id        int64
	ServiceId int64
	Address   string
}

type API struct {
	Id        int64
	ServiceId int64
	Path      string
	Method    string
}
