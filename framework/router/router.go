package router

import (
	co "saiang/framework/core"
	"saiang/framework/handler"
)

func RegisterRouter(core *co.Core) {
	// static router
	core.Get("Foo", handler.FooHandler)
	// router group
	subApi := core.GroupFunc("/sub")
	{
		subApi.Get("/Foo", handler.FooHandler)
	}
}
