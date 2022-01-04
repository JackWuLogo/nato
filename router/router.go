package router

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"micro-libs/utils/errors"
	"micro-libs/utils/pb"
	"net/http"
)

// 游戏协议路由
type Router struct {
	routes map[string]*Route
}

func (r *Router) Routes() map[string]*Route {
	routes := make(map[string]*Route, len(r.routes))
	for _, route := range r.routes {
		routes[route.name] = route
	}
	return r.routes
}

func (r *Router) AddRoute(handles ...interface{}) {
	for _, handler := range handles {
		for _, h := range ParseRoutes(handler) {
			r.routes[h.name] = h
		}
	}
}

// 调用
func (r *Router) Call(ctx context.Context, method string, data *any.Any) (*any.Any, error) {
	route, ok := r.routes[method]
	if !ok {
		return nil, errors.New(http.StatusNotFound, "not found router method %s", method)
	}

	req := route.NewReqValue().Interface()
	if err := pb.UnmarshalAny(data, req.(proto.Message)); err != nil {
		return nil, err
	}

	rsp := route.NewRspValue().Interface()
	if err := route.Call(ctx, req, rsp); err != nil {
		return nil, err
	}

	out, err := pb.MarshalAny(rsp.(proto.Message))
	if err != nil {
		return nil, err
	}

	return out, nil
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}
