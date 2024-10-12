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
 @Time    : 2024/10/12 -- 16:37
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: parse.go
*/

package parse

import (
	"github.com/qiguanzhu/infrastructure/apogo/parse/imp"
	"github.com/qiguanzhu/infrastructure/apogo/parse/json"
	"github.com/qiguanzhu/infrastructure/apogo/parse/yaml"
	"github.com/qiguanzhu/infrastructure/apogo/parse/yml"
)

type ContentParser interface {
	Parse(configContent interface{}) (map[string]interface{}, error)
	GetParserType() string
	Unmarshal(data []byte, val interface{}) error
}

func GetParser(typ string) ContentParser {
	switch typ {
	case "yml":
		return yml.NewParser()
	case "yaml":
		return yaml.NewParser()
	case "json":
		return json.NewParser()
	default:
		return imp.NewParser()
	}
}
