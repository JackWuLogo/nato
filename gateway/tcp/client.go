package tcp

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/google/uuid"
	"io"
	"micro-libs/codec"
	"micro-libs/gateway"
	"micro-libs/utils/color"
	"micro-libs/utils/log"
	"net"
	"sync"
	"time"
)

type Client struct {
	sync.RWMutex
	id     string         // Client ID
	server gateway.Server // 服务器
	conn   net.Conn       // socket连接
	meta   *gateway.Meta  // 客户端上下文
	log    *logger.Helper // 日志对象

	writeChan chan []byte // 写入消息缓冲
	closed    bool        // 是否已关闭
	waitAuth  *time.Timer // 等待认证定时器
}

// 获取Client ID
func (c *Client) Id() string {
	return c.id
}

// 获取关联服务端
func (c *Client) Server() gateway.Server {
	return c.server
}

// 获取客户端元数据
func (c *Client) Meta() *gateway.Meta {
	return c.meta
}

// 日志对象
func (c *Client) Log() *logger.Helper {
	return c.log
}

// 判断是否关闭
func (c *Client) Closed() bool {
	c.RLock()
	defer c.RUnlock()

	return c.closed
}

// 认证成功
func (c *Client) SetAuthState(state bool) {
	if state {
		if c.waitAuth == nil {
			return
		}
		c.Lock()
		c.waitAuth.Stop()
		c.waitAuth = nil
		c.Unlock()
	} else {
		if c.waitAuth != nil {
			return
		}
		c.Lock()
		c.waitAuth = time.AfterFunc(c.server.Opts().WaitAuthTime, c.Close) // 连接成功后, 启动认证超时验证
		c.Unlock()
	}
}

// 发送消息
func (c *Client) Read() (*codec.ClientHead, []byte, error) {
	clientCodec := c.server.Gateway().ClientCodec()
	headLen := clientCodec.HeadLen()
	headBuf := make([]byte, headLen)
	if _, err := io.ReadFull(c.conn, headBuf); err != nil {
		return nil, nil, err
	}

	// 解析消息头
	head, _, err := clientCodec.Unmarshal(headBuf)
	if err != nil {
		return nil, nil, err
	}

	dataBuf := make([]byte, head.DataLen)
	if head.DataLen > 0 {
		if _, err := io.ReadFull(c.conn, dataBuf); err != nil {
			return nil, nil, err
		}
	}

	return head, dataBuf, nil
}

// 发送消息
func (c *Client) Write(b []byte) {
	c.Lock()
	defer c.Unlock()

	if c.closed || b == nil {
		return
	}

	c.doWrite(b)
}

// 执行写入消息
func (c *Client) doWrite(buf []byte) {
	if len(c.writeChan) == cap(c.writeChan) {
		c.log.Warn(color.Warn.Text("close conn: channel full"))
		c.doDestroy()
		return
	}

	c.writeChan <- buf
}

// 关闭连接
func (c *Client) Close() {
	c.Lock()
	defer c.Unlock()

	if c.closed {
		return
	}

	c.doWrite(nil)
	c.closed = true
}

// 销毁连接 (丢弃任何未发送或未确认的数据)
func (c *Client) Destroy() {
	c.Lock()
	defer c.Unlock()

	if c.closed {
		return
	}

	c.doDestroy()
}

// 关闭操作
func (c *Client) doDestroy() {
	_ = c.conn.(*net.TCPConn).SetLinger(0)
	_ = c.conn.Close()

	close(c.writeChan)
	c.closed = true
}

// 实例化新的客户端连接
func NewClient(server gateway.Server, conn net.Conn, ip string) gateway.Client {
	c := &Client{
		id:        uuid.New().String(),
		server:    server,
		conn:      conn,
		writeChan: make(chan []byte, 100),
	}

	// 设置客户端信息
	c.meta = gateway.NewMeta(c.id)
	c.meta.Set(gateway.MetaClientIp, ip)

	c.log = log.Logger.WithFields(map[string]interface{}{
		"client": c.id,
		"ip":     ip,
	})

	// 连接成功后, 启动认证超时验证
	c.SetAuthState(false)

	// 异步处理推送消息
	go func() {
		for b := range c.writeChan {
			if b == nil {
				break
			}

			if _, err := conn.Write(b); err != nil {
				c.log.Warnf(color.Warn.Text("Write client message error: %s", err))
				break
			}
		}

		_ = conn.Close()

		c.Lock()
		c.closed = true
		c.Unlock()

		c.log.Debugf("Client Write Chan is Closed ...")
	}()

	return c
}
