package event

import (
	"micro-libs/utils/errors"
	"micro-libs/utils/log"
	"sync"
)

// 事件管理器 (异步执行)
type Events struct {
	sync.RWMutex
	name   string
	events map[string]*Event // 事件列表
}

func (e *Events) String() string {
	return e.name
}

func (e *Events) Exist(name string) bool {
	e.RLock()
	defer e.RUnlock()
	_, ok := e.events[name]
	return ok
}

// Bind 绑定事件
func (e *Events) Bind(name string, handler Handler, maxCount ...int64) error {
	if e.Exist(name) {
		return errors.Exists("%s is exist", name)
	}

	e.Lock()
	e.events[name] = NewEvent(name, handler, maxCount...)
	e.Unlock()

	return nil
}

// Bind 绑定事件 (仅执行一次)
func (e *Events) Once(name string, handler Handler) error {
	if e.Exist(name) {
		return errors.Exists("%s is exist", name)
	}

	e.Lock()
	e.events[name] = NewEvent(name, handler, 1)
	e.Unlock()

	return nil
}

// Delete 删除事件
func (e *Events) Delete(name string) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.events[name]; ok {
		delete(e.events, name)
	}
}

// Emit 执行事件
func (e *Events) Emit(name string, data interface{}) error {
	e.RLock()
	h, ok := e.events[name]
	e.RUnlock()

	if !ok {
		return errors.NotFound("没有找到指定事件")
	}

	if !h.IsAllowExec() {
		e.Delete(name)
		return errors.Finished("已达最大执行次数")
	}

	return h.Call(data)
}

// EmitAsync 异步执行事件
func (e *Events) EmitAsync(name string, data interface{}) {
	go func() {
		if err := e.Emit(name, data); err != nil {
			log.Error("[AsyncEvent] [%s] run error: %s", name, err)
		}
	}()
}

func NewEvents(name string) *Events {
	return &Events{
		name:   name,
		events: make(map[string]*Event),
	}
}
