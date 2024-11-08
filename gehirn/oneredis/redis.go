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
 @Time    : 2024/11/5 -- 17:50
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: redis.go
*/

package oneredis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// redis://role:pwd@addr:6379/database
// redis-cli -h addr -p 6379 -a role:pwd
// select database

func Connect(dsn string) (*redis.Client, error) {
	options, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(options), nil
}

func MustConnect(dsn string) *redis.Client {
	cli, err := Connect(dsn)
	if err != nil {
		panic(fmt.Sprintf("connect to redis %s error %+v", dsn, err))
	}
	return cli
}
