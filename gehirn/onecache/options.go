/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/11/5 -- 15:49
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: options.go
*/

package onecache

type Option func(*Cache)

// IRateLimiter
// Interface of rate limiter
// TakeAvailable
// is a non-blocking function, it takes up to count immediately available tokens from the bucket.
// It returns the number of tokens removed, or zero if there are no available tokens.
type IRateLimiter interface {
	TakeAvailable(count int64) int64
}

func FallbackWhenError() Option {
	return Option(func(m *Cache) {
		m.fallbackWhenError = true
	})
}

func RateLimiter(limiter IRateLimiter) Option {
	return Option(func(m *Cache) {
		m.limiter = limiter
	})
}
