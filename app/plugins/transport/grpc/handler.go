package grpc

import (
	"runtime/debug"

	"github.com/asim/go-micro/v3/errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/transport"
	"google.golang.org/grpc/peer"
	pb "micro-libs/app/plugins/transport/grpc/proto"
)

// microTransport satisfies the pb.TransportServer inteface
type microTransport struct {
	addr string
	fn   func(transport.Socket)
}

func (m *microTransport) Stream(ts pb.Transport_StreamServer) (err error) {

	sock := &grpcTransportSocket{
		stream: ts,
		local:  m.addr,
	}

	p, ok := peer.FromContext(ts.Context())
	if ok {
		sock.remote = p.Addr.String()
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r, string(debug.Stack()))
			sock.Close()
			err = errors.InternalServerError("go.micro.transport", "panic recovered: %v", r)
		}
	}()

	// execute socket func
	m.fn(sock)

	return err
}
