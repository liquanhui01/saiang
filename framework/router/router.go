package router

import (
	co "saiang/framework/core"
	"saiang/framework/handler"
)

func RegisterRouter(core *co.Core) {
	// static router
	core.Get("Foo", handler.FooHandler)
}
