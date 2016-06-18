// stat
package middleware

import (
	"time"

	"github.com/lvshuchengyin/gosws/context"
	"github.com/lvshuchengyin/gosws/logger"
)

const (
	NAME_STAT = "stat"
)

func init() {
	Register(NAME_STAT, &MiddlewareStat{})
}

type MiddlewareStat struct {
	startTime int64
}

func (self *MiddlewareStat) Name() string {
	return NAME_STAT
}

func (self *MiddlewareStat) ProcessRequest(arg *context.Context) error {
	self.startTime = time.Now().UnixNano()
	return nil
}

func (self *MiddlewareStat) ProcessResponse(arg *context.Context) error {
	logger.Info("[%s:%s:%d:%dms]; trace:%s", arg.Req.Method, arg.Req.RequestURI,
		arg.Status, (time.Now().UnixNano()-self.startTime)/1000000, arg.Log.String())
	return nil
}
