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
 @Time    : 2024/11/1 -- 18:13
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: store.go
*/

package store

import (
	"context"
	"errors"
	"time"
)

var (
	ErrBadSrcType         = errors.New("cache store: src must be a slice or pointer-to-slice")
	ErrKeysLengthNotMatch = errors.New("cache store: keys and src slices have different length")
	ErrBadDstType         = errors.New("cache store: dst must be a pointer")
	ErrBadDstMapType      = errors.New("cache store: dst must be a map or pointer-to-map")
	ErrBadDstMapValue     = errors.New("cache store: dst must not be a nil map")
	ErrSrcDstTypeMismatch = errors.New("cache store: type of fallback result mismatch error")
)

type Store interface {
	// Get from backend and store the value in dst.
	Get(ctx context.Context, key string, dst interface{}) error
	MustGet(ctx context.Context, key string, dst interface{})

	// GetMulti get from backend and store the value in dst,
	// dst must be a map or pointer-to-map in format of
	// map[string]interface{} where key is the passed key if
	// key exists.
	GetMulti(ctx context.Context, keys []string, dstMap interface{}) error
	MustGetMulti(ctx context.Context, keys []string, dstMap interface{})

	// Exists ask for backend whether specified item exists.
	Exists(ctx context.Context, key string) (bool, error)
	MustExists(ctx context.Context, key string) bool

	// ExistsMulti ask for backend whether specified items exists.
	ExistsMulti(ctx context.Context, keys ...string) ([]bool, error)
	MustExistsMulti(ctx context.Context, keys ...string) []bool

	// Set set key and value with timeout.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	MustSet(ctx context.Context, key string, value interface{}, ttl time.Duration)

	// SetMulti set keys with values.
	SetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) error
	MustSetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration)

	// Delete remove the specified item by key.
	Delete(ctx context.Context, keys ...string) error
	MustDelete(ctx context.Context, keys ...string)
}
