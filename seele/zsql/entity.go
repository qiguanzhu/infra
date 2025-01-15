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

package zsql

import (
	"database/sql"
	"github.com/qiguanzhu/infra/pkg/consts"
	"time"
)

// MysqlConf ...
type MysqlConf struct {
	DBName  string   `properties:"dbName" json:"dbName" yaml:"dbName"`
	DBAddrs []string `properties:"dbAddrs" json:"dbAddrs" yaml:"dbAddrs"`

	MaxIdleConns     int    `properties:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns     int    `properties:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
	MaxLifeTimeSec   int    `properties:"maxLifeTimeSec" json:"maxLifeTimeSec" yaml:"maxLifeTimeSec"`
	TimeoutMsec      int    `properties:"timeoutMsec" json:"timeoutMsec" yaml:"timeoutMsec"`
	ReadTimeoutMsec  int    `properties:"readTimeoutMsec" json:"readTimeoutMsec" yaml:"readTimeoutMsec"`
	WriteTimeoutMsec int    `properties:"writeTimeoutMsec" json:"writeTimeoutMsec" yaml:"writeTimeoutMsec"`
	Username         string `properties:"username" json:"username" yaml:"username"`
	Password         string `properties:"password" json:"password" yaml:"password"`
}

func (c *MysqlConf) LoadDefault(insName string) {
	if c.TimeoutMsec == 0 {
		c.TimeoutMsec = consts.DefaultTimeoutSecond
	}
	if c.ReadTimeoutMsec == 0 {
		c.ReadTimeoutMsec = consts.DefaultReadTimeoutSecond
	}
	if c.WriteTimeoutMsec == 0 {
		c.WriteTimeoutMsec = consts.DefaultWriteTimeoutSecond
	}
	if c.MaxLifeTimeSec == 0 {
		c.MaxLifeTimeSec = consts.DefaultMaxLifeTimeSecond
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = consts.DefaultMaxIdleConns
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = consts.DefaultMaxOpenConns
	}

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
