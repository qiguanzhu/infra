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
 @Time    : 2024/9/30 -- 12:21
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: database.go
*/

package xsqlIface

import (
	"database/sql"
	"time"
)

// MysqlConf ...
type MysqlConf struct {
	MaxIdleConns     int    `properties:"maxIdleConns"`
	MaxOpenConns     int    `properties:"maxOpenConns"`
	MaxLifeTimeSec   int    `properties:"maxLifeTimeSec"`
	TimeoutMsec      int    `properties:"timeoutMsec"`
	ReadTimeoutMsec  int    `properties:"readTimeoutMsec"`
	WriteTimeoutMsec int    `properties:"writeTimeoutMsec"`
	Username         string `properties:"username"`
	Password         string `properties:"password"`
}

type Cfg struct {
	ConfMap   map[string]MysqlConf `properties:"conf_map"`
	ProxyHost string               `properties:"proxy_host"`
	ProxyPort int                  `properties:"proxy_port"`
}

func (c *Cfg) IsProxyHostSet() bool {
	return c.ProxyHost != ""
}

func (c *Cfg) IsProxyPortSet() bool {
	return c.ProxyPort != 0
}

func (c *Cfg) GetProxyHost() string {
	return c.ProxyHost
}

func (c *Cfg) GetProxyPort() int {
	return c.ProxyPort
}

// DynamicConf ...
type DynamicConf struct {
	Timeout        time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxLifeTimeSec time.Duration
	MaxIdleConns   int
	MaxOpenConns   int
	Username       string
	Password       string
}

// Option stands for a series of options for creating a DB
type Option struct {
	driver   string
	DbName   string
	User     string
	Password string
	Host     string
	port     int
	settings []Setting
}

type Setting func(string) string

// Port sets the server port,default 3306
func (o *Option) Port(port int) *Option {
	o.port = port
	return o
}

// Driver sets the driver, default mysql
func (o *Option) Driver(driver string) *Option {
	o.driver = driver
	return o
}

// Set receives a series of Set*-like functions
func (o *Option) Set(sets ...Setting) *Option {
	o.settings = append(o.settings, sets...)
	return o
}

type OpenFunc func(option *Option) (*sql.DB, error)

// Open is used for creating a *sql.DB
// Use it at the last
func (o *Option) Open(ping bool, open OpenFunc) (*sql.DB, error) {
	db, err := open(o)
	if nil != err {
		return nil, err
	}
	if ping {
		err = db.Ping()
	}
	return db, err
}

func (o *Option) GetUser() string {
	return o.User
}
func (o *Option) GetPassword() string {
	return o.Password
}
func (o *Option) GetHost() string {
	return o.Host
}
func (o *Option) GetPort() int {
	return o.port
}
func (o *Option) GetDbName() string {
	return o.DbName
}
func (o *Option) GetDriver() string {
	return o.driver
}
func (o *Option) GetSettings() []Setting {
	return o.settings
}

// Config ...
type Config struct {
	DBName   string
	DBType   string
	DBAddr   []string
	UserName string
	PassWord string
}

// ChangeIns ...
type ChangeIns struct {
	InsNames []string
}

// ConfigChange 配置变更
type ConfigChange struct {
	DbInstanceChange map[string][]string
	DbGroups         []string
}
