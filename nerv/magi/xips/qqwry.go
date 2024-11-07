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
 @Time    : 2024/10/25 -- 18:13
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: qqwry.go
*/

package xips

import (
	"context"
	"github.com/sinlov/qqwry-golang/qqwry"
	"log"
	"strings"
	"sync"
)

/**
纯真是国内IP库免费社区的维护者。可以通过官方接口下载qqwry.dat ip数据文件
https://github.com/oschwald/maxminddb-golang
*/

type QQWry struct {
	lock                  sync.Mutex
	qwry                  *qqwry.QQwry
	path                  string
	refreshIntervalSecond int64
}

func NewQQWry(path string, refreshIntervalSecond int64) (*QQWry, error) {
	obj := &QQWry{
		path:                  path,
		refreshIntervalSecond: refreshIntervalSecond,
	}
	err := obj.load()
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (m *QQWry) start() error {
	if err := m.load(); err != nil {
		return err
	}
	if m.refreshIntervalSecond > 0 {
		//	启动更新任务
	}
	return nil
}
func (m *QQWry) load() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	qqwry.DatData.FilePath = m.path
	init := qqwry.DatData.InitDatFile()
	if v, ok := init.(error); ok {
		if v != nil {
			log.Fatal("init InitDatFile error: ", v)
			return v
		}
	}
	m.qwry = qqwry.NewQQwry()
	return nil
}
func (m *QQWry) getDb() *qqwry.QQwry {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.qwry
}

// todo 未实现解析省市
func (m *QQWry) GetRegion(ctx context.Context, ipstr string) IpRegion {
	res := m.qwry.SearchByIPv4(ipstr)
	var pro = res.Country
	var city string
	arrs := strings.Split(pro, "省")
	if len(arrs) > 1 {
		city = arrs[1]
		pro = arrs[0]
	}
	if len(pro) > 3 {
		city = strings.ReplaceAll(pro, "内蒙古", "")
		pro = "内蒙古"
	}
	return IpRegion{
		Country: pro,
		City:    city,
	}
}
