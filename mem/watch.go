package mem

import "sync"

// 模型数据变化监听函数
type WatchHandler func(field string, value interface{})

// 数据变化监听
type Watch struct {
	wg     *sync.WaitGroup
	model  *ModelLocal // 监听的数据模型
	change chan string // 字段数据变化
}

// 字段变化触发器
func (w *Watch) Trigger(field string) {
	if w.change == nil {
		return
	}
	w.change <- field
}

// 监听数据变化
func (w *Watch) Listen(fn WatchHandler) {
	if w.change != nil {
		w.Close()
	}
	w.change = make(chan string, 30)

	go func() {
		for {
			field, ok := <-w.change
			if !ok {
				break
			}

			w.wg.Add(1)
			go func() {
				defer w.wg.Done()
				fn(field, w.model.GetValue(field))
			}()
		}
	}()
}

// 关闭
func (w *Watch) Close() {
	if w.change == nil {
		return
	}

	close(w.change)
	w.wg.Wait()
	w.change = nil
}

func NewWatch(m *ModelLocal) *Watch {
	return &Watch{
		wg:    new(sync.WaitGroup),
		model: m,
	}
}
