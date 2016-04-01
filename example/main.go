// main
package main

import (
	"fmt"
	"gosws"
	"gosws/context"
	"time"
	"views"
)

func test(ctx *context.Context) {
	ctx.WriteString(fmt.Sprintf("%d", time.Now().UnixNano()))
}

func main() {
	err := gosws.Init("conf.xml")
	if err != nil {
		fmt.Println("gosws Init err:", err)
		panic(err)
	}

	httpServer := gosws.NewHttpServer()

	httpServer.AddHandle("/", "get", test)

	err = httpServer.Run()

	fmt.Println("httpserver run err:", err)
}
