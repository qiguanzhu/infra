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
 @Time    : 2024/8/30 -- 16:34
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: cbc.go
*/

package aescipher

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"github.com/qiguanzhu/infra/nerv/magi/xcrypto/base"
	"github.com/qiguanzhu/infra/nerv/magi/xcrypto/padding"
	"github.com/qiguanzhu/infra/nerv/magi/xrand"
	"golang.org/x/crypto/pbkdf2"
	"hash"
)

type cbc struct{}

// aesCBCEncrypt 加密 由key的长度决定是128, 192 还是 256
func aesCBCEncrypt(_ context.Context, origData, key, iv []byte, pad padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为128位，所以blockSize=16，单位字节
	// 此处没有固定key的长度，可以由上层调用再封装常用功能
	blockSize := block.BlockSize()
	origData = pad.Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize]) // 初始向量的长度必须等于块block的长度16字节
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// aesCBCDecrypt 解密
func aesCBCDecrypt(_ context.Context, encrypted, key, iv []byte, pad padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize]) // 初始向量的长度必须等于块block的长度
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData, _ = pad.UnPadding(origData, blockSize)
	return origData, nil
}

// Encrypt 普通加密
func (cbc) Encrypt(ctx context.Context, origData, key, iv []byte, pad padding.Pad) ([]byte, error) {
	return aesCBCEncrypt(ctx, origData, key, iv, pad)
}

// EncryptBase64 加密 结果返回base64编码后的string
func (cbc) EncryptBase64(ctx context.Context, origData, key, iv []byte, pad padding.Pad) (string, error) {
	bs, err := aesCBCEncrypt(ctx, origData, key, iv, pad)
	if err != nil {
		return "", err
	}
	return base.Base64.Encode(bs), nil
}

// Decrypt 普通解密
func (cbc) Decrypt(ctx context.Context, encrypted, key, iv []byte, pad padding.Pad) ([]byte, error) {
	return aesCBCDecrypt(ctx, encrypted, key, iv, pad)
}

// DecryptBase64 base64编码后的加密串，返回原始数据
func (cbc) DecryptBase64(ctx context.Context, baseEncrypted, key, iv []byte, pad padding.Pad) ([]byte, error) {
	encrypted, err := base.Base64.Decode(string(baseEncrypted))
	if err != nil {
		return []byte{}, err
	}
	bs, err := aesCBCEncrypt(ctx, encrypted, key, iv, pad)
	if err != nil {
		return []byte{}, err
	}
	return bs, nil
}

func (cbc) EncryptWithSalt(ctx context.Context, origData, key []byte, iter int, magic string, h func() hash.Hash, pad padding.Pad, f Cipher) ([]byte, error) {
	if iter <= 0 {
		iter = pkc5DefaultIter
	}

	if h == nil {
		h = md5.New
	}

	var salt = xrand.RandString(pkc5SaltLen)
	var sKey = pbkdf2.Key(key, []byte(salt), iter, len(key), h)
	var sIV = pbkdf2.Key(sKey, []byte(salt), iter, maxIvLen, h)

	var encrypted, err = f(ctx, origData, sKey, sIV, pad)

	encrypted = append([]byte(salt), encrypted...)
	encrypted = append([]byte(magic), encrypted...)

	return encrypted, err
}

func (cbc) DecryptWithSalt(ctx context.Context, encrypted, key []byte, iterCount int, magic string, h func() hash.Hash, pad padding.Pad, f Cipher) ([]byte, error) {
	if iterCount <= 0 {
		iterCount = pkc5DefaultIter
	}

	if h == nil {
		h = md5.New
	}

	var salt = encrypted[len(magic) : len(magic)+pkc5SaltLen]
	var sKey = pbkdf2.Key(key, salt, iterCount, len(key), h)
	var sIV = pbkdf2.Key(sKey, salt, iterCount, maxIvLen, h)

	var plaintext, err = f(ctx, encrypted[len(magic)+pkc5SaltLen:], sKey, sIV, pad)
	return plaintext, err
}
