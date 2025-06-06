package models

import "time"

type LimiterInterface interface {
	IsBlackListed(ip string) bool
	AddToBlackList(ip string) bool
	SyncBlackList() error
	AllowRequest(ip string) bool
	StartCacheRefresher(duration time.Duration)
}
