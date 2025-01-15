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
 @Time    : 2024/11/11 -- 12:02
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: repo_request.go
*/

package zsql

import "context"

// QueryRequest
// query related repo requests
type QueryRequest[EntityObj any] interface {
	QueryGetRequest[EntityObj]
	// GetByID func by name
	GetByID(ctx context.Context, id int64, withCache bool) (*EntityObj, error)

	QueryListRequest[EntityObj]
	// List all with offset limit
	List(ctx context.Context, offset, limit int) ([]*EntityObj, int64, error)
	// ListByIDs must be autoincrement id column
	ListByIDs(ctx context.Context, ids []int64, offset, limit int, withCache bool) ([]*EntityObj, int64, error)
}

type QueryListRequest[EntityObj any] interface {
	// ListWithConditions list by conditions with offset limit
	ListWithConditions(ctx context.Context, conditions ConditionsProxy, offset, limit int) ([]*EntityObj, int64, error)
	// ListAllWithConditions list all with conditions
	ListAllWithConditions(ctx context.Context, conditions ConditionsProxy) ([]*EntityObj, int64, error)
}

type QueryGetRequest[EntityObj any] interface {
	// Get by conditions
	Get(ctx context.Context, cond ConditionsProxy, withCache bool) (*EntityObj, error)
	// GetByColumn get obj with one column condition
	GetByColumn(ctx context.Context, column string, val interface{}, withCache bool) (*EntityObj, error)
}

// InsertRequest
// insert related repo requests
type InsertRequest[EntityObj any] interface {
	// Insert simple insert
	Insert(ctx context.Context, req *EntityObj, opAccount string) (int64, error)
}

type UpsertRequest[EntityObj any] interface {
	// Upsert
	// ddl must has unique key : id or other
	// ddl`s unique column must contained in req
	Upsert(ctx context.Context, req *EntityObj, opAccount string) (int64, error)
}

// UpdateRequest
// update related repo requests
type UpdateRequest interface {
	Update(ctx context.Context, conditions, changes ConditionsProxy, opAccount string) error // update
}

type DeleteRequest interface {
	DeleteLogicRequest
	DeletePhysicRequest
}

type DeleteLogicRequest interface {
	// Del logically delete from db. mostly update state || status
	Del(ctx context.Context, cond ConditionsProxy, opAccount string) error
	DelById(ctx context.Context, id int64, opAccount string) error
}

type DeletePhysicRequest interface {
	// RealDel physically delete from db
	RealDel(ctx context.Context, cond ConditionsProxy, opAccount string) error
}

type CommonExecRequest[EntityObj any] interface {
	Total(ctx context.Context) (int, error)
	Count(ctx context.Context, cond ConditionsProxy) (int, error)
}

// ComplexRequest
// usually ToSql will pass from serviceModel through repoModel direct to daoModel and be executed
// if you concern any SQL injection issue, you can easily build your own complex query function at any logical name in your repoModel
type ComplexRequest[EntityObj any] interface {
	ComplexQueryRequest[EntityObj]
	ComplexExecRequest
}

// ComplexQueryRequest ...
type ComplexQueryRequest[EntityObj any] interface {
	ComplexQuery(ctx context.Context, ts ToSql, columns []string, b BindFunc) ([]*EntityObj, error)
}

// ComplexExecRequest ...
type ComplexExecRequest interface {
	ComplexExec(ctx context.Context, ts ToSql) error
}
