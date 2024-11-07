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
 @Time    : 2024/11/4 -- 17:12
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: xsync.go
*/

package xsync

import (
	"sync"
)

type Mutex struct {
	once sync.Once
	mu   chan bool
}

func (m *Mutex) initLock() {
	m.once.Do(func() {
		m.mu = make(chan bool, 1)
	})
}

func (m *Mutex) Lock() {
	m.initLock()
	m.mu <- true
}

func (m *Mutex) Unlock() {
	select {
	case <-m.mu:
	default:
		panic("xsync: unlock of unlocked mutex")
	}
}

func (m *Mutex) Trylock() bool {
	m.initLock()
	select {
	case m.mu <- true:
		return true
	default:
		return false
	}
}
