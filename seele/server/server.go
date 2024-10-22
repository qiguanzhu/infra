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
 @Time    : 2024/10/17 -- 18:07
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: server.go
*/

package server

import (
	"context"
	"github.com/qiguanzhu/infra/seele/xconfigIface"
)

type ServConfType string

const (
	SERV ServConfType = "SERV"
	DB   ServConfType = "DB"
)

type Server[servInfo any] interface {
	ServName(ctx context.Context) string               // eg: user
	ServLocation(ctx context.Context) string           // location
	ServIp(ctx context.Context) string                 // ip
	ServId(ctx context.Context) int                    // e.g: 1
	CopyName(ctx context.Context) string               // eg: trade/points1
	Startup(ctx context.Context) error                 // 启动服务
	Shutdown(ctx context.Context) error                // 关停服务
	AppendShutdownCallback(context.Context, func())    // 添加关停回调函数
	IsLocalRunning(ctx context.Context) bool           // return true if server is local running
	IsStopped(ctx context.Context) bool                // is stopped
	ServInfos(ctx context.Context) map[string]servInfo // serv infos
	ConfigCenter(ctx context.Context) xconfigIface.ConfigCenter
}
