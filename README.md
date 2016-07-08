## gosws

gosws is mean golang simple web server.

## Quick Start

###### Download and install

    go get github.com/lvshuchengyin/gosws

###### Create file `main.go`

```go
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
