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
 @Time    : 2024/8/30 -- 14:54
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: nonce.go
*/

package xcrypto

import "github.com/qiguanzhu/infra/nerv/magi/xrand"

const (
	// NonceSymbols 随机字符串可用字符集
	NonceSymbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// NonceLength 随机字符串的长度
	NonceLength = 32
)

// GenerateNonce 生成一个长度为 NonceLength 的随机字符串（只包含大小写字母与数字）
func GenerateNonce() (string, error) {
	return string(xrand.RandAnythingSomeFrom([]byte(NonceSymbols), NonceLength)), nil
}
