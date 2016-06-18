// main
package main

import (
	"fmt"
	"time"

	"github.com/lvshuchengyin/gosws"
	"github.com/lvshuchengyin/gosws/context"
)

func test(ctx *context.Context, argid int64) {
	msg := fmt.Sprintf("hello world, nowtime:%d, argid:%d", time.Now().Unix(), argid)
	ctx.WriteString(msg)
}

func main() {
	httpServer := gosws.NewHttpServer()
	httpServer.AddHandle(`^/(\d+)$`, "get", test)
	httpServer.Run()
}
