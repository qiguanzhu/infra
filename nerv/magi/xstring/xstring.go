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
 @Time    : 2024/11/4 -- 17:57
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: xstring.go
*/

package xstring

import (
	"bytes"
	"encoding/json"
	"github.com/kaneshin/go-pkg/unicode"
	"unicode/utf8"
	"unsafe"
)

func GetUtf8Chars(s string, num int) string {
	rv := ""
	for i := 0; len(s) > 0 && i < num; i++ {
		_, size := utf8.DecodeRuneInString(s)
		rv += s[:size]
		s = s[size:]
	}

	return rv
}

// GetInvalidUtf8String
// 截取获取合法 num 个unicode字符 的utf8字符串
// num 为0，全部截取
// 返回截取的unicode字符个数，以及字符串
func GetInvalidUtf8String(s string, num int) (string, int) {
	rv := ""
	count := 0
	for i := 0; len(s) > 0; i++ {
		ru, size := utf8.DecodeRuneInString(s)
		if ru != utf8.RuneError {
			rv += s[:size]
			count++

			if num > 0 && count >= num {
				break
			}
		}
		s = s[size:]
	}

	return rv, count
}

func HaveEmoji(s string) bool {
	rs := []rune(s)
	for _, r := range rs {
		if unicode.IsEmoji(r) {
			return true
		}
	}

	return false
}

func ReplaceEmoji(s string, ns string) string {
	var rv []rune
	rs := []rune(s)
	for _, r := range rs {
		if unicode.IsEmoji(r) {
			rv = append(rv, []rune(ns)...)
		} else {
			rv = append(rv, r)
		}
	}

	return string(rv)
}

var (
	sbc2dbcMap = map[rune]rune{
		'＋': '+',
		'－': '-',
		'０': '0',
		'１': '1',
		'２': '2',
		'３': '3',
		'４': '4',
		'５': '5',
		'６': '6',
		'７': '7',
		'８': '8',
		'９': '9',
		'‘': '\'',
		'’': '\'',
		'“': '"',
		'”': '"',
		'，': ',',
		'。': '.',
		'？': '?',
		'×': '*',
		'／': '/',
		'％': '%',
		'＃': '#',
		'＠': '@'}
)

// UCFirst 首写字母大写
func UCFirst(str string) string {
	runes := []rune(str)
	if len(runes) < 1 {
		return str
	}
	if runes[0] >= 97 && runes[0] <= 122 {
		runes[0] -= 32
	}
	return string(runes)
}

func SBC2DBC(str string) string {
	runes := []rune(str)
	var buf bytes.Buffer
	for i := 0; i < len(runes); i++ {
		r, ok := sbc2dbcMap[runes[i]]
		if ok {
			buf.WriteRune(r)
		} else {
			buf.WriteRune(runes[i])
		}
	}
	return buf.String()
}

func Concat(strings ...string) string {
	var buffer bytes.Buffer
	for _, s := range strings {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func StrLen(str string) int {
	runes := []rune(str)
	return len(runes)
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]} // copy data, len, cap
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func FilterEmoji(s string) string {
	dst := ""
	for _, value := range s {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			dst += string(value)
		}
	}
	return dst
}

func Any2String(v any) (str string) {
	if v == nil {
		return ""
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return ""
	}
	str = string(bs)
	return
}
