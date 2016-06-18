## gosws

gosws is mean golang simple web server.

## Quick Start

###### Download and install

    go get github.com/lvshuchengyin/gosws

###### Create file `main.go`

```go
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
