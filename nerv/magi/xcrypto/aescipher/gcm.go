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
 @Time    : 2024/10/11 -- 09:58
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: gcm.go
*/

package aescipher

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"github.com/qiguanzhu/infra/nerv/magi/xcrypto/padding"
	"github.com/qiguanzhu/infra/nerv/magi/xrand"
)

type gcm struct{}

func (gcm) Encrypt(ctx context.Context, origData, key, additional []byte, _ padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := xrand.RandString(blockMode.NonceSize())

	return blockMode.Seal([]byte(nonce), []byte(nonce), origData, additional), nil
}

func (gcm) Decrypt(ctx context.Context, encrypted, key, additional []byte, _ padding.Pad) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := blockMode.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	var nonce []byte
	nonce, encrypted = encrypted[:nonceSize], encrypted[nonceSize:]
	return blockMode.Open(nil, nonce, encrypted, additional)
}

func (gcm) EncryptWithNonce(ctx context.Context, origData, key, nonce, additional []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != blockMode.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size, must contain %d characters", blockMode.NonceSize())
	}

	return blockMode.Seal(nil, nonce, origData, additional), nil
}

func (gcm) DecryptWithNonce(ctx context.Context, encrypted, key, nonce, additional []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != blockMode.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size, must contain %d characters", blockMode.NonceSize())
	}

	return blockMode.Open(nil, nonce, encrypted, additional)
}
