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
 @Time    : 2024/10/11 -- 15:01
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: sqlutil.go
*/

package sqlutils

import (
	"context"
	"fmt"
	"github.com/qiguanzhu/infra/pkg"
	"github.com/qiguanzhu/infra/seele/xsqlIface"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var (
	defaultSortAlgorithm = sort.Strings
)

type insertType string

const (
	CommonInsert  insertType = "INSERT INTO"
	IgnoreInsert  insertType = "INSERT IGNORE INTO"
	ReplaceInsert insertType = "REPLACE INTO"
)

func build(m map[string]interface{}, op string) ([]string, []interface{}) {
	if nil == m || 0 == len(m) {
		return nil, nil
	}
	length := len(m)
	cond := make([]string, length)
	vals := make([]interface{}, length)
	var i int
	for key := range m {
		cond[i] = key
		i++
	}
	defaultSortAlgorithm(cond)
	for i = 0; i < length; i++ {
		vals[i] = m[cond[i]]
		cond[i] = assembleExpression(cond[i], op)
	}
	return cond, vals
}

func assembleExpression(field, op string) string {
	return quoteField(field) + op + "?"
}

// caller ensure that orderMap is not empty
func orderBy(orderMap []eleOrderBy) (string, error) {
	var str strings.Builder
	for _, orderInfo := range orderMap {
		realOrder := strings.ToUpper(orderInfo.order)
		if realOrder != "ASC" && realOrder != "DESC" {
			return "", pkg.ErrBuilderOrderByParam
		}
		str.WriteString(orderInfo.field)
		str.WriteByte(' ')
		str.WriteString(realOrder)
		str.WriteByte(',')
	}
	finalSQL := str.String()
	return finalSQL[:len(finalSQL)-1], nil
}

func resolveKV(m map[string]interface{}) (keys []string, vals []interface{}) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vals = append(vals, m[k])
	}
	return
}

func resolveFields(m map[string]interface{}) []string {
	var fields []string
	for k := range m {
		fields = append(fields, quoteField(k))
	}
	defaultSortAlgorithm(fields)
	return fields
}

func whereConnector(conditions ...xsqlIface.Comparable) (string, []interface{}) {
	if len(conditions) == 0 {
		return "", nil
	}
	var where []string
	var values []interface{}
	for _, cond := range conditions {
		cons, vals := cond.Build()
		if nil == cons {
			continue
		}
		where = append(where, cons...)
		values = append(values, vals...)
	}
	if 0 == len(where) {
		return "", nil
	}
	whereString := strings.Join(where, " AND ")
	return whereString, values
}

func quoteField(field string) string {
	return field
}

func BuildInsert(table string, setMap []map[string]interface{}, insertType insertType) (string, []interface{}, error) {
	format := "%s %s (%s) VALUES %s"
	var fields []string
	var vals []interface{}
	if len(setMap) < 1 {
		return "", nil, pkg.ErrBuilderInsertNullData
	}
	fields = resolveFields(setMap[0])
	placeholder := "(" + strings.TrimRight(strings.Repeat("?,", len(fields)), ",") + ")"
	var sets []string
	for _, mapItem := range setMap {
		sets = append(sets, placeholder)
		for _, field := range fields {
			val, ok := mapItem[strings.Trim(field, "`")]
			if !ok {
				return "", nil, pkg.ErrBuilderInsertDataNotMatch
			}
			vals = append(vals, val)
		}
	}
	return fmt.Sprintf(format, insertType, quoteField(table), strings.Join(fields, ","), strings.Join(sets, ",")), vals, nil
}

func BuildUpdate(table string, update map[string]interface{}, clauses *StatementClauses) (string, []interface{}, error) {
	format := "UPDATE %s SET %s"
	keys, vals := resolveKV(update)
	var sets string
	for _, k := range keys {
		sets += fmt.Sprintf("%s=?,", quoteField(k))
	}
	sets = strings.TrimRight(sets, ",")
	fieldStr := quoteField(table)
	if len(clauses.forceIndex) != 0 {
		fieldStr += fmt.Sprintf(" FORCE INDEX(%s)", clauses.forceIndex)
	}
	cond := fmt.Sprintf(format, fieldStr, sets)
	whereString, whereVals := whereConnector(clauses.conditions...)
	if "" != whereString {
		cond = fmt.Sprintf("%s WHERE %s", cond, whereString)
		vals = append(vals, whereVals...)
	}
	return cond, vals, nil
}

func BuildDelete(table string, conditions ...xsqlIface.Comparable) (string, []interface{}, error) {
	whereString, vals := whereConnector(conditions...)
	if "" == whereString {
		return fmt.Sprintf("DELETE FROM %s", table), nil, nil
	}
	format := "DELETE FROM %s WHERE %s"

	cond := fmt.Sprintf(format, quoteField(table), whereString)
	return cond, vals, nil
}

