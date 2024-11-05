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
 @Time    : 2024/11/4 -- 17:13
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: zatomic.go : easier to use atomic
*/

package zatomic

import "time"

type ZAtomic[Elem any] interface {
	Set(Elem)
	Get() Elem
	CompareAndSwap(Elem, Elem) bool
}

type ZAtomicNumLike[Elem int32 | int64 | float64 | time.Duration] interface {
	Add(Elem) Elem
	ZAtomic[Elem]
}
