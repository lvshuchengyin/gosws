// manager_middleware
package gosws

import (
	"github.com/lvshuchengyin/gosws/context"
	"github.com/lvshuchengyin/gosws/logger"
	"github.com/lvshuchengyin/gosws/middleware"
	"github.com/lvshuchengyin/gosws/util"
)

type ManagerMiddleware struct {
	mdws []middleware.Middleware
}

func NewManagerMiddleware() *ManagerMiddleware {
	return &ManagerMiddleware{
		mdws: make([]middleware.Middleware, 0, 2),
	}
}

func (self *ManagerMiddleware) Add(name string) {
	for _, m := range self.mdws {
		if m.Name() == name {
			logger.Warning("middleware %s already add", name)
			return
		}
	}

	mw := middleware.Get(name)
	self.mdws = append(self.mdws, mw)
}

func (self *ManagerMiddleware) ProcessRequest(ctx *context.Context) error {
	for _, ms := range self.getCopy() {
		err := ms.ProcessRequest(ctx)
		if err == nil {
			continue
		}

		logger.Error("middleware %+v, ProcessRequest error! uri:%s, err:%v",
			ms, ctx.Req.URL.Path, err)

		return err
	}

	return nil
}

func (self *ManagerMiddleware) ProcessResponse(ctx *context.Context) error {
	mdws := self.getCopy()
	for i := len(mdws) - 1; i >= 0; i-- {
		err := mdws[i].ProcessResponse(ctx)
		if err == nil {
			continue
		}

		logger.Error("middleware %+v, ProcessResponse error! uri:%s, err:%v",
			mdws[i], ctx.Req.URL.Path, err)

		return err
	}

	return nil
}

func (self *ManagerMiddleware) getCopy() []middleware.Middleware {
	cpmws := make([]middleware.Middleware, 0, len(self.mdws))
	for _, ms := range self.mdws {
		cpmw := util.CloneValue(ms).(middleware.Middleware)
		cpmws = append(cpmws, cpmw)
	}

	return cpmws
}
