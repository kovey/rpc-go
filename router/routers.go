package router

import (
	"fmt"
	"reflect"
)

var (
	Route *Routers
)

type Routers struct {
	routers map[string]*Router
}

func NewRouters() *Routers {
	return &Routers{routers: make(map[string]*Router, 0)}
}

func (r *Routers) Add(path string, router *Router) *Routers {
	r.routers[getPath(path, router.GetMethod())] = router
	return r
}

func (r *Routers) AddAll(path string, obj interface{}) *Routers {
	val := reflect.TypeOf(obj)
	for i := 0; i < val.NumMethod(); i++ {
		fun := val.Method(i)
		router := NewRouter(path, fun.Name, obj)
		r.Add(path, router)
	}

	return r
}

func (r *Routers) Get(path string, method string) (*Router, error) {
	router, ok := r.routers[getPath(path, method)]
	if !ok {
		return nil, fmt.Errorf("method[%s] of path[%s] is not exists", method, path)
	}

	return router, nil
}

func (r *Routers) IsExists(path string, method string) bool {
	_, ok := r.routers[getPath(path, method)]
	return ok
}

func getPath(path string, method string) string {
	return fmt.Sprintf("%s.%s", path, method)
}

func init() {
	Route = NewRouters()
}
