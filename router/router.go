package router

import (
	"fmt"
	"reflect"

	"github.com/kovey/logger-go/logger"
)

type Router struct {
	handler     string
	method      string
	refType     reflect.Type
	paramsCount int
}

func NewRouter(handler string, method string, hClass interface{}) *Router {
	fun := reflect.ValueOf(hClass).MethodByName(method)
	if fun.IsZero() {
		panic(fmt.Sprintf("method[%s] is not exists", method))
	}

	return &Router{
		handler: handler, refType: reflect.TypeOf(hClass), method: method,
		paramsCount: fun.Type().NumIn(),
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

func (r *Router) Call(params ...interface{}) (interface{}, error) {
	instance := reflect.New(r.refType.Elem())
	var vals []reflect.Value

	for _, val := range params {
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
