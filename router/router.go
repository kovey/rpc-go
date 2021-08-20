package router

import (
	"fmt"
	"reflect"

	"github.com/kovey/logger-go/logger"
	hd "github.com/kovey/rpc-go/handler"
	"github.com/kovey/rpc-go/protocol"
)

type Router struct {
	handler     string
	method      string
	refType     reflect.Type
	paramsCount int
	base        reflect.Type
}

func NewRouter(handler string, method string, hClass interface{}) *Router {
	fun := reflect.ValueOf(hClass).MethodByName(method)
	if fun.IsZero() {
		panic(fmt.Sprintf("method[%s] is not exists", method))
	}

	return &Router{
		handler: handler, refType: reflect.TypeOf(hClass), method: method,
		paramsCount: fun.Type().NumIn(), base: reflect.TypeOf((*hd.HandlerInterface)(nil)).Elem(),
	}
}

func (r *Router) GetHandler() string {
	return r.handler
}

func (r *Router) GetMethod() string {
	return r.method
}

func (r *Router) GetValue() reflect.Value {
	return reflect.ValueOf(r.refType.Elem())
}

func (r *Router) Call(request *protocol.Request, spandId string) (interface{}, error) {
	var instance reflect.Value
	if r.refType.Kind() == reflect.Ptr {
		instance = reflect.New(r.refType.Elem())
	} else {
		instance = reflect.New(r.refType)
	}

	if !r.refType.Implements(r.base) {
		return nil, fmt.Errorf("handler[%s] is not extends handler.Handler", r.handler)
	}

	var base reflect.Value
	if instance.Kind() == reflect.Ptr {
		base = instance.Elem().FieldByName("Handler")
	} else {
		base = instance.FieldByName("Handler")
	}

	if !base.Type().Implements(r.base) {
		return nil, fmt.Errorf("field[Handler] not in [%s]", r.handler)
	}

	base.Set(reflect.ValueOf(hd.NewHandler(request.TraceId, spandId)))

	var vals []reflect.Value

	for _, val := range request.Args {
		vals = append(vals, reflect.ValueOf(val))
	}

	logger.Debug("call vals: %v", vals)

	fun := instance.MethodByName(r.method)
	if fun.IsZero() {
		return nil, fmt.Errorf("method[%s] is not exists in %s", r.method, r.handler)
	}

	result := fun.Call(vals)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		return v.Interface(), nil
	}

	return nil, nil
}
