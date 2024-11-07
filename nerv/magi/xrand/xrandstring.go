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
 @Time    : 2024/7/13 -- 15:50
 @Author  : bishop ❤️ MONEY
 @Description: xrandstring.go
*/

package xrand

import "math/rand"

const letterDigit = "0123456789"

const letterAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int, ctype int) string {

	randSet := letterAlpha
	if ctype == 1 {
		randSet = letterDigit

	} else if ctype == 2 {
		randSet = letterAlpha
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = randSet[rand.Intn(len(randSet))]
	}
	return string(b)
}

func RandString(n int) string {
	return randString(n, 2)
}

func RandDigit(n int) string {
	return randString(n, 1)
}

func RandAnythingOnceFrom[S ~[]T, T any](from S) T {
	length := len(from)
	r := rand.New(rand.NewSource(int64(length)))
	return from[r.Intn(length)]
}

func RandAnythingSomeFrom[S ~[]T, T any](from S, tar int) S {
	length := len(from)
	choice := ChooseMN(length, tar)
	var result S
	for idx := range choice {
		result = append(result, from[idx])
	}
	return result
}
