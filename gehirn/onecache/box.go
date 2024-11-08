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
 @Time    : 2024/11/1 -- 18:12
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: box.go
*/

package onecache

import (
	"context"
	"github.com/qiguanzhu/infra/gehirn"
	"github.com/qiguanzhu/infra/gehirn/onecache/store"
	"github.com/qiguanzhu/infra/nerv/magi/xreflect"
	"reflect"
	"time"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	store             store.Store
	expire            time.Duration
	fallbackWhenError bool
	limiter           IRateLimiter
}

var _ CacheBox = (*Cache)(nil)

func NewCache(store store.Store, expire time.Duration, opts ...Option) *Cache {
	r := &Cache{
		store:  store,
		expire: expire,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

func (c *Cache) canFallbackWhenError() bool {
	if !c.fallbackWhenError {
		return false
	}
	if c.limiter == nil {
		return true
	}
	return c.limiter.TakeAvailable(1) != 0
}

func (c *Cache) get(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, expire time.Duration) error {
	dstPtrV := reflect.ValueOf(dst)
	if dstPtrV.Kind() != reflect.Ptr {
		return ErrBadDstType
	}

	key := keyFunc()
	err := c.store.Get(ctx, key, dst)
	if err == nil {
		return nil
	} else if err == redis.Nil || err == gcache.KeyNotFoundError || c.canFallbackWhenError() {
		fallbackResult, err := fallbackFunc()
		if err != nil {
			return err
		}

		// check nil
		fV := reflect.ValueOf(fallbackResult)
		if !fV.IsValid() || fV.Kind() == reflect.Ptr && fV.IsNil() {
			return nil
		}

		// check fallback type
		dstV := reflect.Indirect(reflect.ValueOf(dst))
		// 此处不需要强一致类型 **struct 和 *struct 是可以的
		left := xreflect.RecursiveIndirectType(dstV.Type())
		right := xreflect.RecursiveIndirectType(fV.Type())
		if left.Kind() != right.Kind() {
			return ErrDstFallbackTypeMismatch
		}
		dstV.Set(fV)

		_ = c.store.Set(ctx, key, dst, expire)
	} else {
		return err
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.get(ctx, keyFunc, fallbackFunc, dst, expire)
}

func (c *Cache) MustGet(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, dst interface{}, ttl *time.Duration) {
	gehirn.PanicIf(c.Get(ctx, keyFunc, fallbackFunc, dst, ttl))
}

func (c *Cache) getMulti(ctx context.Context, ids interface{},
	keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc,
	dstMap interface{}, expire time.Duration) error {
	// dst must be a map or pointer-to-map
	dstPtrV := reflect.ValueOf(dstMap)
	dstV := reflect.Indirect(dstPtrV)
	if dstV.Kind() != reflect.Map {
		return ErrBadDstMapType
	}

	// nil map
	if dstPtrV.Kind() != reflect.Ptr && dstV.IsNil() {
		return ErrBadDstMapValue
	}

	// auto generate map
	if dstPtrV.Kind() == reflect.Ptr && dstV.IsNil() {
		m := reflect.MakeMap(reflect.MapOf(dstV.Type().Key(), dstV.Type().Elem()))
		dstV.Set(m)
	}

	dstType := xreflect.RecursiveIndirectType(reflect.TypeOf(dstMap))
	mapKeyType := dstType.Key()
	mapValueType := dstType.Elem()

	// check slice type
	idsT := reflect.ValueOf(ids)
	if idsT.Kind() != reflect.Slice {
		return ErrBadSrcType
	}

	length := idsT.Len()
	if length == 0 {
		return nil
	}

	// check type of id vs key type of dstMap
	left := idsT.Type().Elem()
	right := mapKeyType
	if left.Kind() != right.Kind() {
		return ErrDstFallbackTypeMismatch
	}

	actualIds := reflect.MakeSlice(reflect.SliceOf(left), length, length)
	for i := 0; i < length; i++ {
		actualIds.Index(i).Set(idsT.Index(i))
	}

	// generate keys
	keys := make([]string, length)
	revertKeyMap := make(map[string]interface{}) // key -> id
	for i := 0; i < length; i++ {
		id := actualIds.Index(i).Interface()
		key := keyFunc(id)
		keys[i] = key
		revertKeyMap[key] = id
	}

	// 这里要处理 *map 的情况
	dstValueT := reflect.TypeOf(dstMap)
	if dstPtrV.Kind() == reflect.Ptr {
		dstValueT = dstValueT.Elem()
	}
	dstValueT = dstValueT.Elem()

	// get from cache
	cacheDstV := reflect.MakeMap(reflect.MapOf(reflect.ValueOf("").Type(), dstValueT))
	err := c.store.GetMulti(ctx, keys, cacheDstV.Interface())
	if err != nil && !c.canFallbackWhenError() {
		return err
	}

	// check miss
	cacheMissIdsV := reflect.MakeSlice(reflect.SliceOf(mapKeyType), 0, 8)
	for _, key := range keys {
		id := revertKeyMap[key]

		// find value from map
		v := cacheDstV.MapIndex(reflect.ValueOf(key))
		if !v.IsValid() {
			cacheMissIdsV = reflect.Append(cacheMissIdsV, reflect.ValueOf(id))
		} else {
			dstV.SetMapIndex(reflect.ValueOf(id), v)
		}
	}

	// fallback && set cache
	if cacheMissIdsV.Len() == 0 {
		return nil
	}

	fallbackResult, err := fallbackFunc(cacheMissIdsV.Interface())
	if err != nil {
		return err
	}

	// check map type
	fallbackResultV := reflect.ValueOf(fallbackResult)
	if fallbackResultV.Kind() != reflect.Map {
		return ErrBadFallbackResultType
	}

	if fallbackResultV.Type().Key().Kind() != mapKeyType.Kind() ||
		fallbackResultV.Type().Elem().Kind() != mapValueType.Kind() {
		return ErrDstFallbackTypeMismatch
	}

	if fallbackResultV.Len() == 0 {
		return nil
	}

	cacheKeys := make([]string, 0)
	cacheSrcV := reflect.MakeSlice(reflect.SliceOf(dstValueT), 0, 8)

	left = xreflect.RecursiveIndirectType(xreflect.RecursiveIndirect(reflect.ValueOf(dstMap)).Type().Elem())
	mapKeys := fallbackResultV.MapKeys()

	for _, keyV := range mapKeys {
		valueV := fallbackResultV.MapIndex(keyV)

		if !valueV.IsValid() || valueV.Kind() == reflect.Ptr && valueV.IsNil() {
			continue
		}

		key := keyFunc(keyV.Interface())
		cacheKeys = append(cacheKeys, key)
		cacheSrcV = reflect.Append(cacheSrcV, valueV)
		dstV.SetMapIndex(keyV, valueV)
	}

	if cacheSrcV.Len() > 0 {
		_ = c.store.SetMulti(ctx, cacheKeys, cacheSrcV.Interface(), expire)
	}

	return nil
}

func (c *Cache) GetMulti(ctx context.Context, ids interface{},
	keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc,
	dstMap interface{},
	ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.getMulti(ctx, ids, keyFunc, fallbackFunc, dstMap, expire)
}

func (c *Cache) MustGetMulti(ctx context.Context, ids interface{},
	keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc,
	dstMap interface{},
	ttl *time.Duration) {
	gehirn.PanicIf(c.GetMulti(ctx, ids, keyFunc, fallbackFunc, dstMap, ttl))
}

func (c *Cache) Evict(ctx context.Context, keyFunc KeyFunc) error {
	return c.store.Delete(ctx, keyFunc())
}

func (c *Cache) MustEvict(ctx context.Context, keyFunc KeyFunc) {
	gehirn.PanicIf(c.Evict(ctx, keyFunc))
}

func (c *Cache) EvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc) error {
	// check slice type
	idsV := reflect.ValueOf(ids)
	if idsV.Kind() != reflect.Slice {
		return ErrBadIdsType
	}

	length := idsV.Len()
	if length == 0 {
		return nil
	}

	actualIds := reflect.MakeSlice(reflect.SliceOf(idsV.Type().Elem()), length, length)
	for i := 0; i < length; i++ {
		actualIds.Index(i).Set(idsV.Index(i))
	}

	keys := make([]string, length)
	for i := 0; i < length; i++ {
		keys[i] = keyFunc(actualIds.Index(i).Interface())
	}

	return c.store.Delete(ctx, keys...)
}

func (c *Cache) MustEvictMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc) {
	gehirn.PanicIf(c.EvictMulti(ctx, ids, keyFunc))
}

func (c *Cache) refresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, p reflect.Type, expire time.Duration) error {
	key := keyFunc()
	err := c.store.Delete(ctx, key)
	if err != nil {
		return err
	}

	fallbackResult, err := fallbackFunc()
	if err != nil {
		return err
	}

	// check nil
	fV := reflect.ValueOf(fallbackResult)
	if !fV.IsValid() || fV.Kind() == reflect.Ptr && fV.IsNil() {
		return nil
	}

	// 右侧 fallback 可以使用指针
	left := p
	right := reflect.Indirect(fV).Type()
	if left.Kind() != right.Kind() {
		return ErrDstFallbackTypeMismatch
	}

	err = c.store.Set(ctx, key, fallbackResult, expire)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Refresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.refresh(ctx, keyFunc, fallbackFunc, reflect.TypeOf(i), expire)
}

