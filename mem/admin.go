package mem

import (
	"micro-libs/store/mongo"
	"micro-libs/utils/debug"
	"micro-libs/utils/log"
	"sync"
	"time"
)

// Admin 内存管理器
type Admin struct {
	sync.RWMutex
	wg     sync.WaitGroup
	dbName string            // 数据库名称
	opts   *Options          // 参数
	client *RemoteClient     // 内存数据远程客户端
	ticker *time.Ticker      // 数据状态检查
	tables map[string]*Table // 数据表集合
}

func (a *Admin) SrvName() string {
	return a.client.name
}

func (a *Admin) DbName() string {
	return a.dbName
}

func (a *Admin) Client() *RemoteClient {
	return a.client
}

func (a *Admin) Options() *Options {
	return a.opts
}

// Init 初始化
func (a *Admin) Init(tables map[string]*mongo.Table) {
	a.Lock()
	defer a.Unlock()

	for nm, tab := range tables {
		a.tables[nm] = NewTable(a, tab)
	}
}

// 开始同步任务
func (a *Admin) SyncStart() {
	a.Lock()
	defer a.Unlock()

	if a.ticker != nil {
		return
	}

	a.ticker = time.NewTicker(a.opts.StateCheckInterval)

	go func() {
		for range a.ticker.C {
			a.SyncAll()
		}
	}()
}

// 结束同步任务
func (a *Admin) SyncStop() {
	a.Lock()
	defer a.Unlock()

	if a.ticker == nil {
		return
	}

	a.ticker.Stop()
	a.ticker = nil
}

// MTable 获取内存模型表
func (a *Admin) Table(name string) (*Table, error) {
	if mt, ok := a.tables[name]; ok {
		return mt, nil
	}
	return nil, ErrTableNotFound
}

// SyncAll 同步所有数据
func (a *Admin) SyncAll() {
	a.RLock()
	tables := make(map[string]*Table, len(a.tables))
	for k, v := range a.tables {
		tables[k] = v
	}
	a.RUnlock()

	var prof *debug.Prof
	if log.IsTrace() {
		prof = debug.NewProf(log.Logger, "[MemSync]")
	}
	log.Debug("[MemSync] check data state starting ...")

	// 同步所有数据
	for _, table := range tables {
		table := table
		a.wg.Add(1)

		go func() {
			defer a.wg.Done()
			table.SyncAll()
		}()
	}

	a.wg.Wait()

	if prof != nil {
		prof.Result()
	}
}

func NewAdmin(srvName string, dbName string, opts ...Option) *Admin {
	ma := &Admin{
		dbName: dbName,
		opts:   newOptions(opts...),
		client: NewRemoteClient(srvName),
		tables: make(map[string]*Table),
	}
	return ma
}
