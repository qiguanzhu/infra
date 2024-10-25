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
 @Time    : 2024/10/11 -- 18:19
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: scan.go
*/

package sqlutils

import (
	"github.com/qiguanzhu/infra/pkg"
	"github.com/qiguanzhu/infra/seele/zsql"
)

func ResolveDataFromRows(rows zsql.Rows) ([]map[string]interface{}, error) {
	if nil == rows {
		return nil, pkg.ErrScannerNilRows
	}
	columns, err := rows.Columns()
	if nil != err {
		return nil, err
	}
	length := len(columns)
	var result []map[string]interface{}
	// unnecessary to put below into rows.Next loop,reduce allocating
	values := make([]interface{}, length)
	for i := 0; i < length; i++ {
		values[i] = new(interface{})
	}
	for rows.Next() {
		err = rows.Scan(values...)
		if nil != err {
			return nil, err
		}
		mp := make(map[string]interface{})
		for idx, name := range columns {
			// mp[name] = reflect.ValueOf(values[idx]).Elem().Interface()
			mp[name] = *(values[idx].(*interface{}))
		}
		result = append(result, mp)
	}
	return result, nil
}
