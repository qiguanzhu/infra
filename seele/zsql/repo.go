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

type DataObj[SvcObj any] interface {
	ToEntity(ctx context.Context) *SvcObj
}

// DaoModel
// interactive between DataObj and db
// reflect db instance rows to DataObj
type DaoModel[SvcObj any, DObj DataObj[SvcObj]] interface {
	Init(cons SqlConstructor, tableName func() string, omits func() []string, b BindFunc)
	TableName() string
	Omits() []string
	GetScanner() ScannerProxy
	GetBuilder() BuilderProxy

	SelectOne(ctx context.Context, db XDB, where map[string]interface{}) (DObj, error)
	SelectMulti(ctx context.Context, db XDB, where map[string]interface{}) ([]DObj, error)
	Insert(ctx context.Context, db XDB, data []map[string]interface{}) (int64, error)
	Update(ctx context.Context, db XDB, where, data map[string]interface{}) (int64, error)
	Delete(ctx context.Context, db XDB, where map[string]interface{}) (int64, error)
	CountOf(ctx context.Context, db XDB, where map[string]interface{}) (count int, err error)
	ToEntity(ctx context.Context, t DObj) *SvcObj
	MultiToEntity(ctx context.Context, ts []DObj) []*SvcObj
}

type ComplexQueryMod[ans any] func(
	ctx context.Context,
	db XDB,
	scanner ScannerProxy,
	f ToSql,
	bind BindFunc,
) (res []ans, err error)

type ComplexExecMod func(
	ctx context.Context,
	db XDB,
	ts ToSql,
) (int64, error)

// RepoModel
// basic function definition of sql repo
// normally we access to these functions from object Service level
// so we define SvcObj as the entity used in Service level
// this model transfer ServiceObj to dataObj as built-in, so users from upper level can ignore whatever in lower level
type RepoModel[SvcObj any] interface {
	QueryRequest[SvcObj]
	InsertRequest[SvcObj]
	UpdateRequest[SvcObj]

	Valid(obj *SvcObj) (bool, error)                                                   // verify SvcObj values
	Audit(ctx context.Context, id int64, action string, remark string, changes ...any) // log repo actions for audit
}

// QueryRequest
// query related repo requests
type QueryRequest[SvcObj any] interface {
	Get(context.Context, ConditionsProxy) (*SvcObj, error)                            // get by changes
	GetByID(context.Context, int64) (*SvcObj, error)                                  // get by id
	GetByColumn(ctx context.Context, column string, val interface{}) (*SvcObj, error) // get by one column

	List(ctx context.Context, offset, limit int) ([]*SvcObj, int64, error)                                           // list all with page
	ListWithConditions(ctx context.Context, conditions ConditionsProxy, offset, limit int) ([]*SvcObj, int64, error) // list by changes
}

// InsertRequest
// insert related repo requests
type InsertRequest[SvcObj any] interface {
	Insert(ctx context.Context, req *SvcObj, opAccount string) (int64, error)
}

// UpdateRequest
// update related repo requests
type UpdateRequest[SvcObj any] interface {
	Del(ctx context.Context, column string, val interface{}, opAccount string) error         // 默认逻辑删除 default del by id
	RealDel(ctx context.Context, column string, val interface{}, opAccount string) error     // 从库表物理删除 default del by id
	Update(ctx context.Context, conditions, changes ConditionsProxy, opAccount string) error // update
}

// ComplexRequest
// usually ToSql will pass from serviceModel through repoModel direct to daoModel and be executed
// if you concern any SQL injection issue, you can easily build your own complex query function at any logical name in your repoModel
type ComplexRequest[SvcObj any] interface {
	ComplexQuery(ctx context.Context, ts ToSql, columns []string, b BindFunc) ([]*SvcObj, error)
	ComplexExec(ctx context.Context, ts ToSql) error
}

type ConditionsProxy interface {
	Init(columns ...string)
	SetCondition(column string, value interface{}) ConditionsProxy
	GetConditions() map[string]interface{}
}
