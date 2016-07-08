// main
package main

import (
	"fmt"
	"time"

	"github.com/lvshuchengyin/gosws"
	"github.com/lvshuchengyin/gosws/controller"
)

type TestController struct {
	controller.Controller
}

func (self *TestController) Get(i int64) {
	self.Ctx.WriteString(fmt.Sprintf("%d, %d", time.Now().UnixNano(), i))
}

func main() {
	httpServer := gosws.NewHttpServer()
	httpServer.AddRoute(`^/(\d+)$`, &TestController{})
	httpServer.Run()
}
