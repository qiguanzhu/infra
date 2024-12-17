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
 @Time    : 2024/12/10 -- 17:40
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: json.go
*/

package xjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func WalkEndReader(data []byte, path string, sep string) (io.Reader, error) {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	keys := splitPath(path, sep)
	value, err := walk(jsonData, keys)
	if err != nil {
		return nil, err
	}

	// 将最终的值转换为 JSON 格式并返回 io.Reader
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(valueBytes), nil
}

// 分割路径
func splitPath(path string, sep string) []string {
	// suit for empty sep and only one level of path
	if sep != "" {
		return []string{path}
	}
	return strings.Split(path, sep)
}

// 遍历 JSON 数据
func walk(data interface{}, keys []string) (interface{}, error) {
	for _, key := range keys {
		switch d := data.(type) {
		case map[string]interface{}:
			// 处理 map
			if value, ok := d[key]; ok {
				data = value
			} else {
				return nil, fmt.Errorf("路径 %s 不存在", key)
			}
		case []interface{}:
			// 处理 list
			index, err := strconv.Atoi(key)
			if err != nil || index < 0 || index >= len(d) {
				return nil, fmt.Errorf("路径 %s 不存在或索引超出范围", key)
			}
			data = d[index]
		default:
			return nil, fmt.Errorf("路径 %s 不存在", key)
		}
	}
	return data, nil
}
