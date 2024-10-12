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
 @Time    : 2024/10/12 -- 16:39
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: imp.go
*/

package imp

type Parser struct {
}

func NewParser() *Parser {
	p := Parser{}
	return &p
}

func (d *Parser) Parse(configContent interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (this *Parser) GetParserType() string {
	return "normal"
}

func (this *Parser) Unmarshal(data []byte, val interface{}) error {
	return nil
}
