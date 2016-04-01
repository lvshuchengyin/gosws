// middleware
package middleware

import "gosws/context"

type Middleware interface {
	Name() string
	ProcessRequest(arg *context.Context) error
	ProcessResponse(arg *context.Context) error
}
