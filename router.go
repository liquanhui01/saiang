package main

import (
	co "saiang/framework/core"
)

func RegisterRouter(core *co.Core) {
	// static router
	core.Get("Foo", FooHandler)
	// router group
	subApi := core.Group("/sub")
	{
		subApi.Delete("/:id", FooHandler)
		subApi.Put("/:id", FooHandler)
		subApi.Get("/:id", FooHandler)
		subApi.Get("/list/all", FooHandler)
	}
}