func (c *Cache) MustRefresh(ctx context.Context, keyFunc KeyFunc, fallbackFunc FallbackFunc, i interface{}, ttl *time.Duration) {
	gehirn.PanicIf(c.Refresh(ctx, keyFunc, fallbackFunc, i, ttl))
}

func (c *Cache) refreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc, dst interface{}, expire time.Duration) error {
	dstV := reflect.ValueOf(dst)
	dstT := dstV.Type()
	// check dst
	if dstV.Kind() != reflect.Map {
		return ErrBadDstType
	}
	if dstV.IsNil() {
		return ErrBadDstMapValue
	}
	dstKeyT := dstT.Key()
	dstValT := dstT.Elem()

	// check ids is slice and ids element type equal to map key type
	idsV := reflect.ValueOf(ids)
	idsT := idsV.Type()
	if idsV.Kind() != reflect.Slice {
		return ErrBadIdsType
	}
	length := idsV.Len()
	if length == 0 {
		return nil
	}
	idsValT := idsT.Elem()
	if idsValT != dstKeyT {
		return ErrSrcDstTypeMismatch
	}

	// gen key
	keys := make([]string, length)
	for i := 0; i < length; i++ {
		id := idsV.Index(i).Interface()
		keys[i] = keyFunc(id)
	}

	// delete cache
	err := c.store.Delete(ctx, keys...)
	if err != nil {
		return err
	}

	// fallback
	r, err := fallbackFunc(ids)
	if err != nil {
		return err
	}

	// check fallback with dst
	rV := reflect.ValueOf(r)
	rT := rV.Type()
	if rV.Kind() != reflect.Map {
		return ErrBadFallbackResultType
	}
	length = rV.Len()
	if length == 0 {
		return nil
	}
	rKeyT := rT.Key()
	rValT := rT.Elem()
	if rKeyT != dstKeyT {
		return ErrDstFallbackTypeMismatch
	}
	if rValT != dstValT {
		return ErrDstFallbackTypeMismatch
	}

	keys = make([]string, length)
	values := reflect.MakeSlice(reflect.SliceOf(rValT), 0, length)
	for i, idV := range rV.MapKeys() {
		id := idV.Interface()
		key := keyFunc(id)
		keys[i] = key
		values = reflect.Append(values, rV.MapIndex(idV))
	}
	return c.store.SetMulti(ctx, keys, values.Interface(), expire)
}

func (c *Cache) RefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc, dst interface{}, ttl *time.Duration) error {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	return c.refreshMulti(ctx, ids, keyFunc, fallbackFunc, dst, expire)
}

func (c *Cache) MustRefreshMulti(ctx context.Context, ids interface{}, keyFunc KeyMultiFunc,
	fallbackFunc FallbackMultiFunc, dst interface{}, ttl *time.Duration) {
	expire := c.expire
	if ttl != nil {
		expire = *ttl
	}
	gehirn.PanicIf(c.refreshMulti(ctx, ids, keyFunc, fallbackFunc, dst, expire))
}
