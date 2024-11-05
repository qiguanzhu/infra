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
 @Time    : 2024/10/25 -- 10:45
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: crond.go
*/

package zcron

import (
	"context"
	"time"
)

// CronRunnerProxy
// entry of cron CronSchedulerModel
type CronRunnerProxy[unit CronUnit] interface {
	Register(remark string, u unit) CronRunnerProxy[unit]
	ForceRegister(remark string, u unit) CronRunnerProxy[unit] // if you want to change one unit already registered in remark
	Run(ctx context.Context)
}

// CronSchedulerModel
// any CronUnit register to CronRunner will run themselves
// you should register CronCore separately
type CronSchedulerModel interface {
	Register(cu CronUnit) error
}

type CronUnit interface {
	Init(ctx context.Context) error
	Name() string
	Do()
	GetTicker() <-chan time.Time
}
