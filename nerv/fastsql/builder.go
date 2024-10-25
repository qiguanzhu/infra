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
 @Time    : 2024/10/11 -- 14:18
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: builder.go
*/

package fastsql

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/qiguanzhu/infra/pkg"
	"github.com/qiguanzhu/infra/seele/zsql"
	"github.com/qiguanzhu/infra/seele/zsql/sqlutils"
	"strings"
)

type fastBuilder struct{}

// buildSelect
// sq.where including in ; not in ; like ; and etc. praser
// normally we support _limit; _orderby; _join; _having; _groupby
//
// where := map[string]interface{}{
// "age >": 100,
// "_orderby": "fieldName asc",
// "_groupby": "fieldName",
// "_having": map[string]interface{}{"foo":"bar",},
// "_limit": []uint{offset, row_count},
// "_forceindex": "PRIMARY",
// }
func (f fastBuilder) buildSelect(tableName string, where map[string]interface{}, selectedField []string) (query string, args []interface{}, err error) {
	field := "*"
	if len(selectedField) > 0 {
		field = strings.Join(selectedField, ",")
	}
	copiedWhere := sqlutils.CopyWhere(where)
	selector := sq.
		Select(field).
		From(tableName)

	for k, v := range where {
		switch k {
		case "_limit":
			_limit := v.([]uint)
			if len(_limit) != 2 && len(_limit) != 1 {
				return "", nil, errors.New("invalid limit")
			}
			if len(_limit) == 2 {
				selector.Limit(uint64(_limit[1])).Offset(uint64(_limit[0]))
			} else {
				selector.Limit(uint64(_limit[0]))
			}
			delete(copiedWhere, "_limit")
		case "_orderby":
			// we use "? desc", count by default
			// eg:  ? desc, 2
			//      select * from table1 ? desc ? desc
			//      and we can inject args to the query
			// here we just use "a asc b desc" without injections of args in order to suit one and more orderby func
			selector = selector.OrderByClause(v.(string))
			delete(copiedWhere, "_orderby")
		case "_having":
			selector = selector.Having(v.(string))
			delete(copiedWhere, "_having")
		case "_groupby":
			selector = selector.GroupBy(v.(string))
			delete(copiedWhere, "_groupby")
		}
	}

	selector, err = formatWhere(copiedWhere, selector)
	if err != nil {
		return "", nil, err
	}

	return selector.ToSql()
}

func (f fastBuilder) BuildSelect(tableName string, where map[string]interface{}, selectedField []string) (query string, args []interface{}, err error) {
	return f.buildSelect(tableName, where, selectedField)
}

func (f fastBuilder) BuildSelectWithContext(ctx context.Context, tableName string, where map[string]interface{}, selectedField []string) (query string, args []interface{}, err error) {
	return f.buildSelect(tableName, where, selectedField)
}

func (fastBuilder) BuildUpdate(tableName string, where map[string]interface{}, update map[string]interface{}) (string, []interface{}, error) {
	return sq.Update(tableName).
		SetMap(update).
		Where(sq.Eq(where)).
		ToSql()
}

func (fastBuilder) BuildDelete(tableName string, where map[string]interface{}) (string, []interface{}, error) {
	return sq.Delete("").
		From(tableName).
		Where(where).
		ToSql()
}

func (f fastBuilder) BuildInsert(tableName string, data []map[string]interface{}) (string, []interface{}, error) {
	return f.buildInsert(tableName, data, sq.Insert(tableName))
}

func (fastBuilder) buildInsert(tableName string, data []map[string]interface{}, inserter sq.InsertBuilder) (string, []interface{}, error) {
	columns := make([]string, 0)
	values := make([][]interface{}, 0)
	for _, d := range data {
		value := make([]interface{}, 0)
		if len(columns) == 0 {
			for k, v := range d {
				columns = append(columns, k)
				value = append(value, v)
			}
		} else {
			for _, k := range columns {
				if v, ok := d[k]; ok {
					value = append(value, v)
				} else {
					value = append(value, nil)
				}
			}
		}
		values = append(values, value)
	}
	fmt.Println(columns, values)
	builder := inserter.Columns(columns...)
	for _, v := range values {
		// values(1,2,3,4)
		builder = builder.Values(v...)
	}

	return builder.ToSql()
}

func (f fastBuilder) BuildInsertIgnore(tableName string, data []map[string]interface{}) (string, []interface{}, error) {
	query, args, _ := f.buildInsert(tableName, data, sq.Insert(tableName))
	return query, args, nil
}

func (f fastBuilder) BuildReplaceIgnore(tableName string, data []map[string]interface{}) (string, []interface{}, error) {
	query, args, _ := f.buildInsert(tableName, data, sq.Replace(tableName))
	return query, args, nil
}

func (fastBuilder) AggregateQuery(ctx context.Context, db zsql.XDB, tableName string, where map[string]interface{}, aggregate zsql.AggregateSymbolBuilder) (zsql.ResultResolver, error) {
	return nil, errors.New("not implemented")
}

func formatWhere(where map[string]interface{}, builder sq.SelectBuilder) (sq.SelectBuilder, error) {
	if len(where) == 0 {
		return builder, nil
	}
	var field, operator string
	var err error
	for k, v := range where {
		field, operator, err = sqlutils.SplitKey(k)
		if !sqlutils.IsStringInSlice(operator, sqlutils.OpOrder) {
			return builder, pkg.ErrBuilderUnsupportedOperator
		}
		if nil != err {
			return builder, err
		}
		var op zsql.ZSqlizer
		op, err = sqlutils.OpOp[operator](field, v)
		if nil != err {
			return builder, err
		}
		builder = builder.Where(op)
	}
	return builder, nil
}
