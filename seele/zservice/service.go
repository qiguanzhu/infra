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
 @Time    : 2024/10/24 -- 17:30
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: service.go
*/

package zservice

import (
	"github.com/qiguanzhu/infra/seele/zconfig"
)

type SvcBase[SvcInfo any] interface {
	// RegisterService key is processor to ServInfo
	RegisterService(svcs map[string]*SvcInfo) error
	RegisterBackDoor(svcs map[string]*SvcInfo) error
	RegisterCrossDCService(svcs map[string]*SvcInfo) error

	SvcName() string
	SvcIp() string
	ServId() int
	// CopyName 服务副本名称, servename + servid
	CopyName() string

	// ServConfig 获取服务的配置
	ServConfig(cfg interface{}) error
	// 任意路径的配置信息
	// ArbiConfig(location string) (string, error)

	// GenSlowId
	// 慢id生成器，适合id产生不是非常快的场景,基于毫秒时间戳，每毫秒最多产生2个id，过快会自动阻塞，直到毫秒递增
	// id表示可以再52bit完成，用double表示不会丢失精度，javascript等弱类型语音可以直接使用
	GenSlowId(tp string) (int64, error)
	GetSlowIdStamp(sid int64) int64
	GetSlowIdWithStamp(stamp int64) int64

	// GenSnowFlakeId 雪花id生成逻辑
	GenSnowFlakeId() (int64, error)
	// GetSnowFlakeIdStamp 获取snowflakeid生成时间戳，单位ms
	GetSnowFlakeIdStamp(sid int64) int64
	// GetSnowFlakeIdWithStamp 按给定的时间点构造一个起始snowflakeid，一般用于区域判断
	GetSnowFlakeIdWithStamp(stamp int64) int64

	GenUuid() (string, error)
	GenUuidSha1() (string, error)
	GenUuidMd5() (string, error)

	// 默认的锁，局部分布式锁，各个服务之间独立不共享

	// Lock 获取到lock立即返回，否则block直到获取到
	Lock(name string) error
	// Unlock 没有lock的情况下unlock，程序会直接panic
	Unlock(name string) error
	// TryLock 立即返回，如果获取到lock返回true，否则返回false
	TryLock(name string) (bool, error)

	// 全局分布式锁，全局只有一个，需要特殊加global说明

	LockGlobal(name string) error
	UnlockGlobal(name string) error
	TryLockGlobal(name string) (bool, error)

	// ConfigCenter ...
	ConfigCenter() zconfig.ConfigCenter

	// RegInfos ...
	RegInfos() map[string]string

	// Stop ...
	Stop()

	// SetOnShutdown set app shutdown hook
	SetOnShutdown(func())
}
