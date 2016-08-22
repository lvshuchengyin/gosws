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
	"github.com/lvshuchengyin/gosws/middleware"
)

type HttpServer struct {
	server            *http.Server
	router            *Router
	middlewareManager *middleware.MiddlewareManager
	secretKey         string
	sessExpire        int64
}

func NewHttpServer() *HttpServer {
	addr := config.ListenAddr()
	if addr == "" {
		addr = ":8000"
	}

	hs := &HttpServer{
		server: &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		},
		router:            NewRouter(),
		middlewareManager: middleware.NewMiddlewareManager(),
		secretKey:         config.SecretKey(),
		sessExpire:        config.SessionLifeTime(),
	}

	// middleware
	hs.middlewareManager.Add(&middleware.MiddlewareStat{})

	// route
	http.HandleFunc("/", hs.Process)

	// static files
	staticDir := config.StaticDir()
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	return hs
}

func (self *HttpServer) AddMiddleware(mdw middleware.Middleware) {
	self.middlewareManager.Add(mdw)
}

func (self *HttpServer) AddController(pattern string, ctrl controller.ControllerInterface) error {
	return self.router.AddController(pattern, ctrl)
}

func (self *HttpServer) AddRoute(pattern, method string, handleFunc interface{}) error {
	// check the handleFunc params
	ht := reflect.TypeOf(handleFunc)
	if ht.Kind() != reflect.Func {
		err := fmt.Errorf("%s not a func type", ht.Name())
		panic(err)
		return err
	}

	for i := 0; i < ht.NumIn(); i++ {
		at := ht.In(i)
		if i == 0 {
			var hap *context.Context
			if at != reflect.TypeOf(hap) {
				err := fmt.Errorf("handle func: %s, first arg type must be *Context, now is %+v", ht.String(), at)
				panic(err)
				return err
			}
			continue
		}

		k := at.Kind()
		if k != reflect.Int64 && k != reflect.Float64 && k != reflect.String {
			err := fmt.Errorf("handle func: %s, args type must be int64, float64, string", ht.String())
			panic(err)
			return err
		}
	}

	return self.router.AddRoute(pattern, method, handleFunc)
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
	ctx := context.NewContext(w, r, self.secretKey, self.sessExpire)

	defer func() {
		if rec := recover(); rec != nil && !ctx.IsAbort() {
			logger.Error("panic! uri:%s, err:%v \n%s", r.URL.Path, rec, string(debug.Stack()))
			self.End(w, 500, "server error 500")
		}
	}()

	// route
	handle, ss := self.router.Route(r)
	if handle == nil {
		logger.Warning("Not found: %s", r.URL.Path)
		self.End(w, 404, "not found")
		return
	}

	// cb
	processResFunc := func() error {
		ctrl, ok := handle.(controller.ControllerInterface)
		if ok {
			self.processController(w, r, ctrl, ctx, ss)
		} else {
			self.processHandleFunc(w, r, handle, ctx, ss)
		}

		return nil
	}

	self.middlewareManager.Process(ctx, processResFunc)
}

func (self *HttpServer) processController(w http.ResponseWriter, r *http.Request,
	ctrl controller.ControllerInterface, ctx *context.Context, ss []string) {

	// new ctrl
	ctrlValue := reflect.ValueOf(ctrl)
	newCtrlValue := reflect.New(ctrlValue.Elem().Type())
	newCtrl := newCtrlValue.Interface().(controller.ControllerInterface)
	newCtrl.Init(ctx)

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
		logger.Error("args num not correct, get args: %v", ss)
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
				self.End(w, 400, "bad request, 40060")
				return
			}
			args = append(args, reflect.ValueOf(v))
		case reflect.Float64:
			v, err := strconv.ParseInt(ss[i], 10, 64)
			if err != nil {
				self.End(w, 400, "bad request, 40061")
				return
			}
			args = append(args, reflect.ValueOf(v))
		case reflect.String:
			args = append(args, reflect.ValueOf(ss[i]))
		default:
			logger.Error("controller:%+v args have invalid type:%+v, must be int64, float64, string", ctrl, fak)
			self.End(w, 400, "bad request, 40062")
			return
		}
	}

	// do
	newCtrlValue.MethodByName("Prepare").Call(nil)
	methodValue.Call(args)
	newCtrlValue.MethodByName("Finish").Call(nil)
	return
}

func (self *HttpServer) processHandleFunc(w http.ResponseWriter, r *http.Request,
	handleFunc interface{}, ctx *context.Context, ss []string) {

	handleType := reflect.TypeOf(handleFunc)
	if handleType.NumIn()-1 != len(ss) {
		logger.Error("args num not correct, get args: %v", ss)
		self.End(w, 400, "bad request, 40076")
		return
	}

	// convert
	args := []reflect.Value{}
	for i := 0; i < len(ss); i++ {
		fak := handleType.In(i + 1).Kind()
		switch fak {
		case reflect.Int64:
			v, err := strconv.ParseInt(ss[i], 10, 64)
			if err != nil {
				self.End(w, 400, "bad request, 40070")
				return
			}
			args = append(args, reflect.ValueOf(v))
		case reflect.Float64:
			v, err := strconv.ParseInt(ss[i], 10, 64)
			if err != nil {
				self.End(w, 400, "bad request, 40071")
				return
			}
			args = append(args, reflect.ValueOf(v))
		case reflect.String:
			args = append(args, reflect.ValueOf(ss[i]))
		default:
			logger.Error("handleFunc:%+v args have invalid type:%+v, must be int64, float64, string", handleFunc, fak)
			self.End(w, 400, "bad request, 40072")
			return
		}
	}

	// add ctx argument
	args = append([]reflect.Value{reflect.ValueOf(ctx)}, args...)

	// do
	reflect.ValueOf(handleFunc).Call(args)
	return
}
