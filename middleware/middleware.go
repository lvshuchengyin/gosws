// middleware
package middleware

import "github.com/lvshuchengyin/gosws/context"

type Middleware interface {
	Process(ctx *context.Context, nextFunc ProcessNextFunc) error
}
