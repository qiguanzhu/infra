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
 @Time    : 2024/9/30 -- 15:14
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: error.go
*/

package pkg

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrBuilderSplitEmptyKey = errors.New("[builder] couldn't split a empty string")
	ErrBuilderSplitOrderBy  = errors.New(`[builder] the value of _orderby should be "fieldName direction [,fieldName direction]"`)
	// ErrBuilderUnsupportedOperator reports there's unsupported operators in where-condition
	ErrBuilderUnsupportedOperator       = errors.New("[builder]: unsupported operator")
	ErrBuilderWhereInType               = errors.New(`[builder]: the value of "xxx in" must be of []interface{} type`)
	ErrBuilderOrderByValueType          = errors.New(`[builder]: the value of "_orderby" must be of string type`)
	ErrBuilderGroupByValueType          = errors.New(`[builder]: the value of "_groupby" must be of string type`)
	ErrBuilderLimitValueType            = errors.New(`[builder]: the value of "_limit" must be of []uint or []int type`)
	ErrBuilderLimitValueLength          = errors.New(`[builder]: the value of "_limit" must contain one or two uint elements`)
	ErrBuilderEmptyINCondition          = errors.New(`[builder]: the value of "in" must contain at least one element`)
	ErrBuilderHavingValueType           = errors.New(`[builder]: the value of "_having" must be of map[string]interface{}`)
	ErrBuilderHavingUnsupportedOperator = errors.New(`[builder]: "_having" contains unsupported operator`)
	ErrBuilderForceIndexType            = errors.New(`[builder]: the value of "_forceindex" must be of []string type`)

	ErrBuilderInsertDataNotMatch = errors.New("[builder]: [dao] insert data not match")
	ErrBuilderInsertNullData     = errors.New("[builder]: [dao] insert null data")
	ErrBuilderOrderByParam       = errors.New("[builder]: [dao] order param only should be ASC or DESC")
)

var (
	// ErrScannerTargetNotSettable means the second param of Bind is not settable
	ErrScannerTargetNotSettable = errors.New("[scanner]: target is not settable! a pointer is required")
	// ErrScannerNilRows means the first param can't be a nil
	ErrScannerNilRows = errors.New("[scanner]: rows can't be nil")
	// ErrScannerSliceToString means only []uint8 can be transmuted into string
	ErrScannerSliceToString = errors.New("[scanner]: can't transmute a non-uint8 slice to string")
	// ErrScannerEmptyResult occurs when target of Scan isn't slice and the result of the query is empty
	ErrScannerEmptyResult     = errors.New(`[scanner]: empty result`)
	ErrScannerFromNilFunction = errors.New("[scanner]: function can't be nil")
)

var (
	ErrRepoUpdateNothingChanged = errors.New("[repo]: nothing changed")
)

var (
	ErrNotASliceValueForBetweenStatement = errors.New(`[sqlutils]: the value of "between" must be of []interface{} and have only 2 values`)
	ErrNotASliceValueForInStatement      = errors.New(`[sqlutils]: the value of "in" must be of []interface{}`)
)

// ScanErr will be returned if an underlying type couldn't be AssignableTo type of target field
type ScanErr struct {
	structName, fieldName string
	from, to              reflect.Type
}

func (s ScanErr) Error() string {
	return fmt.Sprintf("[scanner]: %s.%s is %s which is not AssignableBy %s", s.structName, s.fieldName, s.to, s.from)
}

func NewScanErr(structName, fieldName string, from, to reflect.Type) ScanErr {
	return ScanErr{structName, fieldName, from, to}
}

// CloseErr is the error occurs when rows.Close()
type CloseErr struct {
	err error
}

func (e CloseErr) Error() string {
	return e.err.Error()
}

func NewCloseErr(err error) error {
	if err == nil {
		return nil
	}
	return CloseErr{err}
}

var (
	ErrNoneStructTarget = errors.New("target must be a struct type")
)
