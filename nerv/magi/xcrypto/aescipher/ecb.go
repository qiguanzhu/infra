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
 @Time    : 2024/10/11 -- 09:56
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: ecb.go
*/

package aescipher

import (
	"context"
	"crypto/aes"
	"github.com/qiguanzhu/infra/nerv/magi/xcrypto/padding"
)

type ecb struct{}

func (ecb) Encrypt(ctx context.Context, origData, key, _ []byte, pad padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = pad.Padding(origData, blockSize)
	encrypted := make([]byte, len(origData))
	var start = 0
	var end = blockSize
	for start < len(origData) {
		block.Encrypt(encrypted[start:end], origData[start:end])
		start += blockSize
		end += blockSize
	}
	return encrypted, nil
}

func (ecb) Decrypt(ctx context.Context, encrypted, key, _ []byte, pad padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData := make([]byte, len(encrypted))
	var start = 0
	var end = blockSize

	for start < len(encrypted) {
		block.Decrypt(origData[start:end], encrypted[start:end])
		start = start + blockSize
		end = end + blockSize
	}
	return pad.UnPadding(origData, blockSize)
}
