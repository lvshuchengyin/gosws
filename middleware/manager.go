// manager
package middleware

import (
	"fmt"

	"github.com/lvshuchengyin/gosws/config"
	"github.com/lvshuchengyin/gosws/session"
)

var (
	middlewares map[string]Middleware = map[string]Middleware{}
)

func Register(name string, sf Middleware) error {
	middlewares[name] = sf
	return nil
}

func Get(name string) Middleware {
	mw, ok := middlewares[name]
	if !ok {
		panic(fmt.Sprintf("not found %s middleware", name))
	}
	return mw
}

func Init() error {
	confSession := config.Session()
	if confSession.Sessname == "" {
		return nil
	}

	name := confSession.Sessname
	secretKey := config.SecretKey()

	sessionFactory := session.Get(name)
	mw := Get(NAME_SESS)
	mw_sess := mw.(*MiddlewareSession)
	mw_sess.Init(sessionFactory, secretKey)

	return nil
}
