package protocol

import (
	"github.com/golang/protobuf/proto"
	"micro-libs/gateway"
	"micro-libs/utils/errors"
	"net/http"
)

// 游戏协议路由
type Router struct {
	routes map[uint32]*Route
}

// 注册
func (r *Router) Routes() map[uint32]*Route {
	routes := make(map[uint32]*Route, len(r.routes))
	for _, route := range r.routes {
		routes[route.cmd] = route
	}
	return routes
}

// 注册
func (r *Router) AddRoute(handles ...interface{}) {
	for _, handler := range handles {
		for _, h := range ParseRoutes(handler) {
			r.routes[h.cmd] = h
		}
	}
}

// 调用
func (r *Router) Call(gmt *gateway.Meta, cmd uint32, req []byte) (rsp []byte, err error) {
	route, ok := r.routes[cmd]
	if !ok {
		return nil, errors.New(http.StatusNotFound, "not found protocol %d", cmd)
	}

	c2s := route.NewReqValue().Interface()
	if err := proto.Unmarshal(req, c2s.(proto.Message)); err != nil {
		return nil, err
	}

	s2c := route.NewRspValue().Interface()
	if err := route.Call(gmt, c2s, s2c); err != nil {
		return nil, err
	}

	return proto.Marshal(s2c.(proto.Message))
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[uint32]*Route),
	}
}
