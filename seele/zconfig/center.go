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
 @Time    : 2024/9/30 -- 12:02
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: interface.go
*/

package zconfig

import (
	"context"
	"github.com/qiguanzhu/infra/seele/zconfig/zobserver"
)

// Option ...
type Option func(ConfigCenter)

// ConfigureType ...
type ConfigureType string

type ConfigCenter interface {
	// RegisterObserver register observer return recall func to cancel observer
	RegisterObserver(ctx context.Context, observer *zobserver.ConfigObserver) (recall func())
	// Stop stop client include cancel client ctx, cancel long poller ctx, close updateChan
	Stop(ctx context.Context) error
	// SubscribeNamespaces subscribe new namespaces if init not set
	SubscribeNamespaces(ctx context.Context, namespaceNames []string) error
	// GetString get string value form default namespace application
	GetString(ctx context.Context, key string) (string, bool)
	// GetStringWithNamespace get string value form specified namespace
	GetStringWithNamespace(ctx context.Context, namespace, key string) (string, bool)
	// GetBool get bool value form default namespace application
	GetBool(ctx context.Context, key string) (bool, bool)
	// GetBoolWithNamespace get bool value form specified namespace
	GetBoolWithNamespace(ctx context.Context, namespace, key string) (bool, bool)
	// GetInt get int value form default namespace application
	GetInt(ctx context.Context, key string) (int, bool)
	// GetIntWithNamespace get int value form specified namespace
	GetIntWithNamespace(ctx context.Context, namespace, key string) (int, bool)
	// GetInt64 get int64 value form default namespace application
	GetInt64(ctx context.Context, key string) (int64, bool)
	// GetInt64WithNamespace get int64 value form specified namespace
	GetInt64WithNamespace(ctx context.Context, namespace, key string) (int64, bool)
	// GetInt32 get int32 value form default namespace application
	GetInt32(ctx context.Context, key string) (int32, bool)
	// GetInt32WithNamespace get int32 value form specified namespace
	GetInt32WithNamespace(ctx context.Context, namespace, key string) (int32, bool)
	// GetFloat64 get float64 value form default namespace application
	GetFloat64(ctx context.Context, key string) (float64, bool)
	// GetFloat64WithNamespace get float64 value form specified namespace
	GetFloat64WithNamespace(ctx context.Context, namespace, key string) (float64, bool)
	// GetIntSlice get []int value form default namespace application
	GetIntSlice(ctx context.Context, keyPrefix string) ([]int, bool)
	// GetIntSliceWithNamespace get []int value form specified namespace
	GetIntSliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]int, bool)
	// GetInt64Slice get []int value form default namespace application
	GetInt64Slice(ctx context.Context, keyPrefix string) ([]int64, bool)
	// GetInt64SliceWithNamespace get []int value form specified namespace
	GetInt64SliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]int64, bool)
	// GetInt32Slice get []int value form default namespace application
	GetInt32Slice(ctx context.Context, keyPrefix string) ([]int32, bool)
	// GetInt32SliceWithNamespace get []int value form specified namespace
	GetInt32SliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]int32, bool)
	// GetStringSlice get []int value form default namespace application
	GetStringSlice(ctx context.Context, keyPrefix string) ([]string, bool)
	// GetStringSliceWithNamespace get []int value form specified namespace
	GetStringSliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]string, bool)
	// GetAllKeys get all keys from default namespace application
	GetAllKeys(ctx context.Context) []string
	// GetAllKeysWithNamespace get all keys form specified namespace
	GetAllKeysWithNamespace(ctx context.Context, namespace string) []string
	// Unmarshal unmarshal from default namespace application
	Unmarshal(ctx context.Context, v interface{}) error
	// UnmarshalWithNamespace unmarshal form specified namespace
	UnmarshalWithNamespace(ctx context.Context, namespace string, v interface{}) error
	// UnmarshalKey unmarshal key from default namespace application
	UnmarshalKey(ctx context.Context, key string, v interface{}) error
	// UnmarshalKeyWithNamespace unmarshal key form specified namespace
	UnmarshalKeyWithNamespace(ctx context.Context, namespace string, key string, v interface{}) error
	// SetCluster set center cluster
	SetCluster(cluster string)
	// SetCacheDir set center cache dir
	SetCacheDir(cacheDir string)
	// SetIPHost set center remote host
	SetIPHost(ipHost string)
}
