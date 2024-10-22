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
 @Time    : 2024/10/11 -- 15:59
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: scanner.go
*/

package fastsql

import (
	"github.com/qiguanzhu/infra/pkg"
	"github.com/qiguanzhu/infra/seele/xsqlIface"
	"github.com/qiguanzhu/infra/seele/xsqlIface/sqlutils"
	"reflect"
	"strconv"
)

type fastScanner struct{}

func (s fastScanner) Scan(rows xsqlIface.Rows, target any, b xsqlIface.Bind) error {
	if nil == target || reflect.ValueOf(target).IsNil() || reflect.TypeOf(target).Kind() != reflect.Ptr {
		return pkg.ErrScannerTargetNotSettable
	}
	if nil == b {
		return pkg.ErrScannerFromNilFunction
	}

	if reflect.TypeOf(target).Elem().Kind() == reflect.Slice {
		// reflect.TypeOf(target).Elem() -> []DObj
		// reflect.TypeOf(target).Elem().Elem() -> DObj
		// res -> make([]DObj, 0)
		res := reflect.MakeSlice(reflect.TypeOf(target).Elem(), 0, 0)
		for rows.Next() {
			r, err := b(rows)
			if err != nil {
				return err
			}
			// append([]DObj, r.(DObj))
			res = reflect.Append(res, reflect.ValueOf(r).Elem())
		}
		// res.Interface() []interface{} -> []DObj{}
		resI := res.Interface()
		// reflect.ValueOf -> reflect.ValueOf([]DObg{})
		value := reflect.ValueOf(resI)
		// convert -> []DObj{}, call Interface() to get real value
		convert := value.Convert(reflect.TypeOf(target).Elem())
		return Mapping(convert.Interface(), target)
	} else {
		for rows.Next() {
			res, err := b(rows)
			if err != nil {
				return err
			}
			if nil == res {
				return pkg.ErrScannerEmptyResult
			}
			return Mapping(res, target)
		}
		return pkg.ErrScannerEmptyResult
	}
}

func (s fastScanner) ScanMap(rows xsqlIface.Rows) ([]map[string]interface{}, error) {
	return sqlutils.ResolveDataFromRows(rows)
}

func (s fastScanner) ScanMapDecode(rows xsqlIface.Rows) ([]map[string]interface{}, error) {
	results, err := sqlutils.ResolveDataFromRows(rows)
	if nil != err {
		return nil, err
	}
	for i := 0; i < len(results); i++ {
		for k, v := range results[i] {
			rv, ok := v.([]uint8)
			if !ok {
				continue
			}
			s := string(rv)
			// convert to int
			intVal, err := strconv.Atoi(s)
			if err == nil {
				results[i][k] = intVal
				continue
			}
			// convert to float64
			floatVal, err := strconv.ParseFloat(s, 64)
			if err == nil {
				results[i][k] = floatVal
				continue
			}
			// convert to string
			results[i][k] = s
		}
	}
	return results, nil
}

func (s fastScanner) ScanMapDecodeClose(rows xsqlIface.Rows) ([]map[string]interface{}, error) {
	result, err := s.ScanMapDecode(rows)
	if nil != rows {
		errClose := rows.Close()
		if err == nil {
			err = pkg.NewCloseErr(errClose)
		}
	}
	return result, err
}

// ScanMapClose is the same as ScanMap and close the rows
func (s fastScanner) ScanMapClose(rows xsqlIface.Rows) ([]map[string]interface{}, error) {
	result, err := s.ScanMap(rows)
	if nil != rows {
		errClose := rows.Close()
		if err == nil {
			err = pkg.NewCloseErr(errClose)
		}
	}
	return result, err
}

// ScanClose is the same as Scan and helps you Close the rows
// Don't exec the rows.Close after calling this
func (s fastScanner) ScanClose(rows xsqlIface.Rows, target any, b xsqlIface.Bind) error {
	err := s.Scan(rows, target, b)
	if nil != rows {
		errClose := rows.Close()
		if err == nil {
			err = pkg.NewCloseErr(errClose)
		}
	}
	return err
}