func splitCondition(conditions []xsqlIface.Comparable) ([]xsqlIface.Comparable, []xsqlIface.Comparable) {
	var having []xsqlIface.Comparable
	var i int
	for i = len(conditions) - 1; i >= 0; i-- {
		if _, ok := conditions[i].(nilComparable); ok {
			break
		}
	}
	if i >= 0 && i < len(conditions)-1 {
		having = conditions[i+1:]
		return conditions[:i], having
	}
	return conditions, nil
}

func BuildSelect(table string, ufields []string, clauses *StatementClauses) (string, []interface{}, error) {
	fields := "*"
	if len(ufields) > 0 {
		for i := range ufields {
			ufields[i] = quoteField(ufields[i])
		}
		fields = strings.Join(ufields, ",")
	}
	bd := strings.Builder{}
	bd.WriteString("SELECT ")
	bd.WriteString(fields)
	bd.WriteString(" FROM ")
	bd.WriteString(table)
	if len(clauses.forceIndex) != 0 {
		bd.WriteString(fmt.Sprintf(" FORCE INDEX(%s)", clauses.forceIndex))
	}
	where, having := splitCondition(clauses.conditions)
	whereString, vals := whereConnector(where...)
	if "" != whereString {
		bd.WriteString(" WHERE ")
		bd.WriteString(whereString)
	}
	if "" != clauses.groupBy {
		bd.WriteString(" GROUP BY ")
		bd.WriteString(clauses.groupBy)
	}
	if nil != having {
		havingString, havingVals := whereConnector(having...)
		bd.WriteString(" HAVING ")
		bd.WriteString(havingString)
		vals = append(vals, havingVals...)
	}
	if "" != clauses.orderBy {
		bd.WriteString(" ORDER BY ")
		bd.WriteString(clauses.orderBy)
	}
	if nil != clauses.limit {
		bd.WriteString(" LIMIT ?,?")
		vals = append(vals, clauses.limit.begin, clauses.limit.step)
	}
	return bd.String(), vals, nil
}

func BuildSelectWithContext(ctx context.Context, table string, ufields []string, clauses *StatementClauses) (string, []interface{}, error) {
	fields := "*"
	if len(ufields) > 0 {
		for i := range ufields {
			ufields[i] = quoteField(ufields[i])
		}
		fields = strings.Join(ufields, ",")
	}
	bd := strings.Builder{}
	bd.WriteString("SELECT ")
	bd.WriteString(fields)
	bd.WriteString(" FROM ")
	bd.WriteString(table)
	if len(clauses.forceIndex) != 0 {
		bd.WriteString(fmt.Sprintf(" FORCE INDEX(%s)", clauses.forceIndex))
	}
	where, having := splitCondition(clauses.conditions)
	whereString, vals := whereConnector(where...)
	if force, ok := ctx.Value(ContextKeyForceIndex).([]string); ok && len(force) > 0 && len(clauses.forceIndex) == 0 {
		bd.WriteString(" FORCE INDEX ")
		index := strings.Join(force, ",")
		indexStr := fmt.Sprintf("(%s)", index)
		bd.WriteString(indexStr)
	}
	if "" != whereString {
		bd.WriteString(" WHERE ")
		bd.WriteString(whereString)
	}
	if "" != clauses.groupBy {
		bd.WriteString(" GROUP BY ")
		bd.WriteString(clauses.groupBy)
	}
	if nil != having {
		havingString, havingVals := whereConnector(having...)
		bd.WriteString(" HAVING ")
		bd.WriteString(havingString)
		vals = append(vals, havingVals...)
	}
	if "" != clauses.orderBy {
		bd.WriteString(" ORDER BY ")
		bd.WriteString(clauses.orderBy)
	}
	if nil != clauses.limit {
		bd.WriteString(" LIMIT ?,?")
		vals = append(vals, int(clauses.limit.begin), int(clauses.limit.step))
	}
	return bd.String(), vals, nil
}

func CopyWhere(src map[string]interface{}) (target map[string]interface{}) {
	target = make(map[string]interface{})
	for k, v := range src {
		target[k] = v
	}
	return
}

