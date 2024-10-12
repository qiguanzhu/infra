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
 @Time    : 2024/10/12 -- 16:35
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: j.go
*/

package json

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
)

// Parser json 转换器
type Parser struct {
	Vp *viper.Viper
}

func NewParser() *Parser {
	p := Parser{
		Vp: viper.New(),
	}
	p.Vp.SetConfigType("json")
	return &p
}

// Parse 内存内容 => json 数据格式转换器
func (p *Parser) Parse(configContent interface{}) (map[string]interface{}, error) {
	content, ok := configContent.(string)
	if !ok {
		return nil, nil
	}
	if content == "" {
		return nil, nil
	}

	buffer := bytes.NewBufferString(content)

	err := p.Vp.ReadConfig(buffer)
	if err != nil {
		return nil, err
	}

	return p.convertToMap(), nil
}

func (p *Parser) convertToMap() map[string]interface{} {
	if p.Vp == nil {
		return nil
	}

	m := make(map[string]interface{})
	for _, key := range p.Vp.AllKeys() {
		m[key] = p.Vp.Get(key)
	}
	return m
}

func (p *Parser) GetParserType() string {
	return "json"
}

func (p *Parser) Unmarshal(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}
