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
 @Time    : 2024/10/21 -- 09:54
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: fastsql_test.go
*/

package onesql

import (
	"github.com/qiguanzhu/infra/nerv/magi/xreflect"
	"testing"
)

func TestMirroring(t *testing.T) {
	type Sa struct {
		A string
	}
	type Sb struct {
		B string
	}

	type Ss struct {
		S int
	}

	src := &Sa{A: "a"}
	sb := new(Sb)
	ss := new(Ss)
	sa := new(Sa)
	sa2 := Sa{}

	err := xreflect.Mirroring(src, sb)
	t.Log(src, sb, err)
	// &{a} &{} src and tar must have the same struct type

	err = xreflect.Mirroring(src, ss)
	t.Log(src, ss, err)
	// &{a} &{0} src and tar must have the same struct type

	err = xreflect.Mirroring(src, sa)
	t.Log(src, sa, err)
	// &{a} &{a} <nil>
	// Success!!

	err = xreflect.Mirroring(src, sa2)
	t.Log(src, sa2, err)
	// &{a} {} src must be a pointer to a struct

}
