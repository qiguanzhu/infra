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
 @Time    : 2024/10/12 -- 14:17
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: manager.go
*/

package fastsql

import (
	"database/sql"
	"fmt"
	"github.com/qiguanzhu/infra/seele/xsqlIface"
	"time"
)

const (
	// mysql pool
	mysqlConnMaxIdleTime = 10 * time.Minute
	mysqlConnMaxLifetime = 30 * time.Minute
	mysqlMaxIdleConns    = 128
	mysqlMaxOpenConns    = 1024 * 16
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(mysqlConnMaxIdleTime)
	db.SetConnMaxLifetime(mysqlConnMaxLifetime)
	db.SetMaxIdleConns(mysqlMaxIdleConns)
	db.SetMaxOpenConns(mysqlMaxOpenConns)

	return db, nil
}

func MustConnect(dsn string) xsqlIface.XDB {
	db, err := Connect(dsn)
	if err != nil {
		panic(fmt.Sprintf("connect to mysql %s error %+v", dsn, err))
	}
	return db
}
