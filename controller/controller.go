// controller
package controller

import "github.com/lvshuchengyin/gosws/context"

type ControllerInterface interface {
	Init(ctx *context.Context)
	Prepare()
	Finish()
}

type Controller struct {
	Ctx *context.Context
}

func (self *Controller) Init(ctx *context.Context) {
	self.Ctx = ctx
}

func (self *Controller) Prepare() {

}

func (self *Controller) Finish() {

}

func (self *Controller) Get() {
	self.Ctx.Error(405, "Method Not Allowed")
}

func (self *Controller) Post() {
	self.Ctx.Error(405, "Method Not Allowed")
}

func (self *Controller) Delete() {
	self.Ctx.Error(405, "Method Not Allowed")
}

func (self *Controller) Put() {
	self.Ctx.Error(405, "Method Not Allowed")
}

func (self *Controller) Head() {
	self.Ctx.Error(405, "Method Not Allowed")
}

func (self *Controller) Patch() {
	self.Ctx.Error(405, "Method Not Allowed")
}

func (self *Controller) Options() {
	self.Ctx.Error(405, "Method Not Allowed")
}
