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
 @Time    : 2024/8/30 -- 15:05
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: cipher.go
*/

package aescipher

import (
	"context"
	"github.com/qiguanzhu/infra/nerv/magi/xcrypto/padding"
	"hash"
)

const (
	pkc5SaltLen          = 8
	pkc5DefaultIter      = 2048
	pkc5DefaultMagicWord = "Salted__"
	maxIvLen             = 16
)

type Cipher func(ctx context.Context, data, key, iv []byte, pad padding.Pad) ([]byte, error)

// AesCipher 标准参数加密解密功能入口
type AesCipher interface {
	Encrypt(ctx context.Context, plainText, secretKey, ivAes []byte, pad padding.Pad) (cipherText []byte, err error)
	Decrypt(ctx context.Context, cipherText, secretKey, ivAes []byte, pad padding.Pad) (plainText []byte, err error)
}

type AesCipherByBase64 interface {
	EncryptBase64(ctx context.Context, plainText, secretKey, ivAes []byte, pad padding.Pad) (cipherText string, err error)
	DecryptBase64(ctx context.Context, cipherText, secretKey, ivAes []byte, pad padding.Pad) (plainText []byte, err error)
}

type AesCipherBySalt interface {
	EncryptWithSalt(ctx context.Context, origData, key []byte, iterCount int, magic string, h func() hash.Hash, pad padding.Pad, f Cipher) ([]byte, error)
	DecryptWithSalt(ctx context.Context, encrypted, key []byte, iterCount int, magic string, h func() hash.Hash, pad padding.Pad, f Cipher) ([]byte, error)
}

type AesCipherByNonce interface {
	EncryptWithNonce(ctx context.Context, origData, key, nonce, additional []byte) ([]byte, error)
	DecryptWithNonce(ctx context.Context, encrypted, key, nonce, additional []byte) ([]byte, error)
}

func CBC() AesCipher {
	return cbc{}
}
func CBCBase64() AesCipherByBase64 {
	return cbc{}
}
func CBCBySalt() AesCipherBySalt {
	return cbc{}
}

func CFB() AesCipher {
	return cfb{}
}

func ECB() AesCipher {
	return ecb{}
}

func GCM() AesCipher {
	return gcm{}
}
func GCMByNonce() AesCipherByNonce {
	return gcm{}
}
