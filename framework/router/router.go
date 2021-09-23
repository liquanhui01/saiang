package router

import (
	co "saiang/framework/core"
	"saiang/framework/handler"
)

func RegisterRouter(core *co.Core) {
	core.Get("foo", handler.FooHandler)
}