func ParseWhere(where map[string]interface{}) (statement *StatementClauses, err error) {
	var release func()
	var having map[string]interface{}
	statement = new(StatementClauses)
	copiedWhere := CopyWhere(where)
	if val, ok := copiedWhere["_orderby"]; ok {
		s, ok := val.(string)
		if !ok {
			err = pkg.ErrBuilderOrderByValueType
			return
		}
		statement.orderBy = strings.TrimSpace(s)
		delete(copiedWhere, "_orderby")
	}
	if val, ok := where["_forceindex"]; ok {
		switch result := val.(type) {
		case []string:
			statement.forceIndex = strings.Join(result, ",")
		case string:
			statement.forceIndex = result
		default:
			err = pkg.ErrBuilderForceIndexType
			return
		}
		delete(copiedWhere, "_forceindex")
	}
	if val, ok := copiedWhere["_groupby"]; ok {
		s, ok := val.(string)
		if !ok {
			err = pkg.ErrBuilderGroupByValueType
			return
		}
		statement.groupBy = s
		delete(copiedWhere, "_groupby")
		if h, ok := copiedWhere["_having"]; ok {
			having, err = resolveHaving(h)
			if nil != err {
				return
			}
		}
	}
	if _, ok := copiedWhere["_having"]; ok {
		delete(copiedWhere, "_having")
	}
	if val, ok := copiedWhere["_limit"]; ok {
		arr := make([]int, 2)
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.Slice {
			err = pkg.ErrBuilderLimitValueType
			return
		}
		if v.Len() != 2 && v.Len() != 1 {
			err = pkg.ErrBuilderLimitValueLength
			return
		}
		arr, err = arrayInterfaceToArrayInt(val)
		if err != nil {
			return
		}
		begin, step := arr[0], arr[1]
		statement.limit = &eleLimit{
			begin: begin,
			step:  step,
		}
		delete(copiedWhere, "_limit")
	}
	conditions, release, err := GetWhereConditions(copiedWhere)
	if nil != err {
		return
	}
	defer release()
	if having != nil {
		havingCondition, release1, err1 := GetWhereConditions(having)
		if nil != err1 {
			err = err1
			return
		}
		defer release1()
		conditions = append(conditions, nilComparable(0))
		conditions = append(conditions, havingCondition...)
	}
	statement.conditions = conditions
	return
}

func ParseDMLWhere(where map[string]interface{}) (clauses *StatementClauses, release func(), err error) {
	clauses = new(StatementClauses)
	if val, ok := where["_forceindex"]; ok {
		switch result := val.(type) {
		case []string:
			clauses.forceIndex = strings.Join(result, ",")
		case string:
			clauses.forceIndex = result
		default:
			err = pkg.ErrBuilderForceIndexType
			return
		}
		delete(where, "_forceindex")
	}
	clauses.conditions, release, err = GetWhereConditions(where)
	return
}

func resolveHaving(having interface{}) (map[string]interface{}, error) {
	var havingMap map[string]interface{}
	var ok bool
	if havingMap, ok = having.(map[string]interface{}); !ok {
		return nil, pkg.ErrBuilderHavingValueType
	}
	copiedMap := make(map[string]interface{})
	for key, val := range havingMap {
		_, operator, err := SplitKey(key)
		if nil != err {
			return nil, err
		}
		if !IsStringInSlice(operator, OpOrder) {
			return nil, pkg.ErrBuilderHavingUnsupportedOperator
		}
		copiedMap[key] = val
	}
	return copiedMap, nil
}

var (
	cpPool = sync.Pool{
		New: func() interface{} {
			return make([]xsqlIface.Comparable, 0)
		},
	}
)

func getCpPool() ([]xsqlIface.Comparable, func()) {
	obj := cpPool.Get().([]xsqlIface.Comparable)
	return obj[:0], func() { cpPool.Put(obj) }
}

func emptyFunc() {}

func IsStringInSlice(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func GetWhereConditions(where map[string]interface{}) ([]xsqlIface.Comparable, func(), error) {
	if len(where) == 0 {
		return nil, emptyFunc, nil
	}
	wms := &whereMapSet{}
	var field, operator string
	var err error
	for key, val := range where {
		field, operator, err = SplitKey(key)
		if !IsStringInSlice(operator, OpOrder) {
			return nil, emptyFunc, pkg.ErrBuilderUnsupportedOperator
		}
		if nil != err {
			return nil, emptyFunc, err
		}
		if _, ok := val.(NullType); ok {
			operator = opNull
		}
		wms.add(operator, field, val)
	}

	return buildWhereCondition(wms)
}

const (
	opEq         = "="
	opNe1        = "!="
	opNe2        = "<>"
	opIn         = "in"
	opNotIn      = "not in"
	opGt         = ">"
	opGte        = ">="
	opLt         = "<"
	opLte        = "<="
	opLike       = "like"
	opNotLike    = "not like"
	opBetween    = "between"
	opNotBetween = "not between"
	// special
	opNull = "null"
)

func buildWhereCondition(mapSet *whereMapSet) ([]xsqlIface.Comparable, func(), error) {
	cpArr, release := getCpPool()
	for _, operator := range OpOrder {
		whereMap, ok := mapSet.set[operator]
		if !ok {
			continue
		}
		f, ok := op2Comparable[operator]
		if !ok {
			release()
			return nil, emptyFunc, pkg.ErrBuilderUnsupportedOperator
		}
		cp, err := f(whereMap)
		if nil != err {
			release()
			return nil, emptyFunc, err
		}
		cpArr = append(cpArr, cp)
	}
	return cpArr, release, nil
}

func convertWhereMapToWhereMapSlice(where map[string]interface{}) (map[string][]interface{}, error) {
	result := make(map[string][]interface{})
	for key, val := range where {
		vals, ok := convertInterfaceToMap(val)
		if !ok {
			return nil, pkg.ErrBuilderWhereInType
		}
		if 0 == len(vals) {
			return nil, pkg.ErrBuilderEmptyINCondition
		}
		result[key] = vals
	}
	return result, nil
}

func convertInterfaceToMap(val interface{}) ([]interface{}, bool) {
	s := reflect.ValueOf(val)
	if s.Kind() != reflect.Slice {
		return nil, false
	}
	interfaceSlice := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		interfaceSlice[i] = s.Index(i).Interface()
	}
	return interfaceSlice, true
}

