package app

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/cmd"
	"github.com/asim/go-micro/v3/transport"
	"github.com/google/gops/agent"
	"micro-libs/utils/errors"
	"micro-libs/utils/tool"
	"path/filepath"
	"sync"

	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/micro/cli/v2"

	cgrpc "micro-libs/app/plugins/client/grpc"
	sgrpc "micro-libs/app/plugins/server/grpc"
	tgrpc "micro-libs/app/plugins/transport/grpc"

	_ "micro-libs/app/plugins/broker/nats"
	_ "micro-libs/app/plugins/broker/nsq"
	_ "micro-libs/app/plugins/broker/rabbitmq"

	_ "micro-libs/app/plugins/registry/etcd"
	_ "micro-libs/app/plugins/registry/kubernetes"

	"micro-libs/compile"
	"micro-libs/utils/log"
	"micro-libs/utils/pb"
)

var (
	appInstance *app
	srvOnce     sync.Once
)

// 服务单例
func APP() *app {
	srvOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		appInstance = &app{
			cancel:    cancel,
			ctx:       ctx,
			publisher: make(map[string]micro.Publisher),
		}
	})
	return appInstance
}

// 当前服务
type app struct {
	sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	cmd       cmd.Cmd
	srv       micro.Service              // micro服务
	publisher map[string]micro.Publisher // 订阅
}

func (a *app) Context() context.Context {
	return a.ctx
}

func (a *app) Cancel() {
	a.cancel()
}

func (a *app) Cmd() cmd.Cmd {
	return a.cmd
}

func (a *app) Service() micro.Service {
	return a.srv
}

func (a *app) AddCommands(commands ...*cli.Command) {
	a.cmd.App().Commands = append(a.cmd.App().Commands, commands...)
}

func (a *app) newService(srvName string, flags ...[]cli.Flag) {
	// 合并flags
	if len(flags) > 0 {
		for _, fg := range flags {
			Flags = append(Flags, fg...)
		}
	}

	// 创角命令行
	a.cmd = cmd.NewCmd(
		cmd.Name(srvName),
		cmd.Description(fmt.Sprintf("a %s service", srvName)),
		cmd.Version(compile.Version()),
	)

	// 重载Cmd参数, 省略用不上的参数
	a.cmd.App().Flags = Flags

	a.AddCommands(compile.Cmd...)

	var opts = []micro.Option{
		micro.Context(a.ctx),
		micro.Cmd(a.cmd),
		micro.Client(cgrpc.NewClient()),
		micro.Server(sgrpc.NewServer(
			server.Name(srvName),
			server.Version(compile.Version()),
		)),
		micro.Transport(tgrpc.NewTransport(transport.Secure(true))),
		micro.WrapHandler(serverWrapper),
		micro.WrapSubscriber(subscriberWrapper),
		micro.BeforeStart(func() error {
			// 启动时输出版本信息
			compile.EchoVersion(a.srv)

			// 检查集群ID取值范围
			if Opts.ClusterId < 1000 || Opts.ClusterId > 9999 {
				return errors.Invalid("[ClusterId] 的取值范围为: 1000 < ClusterId < 9999")
			}

			// 检查目录是否存在
			if err := tool.InitFolder(Opts.StoreRoot, 0755); err != nil {
				return err
			}

			if Opts.Dev {
				log.Warn("the current startup mode is development ...")
			}
			if Opts.PsAddr != "" || Opts.Dev {
				if err := agent.Listen(agent.Options{Addr: Opts.PsAddr, ConfigDir: filepath.Join(Opts.StoreRoot, "gops")}); err != nil {
					return err
				}
			}

			return nil
		}),
	}

	a.srv = micro.NewService(opts...)
}

func (a *app) initService(opts ...micro.Option) {
	a.srv.Init(opts...)
}

// New 创建服务
func New(srvName string, flags ...[]cli.Flag) {
	if APP().srv != nil {
		return
	}

	// 设置当前服务名称
	compile.SetName(srvName)

	// 创建服务
	APP().newService(srvName, flags...)

	// 关闭事件
	_ = AddSub("cancel", func(_ context.Context, c *pb.Cancel) error {
		if c.Name != Name() || (c.NodeId != "" && c.NodeId != Id()) {
			return nil
		}

		log.Debug("[ServiceCancel] service [%s][%s] is being stopped ...", Name(), Id())

		// 停止服务
		APP().Cancel()

		return nil
	})
}

