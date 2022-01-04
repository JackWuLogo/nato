package app

import (
	"github.com/micro/cli/v2"
	"os"
)

var (
	Opts = new(struct {
		Dev       bool   // 开发模式
		PsAddr    string // gops分析
		ClusterId int64  // 集群ID(取值范围100-999)
		StoreRoot string // 数据保存根路径
		Lang      string // 设置语言环境
	})

	Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "dev",
			Usage:       "设置是否为开发模式",
			EnvVars:     []string{"GAME_DEV"},
			Destination: &Opts.Dev,
		},
		&cli.StringFlag{
			Name:        "gops",
			Usage:       "设置gops分析端口",
			EnvVars:     []string{"GAME_GOPS"},
			Destination: &Opts.PsAddr,
		},
		&cli.StringFlag{
			Name:        "lang",
			Value:       "zh-cn",
			Usage:       "设置当前语言环境.",
			EnvVars:     []string{"GAME_LANG"},
			Destination: &Opts.Lang,
		},
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Debug profiler for cpu and memory stats",
			EnvVars: []string{"GAME_DEBUG_PROFILE"},
		},
		&cli.StringFlag{
			Name:    "client_request_timeout",
			EnvVars: []string{"GAME_CLIENT_REQUEST_TIMEOUT"},
			Usage:   "设置客户端请求超时. e.g 500ms, 5s, 1m. Default: 5s",
		},
		&cli.IntFlag{
			Name:    "client_pool_size",
			EnvVars: []string{"GAME_CLIENT_POOL_SIZE"},
			Usage:   "设置客户端连接池大小. Default: 1",
		},
		&cli.StringFlag{
			Name:    "client_pool_ttl",
			EnvVars: []string{"GAME_CLIENT_POOL_TTL"},
			Usage:   "设置客户端连接TTL. e.g 500ms, 5s, 1m. Default: 1m",
		},
		&cli.StringFlag{
			Name:    "registry",
			EnvVars: []string{"GAME_REGISTRY"},
			Value:   "etcd",
			Usage:   "服务发现类型. etcd, consul, kubernetes",
		},
		&cli.StringFlag{
			Name:    "registry_address",
			EnvVars: []string{"GAME_REGISTRY_ADDRESS"},
			Usage:   "服务发现地址. 以逗号分隔",
		},
		&cli.IntFlag{
			Name:    "register_ttl",
			EnvVars: []string{"GAME_REGISTER_TTL"},
			Value:   60,
			Usage:   "服务发现注册TTL",
		},
		&cli.IntFlag{
			Name:    "register_interval",
			EnvVars: []string{"GAME_REGISTER_INTERVAL"},
			Value:   30,
			Usage:   "服务发现注册时间间隔",
		},
		&cli.StringFlag{
			Name:    "broker",
			EnvVars: []string{"GAME_BROKER"},
			Value:   "nats",
			Usage:   "消息队列类型. nats, nsq, rabbitmq",
		},
		&cli.StringFlag{
			Name:    "broker_address",
			EnvVars: []string{"GAME_BROKER_ADDRESS"},
			Usage:   "消息队列地址. 以逗号分隔",
		},
		&cli.StringFlag{
			Name:    "selector",
			EnvVars: []string{"GAME_SELECTOR"},
			Usage:   "设置用于选择节点进行查询的选择器",
		},
		&cli.StringFlag{
			Name:    "tracer",
			EnvVars: []string{"GAME_TRACER"},
			Usage:   "设置分布式跟踪的跟踪器, e.g. memory, jaeger",
		},
		&cli.StringFlag{
			Name:    "tracer_address",
			EnvVars: []string{"GAME_TRACER_ADDRESS"},
			Usage:   "设置分布式跟踪的跟踪器地址. 以逗号分隔",
		},
		&cli.Int64Flag{
			Name:        "cluster_id",
			Usage:       "设置当前集群ID. 取值范围: 1000 ~ 9999",
			EnvVars:     []string{"GAME_CLUSTER_ID"},
			Destination: &Opts.ClusterId,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "store_root",
			Value:       os.TempDir(),
			Usage:       "设置数据保存磁盘路径, 默认临时目录",
			EnvVars:     []string{"GAME_STORE_ROOT"},
			Destination: &Opts.StoreRoot,
		},
	}
)
