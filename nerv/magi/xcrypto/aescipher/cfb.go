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
 @Time    : 2024/10/11 -- 09:55
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: cfb.go
*/

package aescipher

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"github.com/qiguanzhu/infra/nerv/magi/xcrypto/padding"
)

type cfb struct{}

func (cfb) Encrypt(ctx context.Context, origData, key, iv []byte, pad padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = pad.Padding(origData, blockSize)
	blockMode := cipher.NewCFBEncrypter(block, iv[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.XORKeyStream(encrypted, origData)
	return encrypted, nil
}

func (cfb) Decrypt(ctx context.Context, encrypted, key, iv []byte, pad padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCFBDecrypter(block, iv[:blockSize]) // 初始向量的长度必须等于块block的长度
	origData := make([]byte, len(encrypted))
	blockMode.XORKeyStream(origData, encrypted)
	return pad.UnPadding(origData, blockSize)
}
