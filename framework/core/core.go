package core

import (
	"log"
	"net/http"

	ctx "saiang/framework/context"
)

// Core defines the framework's core struct
type Core struct {
	router map[string]ctx.HandlerFunc
}

// NewCore initialize the Core struct
func NewCore() *Core {
	return &Core{router: make(map[string]ctx.HandlerFunc)}
}

func (c *Core) Get(path string, handler ctx.HandlerFunc) {
	c.router[path] = handler
}

// Run defines a method to run and listen the server's port specified
func (c *Core) Run(addr string) {
	http.ListenAndServe(addr, c)
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	log.Println("core.ServeHTTP")
	ctx := ctx.NewContext(r, w)

	router := c.router["foo"]
	if router == nil {
		return
	}
	log.Println("core.router")
	router(ctx)
}
