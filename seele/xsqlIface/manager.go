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
 @Time    : 2024/9/30 -- 15:27
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: manager.go
*/

package xsqlIface

import (
	"context"
	"database/sql"
	"github.com/qiguanzhu/infra/seele/xconfigIface"
)

// ManagerProxy DBManagerProxy 定义 manager 的操作。涉及不同的库 or 配置可以维护不同的 manager。
type ManagerProxy interface {
	InitConf(ctx context.Context, config xconfigIface.ConfigCenter) error
	GetDB(ctx context.Context, insName, dbName string) (XDBWrapper, error)
	ReloadConf(ctx context.Context, config xconfigIface.ConfigCenter, event xconfigIface.ChangeEvent) error
	GetInstance(insName, dbName string) (DBInstanceProxy, error)
}

type InstanceProxy interface {
	String() string
	GetInstanceName() string
	GetDbName() string
}

type DBInstanceProxy interface {
	GetType() string
	Close() error
	Reload() error
	GetDbName() string
	GetDB() *sql.DB
}
