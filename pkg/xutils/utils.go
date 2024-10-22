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
 @Time    : 2024/9/30 -- 15:41
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: utils.go
*/

package xutils

import (
	"bytes"
	"github.com/qiguanzhu/infra/pkg"
	"reflect"
	"strings"
)

// Map converts a struct to a map
// type for each field of the struct must be built-in type
func Map(target interface{}, useTag string) (map[string]interface{}, error) {
	if nil == target {
		return nil, nil
	}
	v := reflect.ValueOf(target)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, pkg.ErrNoneStructTarget
	}
	t := v.Type()
	result := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		keyName := getKey(t.Field(i), useTag)
		if "" == keyName {
			continue
		}
		result[keyName] = v.Field(i).Interface()
	}
	return result, nil
}

func isExportedField(name string) bool {
	return strings.ToUpper(name) == name
}

func getKey(field reflect.StructField, useTag string) string {
	if !isExportedField(field.Name) {
		return ""
	}
	if field.Type.Kind() == reflect.Ptr {
		return ""
	}
	if "" == useTag {
		return field.Name
	}
	tag, ok := field.Tag.Lookup(useTag)
	if !ok {
		return ""
	}
	return ResolveTagName(tag)
}

func ResolveTagName(tag string) string {
	idx := strings.IndexByte(tag, ',')
	if -1 == idx {
		return tag
	}
	return tag[:idx]
}

func Concat(strings ...string) string {
	var buffer bytes.Buffer
	for _, s := range strings {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func Mirror(src, tar any) {
	srcValue := reflect.ValueOf(src)
	tarValue := reflect.ValueOf(tar)
	if srcValue.Type().Kind() != reflect.Ptr {
		tarValue.Elem().Set(srcValue)
	} else {
		tarValue.Elem().Set(srcValue.Elem())
	}
}
