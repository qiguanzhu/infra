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
 @Time    : 2024/11/1 -- 18:01
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: cache.go
*/

package onecache

import (
	"context"
	"errors"
	"time"
)

type (
	KeyFunc      func() string
	KeyMultiFunc func(interface{}) string
)

type (
	FallbackFunc      func() (interface{}, error)
	FallbackMultiFunc func(interface{}) (interface{}, error)
)

var (
	ErrBadIdsType              = errors.New("cache: ids must be slice")
	ErrBadSrcType              = errors.New("cache: src must be a slice or pointer-to-slice")
	ErrBadDstType              = errors.New("cache: dst must be a pointer")
	ErrBadDstMapType           = errors.New("cache: dst must be a map or pointer-to-map")
	ErrBadDstMapValue          = errors.New("cache: dst must not be a nil map")
	ErrSrcDstTypeMismatch      = errors.New("cache: type of fallback result mismatch error")
	ErrBadFallbackResultType   = errors.New("cache: type of fallbackResultV must be map")
	ErrDstFallbackTypeMismatch = errors.New("cache: key type and value type of fallbackResult is not equal to map")
)

type CacheBox interface {
	Get(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration) error
	MustGet(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration)

	GetMulti(ctx context.Context, ids interface{},
		keyFunc KeyMultiFunc,
		fallbackFunc FallbackMultiFunc,
		dstMap interface{},
		ttl *time.Duration) error
	MustGetMulti(ctx context.Context, ids interface{},
		keyFunc KeyMultiFunc,
		fallbackFunc FallbackMultiFunc,
		dstMap interface{},
		ttl *time.Duration)

	Evict(ctx context.Context, keyFunc KeyFunc) error
	MustEvict(ctx context.Context, keyFunc KeyFunc)

	EvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc) error
	MustEvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc)

	Refresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration) error
	MustRefresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration)

	RefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc, fallbackFunc FallbackMultiFunc, i interface{}, ttl *time.Duration) error
	MustRefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc, fallbackFunc FallbackMultiFunc, i interface{}, ttl *time.Duration)
}
