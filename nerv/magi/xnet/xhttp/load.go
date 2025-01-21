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
 @Time    : 2024/8/27 -- 18:47
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: load.go
*/

package xhttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

var _ResContentTypeLoader = map[string]loadFunc{
	ResTypeJSON: loadJson(),
	ResTypeXML:  loadXml(),
}

type loadFunc func(bs []byte, v any) error

func loadJson() loadFunc {
	return func(bs []byte, v any) error {
		err := json.Unmarshal(bs, &v)
		return fmt.Errorf("json.Unmarshal(%s, %+v)：%w", string(bs), v, err)
	}
}

func loadXml() loadFunc {
	return func(bs []byte, v any) error {
		err := xml.Unmarshal(bs, &v)
		return fmt.Errorf("xml.Unmarshal(%s, %+v)：%w", string(bs), v, err)
	}
}

func defaultLoader() CfgOp {
	return Res(ResTypeJSON)
}
