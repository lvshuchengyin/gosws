// router
package gosws

import (
	"fmt"
	"gosws/context"
	"gosws/logger"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type HandleInfo struct {
	re         *regexp.Regexp
	method     string
	handleFunc interface{}
}

type Router struct {
	handleInfos []*HandleInfo
}

func NewRouter() *Router {
	return &Router{
		handleInfos: make([]*HandleInfo, 0, 2),
	}
}

// add handle, must add before sws start
func (self *Router) AddHandle(pattern, method string, handleFunc interface{}) error {
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
				err := fmt.Errorf("handle func: %s, first arg type must be *Context", ht.String())
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

	info := &HandleInfo{
		re:         regexp.MustCompile(pattern),
		method:     strings.ToUpper(method),
		handleFunc: handleFunc,
	}

	self.handleInfos = append(self.handleInfos, info)
	return nil
}

func (self *Router) Route(r *http.Request) (handleFunc interface{}, args []reflect.Value) {
	for _, handleInfo := range self.handleInfos {
		if handleInfo.method != "*" && handleInfo.method != strings.ToUpper(r.Method) {
			continue
		}

		// match route
		ss := handleInfo.re.FindStringSubmatch(r.URL.Path)
		if len(ss) <= 0 {
			continue
		}

		handleType := reflect.TypeOf(handleInfo.handleFunc)
		if handleType.NumIn() != len(ss) {
			continue
		}

		argsOk := true
		for i := 1; i < len(ss); i++ {
			fak := handleType.In(i).Kind()
			// 类型转换
			if fak == reflect.Int64 {
				v, err := strconv.ParseInt(ss[i], 10, 64)
				if err != nil {
					argsOk = false
					break
				}
				args = append(args, reflect.ValueOf(v))

			} else if fak == reflect.Float64 {
				v, err := strconv.ParseFloat(ss[i], 64)
				if err != nil {
					argsOk = false
					break
				}
				args = append(args, reflect.ValueOf(v))

			} else if fak == reflect.String {
				args = append(args, reflect.ValueOf(ss[i]))

			} else {
				logger.Error("handle func have invalid type, must be int64, float64, string")
				argsOk = false
				break
			}
		}

		if !argsOk {
			continue
		}

		logger.Swsd("route match, pattern[%s], method[%s], url[%s], args%v", handleInfo.re.String(), r.Method, r.URL.Path, ss)

		handleFunc = handleInfo.handleFunc
		return
	}

	return
}
