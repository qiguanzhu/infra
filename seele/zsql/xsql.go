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
 @Time    : 2024/9/30 -- 15:13
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: xsql.go
*/

package zsql

import (
	"context"
	"database/sql"
)

type XDB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type XDBWrapper interface {
	XDB
	GetDb() *sql.DB
}

type DynamicConfigureProxy[Config any, DynamicConf any] interface {
	SetConf(Config)
	LoadDynamicConf(string, DynamicConf)
	IsInit() bool
}

type SqlConfigProxy interface {
	IsProxyHostSet() bool
	IsProxyPortSet() bool
	GetProxyHost() string
	GetProxyPort() int
}

// ConfigureProxy 配置接口抽象
type ConfigureProxy[Config any] interface {
	GetInstanceName(ctx context.Context, cluster, table string) string
	GetInstanceConfig(ctx context.Context, instance, group string) Config
	GetAllGroups(ctx context.Context) []string
}

type SqlConstructor interface {
	GetBuilder() Builder
	GetScanner() Scanner
}
