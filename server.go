// server
package gosws

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/lvshuchengyin/gosws/config"
	"github.com/lvshuchengyin/gosws/context"
	"github.com/lvshuchengyin/gosws/controller"
	"github.com/lvshuchengyin/gosws/logger"
)

type HttpServer struct {
	server            *http.Server
	router            *Router
	managerMiddleware *ManagerMiddleware
}

func NewHttpServer() *HttpServer {
	addr := config.ListenAddr()
	if addr == "" {
		addr = ":8000"
	}
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
	staticDir := config.StaticDir()
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	return hs
}

func (self *HttpServer) AddRoute(pattern string, ctrl controller.ControllerInterface) error {
	return self.router.AddRoute(pattern, ctrl)
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
	ctx := context.NewContext(w, r)

	defer func() {
		if rec := recover(); rec != nil && !ctx.IsAbort() {
			logger.Error("panic! uri:%s, err:%v \n%s", r.URL.Path, rec, string(debug.Stack()))
			self.End(w, 500, "server error 500")
		}
	}()

	// route
	ctrl, ss := self.router.Route(r)
	if ctrl == nil {
		logger.Warning("Not found: %s", r.URL.Path)
		self.End(w, 404, "not found")
		return
	}

	// new ctrl
	ctrlValue := reflect.ValueOf(ctrl)
	newCtrlValue := reflect.New(ctrlValue.Elem().Type())
	newCtrl := newCtrlValue.Interface().(controller.ControllerInterface)
	newCtrl.Init(ctx)

	// mw process request
	err := self.managerMiddleware.ProcessRequest(ctx)
	if err != nil {
		logger.Error("managerMiddleware.ProcessRequest error, uri:%s, err:%v", r.URL.Path, err)
		self.End(w, 500, "server error 50010")
		return
	}

	ctrlMethodName := ""
	switch strings.ToUpper(r.Method) {
	case "GET":
		ctrlMethodName = "Get"
	case "POST":
		ctrlMethodName = "Post"
	case "DELETE":
		ctrlMethodName = "Delete"
	case "PUT":
		ctrlMethodName = "Put"
	case "HEAD":
		ctrlMethodName = "Head"
	case "PATCH":
		ctrlMethodName = "Patch"
	case "OPTIONS":
		ctrlMethodName = "Options"
	default:
		logger.Warning("unknow method: %s", r.Method)
		self.End(w, 400, fmt.Sprintf("unknow method: %s, 40050", r.Method))
		return
	}

	methodValue := newCtrlValue.MethodByName(ctrlMethodName)
	if !methodValue.IsValid() {
		logger.Error("can't find controller method:%s", ctrlMethodName)
		self.End(w, 404, "not found method, 40451")
		return
	}

	handleType := methodValue.Type()
	if handleType.NumIn() != len(ss) {
		logger.Error("args num not correct")
		self.End(w, 400, "bad request, 40052")
		return
	}

	// convert
	args := []reflect.Value{}
	for i := 0; i < len(ss); i++ {
		fak := handleType.In(i).Kind()
		switch fak {
		case reflect.Int64:
			v, err := strconv.ParseInt(ss[i], 10, 64)
			if err != nil {
				self.End(w, 500, "bad request, 50060")
				break
			}
			args = append(args, reflect.ValueOf(v))
		case reflect.Float64:
			v, err := strconv.ParseInt(ss[i], 10, 64)
			if err != nil {
				self.End(w, 500, "bad request, 50061")
				break
			}
			args = append(args, reflect.ValueOf(v))
		case reflect.String:
			args = append(args, reflect.ValueOf(ss[i]))
		default:
			logger.Error("controller:%+v args have invalid typev:%+v, must be int64, float64, string", ctrl, fak)
			self.End(w, 500, "bad request, 50062")
			return
		}
	}

	// do
	newCtrlValue.MethodByName("Prepare").Call(nil)
	methodValue.Call(args)
	newCtrlValue.MethodByName("Finish").Call(nil)

	// mw process reponse
	err = self.managerMiddleware.ProcessResponse(ctx)
	if err != nil {
		logger.Error("managerMiddleware.ProcessResponse error, uri:%s, err:%v", r.URL.Path, err)
		self.End(w, 500, "server error 50011")
		return
	}
}
