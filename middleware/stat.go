// middleware_stat
package middleware

import (
	"time"

	"github.com/lvshuchengyin/gosws/context"
	"github.com/lvshuchengyin/gosws/logger"
)

type MiddlewareStat struct {
}

func (self *MiddlewareStat) Process(ctx *context.Context, nextFunc ProcessNextFunc) (err error) {
	startTime := time.Now().UnixNano()

	err = nextFunc()

	logger.Info("[%s:%s:%d:%dms] trace:%s", ctx.Req.Method, ctx.Req.RequestURI,
		ctx.Status, (time.Now().UnixNano()-startTime)/1000000, ctx.Log.String())

	return
}
