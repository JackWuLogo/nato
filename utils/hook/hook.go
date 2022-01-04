package hook

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found hook")
	ErrExists   = errors.New("hook exists")
)

// Handler 钩子函数
type Handler func(data interface{}) (interface{}, error)

// 钩子
type Hook struct {
	sync.Mutex
	name    string  // 钩子名称
	count   int64   // 执行次数
	handler Handler // 钩子函数
}

// 钩子名称
func (h *Hook) Name() string {
	return h.name
}

// 执行次数
func (h *Hook) Count() int64 {
	return h.count
}

// 执行钩子
func (h *Hook) Call(obj interface{}) (interface{}, error) {
	h.Lock()
	defer h.Unlock()

	h.count++

	return h.handler(obj)
}

// NewHook 实例化钩子对象
func NewHook(name string, handler Handler) *Hook {
	h := &Hook{
		name:    name,
		handler: handler,
	}
	return h
}
