package main

import (
	"net/http"

	"saiang/framework"
)

func main() {
	core := framework.NewCore()
	framework.RegisterRouter(core)
	server := http.Server{
		Handler: framework.NewCore(),
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
