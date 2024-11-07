package xcrypto

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/qiguanzhu/infra/nerv/magi/xstring"
	"hash"
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
@Time    : 2024/7/25 -- 14:00
@Author  : bishop
@Description: hash 加密算法集合。默认都返回 16 进制字符串
*/

type Hash interface {
	DoString(string) string
	DoBytes([]byte) string
}

var MD5 _MD5
var SHA1 _SHA1
var SHA256 _SHA256
var SHA512 _SHA512

type _MD5 struct{}

func (h _MD5) DoString(src string) string {
	return h.DoBytes(xstring.Str2Bytes(src))
}
func (h _MD5) DoBytes(bs []byte) string {
	sum := md5.Sum(bs)
	return hex.EncodeToString(sum[:])
}

type _SHA1 struct{}

func (h _SHA1) DoString(src string) string {
	return h.DoBytes(xstring.Str2Bytes(src))
}
func (h _SHA1) DoBytes(bs []byte) string {
	sum := sha1.Sum(bs)
	return hex.EncodeToString(sum[:])
}

type _SHA256 struct{}

func (h _SHA256) DoString(src string) string {
	return h.DoBytes(xstring.Str2Bytes(src))
}
func (h _SHA256) DoBytes(bs []byte) string {
	sum := sha256.Sum256(bs)
	return hex.EncodeToString(sum[:])
}

type _SHA512 struct{}

func (h _SHA512) DoString(src string) string {
	return h.DoBytes(xstring.Str2Bytes(src))
}
func (h _SHA512) DoBytes(bs []byte) string {
	sum := sha512.Sum512(bs)
	return hex.EncodeToString(sum[:])
}

func HMacSHA1(key []byte, bs []byte) string {
	return HMacEx(sha1.New, key, bs)
}

func HMacSHA256(key []byte, bs []byte) string {
	return HMacEx(sha256.New, key, bs)
}

func HMacMD5(key []byte, bs []byte) string {
	return HMacEx(md5.New, key, bs)
}

func HMacEx(h func() hash.Hash, key []byte, bs []byte) string {
	mac := hmac.New(h, key)
	mac.Write(bs)
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum[:])
}

type HashMethod struct {
	h crypto.Hash
}

func NewHashMethod(h crypto.Hash) *HashMethod {
	var hm = &HashMethod{}
	hm.h = h
	return hm
}

func (hm *HashMethod) Sign(data []byte) ([]byte, error) {
	var h = hm.h.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (hm *HashMethod) Verify(data []byte, signature []byte) error {
	var h = hm.h.New()
	if _, err := h.Write(data); err != nil {
		return err
	}
	var hashed = h.Sum(nil)
	if bytes.Compare(hashed, signature) == 0 {
		return nil
	}
	return ErrVerification
}
