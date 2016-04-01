// manager_middleware
package gosws

import (
	"gosws/context"
	"gosws/logger"
	"gosws/middleware"
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
	for _, ms := range self.mdws {
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
	for i := len(self.mdws) - 1; i >= 0; i-- {
		err := self.mdws[i].ProcessResponse(ctx)
		if err == nil {
			continue
		}

		logger.Error("middleware %+v, ProcessResponse error! uri:%s, err:%v",
			self.mdws[i], ctx.Req.URL.Path, err)

		return err
	}

	return nil
}
