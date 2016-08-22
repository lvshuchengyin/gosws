// context
package context

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lvshuchengyin/gosws/cookie_session"
	"github.com/lvshuchengyin/gosws/logger"
)

type Context struct {
	Req           *http.Request
	Res           http.ResponseWriter
	Status        int
	abort         bool
	Log           *logger.LogTrace
	cookieSession *cookie_session.CookieSession
	ctxData       map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request, secretKey string, sessExpire int64) *Context {
	ctx := &Context{
		Res:           w,
		Req:           r,
		Status:        200,
		Log:           logger.NewLogTrace(),
		cookieSession: cookie_session.NewCookieSession(secretKey, r, w, sessExpire),
		ctxData:       make(map[string]interface{}, 1),
	}

	err := ctx.cookieSession.Parse()
	if err != nil {
		logger.Info("cookieSession.Parse err: %s", err.Error())
	}

	return ctx
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

func (self *Context) WriteString(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
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

//--------------------session-------------------------
func (self *Context) SessGetInt(key string) int64 {
	return self.cookieSession.GetInt(key)
}

func (self *Context) SessGetString(key string) string {
	return self.cookieSession.GetString(key)
}

func (self *Context) SessSet(key string, val interface{}) {
	self.cookieSession.Set(key, val)
}

func (self *Context) SessDel(key string) {
	self.cookieSession.Del(key)
}

//--------------------ctxData-------------------------
func (self *Context) DataGet(key string) (val interface{}) {
	val, _ = self.ctxData[key]
	return
}

func (self *Context) DataSet(key string, val interface{}) {
	self.ctxData[key] = val
}

//--------------------logtrace-------------------------
func (self *Context) SetTraceStatus(name string, status interface{}) {
	self.Log.SetTrace(name, status)
}

func (self *Context) SetTrace(name string, err interface{}) {
	if err != nil {
		self.SetTraceStatus(name, "err")
	} else {
		self.SetTraceStatus(name, "ok")
	}

}
