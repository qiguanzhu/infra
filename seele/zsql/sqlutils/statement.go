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
 @Time    : 2024/10/11 -- 15:03
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: statement.go
*/

package sqlutils

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/qiguanzhu/infra/pkg"
	"github.com/qiguanzhu/infra/seele/zsql"
	"reflect"
	"strings"
)

type eleOrderBy struct {
	field, order string
}

type StatementClauses struct {
	orderBy    string
	limit      *eleLimit
	groupBy    string
	forceIndex string
	conditions []zsql.Comparable
}

type whereMapSet struct {
	set map[string]map[string]interface{}
}

func (w *whereMapSet) add(op, field string, val interface{}) {
	if nil == w.set {
		w.set = make(map[string]map[string]interface{})
	}
	s, ok := w.set[op]
	if !ok {
		s = make(map[string]interface{})
		w.set[op] = s
	}
	s[field] = val
}

type eleLimit struct {
	begin, step int
}

// NullType is the NULL type in mysql
type NullType byte

func (nt NullType) String() string {
	if nt == IsNull {
		return "IS NULL"
	}
	return "IS NOT NULL"
}

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull

	// ContextKeyForceIndex by context use force index key
	ContextKeyForceIndex = "ForceIndex"
)

type nullComparable map[string]interface{}

func (n nullComparable) Build() ([]string, []interface{}) {
	length := len(n)
	if nil == n || 0 == length {
		return nil, nil
	}
	sortedKey := make([]string, 0, length)
	cond := make([]string, 0, length)
	for k := range n {
		sortedKey = append(sortedKey, k)
	}
	defaultSortAlgorithm(sortedKey)
	for _, field := range sortedKey {
		v, ok := n[field]
		if !ok {
			continue
		}
		rv, ok := v.(NullType)
		if !ok {
			continue
		}
		cond = append(cond, field+" "+rv.String())
	}
	return cond, nil
}

type nilComparable byte

func (n nilComparable) Build() ([]string, []interface{}) {
	return nil, nil
}

// Like means like
type Like map[string]interface{}

// Build implements the Comparable interface
func (l Like) Build() ([]string, []interface{}) {
	if nil == l || 0 == len(l) {
		return nil, nil
	}
	var cond []string
	var vals []interface{}
	for k := range l {
		cond = append(cond, k)
	}
	defaultSortAlgorithm(cond)
	for j := 0; j < len(cond); j++ {
		val := l[cond[j]]
		cond[j] = cond[j] + " LIKE ?"
		vals = append(vals, val)
	}
	return cond, vals
}

// NotLike ...
type NotLike map[string]interface{}

// Build implements the Comparable interface
func (l NotLike) Build() ([]string, []interface{}) {
	if nil == l || 0 == len(l) {
		return nil, nil
	}
	var cond []string
	var vals []interface{}
	for k := range l {
		cond = append(cond, k)
	}
	defaultSortAlgorithm(cond)
	for j := 0; j < len(cond); j++ {
		val := l[cond[j]]
		cond[j] = cond[j] + " NOT LIKE ?"
		vals = append(vals, val)
	}
	return cond, vals
}

// Eq means equal(=)
type Eq map[string]interface{}

// Build implements the Comparable interface
func (e Eq) Build() ([]string, []interface{}) {
	return build(e, "=")
}

// Ne means Not Equal(!=)
type Ne map[string]interface{}

// Build implements the Comparable interface
func (n Ne) Build() ([]string, []interface{}) {
	return build(n, "!=")
}

// Lt means less than(<)
type Lt map[string]interface{}

// Build implements the Comparable interface
func (l Lt) Build() ([]string, []interface{}) {
	return build(l, "<")
}

// Lte means less than or equal(<=)
type Lte map[string]interface{}

// Build implements the Comparable interface
func (l Lte) Build() ([]string, []interface{}) {
	return build(l, "<=")
}

// Gt means greater than(>)
type Gt map[string]interface{}

// Build implements the Comparable interface
func (g Gt) Build() ([]string, []interface{}) {
	return build(g, ">")
}

// Gte means greater than or equal(>=)
type Gte map[string]interface{}

// Build implements the Comparable interface
func (g Gte) Build() ([]string, []interface{}) {
	return build(g, ">=")
}

// In means in
type In map[string][]interface{}

// Build implements the Comparable interface
func (i In) Build() ([]string, []interface{}) {
	if nil == i || 0 == len(i) {
		return nil, nil
	}
	var cond []string
	var vals []interface{}
	for k := range i {
		cond = append(cond, k)
	}
	defaultSortAlgorithm(cond)
	for j := 0; j < len(cond); j++ {
		val := i[cond[j]]
		cond[j] = buildIn(cond[j], val)
		vals = append(vals, val...)
	}
	return cond, vals
}

func buildIn(field string, vals []interface{}) (cond string) {
	cond = strings.TrimRight(strings.Repeat("?,", len(vals)), ",")
	cond = fmt.Sprintf("%s IN (%s)", quoteField(field), cond)
	return
}

// NotIn means not in
type NotIn map[string][]interface{}

// Build implements the Comparable interface
func (i NotIn) Build() ([]string, []interface{}) {
	if nil == i || 0 == len(i) {
		return nil, nil
	}
	var cond []string
	var vals []interface{}
	for k := range i {
		cond = append(cond, k)
	}
	defaultSortAlgorithm(cond)
	for j := 0; j < len(cond); j++ {
		val := i[cond[j]]
		cond[j] = buildNotIn(cond[j], val)
		vals = append(vals, val...)
	}
	return cond, vals
}

