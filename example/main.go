// main
package main

import (
	"time"

	"github.com/lvshuchengyin/gosws"
	"github.com/lvshuchengyin/gosws/context"
	"github.com/lvshuchengyin/gosws/controller"
)

// support view function with arguments parse
func Test(ctx *context.Context, s string) {
	ctx.WriteString(s)
	ctx.WriteString("\n now is %d", time.Now().Unix())
}

// support restful api with arguments parse
type TestController struct {
	controller.Controller
}

func (self *TestController) Get(i int64) {
	self.Ctx.WriteString("%d, %d", time.Now().Unix(), i)
}

func main() {
	httpServer := gosws.NewHttpServer()
	httpServer.AddController(`^/(\d+)$`, &TestController{})
	httpServer.AddRoute(`^/(\w+)$`, "get", Test)
	httpServer.Run()
}
