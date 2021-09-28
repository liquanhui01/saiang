package framework

import (
	"fmt"
	"net/http"
)

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("进入服务")
}
