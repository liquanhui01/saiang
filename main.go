package main

import (
<<<<<<< HEAD
	"os"
	"os/signal"
	"syscall"

	co "saiang/framework/core"
)

func main() {
	core := co.NewCore()
	RegisterRouter(core)
	go func() {
		core.Run(":8080")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
=======
	"net/http"

	"saiang/framework"
)

func main() {
	server := http.Server{
		Handler: framework.NewCore(),
		Addr:    ":8080",
	}
	server.ListenAndServe()
>>>>>>> af539f4... 重写框架，提交context和main.go部分
}