func buildNotIn(field string, vals []interface{}) (cond string) {
	cond = strings.TrimRight(strings.Repeat("?,", len(vals)), ",")
	cond = fmt.Sprintf("%s NOT IN (%s)", quoteField(field), cond)
	return
}

// Between ...
type Between map[string][]interface{}

// Build ...
func (bt Between) Build() ([]string, []interface{}) {
	return betweenBuilder(bt, false)
}

func betweenBuilder(bt map[string][]interface{}, notBetween bool) ([]string, []interface{}) {
	if bt == nil || len(bt) == 0 {
		return nil, nil
	}
	var cond []string
	var vals []interface{}
	for k := range bt {
		cond = append(cond, k)
	}
	defaultSortAlgorithm(cond)
	for j := 0; j < len(cond); j++ {
		val := bt[cond[j]]
		condJ, err := buildBetween(notBetween, cond[j], val)
		if nil != err {
			continue
		}
		cond[j] = condJ
		vals = append(vals, val...)
	}
	return cond, vals
}

// NotBetween ...
type NotBetween map[string][]interface{}

// Build ...
func (nbt NotBetween) Build() ([]string, []interface{}) {
	return betweenBuilder(nbt, true)
}

func buildBetween(notBetween bool, key string, vals []interface{}) (string, error) {
	if len(vals) != 2 {
		return "", errors.New("vals of between must be a slice with two elements")
	}
	var operator string
	if notBetween {
		operator = "NOT BETWEEN"
	} else {
		operator = "BETWEEN"
	}
	return fmt.Sprintf("(%s %s ? AND ?)", key, operator), nil
}

type compareProducer func(m map[string]interface{}) (zsql.Comparable, error)

var op2Comparable = map[string]compareProducer{
	opEq: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Eq(m), nil
	},
	opNe1: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Ne(m), nil
	},
	opNe2: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Ne(m), nil
	},
	opIn: func(m map[string]interface{}) (zsql.Comparable, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return In(wp), nil
	},
	opNotIn: func(m map[string]interface{}) (zsql.Comparable, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return NotIn(wp), nil
	},
	opBetween: func(m map[string]interface{}) (zsql.Comparable, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return Between(wp), nil
	},
	opNotBetween: func(m map[string]interface{}) (zsql.Comparable, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return NotBetween(wp), nil
	},
	opGt: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Gt(m), nil
	},
	opGte: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Gte(m), nil
	},
	opLt: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Lt(m), nil
	},
	opLte: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Lte(m), nil
	},
	opLike: func(m map[string]interface{}) (zsql.Comparable, error) {
		return Like(m), nil
	},
	opNotLike: func(m map[string]interface{}) (zsql.Comparable, error) {
		return NotLike(m), nil
	},
	opNull: func(m map[string]interface{}) (zsql.Comparable, error) {
		return nullComparable(m), nil
	},
}

var OpOrder = []string{opEq, opIn, opNe1, opNe2, opNotIn, opGt, opGte, opLt, opLte, opLike, opNotLike, opBetween, opNotBetween, opNull}

type Statement func(field string, value interface{}) (zsql.ZSqlizer, error)

var OpOp = map[string]Statement{
	opEq: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.Eq{field: value}, nil
	},
	opNe1: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.NotEq{field: value}, nil
	},
	opNe2: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.NotEq{field: value}, nil
	},
	opIn: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			return squirrel.Eq{field: value}, nil
		}
		return squirrel.Eq{}, pkg.ErrNotASliceValueForInStatement
	},
	opNotIn: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			return squirrel.NotEq{field: value}, nil
		}
		return squirrel.NotEq{}, pkg.ErrNotASliceValueForInStatement
	},
	opBetween: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		if reflect.ValueOf(value).Kind() == reflect.Slice && reflect.ValueOf(value).Len() == 2 {
			return squirrel.Expr(fmt.Sprintf("%s BETWEEN ? AND ?", field),
				reflect.ValueOf(value).Index(0).Interface(),
				reflect.ValueOf(value).Index(1).Interface(),
			), nil
		}
		return squirrel.Eq{}, pkg.ErrNotASliceValueForBetweenStatement
	},
	opNotBetween: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		if reflect.ValueOf(value).Kind() == reflect.Slice && reflect.ValueOf(value).Len() == 2 {
			return squirrel.Expr(fmt.Sprintf("%s NOT BETWEEN ? AND ?", field),
				reflect.ValueOf(value).Index(0).Interface(),
				reflect.ValueOf(value).Index(1).Interface(),
			), nil
		}
		return squirrel.Eq{}, pkg.ErrNotASliceValueForBetweenStatement
	},
	opGt: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.Gt{field: value}, nil
	},
	opGte: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.GtOrEq{field: value}, nil
	},
	opLt: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.Lt{field: value}, nil
	},
	opLte: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.LtOrEq{field: value}, nil
	},
	opLike: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.Like{field: value}, nil
	},
	opNotLike: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		return squirrel.NotLike{field: value}, nil
	},
	opNull: func(field string, value interface{}) (zsql.ZSqlizer, error) {
		fmt.Println(field, value)
		return squirrel.Expr(fmt.Sprintf("%s %s", field, value.(NullType).String())), nil
	},
}
