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
 @Time    : 2024/10/21 -- 16:38
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: config.go
*/

package config

import "github.com/qiguanzhu/infra/seele/xsqlIface"

// var Center xconfig.ConfigCenter
// var RedisClient *redisext.RedisExt
// var PoolClient *xutil.Pool
// var PulsarManager *pulsar.Manager
// var RandClient *rand.Rand
// var MApolloConfig = &ApolloConfig{}

var DefaultDB xsqlIface.XDB

type Config struct {
	Debug struct {
		// false: Production   true: Development
		Enable bool
	}
}
