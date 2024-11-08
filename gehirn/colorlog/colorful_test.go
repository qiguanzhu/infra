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
 @Time    : 2024/10/25 -- 18:55
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: xlog_test.go
*/

package colorlog

import (
	"context"
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	// Cyan, "青色 debug", Reset, CyanBright, "高亮 debug", Reset, "恢复默认颜色", CyanDelLine, "删除线", Reset, CyanUnderLine, "下划线", Reset, CyanBevel, "斜体 debug", Reset, CyanBg, "背景", Reset
	// default log
	SetLevel(PanicLevel)
	SetColorful(true)
	ctx := context.Background()

	Debug(ctx, "Debug 日志")
	Info(ctx, "Info 日志")
	Warn(ctx, "Warn 日志")
	Error(ctx, "Error 日志")
	Panic(ctx, "Panic 日志")

	log.Println("================")
	Debugf(ctx, "Debugf %s, %s, %d", "日志", "asdasda", 4)
	Infof(ctx, "Infof %s", "日志")
	Warnf(ctx, "Warnf %s", "日志")
	Errorf(ctx, "Errorf %s", "日志")
	Panicf(ctx, "Panicf %s", "日志")

	log.Println("================")

	xLogger := NewLogger()
	// xLogger.Color(Red)
	xLogger.SetLevel(WarnLevel)

	xLogger.Debug(ctx, "new logger Debug 日志")
	xLogger.Info(ctx, "new logger Info 日志")
	xLogger.Warn(ctx, "new logger Warn 日志")
	xLogger.Error(ctx, "new logger Error 日志")

	log.Println("================")

	xLogger.Debugf(ctx, "new logger Debugf %s", "日志")
	xLogger.Infof(ctx, "new logger Infof %s", "日志")
	xLogger.Warnf(ctx, "new logger Warnf %s", "日志")
	xLogger.Errorf(ctx, "new logger Errorf %s", "日志")

	// color log

}
