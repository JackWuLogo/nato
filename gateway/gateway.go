package gateway

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"micro-libs/codec"
	"micro-libs/utils/color"
	"micro-libs/utils/errors"
	"sync"
)

var (
	RegisterGateway = map[string]func(gateway *Gateway) Server{}
)

// 网关网络服务接口
type Server interface {
	Name() string      // 服务器名称
	Gateway() *Gateway // 关联服务端
	Opts() *Options    // 参数
	Port() int         // 监听端口
	Run() error        // 启动服务
	Close()            // 停止服务
}

// 客户端
type Client interface {
	Id() string                               // Client Id
	Server() Server                           // 关联服务端
	Meta() *Meta                              // 客户端上下文
	Log() *logger.Helper                      // 日志对象
	Closed() bool                             // 判断是否关闭
	Read() (*codec.ClientHead, []byte, error) // 读取消息
	Write([]byte)                             // 发送消息
	Close()                                   // 关闭连接
	Destroy()                                 // 销毁连接 (丢弃任何未发送或未确认的数据)
	SetAuthState(state bool)                  // 设置认证状态 (建立Socket连接后, 需要发送Token进行认证)
}

type Gateway struct {
	sync.RWMutex
	wg          *sync.WaitGroup
	opts        *Options
	server      Server
	clientCodec *codec.Client
	serverCodec *codec.Server
	clients     map[string]Client

	OnReceive    func(Client, *codec.ClientHead, []byte) (*codec.ServerHead, []byte, error) // 收到数据调用
	OnDisconnect func(Client)                                                               // 连接断开时调用
}

func (g *Gateway) Opts() *Options {
	return g.opts
}

func (g *Gateway) Server() Server {
	return g.server
}

// Address 网关地址
func (g *Gateway) Address() string {
	return fmt.Sprintf("%s:%d", Opts.Host, g.server.Port())
}

func (g *Gateway) ClientCodec() *codec.Client {
	return g.clientCodec
}

func (g *Gateway) ServerCodec() *codec.Server {
	return g.serverCodec
}

func (g *Gateway) SetOnReceive(fn func(Client, *codec.ClientHead, []byte) (*codec.ServerHead, []byte, error)) {
	g.OnReceive = fn
}

func (g *Gateway) SetOnDisconnect(fn func(Client)) {
	g.OnDisconnect = fn
}

func (g *Gateway) StartClient(client Client) {
	g.wg.Add(1)
	defer g.wg.Done()

	g.Lock()
	g.clients[client.Id()] = client
	g.Unlock()

	client.Log().Debugf("connected ...")

	// 接收消息处理
	for {
		// 接收消息
		cHead, cData, err := client.Read()
		if err != nil {
			client.Log().Warn(color.Warn.Text("read data error: %s", err))
			break
		}

		// 处理接收的消息
		sHead, sData, err := g.OnReceive(client, cHead, cData)
		if err != nil {
			break
		}

		// 响应消息
		if sHead.Code > 0 || len(sData) > 0 {
			b, err := g.ServerCodec().Marshal(sHead, sData)
			if err != nil {
				client.Log().Warn(color.Warn.Text("marshal data error: %s", err))
				break
			}

			client.Write(b)
		}
	}

	// 销毁客户端
	client.Close()

	g.Lock()
	delete(g.clients, client.Id())
	g.Unlock()

	// 连接断开处理
	g.OnDisconnect(client)

	client.Log().Debug("disconnected ...")
}

// 获取客户端连接对象
func (g *Gateway) GetClient(val string, by ...string) Client {
	g.RLock()
	defer g.RUnlock()

	if len(by) > 0 {
		byKey := by[0]
		for _, client := range g.clients {
			if client.Meta().Get(byKey) == val {
				return client
			}
		}
	} else {
		if client, ok := g.clients[val]; ok {
			return client
		}
	}

	return nil
}

// 连接数
func (g *Gateway) Count() int {
	g.RLock()
	defer g.RUnlock()

	return len(g.clients)
}

// 全部客户端
func (g *Gateway) All() map[string]Client {
	g.RLock()
	defer g.RUnlock()

	clients := make(map[string]Client, len(g.clients))
	for k, v := range g.clients {
		clients[k] = v
	}

	return clients
}

// 迭代客户端
func (g *Gateway) ForEach(fn func(client Client) bool) {
	g.RLock()
	defer g.RUnlock()

	for _, client := range g.clients {
		if !fn(client) {
			break
		}
	}
}

// 过滤角色对象
func (g *Gateway) Filter(filter func(client Client) bool) map[string]Client {
	g.RLock()
	defer g.RUnlock()

	var clients = make(map[string]Client)
	for id, client := range g.clients {
		if filter == nil || filter(client) {
			clients[id] = client
		}
	}

	return clients
}

// 常规指定客户端
func (g *Gateway) BroadcastNormal(b []byte, key string, values map[string]struct{}, online ...bool) {
	g.Broadcast(1, b, key, values, online...)
}

// 广播排除客户端
func (g *Gateway) BroadcastExclude(b []byte, key string, values map[string]struct{}, online ...bool) {
	g.Broadcast(2, b, key, values, online...)
}

// 广播所有
func (g *Gateway) BroadcastAll(b []byte, online ...bool) {
	g.Broadcast(0, b, "", nil, online...)
}

// 广播
func (g *Gateway) Broadcast(mode int, b []byte, key string, values map[string]struct{}, online ...bool) {
	if b == nil {
		return
	}

	var isOnline bool
	if len(online) > 0 && online[0] {
		isOnline = online[0]
	}

	g.RLock()
	defer g.RUnlock()

	// 广播
	for _, c := range g.clients {
		if isOnline && !c.Meta().IsOnline() {
			continue
		}

		switch mode {
		case 0:
			go c.Write(b)
		case 1:
			if _, ok := values[c.Meta().Get(key)]; ok {
				go c.Write(b)
			}
		case 2:
			if _, ok := values[c.Meta().Get(key)]; !ok {
				go c.Write(b)
			}
		}
	}
}

// 启动网关服务
func (g *Gateway) Run() error {
	g.Lock()
	defer g.Unlock()

	if g.server != nil {
		return errors.Server("gateway server is set %s ...", g.server.Name())
	}

	if Opts.Type == "ws" || Opts.Type == "wss" {
		Opts.Type = "websocket"
	}

	if s, ok := RegisterGateway[Opts.Type]; ok {
		g.server = s(g)
	} else {
		return errors.Server("Unsupported gateway server type: %s", Opts.Type)
	}

	return g.server.Run()
}

// 关闭网关服务
func (g *Gateway) Close() {
	if g.server != nil {
		g.server.Close()
	}

	g.Lock()
	for _, client := range g.clients {
		client.Close()
	}
	g.Unlock()

	g.wg.Wait()
}

// NewGateway
func NewGateway(codecMix []uint8, opts ...Option) *Gateway {
	g := &Gateway{
		wg:          new(sync.WaitGroup),
		opts:        NewOptions(opts...),
		clientCodec: codec.NewClient(codecMix...),
		serverCodec: codec.NewServer(codecMix...),
		clients:     make(map[string]Client),
	}
	return g
}
