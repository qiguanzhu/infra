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
 @Description: register.go
*/

package zregistry

import (
	"context"
	"time"
)

// Register 注册中心相关接口
type Register interface {
	// Get 获取指定key的值
	Get(ctx context.Context, path string) (string, error)
	// GetNode 获取指定key对应的节点，会将节点及子节点返回
	GetNode(ctx context.Context, path string) (Node, error)
	// Set 设置key，永不过期
	Set(ctx context.Context, path, val string) error
	// CreateDir 创建一个目录节点，不设置值
	CreateDir(ctx context.Context, path string) error
	// SetTtl 设置key值，并指定过期时间
	SetTtl(ctx context.Context, path, val string, ttl time.Duration) error
	// RefreshTtl 刷新节点的过期时间，不更新值
	RefreshTtl(ctx context.Context, path string, ttl time.Duration) error
	// SetNx 设置key值，只有不存在时才设置，否则失败
	SetNx(ctx context.Context, path, val string) error
	// Regist 执行注册，注册后会一直维持心跳。调用注册时会将当前值设置到节点上
	Regist(ctx context.Context, path, val string, heatBeat time.Duration, ttl time.Duration) error
	// Watch 监听一个节点
	Watch(ctx context.Context, path string, handler func(Handler))
	//	锁相关接口
	// id相关接口
	// 分布式选主接口

}
