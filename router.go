// router
package gosws

import (
	"net/http"
	"regexp"

	"github.com/lvshuchengyin/gosws/controller"
	"github.com/lvshuchengyin/gosws/logger"
)

type ControllerInfo struct {
	re   *regexp.Regexp
	ctrl controller.ControllerInterface
}

type Router struct {
	ctrlInfos []*ControllerInfo
}

func NewRouter() *Router {
	return &Router{
		ctrlInfos: make([]*ControllerInfo, 0, 2),
	}
}

// add handle, must add before sws start
func (self *Router) AddRoute(pattern string, ctrl controller.ControllerInterface) error {
	info := &ControllerInfo{
		re:   regexp.MustCompile(pattern),
		ctrl: ctrl,
	}

	self.ctrlInfos = append(self.ctrlInfos, info)
	return nil
}

func (self *Router) Route(r *http.Request) (ctrl controller.ControllerInterface, args []string) {
	for _, ctrlInfo := range self.ctrlInfos {
		ss := ctrlInfo.re.FindStringSubmatch(r.URL.Path)
		if len(ss) <= 0 {
			continue
		}

		logger.Debug("route match, pattern[%s], method[%s], url[%s], args:%v", ctrlInfo.re.String(), r.Method, r.URL.Path, ss)

		ctrl = ctrlInfo.ctrl
		args = ss[1:]

		return
	}

	return
}
