package i

import "time"

type ILimiter interface {
	IsBlackListed(ip string) bool
	AddToBlackList(ip string) bool
	SyncBlackList() error
	AllowIp(ip string) bool
	StartCacheRefresher(duration time.Duration)
	StartIpRateRecorderFlusher()
	IsWhiteListed(ip string) bool
	AddToWhiteList(ip string) bool
}
