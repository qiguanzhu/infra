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
 @Time    : 2024/9/30 -- 15:28
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: scanner.go
*/

package zsql

// ScannerProxy 读取查询结果 or 执行结果 到 目标结构体
type ScannerProxy interface {
	Scan(rows Rows, target any, f BindFunc) error
	ScanClose(rows Rows, target any, f BindFunc) error
	ScanMap(rows Rows) ([]map[string]any, error)
	ScanMapDecode(rows Rows) ([]map[string]any, error)
	ScanMapClose(rows Rows) ([]map[string]any, error)
}

type Rows interface {
	Close() error
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
}

// ByteUnmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalByte must copy the JSON data
// if it wishes to retain the data after returning.
//
// By convention, to approximate the behavior of Unmarshal itself,
// ByteUnmarshaler implement UnmarshalByte([]byte("null")) as a no-op.
type ByteUnmarshaler interface {
	UnmarshalByte(data []byte) error
}

type BindFunc func(rows Rows) (to any, err error)
