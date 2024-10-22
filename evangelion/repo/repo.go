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
 @Time    : 2024/10/17 -- 12:26
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: repo.go
*/

package repo

import "github.com/qiguanzhu/infra/seele/xsqlIface"

type Condition struct {
	cons map[string]interface{}
}

func NewCondition() xsqlIface.ConditionsProxy {
	return &Condition{
		cons: make(map[string]interface{}),
	}
}

// NewChanges is same as NewCondition just for difference
func NewChanges() xsqlIface.ConditionsProxy {
	return &Condition{
		cons: make(map[string]interface{}),
	}
}

func (c *Condition) Init(columns ...string) {}
func (c *Condition) SetCondition(column string, value interface{}) xsqlIface.ConditionsProxy {
	c.cons[column] = value
	return c
}
func (c *Condition) GetConditions() map[string]interface{} {
	return c.cons
}
