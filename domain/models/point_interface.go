package models

import "github.com/trancecho/mundo-gateway/proto/point/v1"

// PointClientInterface 定义积分客户端的接口
type PointClientInterface interface {
	DialToPointsSystem(token string, userID string, points int64, experience int64, reason string) (string, error)
	SignToPointsSystem(token string, userID string) (*point.CommonResponse, error)
	GetPointsInfo(token string, userID string) (*point.UserInfo, error)
	AdminStats(token string, userID string) (*point.AdminStats, error)
}
