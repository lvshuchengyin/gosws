// manager
package middleware

import "github.com/lvshuchengyin/gosws/context"

type ProcessNextFunc func() error
type ProcessResFunc func() error

type MiddlewareManager struct {
	middlewares []Middleware
}

func NewMiddlewareManager() *MiddlewareManager {
	m := &MiddlewareManager{
		middlewares: make([]Middleware, 0, 1),
	}

	return m
}

func (self *MiddlewareManager) Add(mdw Middleware) {
	self.middlewares = append(self.middlewares, mdw)
}

func (self *MiddlewareManager) Process(ctx *context.Context, processResFunc ProcessResFunc) (err error) {
	return self.processNext(0, ctx, processResFunc)
}

func (self *MiddlewareManager) processNext(index int, ctx *context.Context, processResFunc ProcessResFunc) (err error) {
	if index >= len(self.middlewares) {
		return processResFunc()
	}

	mdw := self.middlewares[index]
	cb := func() error {
		return self.processNext(index+1, ctx, processResFunc)
	}

	err = mdw.Process(ctx, cb)
	return
}
