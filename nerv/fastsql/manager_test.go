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
 @Time    : 2024/10/18 -- 15:20
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: manager_test.go
*/

package fastsql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 驱动
	"testing"
)

func TestConnect(t *testing.T) {
	dsn := "root:w12w23B(@tcp(127.0.0.1:3306)/iap?parseTime=true&charset=utf8mb4&loc=Asia%2FShanghai"
	client := MustConnect(dsn)
	fmt.Println(client)
}
