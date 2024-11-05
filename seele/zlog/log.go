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
 @Time    : 2024/10/28 -- 12:06
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: log.go
*/

package zlog

import "context"

type ZLogger interface {
	Debug(ctx context.Context, args ...any)
	Info(ctx context.Context, args ...any)
	Warn(ctx context.Context, args ...any)
	Error(ctx context.Context, args ...any)
	Panic(ctx context.Context, args ...any)

	Debugf(ctx context.Context, fmt string, args ...any)
	Infof(ctx context.Context, fmt string, args ...any)
	Warnf(ctx context.Context, fmt string, args ...any)
	Errorf(ctx context.Context, fmt string, args ...any)
	Panicf(ctx context.Context, fmt string, args ...any)
}
