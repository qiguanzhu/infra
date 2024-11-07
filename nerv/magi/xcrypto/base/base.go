package base

import (
	"encoding/base64"
	"github.com/Bishoptylaor/go-toolkit/xcrypto/base/base62"
	"github.com/Bishoptylaor/go-toolkit/xutils"
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
@Description: base64 and base62
*/

type b interface {
	SEncode(in string) string
	SDecode(in string) (string, error)

	Encode(in []byte) string
	Decode(in string) ([]byte, error)
}

var Base64 b64
var Base64Url b64u
var Base62 b62

type b64 struct{}

func (b b64) SEncode(in string) string {
	return b.Encode(xutils.Str2Bytes(in))
}

func (b b64) SDecode(in string) (string, error) {
	bs, err := b.Decode(in)
	return xutils.Bytes2Str(bs), err
}
func (b b64) Encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func (b b64) Decode(in string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(in)
}

type b64u struct{}

func (b b64u) SEncode(in string) string {
	return b.Encode(xutils.Str2Bytes(in))
}

func (b b64u) SDecode(in string) (string, error) {
	bs, err := b.Decode(in)
	return xutils.Bytes2Str(bs), err
}
func (b b64u) Encode(in []byte) string {
	return base64.URLEncoding.EncodeToString(in)
}

func (b b64u) Decode(in string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(in)
}

type b62 struct{}

func (b b62) SEncode(in string) string {
	return b.Encode(xutils.Str2Bytes(in))
}

func (b b62) SDecode(in string) (string, error) {
	bs, err := b.Decode(in)
	return xutils.Bytes2Str(bs), err
}
func (b b62) Encode(in []byte) string {
	return base62.B62StdEncoding.EncodeToString(in)
}

func (b b62) Decode(in string) ([]byte, error) {
	return base62.B62StdEncoding.DecodeString(in)
}
