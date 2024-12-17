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
 @Time    : 2024/12/10 -- 17:42
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: json_test.go
*/

package xjson

import (
	"fmt"
	"io"
	"log"
	"testing"
)

func TestWalkEndReader(t *testing.T) {
	jsonData := []byte(`{"key1": {"key2": [{"key3": "值1"}, {"key3": "值2"}, {"key3": "最终的值"}]}}`)
	path := "key1.key2.2.key3" // 访问 key2 中的第三个元素
	sep := "."

	reader, err := WalkEndReader(jsonData, path, sep)
	if err != nil {
		log.Fatalf("错误: %v", err)
	}

	// 读取并打印最终的值
	value, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("读取错误: %v", err)
	}

	fmt.Println(string(value)) // 输出: "最终的值"
}
