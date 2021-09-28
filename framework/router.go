package framework

import "fmt"

func RegisterRouter(co *Core) {
	co.Get("/user/login", UserLoginController)
}

func UserLoginController(c *Context) error {
	fmt.Println("登陆成功")
	return nil
}
