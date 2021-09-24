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
	router map[string]*Tree
}

// NewCore initialize the Core struct
func NewCore() *Core {
	// define the second map
	// getRouter := map[string]ctx.HandlerFunc{}
	// postRouter := map[string]ctx.HandlerFunc{}
	// putRouter := map[string]ctx.HandlerFunc{}
	// deleteRouter := map[string]ctx.HandlerFunc{}
	// write the second layer map to the first layer map
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

func (c *Core) Get(path string, handler ctx.HandlerFunc) {
	// fmt.Println("路径为：", path)
	// urlPath := strings.ToUpper(path)
	// c.router["GET"][urlPath] = handler
	if err := c.router["GET"].AddRouter(path, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Post(path string, handler ctx.HandlerFunc) {
	// urlPath := strings.ToUpper(path)
	// c.router["POST"][urlPath] = handler
	if err := c.router["POST"].AddRouter(path, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Put(path string, handler ctx.HandlerFunc) {
	// urlPath := strings.ToUpper(path)
	// c.router["PUT"][urlPath] = handler
	if err := c.router["PUT"].AddRouter(path, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Delete(path string, handler ctx.HandlerFunc) {
	// urlPath := strings.ToUpper(path)
	// c.router["DELETE"][urlPath] = handler
	if err := c.router["DELETE"].AddRouter(path, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// FindRouterByRequest find request's handler according to request's method and url
func (c *Core) FindRouterByRequest(req *http.Request) ctx.HandlerFunc {
	upperPath := strings.ToUpper(req.URL.Path)
	fmt.Println("查询的路由为：", upperPath)
	upperMethod := strings.ToUpper(req.Method)
	if methodHandler, ok := c.router[upperMethod]; ok {
		// if handler, ok := methodHandler[upperPath]; ok {
		// 	return handler
		// }
		return methodHandler.FindHandler(upperPath)
	}
	return nil
}

func (c *Core) Group(prefix string) GroupInterface {
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
