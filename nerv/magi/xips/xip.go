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
 @Time    : 2024/10/25 -- 18:11
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: xip.go
*/

package xips

import (
	"context"
	"fmt"
	"strings"
)

type IpRegion struct {
	Country string `json:"country"`
	// 省
	Province string `json:"province"`
	// 市
	City string `json:"city"`
	// 区
	Region string `json:"region"`
	// 网络运营商  联通 移动...
	Isp string `json:"isp"`
}

// ToCountryProvCity
// 返回格式为 "国家/地区-省-市"，有可能解析不出来，解析不出来为空
// 旧版的ip库解析的结果为该格式
func (m IpRegion) ToCountryProvCity(ctx context.Context) string {
	counry := m.Country
	prov := m.FormatProvince(ctx)
	city := m.FormatCity(ctx)
	return fmt.Sprintf("%s-%s-%s", counry, prov, city)
}
func (m IpRegion) FormatProvince(ctx context.Context) string {
	return trimProCityRegion(ctx, m.Region, "省")
}
func (m IpRegion) FormatCity(ctx context.Context) string {
	return trimProCityRegion(ctx, m.Region, "市")
}
func (m IpRegion) FormatRegion(ctx context.Context) string {
	return trimProCityRegion(ctx, m.Region, "区")
}
func trimProCityRegion(ctx context.Context, origin, suffix string) string {
	if strings.HasSuffix(origin, suffix) {
		runes := []rune(origin)
		return string(runes[0 : len(runes)-1])
	}
	return origin
}
