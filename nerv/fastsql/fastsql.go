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
 @Time    : 2024/10/11 -- 18:31
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: fastsql.go
*/

package fastsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qiguanzhu/infra/seele/xsqlIface"
	"reflect"
)

var Constructor *constructor

type constructor struct {
	_scanner *fastScanner
	_builder *fastBuilder
}

func init() {
	Constructor = &constructor{
		_scanner: &fastScanner{},
		_builder: &fastBuilder{},
	}
}

func (c *constructor) GetBuilder() xsqlIface.BuilderProxy {
	return c._builder
}

func (c *constructor) GetScanner() xsqlIface.ScannerProxy {
	return c._scanner
}

func (c *constructor) ComplexSelect(dbx *sql.DB, builder xsqlIface.XSqlizer, target any, bind xsqlIface.Bind) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		query, args, err := builder.ToSql()
		if err != nil {
			return err
		}
		rows, err := dbx.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		if _, err := bind(rows); err != nil {
			return err
		}
		return err
	}
}

func (c *constructor) ComplexExec(dbx *sql.DB, builder xsqlIface.XSqlizer) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		query, args, err := builder.ToSql()
		if err != nil {
			return err
		}
		_, err = dbx.ExecContext(ctx, query, args...)
		return err
	}
}

// Mapping mapping what`s in src to tar
func Mapping(src, tar any) error {
	srcType := reflect.TypeOf(src)
	tarType := reflect.TypeOf(tar)
	if tarType.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("tar:%s must be a pointer to a struct", tarType.Kind()))
	}

	srcValue := reflect.ValueOf(src)
	tarValue := reflect.ValueOf(tar)
	// recalculate type based on Value`s type kind
	// if Value is ptr, we need to get it`s real value by calling Elem()
	if srcValue.Type().Kind() == reflect.Ptr {
		srcType = srcType.Elem()
		srcValue = srcValue.Elem()
	}
	if tarValue.Type().Kind() == reflect.Ptr {
		tarType = tarType.Elem()
		tarValue = tarValue.Elem()
	}

	if !match(srcType, tarType) {
		return errors.New(fmt.Sprintf("src:%s and tar:%s must have the same struct type", srcType, tarType))
	}

	tarValue.Set(srcValue)
	return nil
}

func match(srcT, tarT reflect.Type) bool {
	return srcT == tarT
}
