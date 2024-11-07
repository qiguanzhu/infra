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
 @Time    : 2024/10/25 -- 18:10
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: ipdb.go
*/

package xips

import (
	"context"
	"github.com/ipipdotnet/ipdb-go"
	"sync"
)

/**
ipdb格式的ip
纯真提供的ipdb格式
https://github.com/metowolf/qqwry.ipdb
*/

type Ipdb struct {
	lock                  sync.Mutex
	db                    *ipdb.City
	path                  string
	refreshIntervalSecond int64
}

func NewIpdb(path string, refreshIntervalSecond int64) (*Ipdb, error) {
	obj := &Ipdb{
		path:                  path,
		refreshIntervalSecond: refreshIntervalSecond,
	}
	err := obj.load()
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (m *Ipdb) start() error {
	if err := m.load(); err != nil {
		return err
	}
	if m.refreshIntervalSecond > 0 {
		//	启动更新任务
	}
	return nil
}
func (m *Ipdb) load() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.db != nil {
		return m.db.Reload(m.path)
	}
	db, err := ipdb.NewCity(m.path)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}
func (m *Ipdb) resetDb(db *ipdb.City) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 原有的db 释放掉
	// if m.db != nil {
	//	m.db.Close()
	// }
	m.db = db
}
func (m *Ipdb) getDb() *ipdb.City {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.db
}

func (m *Ipdb) GetRegion(ctx context.Context, ipstr string) IpRegion {
	res, err := m.db.FindInfo(ipstr, "CN")
	if err != nil || res == nil {
		return IpRegion{}
	}
	return IpRegion{
		Country:  res.CountryName,
		Province: res.RegionName,
		City:     res.CityName,
		Region:   "",
		Isp:      res.IspDomain,
	}
}
