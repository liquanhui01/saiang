package core

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	ctx "saiang/framework/context"
)

// Core defines the framework's core struct. router is a double-layer map,
// the first layer key is the request method, the second layer key is the request
// url path, value is a handler function
type Core struct {
	router map[string]map[string]ctx.HandlerFunc
}

// NewCore initialize the Core struct
func NewCore() *Core {
	// define the second map
	getRouter := map[string]ctx.HandlerFunc{}
	postRouter := map[string]ctx.HandlerFunc{}
	putRouter := map[string]ctx.HandlerFunc{}
	deleteRouter := map[string]ctx.HandlerFunc{}
	// write the second layer map to the first layer map
	router := map[string]map[string]ctx.HandlerFunc{}
	router["GET"] = getRouter
	router["POST"] = postRouter
	router["PUT"] = putRouter
	router["DELETE"] = deleteRouter
	return &Core{router: router}
}

func (c *Core) Get(path string, handler ctx.HandlerFunc) {
	urlPath := strings.ToUpper(path)
	fmt.Println(urlPath)
	c.router["GET"][urlPath] = handler
}

func (c *Core) Post(path string, handler ctx.HandlerFunc) {
	urlPath := strings.ToUpper(path)
	c.router["POST"][urlPath] = handler
}

func (c *Core) Put(path string, handler ctx.HandlerFunc) {
	urlPath := strings.ToUpper(path)
	c.router["PUT"][urlPath] = handler
}

func (c *Core) Delete(path string, handler ctx.HandlerFunc) {
	urlPath := strings.ToUpper(path)
	c.router["DELETE"][urlPath] = handler
}

// FindRouterByRequest find request's handler according to request's method and url
func (c *Core) FindRouterByRequest(req *http.Request) ctx.HandlerFunc {
	upperPath := "/SUB" + strings.ToUpper(req.URL.Path)
	fmt.Println(upperPath)
	upperMethod := strings.ToUpper(req.Method)
	if methodHandler, ok := c.router[upperMethod]; ok {
		if handler, ok := methodHandler[upperPath]; ok {
			return handler
		}
	}
	return nil
}

func (c *Core) GroupFunc(prefix string) GroupInterface {
	return NewGroup(c, prefix)
}

// Run defines a method to run and listen the server's port specified
func (c *Core) Run(addr string) {
	http.ListenAndServe(addr, c)
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	log.Println("core.ServeHTTP")
	ctx := ctx.NewContext(r, w)

	router := c.FindRouterByRequest(r)
	if router == nil {
		ctx.Json(404, "not found")
		return
	}
	log.Println("core.router")
	router(ctx)
}
