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
 @Time    : 2024/10/25 -- 14:21
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: dao.go
 @ currently a copy of xsql.xdao. trying to district behavior
*/

package fastsql

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qiguanzhu/infra/nerv/xlog"
	"github.com/qiguanzhu/infra/seele/zsql"
	"github.com/qiguanzhu/infra/seele/zsql/sqlutils"
)

type FastDao[SvcObj any, DObj zsql.DataObj[SvcObj]] interface {
	zsql.DaoModel[SvcObj, DObj]
}

type DefaultDao[SvcObj any, DObj zsql.DataObj[SvcObj]] struct {
	tableName func() string
	omits     func() []string
	_scanner  zsql.ScannerProxy
	_builder  zsql.BuilderProxy
	bind      zsql.BindFunc
}

// TableName anyone wrapper *DefaultDao should write their own
func (dao *DefaultDao[SvcObj, DObj]) TableName() string { return dao.tableName() }

// Omits anyone wrapper *DefaultDao should write their own
func (dao *DefaultDao[SvcObj, DObj]) Omits() []string { return dao.omits() }

func (dao *DefaultDao[SvcObj, DObj]) Init(cons zsql.SqlConstructor, tableName func() string, omits func() []string, b zsql.BindFunc) {
	dao._builder = cons.GetBuilder()
	dao._scanner = cons.GetScanner()
	dao.omits = omits
	dao.tableName = tableName
	dao.bind = b
}

func (dao *DefaultDao[SvcObj, DObj]) GetScanner() zsql.ScannerProxy {
	return dao._scanner
}

func (dao *DefaultDao[SvcObj, DObj]) GetBuilder() zsql.BuilderProxy {
	return dao._builder
}

func (dao *DefaultDao[SvcObj, DObj]) SelectOne(ctx context.Context, db zsql.XDB, where map[string]interface{}) (res DObj, err error) {
	if nil == db {
		return res, errors.New("manager.XDB object couldn't be nil")
	}
	tar := sqlutils.CopyWhere(where)
	if _, ok := tar["_limit"]; !ok {
		tar["_limit"] = []uint{0, 1}
	}
	cond, vals, err := dao._builder.BuildSelectWithContext(ctx, dao.TableName(), tar, dao.Omits())
	if nil != err {
		return res, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	rows, err := db.QueryContext(ctx, cond, vals...)
	if nil != err || nil == rows {
		return res, err
	}
	defer rows.Close()
	err = dao._scanner.Scan(rows, &res, dao.bind)
	fmt.Println("res", res)
	return res, err
}

func (dao *DefaultDao[SvcObj, DObj]) SelectMulti(ctx context.Context, db zsql.XDB, where map[string]interface{}) (res []DObj, err error) {
	if nil == db {
		return res, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildSelectWithContext(ctx, dao.TableName(), where, dao.Omits())
	if nil != err {
		return nil, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	rows, err := db.QueryContext(ctx, cond, vals...)
	if nil != err || nil == rows {
		return nil, err
	}
	defer rows.Close()
	err = dao._scanner.Scan(rows, &res, dao.bind)
	return res, err
}

func (dao *DefaultDao[SvcObj, DObj]) Insert(ctx context.Context, db zsql.XDB, data []map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildInsert(dao.TableName(), data)
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err || nil == result {
		return 0, err
	}
	return result.LastInsertId()
}

func (dao *DefaultDao[SvcObj, DObj]) Update(ctx context.Context, db zsql.XDB, where, data map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildUpdate(dao.TableName(), where, data)
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err {
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *DefaultDao[SvcObj, DObj]) Delete(ctx context.Context, db zsql.XDB, where map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildDelete(dao.TableName(), where)
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err {
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *DefaultDao[SvcObj, DObj]) CountOf(ctx context.Context, db zsql.XDB, where map[string]interface{}) (count int, err error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildSelect(dao.TableName(), where, []string{"count(*) as count"})
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	rows, err := db.QueryContext(ctx, cond, vals...)
	if nil != err {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}
	return
}

func (dao *DefaultDao[SvcObj, DObj]) ToEntity(ctx context.Context, t DObj) *SvcObj {
	return t.ToEntity(ctx)
}

func (dao *DefaultDao[SvcObj, DObj]) MultiToEntity(ctx context.Context, ts []DObj) []*SvcObj {
	var res []*SvcObj
	for _, t := range ts {
		res = append(res, t.ToEntity(ctx))
	}
	return res
}

// ComplexQuery
// you can use this default logic or
// you can build your own query logic with or without tableName or columns
// depends on your ToSql func
func ComplexQuery[ans any](tableName string, columns ...string) zsql.ComplexQueryMod[ans] {
	return func(
		ctx context.Context,
		db zsql.XDB,
		scanner zsql.ScannerProxy,
		f zsql.ToSql,
		bind zsql.BindFunc,
	) (res []ans, err error) {
		if nil == db {
			return nil, errors.New("manager.XDB object couldn't be nil")
		}
		cond, vals, err := f(tableName, columns...)
		if nil != err {
			return nil, err
		}
		xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
		xlog.Infof(ctx, "build cond: %s vals: %v", cond, vals)
		rows, err := db.QueryContext(ctx, cond, vals...)
		if nil != err || nil == rows {
			return nil, err
		}
		defer rows.Close()
		err = scanner.Scan(rows, &res, bind)
		return res, err
	}
}

func ComplexExec(tableName string) zsql.ComplexExecMod {
	return func(
		ctx context.Context,
		db zsql.XDB,
		f zsql.ToSql,
	) (int64, error) {
		if nil == db {
			return 0, errors.New("manager.XDB object couldn't be nil")
		}
		cond, vals, err := f(tableName)
		if nil != err {
			return 0, err
		}
		xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
		result, err := db.ExecContext(ctx, cond, vals...)
		if nil != err {
			return 0, err
		}
		return result.RowsAffected()
	}
}
