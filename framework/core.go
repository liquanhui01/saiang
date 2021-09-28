package framework

import (
	"fmt"
	"net/http"
	"strings"
)

type Core struct {
	router map[string]map[string]HandlerFunc
}

const (
	CONNECT = http.MethodConnect
	GET     = http.MethodGet
	POST    = http.MethodPost
	PUT     = http.MethodPut
	PATCH   = http.MethodPatch
	DELETE  = http.MethodDelete
	HEAD    = http.MethodHead
	OPTIONS = http.MethodOptions
	TRACE   = http.MethodTrace
)

func NewCore() *Core {
	// define the second layer map
	getRouter := map[string]HandlerFunc{}
	postRouter := map[string]HandlerFunc{}
	putRouter := map[string]HandlerFunc{}
	deleteRouter := map[string]HandlerFunc{}
	// write second layer map in first layer map
	router := map[string]map[string]HandlerFunc{}
	router[GET] = getRouter
	router[POST] = postRouter
	router[PUT] = putRouter
	router[DELETE] = deleteRouter
	return &Core{router: router}
}

func (c *Core) Get(url string, handler HandlerFunc) {
	c.add(GET, url, handler)
}

func (c *Core) Post(url string, handler HandlerFunc) {
	c.add(POST, url, handler)
}

func (c *Core) Put(url string, handler HandlerFunc) {
	c.add(PUT, url, handler)
}

func (c *Core) Delete(url string, handler HandlerFunc) {
	c.add(DELETE, url, handler)
}

func (c *Core) Patch(url string, handler HandlerFunc) {
	c.add(PATCH, url, handler)
}

func (c *Core) add(method, url string, handler HandlerFunc) {
	upperUrl := strings.ToUpper(url)
	c.router[method][upperUrl] = handler
}

func (c *Core) FindRouteByRequest(request *http.Request) HandlerFunc {
	url := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)
	upperUrl := strings.ToUpper(url)

	if methodHandlers, ok := c.router[upperMethod]; ok {
		if handler, ok := methodHandlers[upperUrl]; ok {
			return handler
		}
	}
	return nil
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Servering...")
	ctx := NewContext(w, r)
	router := c.FindRouteByRequest(r)
	if router == nil {
		ctx.Json("not found", 404)
		return
	}

	if err := router(ctx); err != nil {
		ctx.Json("inner error", 500)
		return
	}
}
