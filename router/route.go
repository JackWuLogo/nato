package router

import (
	"context"
	"fmt"
	"reflect"
)

type Route struct {
	name    string
	hdlr    reflect.Value
	method  reflect.Method
	reqType reflect.Type
	rspType reflect.Type
}

func (h *Route) Name() string {
	return h.name
}

func (h *Route) NewReqValue() reflect.Value {
	return reflect.New(h.reqType)
}

func (h *Route) NewRspValue() reflect.Value {
	return reflect.New(h.rspType)
}

func (h *Route) Call(ctx context.Context, req, rsp interface{}) error {
	values := h.method.Func.Call([]reflect.Value{
		h.hdlr,
		reflect.ValueOf(ctx),
		reflect.ValueOf(req),
		reflect.ValueOf(rsp),
	})
	if err := values[0].Interface(); err != nil {
		return err.(error)
	}
	return nil
}

func ParseRoutes(handler interface{}) []*Route {
	var handlers []*Route

	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)

		mtype := method.Type
		reqType := mtype.In(2)
		rspType := mtype.In(3)

		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
		if rspType.Kind() == reflect.Ptr {
			rspType = rspType.Elem()
		}

		handlers = append(handlers, &Route{
			name:    fmt.Sprintf("%s.%s", name, method.Name),
			hdlr:    hdlr,
			method:  method,
			reqType: reqType,
			rspType: rspType,
		})
	}

	return handlers
}
