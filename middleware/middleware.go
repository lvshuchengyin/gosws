// middleware
package middleware

import "github.com/lvshuchengyin/gosws/context"

type Middleware interface {
	Name() string
	ProcessRequest(arg *context.Context) error
	ProcessResponse(arg *context.Context) error
}
