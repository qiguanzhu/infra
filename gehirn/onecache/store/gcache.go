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
 @Time    : 2024/11/5 -- 16:40
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: gcache.go
*/

package store

import (
	"context"
	"reflect"
	"time"

	"github.com/bluele/gcache"
	"github.com/qiguanzhu/infra/gehirn"
	"github.com/qiguanzhu/infra/nerv/magi/xreflect"
)

type LocalStore struct {
	store gcache.Cache
}

func NewLocalStore(size int) *LocalStore {
	return &LocalStore{
		store: gcache.New(size).LRU().Build(),
	}
}

func (ls *LocalStore) Get(_ context.Context, key string, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return ErrBadDstType
	}

	val, err := ls.store.Get(key)
	if err != nil {
		return err
	}

	dstV := reflect.Indirect(reflect.ValueOf(dst))

	newVal := reflect.Indirect(reflect.ValueOf(val))

	left := xreflect.RecursiveIndirectType(dstV.Type())
	right := xreflect.RecursiveIndirectType(newVal.Type())
	if left.Kind() != right.Kind() {
		return ErrSrcDstTypeMismatch
	}

	dstV.Set(newVal)

	return nil
}

func (ls *LocalStore) MustGet(ctx context.Context, key string, dst interface{}) {
	gehirn.PanicIf(ls.Get(ctx, key, dst))
}

func (ls *LocalStore) GetMulti(ctx context.Context, keys []string, dstMap interface{}) error {
	dstPtrV := reflect.ValueOf(dstMap)
	dstV := reflect.Indirect(dstPtrV)
	if dstV.Kind() != reflect.Map {
		return ErrBadDstMapType
	}

	// nil map
	if dstPtrV.Kind() != reflect.Ptr && dstV.IsNil() {
		return ErrBadDstMapValue
	}

	if dstPtrV.Kind() == reflect.Ptr && dstV.IsNil() {
		m := reflect.MakeMap(reflect.MapOf(dstV.Type().Key(), dstV.Type().Elem()))
		dstV.Set(m)
	}

	for i, key := range keys {
		v := reflect.New(dstV.Type().Elem())
		if v.Kind() != reflect.Ptr {
			v = v.Addr()
		}

		val, err := ls.store.Get(key)
		if err != nil {
			continue
		}

		dstV.SetMapIndex(reflect.ValueOf(keys[i]), reflect.ValueOf(val))
	}
	return nil
}

func (ls *LocalStore) MustGetMulti(ctx context.Context, keys []string, dstMap interface{}) {
	gehirn.PanicIf(ls.GetMulti(ctx, keys, dstMap))
}

func (ls *LocalStore) Exists(_ context.Context, key string) (bool, error) {
	return ls.store.Has(key), nil
}

func (ls *LocalStore) MustExists(ctx context.Context, key string) bool {
	ok, err := ls.Exists(ctx, key)
	if err != nil {
		panic(err)
	}
	return ok
}

func (ls *LocalStore) ExistsMulti(ctx context.Context, keys ...string) ([]bool, error) {
	if len(keys) == 0 {
		return []bool{}, nil
	}

	var results []bool
	for _, key := range keys {
		ok, err := ls.Exists(ctx, key)
		if err != nil {
			return nil, err
		}
		results = append(results, ok)
	}

	return results, nil
}

func (ls *LocalStore) MustExistsMulti(ctx context.Context, keys ...string) []bool {
	ret, err := ls.ExistsMulti(ctx, keys...)
	gehirn.PanicIf(err)
	return ret
}

func (ls *LocalStore) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	err := ls.store.SetWithExpire(key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LocalStore) MustSet(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	gehirn.PanicIf(ls.Set(ctx, key, value, ttl))
}

func (ls *LocalStore) SetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) error {
	srcV := reflect.Indirect(reflect.ValueOf(values))

	if srcV.Kind() != reflect.Slice {
		return ErrBadSrcType
	}

	if srcV.Len() != len(keys) {
		return ErrKeysLengthNotMatch
	}

	for index, key := range keys {
		v := srcV.Index(index)

		err := ls.Set(ctx, key, v.Interface(), ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ls *LocalStore) MustSetMulti(ctx context.Context, keys []string, values interface{}, ttl time.Duration) {
	gehirn.PanicIf(ls.SetMulti(ctx, keys, values, ttl))
}

func (ls *LocalStore) Delete(_ context.Context, keys ...string) error {
	for _, key := range keys {
		ls.store.Remove(key)
	}

	return nil
}

func (ls *LocalStore) MustDelete(ctx context.Context, keys ...string) {
	err := ls.Delete(ctx, keys...)
	if err != nil {
		panic(err)
	}
}
