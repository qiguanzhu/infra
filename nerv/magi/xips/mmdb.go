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
 @Time    : 2024/10/25 -- 18:12
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: mmdb.go
*/

package xips

import (
	"context"
	"github.com/oschwald/maxminddb-golang"
	"net"
	"sync"
)

/**
mmdb格式的ip
github 上有每日更新并免费分享的GeoIP 库。https://github.com/P3TERX/GeoLite.mmdb/tree/download
只能解析到国家,
是海外的ip库，理论上解析海外的国家更合适
*/

type Mmdb struct {
	lock                  sync.Mutex
	db                    *maxminddb.Reader
	path                  string
	refreshIntervalSecond int64
}

func NewMmdb(path string, refreshIntervalSecond int64) (*Mmdb, error) {
	obj := &Mmdb{
		path:                  path,
		refreshIntervalSecond: refreshIntervalSecond,
	}
	err := obj.load()
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (m *Mmdb) start() error {
	if err := m.load(); err != nil {
		return err
	}
	if m.refreshIntervalSecond > 0 {
		//	启动更新任务
	}
	return nil
}
func (m *Mmdb) load() error {
	db, err := maxminddb.Open(m.path)
	if err != nil {
		return err
	}
	m.resetDb(db)
	return nil
}
func (m *Mmdb) resetDb(db *maxminddb.Reader) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 原有的db 释放掉
	if m.db != nil {
		m.db.Close()
	}
	m.db = db
}
func (m *Mmdb) getDb() *maxminddb.Reader {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.db
}

// 参考https://dev.maxmind.com/minfraud/api-documentation/responses?lang=en
type mmdbRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
		Names   struct {
			Zh string `maxminddb:"zh-CN"`
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"country"`
	City struct {
		ISOCode string `maxminddb:"iso_code"`
		Names   struct {
			Zh string `maxminddb:"zh-CN"`
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`
	Represented_country struct {
		ISOCode string `maxminddb:"iso_code"`
		Names   struct {
			Zh string `maxminddb:"zh-CN"`
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"represented_country"`
	Registered_country struct {
		ISOCode string `maxminddb:"iso_code"`
		Names   struct {
			Zh string `maxminddb:"zh-CN"`
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"registered_country"`
	Subdivisions struct {
		ISOCode string `maxminddb:"iso_code"`
		Names   struct {
			Zh string `maxminddb:"zh-CN"`
			En string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
} // Or any appropriate struct

func (m *Mmdb) GetRegion(ctx context.Context, ipstr string) IpRegion {
	var r = new(mmdbRecord)
	err := m.UnmarshalRecord(ctx, ipstr, r)
	if err != nil {
		return IpRegion{}
	}
	var country string
	if r != nil && r.Country.Names.Zh != "" {
		country = r.Country.Names.Zh
	}
	return IpRegion{
		Country: country,
	}
}
func (m *Mmdb) UnmarshalRecord(ctx context.Context, ipstr string, record interface{}) error {
	ip := net.ParseIP(ipstr)
	db := m.getDb()
	err := db.Lookup(ip, &record)
	return err
}
