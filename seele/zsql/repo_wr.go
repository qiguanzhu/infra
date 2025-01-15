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
 @Time    : 2024/11/11 -- 12:01
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: repo_wr.go
*/

package zsql

// ReadableRequest group by sql action. functions read from db
type ReadableRequest[EntityObj any] interface {
	QueryListRequest[EntityObj]
	QueryGetRequest[EntityObj]
	ComplexQueryRequest[EntityObj]
}

// WriteableRequest group by sql action. functions write to db
type WriteableRequest[EntityObj any] interface {
	InsertRequest[EntityObj]
	UpsertRequest[EntityObj]
	UpdateRequest
	DeleteRequest
	ComplexExecRequest
}
