package main

import (
	co "saiang/framework/core"
	"saiang/framework/router"
)

func main() {
	core := co.NewCore()
	router.RegisterRouter(core)
	// server := &http.Server{
	// 	Handler: core,
	// 	Addr:    ":8080",
	// }
	// server.ListenAndServe()
	core.Run(":8080")
}