// 启动服务
func Init(opts ...micro.Option) {
	APP().initService(opts...)
}

// AddCmd 添加CMD
func AddCmd(commands ...*cli.Command) {
	APP().AddCommands(commands...)
}

// Cmd 获取命令行对象
func Cmd() cmd.Cmd {
	return APP().cmd
}

// Service 获取服务对象
func Srv() micro.Service {
	return APP().srv
}

// Client 获取客户端
func Client() client.Client {
	return Srv().Client()
}

// Server 获取服务端
func Server() server.Server {
	return Srv().Server()
}

// 启动服务
func Run() error {
	return Srv().Run()
}

// 关闭服务
func Close() {
	APP().Cancel()
}

// 取消服务
func Cancel(name string, id ...string) error {
	in := &pb.Cancel{Name: name}
	if len(id) > 0 {
		in.NodeId = id[0]
	}
	return Pub("cancel", in)
}

// Version 获取服务版本
func Version() string {
	return Server().Options().Version
}

// Id 获取服务UUID
func Id() string {
	return Server().Options().Id
}

// Name 获取服务名称
func Name() string {
	return Server().Options().Name
}

// NameId 获取服务节点ID
func NameId() string {
	return fmt.Sprintf("%s-%s", Name(), Id())
}

// 获取服务发现注册信息
func Registry() registry.Registry {
	return APP().Service().Options().Registry
}

// 获取服务节点列表
func GetServices(name string) []*registry.Service {
	res, _ := Registry().GetService(name)
	return res
}

// 注册发布者
func AddPub(names ...string) {
	for _, name := range names {
		APP().publisher[name] = micro.NewEvent(name, Client())
	}
}

// 获取发布者
func GetPub(name string) micro.Publisher {
	return APP().publisher[name]
}

// 发布消息
func Pub(name string, msg interface{}, opts ...client.PublishOption) error {
	return PubCtx(context.TODO(), name, msg, opts...)
}
func PubCtx(ctx context.Context, name string, msg interface{}, opts ...client.PublishOption) error {
	return GetPub(name).Publish(ctx, msg, opts...)
}

// 注册订阅
func AddSub(name string, h interface{}, queue ...bool) error {
	var opts = []server.SubscriberOption{
		server.InternalSubscriber(true),
	}
	if len(queue) > 0 && queue[0] {
		opts = append(opts, server.SubscriberQueue(name))
	}

	return micro.RegisterSubscriber(name, Server(), h, opts...)
}

func AddSubQueue(name string, h interface{}) error {
	return AddSub(name, h, true)
}

// 注册RPC服务
func AddHandler(handles ...interface{}) {
	opts := []server.HandlerOption{
		server.InternalHandler(true),
	}

	for _, h := range handles {
		if err := micro.RegisterHandler(Server(), h, opts...); err != nil {
			log.Error("RegisterHandler Error: %s", err.Error())
		}
	}
}

// 增加meta信息
func AddMetadata(meta map[string]string) {
	for k, v := range meta {
		Server().Options().Metadata[k] = v
	}
}

// GenClusterName 生成集群专属名称
func GenClusterName(name string) string {
	return fmt.Sprintf("%d_%s", Opts.ClusterId, name)
}

// Call 通过名称调用RPC
// 	@name 服务名称
// 	@method rpc方法名称. 即 serviceName.rpcName
// 	@in 请求参数
// 	@out 返回数据
// 	@nodeId 指定节点ID
func CallCtx(ctx context.Context, srvName string, method string, in interface{}, out interface{}, opts ...client.CallOption) error {
	req := Client().NewRequest(srvName, method, in)
	return Client().Call(ctx, req, out, opts...)
}

// CallNode 请求指定节点RPC
// 	@name 服务名称
// 	@method rpc方法名称. 即 serviceName.rpcName
// 	@in 请求参数
// 	@out 返回数据
// 	@nodeId 指定节点ID
func CallNode(ctx context.Context, srvName string, method string, in interface{}, out interface{}, nodeId ...string) error {
	var opts []client.CallOption
	if len(nodeId) > 0 {
		opts = append(opts, FilterServiceNode(srvName, nodeId[0]))
	}
	return CallCtx(ctx, srvName, method, in, out, opts...)
}
