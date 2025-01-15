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
 @Time    : 2024/10/9 -- 15:13
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: repo.go
*/

package zsql

import (
	"context"
)

// RepoModel
// common function definition of sql repo
// normally we access to these functions from object Service level
// so we define EntityObj as the entity used in Service level
// this model transfer ServiceObj to dataObj as built-in, so users from upper level can ignore anything happened in lower level
// sometimes we do not need these much of functions defined in interface
// you can use dismembered interfaces from repo_request.go and build your own RepoModel anywhere
type RepoModel[EntityObj any] interface {
	QueryRequest[EntityObj]
	InsertRequest[EntityObj]
	UpsertRequest[EntityObj]
	UpdateRequest
	DeleteRequest

	CommonExecRequest[EntityObj]

	Valid(obj *EntityObj) (bool, error)                                                // verify EntityObj values
	Audit(ctx context.Context, id int64, action string, remark string, changes ...any) // log repo actions for audit
}

type ConditionsProxy interface {
	// Set set kv
	Set(column string, value interface{}) ConditionsProxy
	// Export whole map
	Export() map[string]interface{}
	// ToString conditions to string
	ToString() string
}

type ComplexQueryMod[ans any] func(
	ctx context.Context,
	db XDB,
	scanner Scanner,
	f ToSql,
	bind BindFunc,
) (res []ans, err error)

type ComplexExecMod func(
	ctx context.Context,
	db XDB,
	ts ToSql,
) (int64, error)
