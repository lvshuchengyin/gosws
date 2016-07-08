// context
package context

import (
	"encoding/json"
	"net/http"

	"github.com/lvshuchengyin/gosws/logger"
	"github.com/lvshuchengyin/gosws/session"
)

type Context struct {
	Req     *http.Request
	Res     http.ResponseWriter
	Status  int
	Session session.Session
	abort   bool
	Log     *logger.LogTrace
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Res:    w,
		Req:    r,
		Status: 200,
		Log:    logger.NewLogTrace(),
	}
}

func (self *Context) Query(key string) string {
	return self.Req.FormValue(key)
}

func (self *Context) Abort(status int, msg string) {
	self.abort = true
	self.Status = status
	self.Res.WriteHeader(status)
	_, err := self.Res.Write([]byte(msg))
	if err != nil {
		logger.Error("Abort write err: %v", err)
	}
	panic(msg)
}

func (self *Context) IsAbort() bool {
	return self.abort
}

func (self *Context) Error(status int, msg string) {
	self.Status = status
	http.Error(self.Res, msg, status)
}

func (self *Context) SetStatus(status int) {
	self.Status = status
	self.Res.WriteHeader(status)
}

func (self *Context) WriteString(msg string) error {
	_, err := self.Res.Write([]byte(msg))
	if err != nil {
		logger.Error("WriteString err: %v", err)
	}
	return err
}

func (self *Context) ServeJson(code int, msg string, data interface{}) error {
	jsmap := map[string]interface{}{"code": code, "msg": msg, "data": data}
	bs, _ := json.Marshal(jsmap)
	return self.WriteString(string(bs))
}

func (self *Context) Redirect(uri string) error {
	http.Redirect(self.Res, self.Req, uri, 302)
	return nil
}

//--------------------logtrace-------------------------

func (self *Context) Module(m string) {
	self.Log.Moudle(m)
}

func (self *Context) Start(m string) {
	self.Log.Start(m)
}

func (self *Context) End(m string) {
	self.Log.End(m)
}
