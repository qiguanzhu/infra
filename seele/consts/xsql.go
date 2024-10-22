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
 @Time    : 2024/9/30 -- 12:28
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: database.go
*/

package consts

import "time"

const (
	// MysqlConfNamespace mysql apollo conf namespace
	MysqlConfNamespace = "mysql"
	// MaxIdleConnsKey ...
	MaxIdleConnsKey = "maxIdleConns"
	// MaxOpenConnsKey ...
	MaxOpenConnsKey = "maxOpenConns"
	// MaxLifeTimeSecKey ...
	MaxLifeTimeSecKey = "maxLifeTimeSec"
	// TimeoutMsecKey ...
	TimeoutMsecKey = "timeoutMsec"
	// ReadTimeoutMsecKey ...
	ReadTimeoutMsecKey = "readTimeoutMsec"
	// WriteTimeoutMsecKey ...
	WriteTimeoutMsecKey = "writeTimeoutMsec"
	UserNameKey         = "username"
	PasswordKey         = "password"
	// KeySep
	KeySep = "."

	DefaultMaxIdleConns = 64
	DefaultMaxOpenConns = 128
	DefaultReadTimeout  = time.Second * 10
	DefaultWriteTimeout = time.Second * 10
	DefaultMaxLifeTime  = time.Hour * 6
	DefaultTimeout      = time.Second * 3
)

const (
	// DefaultTagName is the default struct tag name
	DefaultTagName = "bdb"
	CTimeFormat    = "2006-01-02 15:04:05"
)

const (
	DefaultDriver = "mysql"
	DefaultPort   = 3306
	CDSNFormat    = "%s%s=%s&"
)

const (
	WeirProxyHost = "weirproxy.service.svc.cluster.local"
	WeirProxyPort = 9021
	DefaultDbType = "weir_proxy"
)

const (
	WeirProxyHostEnv = "WEIR_PROXY_HOST"
	WeirProxyPortEnv = "WEIR_PROXY_PORT"
)