func SplitKey(key string) (field string, operator string, err error) {
	key = strings.Trim(key, " ")
	if "" == key {
		err = pkg.ErrBuilderSplitEmptyKey
		return
	}
	idx := strings.IndexByte(key, ' ')
	if idx == -1 {
		field = key
		operator = "="
	} else {
		field = key[:idx]
		operator = strings.Trim(key[idx+1:], " ")
		operator = removeInnerSpace(operator)
	}
	return
}

func removeInnerSpace(operator string) string {
	n := len(operator)
	firstSpace := strings.IndexByte(operator, ' ')
	if firstSpace == -1 {
		return operator
	}
	lastSpace := firstSpace
	for i := firstSpace + 1; i < n; i++ {
		if operator[i] == ' ' {
			lastSpace = i
		} else {
			break
		}
	}
	return operator[:firstSpace] + operator[lastSpace:]
}

func splitOrderBy(orderby string) ([]eleOrderBy, error) {
	var err error
	var eleOrder []eleOrderBy
	for _, val := range strings.Split(orderby, ",") {
		val = strings.Trim(val, " ")
		idx := strings.IndexByte(val, ' ')
		if idx == -1 {
			err = pkg.ErrBuilderSplitOrderBy
			return eleOrder, err
		}
		field := val[:idx]
		direction := strings.Trim(val[idx+1:], " ")
		eleOrder = append(eleOrder, eleOrderBy{
			field: field,
			order: direction,
		})
	}
	return eleOrder, err
}

const (
	paramPlaceHolder = "?"
)

var searchHandle = regexp.MustCompile(`{{\S+?}}`)

// NamedQuery is used for expressing complex query
func NamedQuery(sql string, data map[string]interface{}) (string, []interface{}, error) {
	length := len(data)
	if length == 0 {
		return sql, nil, nil
	}
	vals := make([]interface{}, 0, length)
	var err error
	cond := searchHandle.ReplaceAllStringFunc(sql, func(paramName string) string {
		paramName = strings.TrimRight(strings.TrimLeft(paramName, "{"), "}")
		val, ok := data[paramName]
		if !ok {
			err = fmt.Errorf("%s not found", paramName)
			return ""
		}
		v := reflect.ValueOf(val)
		if v.Type().Kind() != reflect.Slice {
			vals = append(vals, val)
			return paramPlaceHolder
		}
		length := v.Len()
		for i := 0; i < length; i++ {
			vals = append(vals, v.Index(i).Interface())
		}
		return createMultiPlaceholders(length)
	})
	if nil != err {
		return "", nil, err
	}
	return cond, vals, nil
}

func createMultiPlaceholders(num int) string {
	if 0 == num {
		return ""
	}
	length := (num << 1) | 1
	buff := make([]byte, length)
	buff[0], buff[length-1] = '(', ')'
	ll := length - 2
	for i := 1; i <= ll; i += 2 {
		buff[i] = '?'
	}
	ll = length - 3
	for i := 2; i <= ll; i += 2 {
		buff[i] = ','
	}
	return string(buff)
}

func arrayInterfaceToArrayInt(val interface{}) (arr []int, err error) {
	arr = make([]int, 2)
	switch vals := val.(type) {
	case []int:
		if len(vals) == 1 {
			arr[1] = vals[0]
		} else {
			arr = vals
		}
	case []uint:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []int64:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []int32:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []int16:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []int8:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []uint64:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []uint32:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []uint16:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	case []uint8:
		if len(vals) == 1 {
			arr[1] = int(vals[0])
		} else {
			arr[0] = int(vals[0])
			arr[1] = int(vals[1])
		}
	default:
		err = pkg.ErrBuilderLimitValueType
		return
	}
	return
}
