// router
package gosws

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/lvshuchengyin/gosws/controller"
	"github.com/lvshuchengyin/gosws/logger"
)

type HandleInfo struct {
	re     *regexp.Regexp
	method string
	handle interface{}
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
func (self *Router) AddController(pattern string, ctrl controller.ControllerInterface) error {
	info := &HandleInfo{
		re:     regexp.MustCompile(pattern),
		handle: ctrl,
	}

	self.handleInfos = append(self.handleInfos, info)
	return nil
}

func (self *Router) AddRoute(pattern, method string, handleFunc interface{}) error {
	method = strings.ToUpper(method)
	info := &HandleInfo{
		re:     regexp.MustCompile(pattern),
		method: method,
		handle: handleFunc,
	}

	self.handleInfos = append(self.handleInfos, info)
	return nil
}

func (self *Router) Route(r *http.Request) (handle interface{}, args []string) {
	for _, handleInfo := range self.handleInfos {
		ss := handleInfo.re.FindStringSubmatch(r.URL.Path)
		if len(ss) <= 0 {
			continue
		}

		if handleInfo.method != "" && handleInfo.method != strings.ToUpper(r.Method) {
			continue
		}

		handle = handleInfo.handle
		args = ss[1:]

		logger.Debug("route match, pattern[%s], method[%s], url[%s], args:%v", handleInfo.re.String(), r.Method, r.URL.Path, args)

		return
	}

	return
}
