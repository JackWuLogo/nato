package hook

import "sync"

// 钩子 (同步执行)
type Hooks struct {
	sync.RWMutex
	name  string
	hooks map[string]*Hook // 事件列表
}

func (e *Hooks) String() string {
	return e.name
}

func (e *Hooks) Exist(name string) bool {
	e.RLock()
	defer e.Unlock()
	_, ok := e.hooks[name]
	return ok
}

// Add 添加钩子
func (e *Hooks) Add(name string, handler Handler) error {
	if e.Exist(name) {
		return ErrExists
	}

	e.Lock()
	e.hooks[name] = NewHook(name, handler)
	e.Unlock()

	return nil
}

// Del 删除钩子
func (e *Hooks) Del(name string) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.hooks[name]; ok {
		delete(e.hooks, name)
	}
}

// Call 执行钩子
func (e *Hooks) Call(name string, data interface{}) (interface{}, error) {
	e.RLock()
	h, ok := e.hooks[name]
	e.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	return h.Call(data)
}

func NewHooks(name string) *Hooks {
	return &Hooks{
		name:  name,
		hooks: make(map[string]*Hook),
	}
}
