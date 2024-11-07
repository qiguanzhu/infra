package padding

import (
	"bytes"
	"fmt"
)

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
 @Time    : 2024/7/26 -- 11:38
 @Author  : bishop
 @Description: padding
*/

type Pad interface {
	// Padding 填充明文
	Padding(origData []byte, blockSize int) []byte
	// UnPadding 去除填充数据
	UnPadding(padded []byte, blockSize int) ([]byte, error)
}

var NoPad noPad
var PKCS5 pkcs7
var PKCS7 pkcs7
var Zero zero
var AnsiX923 ansiX923
var Iso7816 iso7816

type noPad struct{}

func (noPad) Padding(origData []byte, blockSize int) []byte {
	return origData
}

func (noPad) UnPadding(padded []byte, blockSize int) ([]byte, error) {
	return padded, nil
}

// PKCS5Padding（也称为 PKCS7Padding）是一种基于 PKCS#5 标准的填充方案，用于加密数据以确保其长度符合加密算法的块大小要求。
// 这种填充方法主要用于块加密模式，如 CBC（Cipher Block Chaining，密码块链接模式）。
type pkcs7 struct{}

func (pkcs7) Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padText...)
}

func (pkcs7) UnPadding(padded []byte, blockSize int) ([]byte, error) {
	length := len(padded)
	unPadding := int(padded[length-1])
	if unPadding > blockSize || unPadding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return padded[:(length - unPadding)], nil
}

// 也称为 Null Padding，简单地在数据末尾添加零直到满足块大小。
type zero struct{}

func (zero) Padding(origData []byte, blockSize int) []byte {
	paddingLength := blockSize - len(origData)%blockSize
	padText := bytes.Repeat([]byte{0x00}, paddingLength)
	return append(origData, padText...)
}

func (zero) UnPadding(padded []byte, blockSize int) ([]byte, error) {
	length := len(padded)
	if length == 0 || padded[length-1] != 0x00 {
		return nil, fmt.Errorf("invalid zero padding")
	}
	unPadding := int(padded[length-1])
	return padded[:length-unPadding], nil
}

// ANSI X.923 Padding：
// 与 PKCS#5 填充类似，但填充值是随机的，而不是固定模式。
type ansiX923 struct{}

func (ansiX923) Padding(origData []byte, blockSize int) []byte {
	paddingLength := blockSize - len(origData)%blockSize
	padText := make([]byte, paddingLength)
	for i := range padText {
		padText[i] = byte(paddingLength)
	}
	return append(origData, padText...)
}

func (ansiX923) UnPadding(padded []byte, blockSize int) ([]byte, error) {
	length := len(padded)
	if length == 0 {
		return nil, fmt.Errorf("invalid X.923 padding")
	}
	unPadding := int(padded[length-1])
	if unPadding > length {
		return nil, fmt.Errorf("invalid X.923 padding")
	}
	return padded[:length-unPadding], nil
}

// 用于智能卡和金融服务，填充模式为 0x80 后跟一系列 0x00。
type iso7816 struct{}

func (iso7816) Padding(origData []byte, blockSize int) []byte {
	paddingLength := blockSize - len(origData)%blockSize
	padText := []byte{0x80}
	if paddingLength > 1 {
		padText = append(padText, bytes.Repeat([]byte{0x00}, paddingLength-1)...)
	}
	return append(origData, padText...)
}

func (iso7816) UnPadding(padded []byte, blockSize int) ([]byte, error) {
	length := len(padded)
	if length == 0 || padded[0] != 0x80 {
		return nil, fmt.Errorf("invalid ISO/IEC 7816-4 padding")
	}
	for i := 1; i < length; i++ {
		if padded[i] != 0x00 {
			return padded[:i], nil
		}
	}
	return nil, fmt.Errorf("invalid ISO/IEC 7816-4 padding")
}
