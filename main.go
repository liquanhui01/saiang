package main

import (
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
}
