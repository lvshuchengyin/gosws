// server
package gosws

import (
	"gosws/config"
	"gosws/context"
	"gosws/logger"
	"net/http"
	"reflect"
	"runtime/debug"
)

type HttpServer struct {
	server            *http.Server
	router            *Router
	managerMiddleware *ManagerMiddleware
}

func NewHttpServer() *HttpServer {
	addr := config.ListenAddr()
	middlewares := config.MiddlewareNames()

	hs := &HttpServer{
		server: &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		},
		router:            NewRouter(),
		managerMiddleware: NewManagerMiddleware(),
	}

	// middleware
	for _, mwName := range middlewares {
		hs.managerMiddleware.Add(mwName)
	}

	// route
	http.HandleFunc("/", hs.Process)

	// static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	return hs
}

func (self *HttpServer) AddHandle(pattern, method string, handleFunc interface{}) error {
	return self.router.AddHandle(pattern, method, handleFunc)
}

// listen and serve
func (self *HttpServer) Run() (err error) {
	logger.Info("sws server ListenAndServe on %s", self.server.Addr)
	return self.server.ListenAndServe()
}

func (self *HttpServer) End(w http.ResponseWriter, status int, data string) {
	w.WriteHeader(status)
	w.Write([]byte(data))
}

func (self *HttpServer) Process(w http.ResponseWriter, r *http.Request) {
	handleFunc, args := self.router.Route(r)
	if handleFunc == nil {
		logger.Warning("Not found: %s", r.URL.Path)
		self.End(w, 404, "not found")
		return
	}

	ctx := &context.Context{
		Res:    w,
		Req:    r,
		Status: 200,
		Log:    logger.NewLogTrace(),
	}

	err := self.managerMiddleware.ProcessRequest(ctx)
	if err != nil {
		logger.Error("managerMiddleware.ProcessRequest error, uri:%s, err:%v", r.URL.Path, err)
		self.End(w, 500, "server error 50010")
		return
	}

	func() {
		defer func() {
			if rec := recover(); rec != nil && !ctx.IsAbort() {
				logger.Error("panic! uri:%s, err:%v \n%s", r.URL.Path, rec, string(debug.Stack()))
				self.End(w, 500, "server error 500")
			}
		}()

		args = append([]reflect.Value{reflect.ValueOf(ctx)}, args...)

		reflect.ValueOf(handleFunc).Call(args)
	}()

	err = self.managerMiddleware.ProcessResponse(ctx)
	if err != nil {
		logger.Error("managerMiddleware.ProcessResponse error, uri:%s, err:%v", r.URL.Path, err)
		self.End(w, 500, "server error 50011")
		return
	}
}
