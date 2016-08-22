## gosws

gosws means golang simple web server.

## feature
* support view function and restful api
* middleware
* logger
* mysql dbpool
* secure cookie session

## Quick Start

###### Download and install

    go get github.com/lvshuchengyin/gosws

###### Create file `main.go`

```go
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
```

######Build and run
```bash
    go build main.go
    ./main
```
######Congratulations! 
You just built your first app.
Open your browser and visit `http://localhost:8000/123`.

## more config

	see example/conf.xml

## LICENSE

gosws source code is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
