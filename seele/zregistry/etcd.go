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
 @Time    : 2024/10/24 -- 13:47
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: etcd.go
*/

package zregistry

import "time"

// Node 注册中心节点信息
type Node interface {
	IsDir() bool
	Key() string
	Value() string
	Ttl() int64
	Expiration() *time.Time
	Children() []Node
}

// Handler 监听回调处理器
type Handler interface {
	// Action 事件
	Action() string
	Node() Node
}
