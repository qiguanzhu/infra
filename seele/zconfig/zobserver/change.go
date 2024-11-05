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
 @Time    : 2024/9/30 -- 12:00
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: config.go
*/

package zobserver

// ChangeType ...
type ChangeType int

// ChangeEventSource ...
type ChangeEventSource int

const (
	// ADD change type:add
	ADD ChangeType = iota
	// MODIFY change type:modify
	MODIFY
	// DELETE change type:delete
	DELETE
)

const (
	// Apollo change event source:apollo
	Apollo ChangeEventSource = iota
)

// String ...
func (c ChangeType) String() string {
	switch c {
	case ADD:
		return "ADD"
	case MODIFY:
		return "MODIFY"
	case DELETE:
		return "DELETE"
	}

	return "UNKOWN"
}

// ChangeEvent ...
type ChangeEvent struct {
	Source    ChangeEventSource
	Namespace string
	Changes   map[string]*Change
}

// Change ...
type Change struct {
	OldValue   string
	NewValue   string
	ChangeType ChangeType
}
