package middleware

import (
	"context"
	"net/http"

	"github.com/ninthsoft/seed"
)

// New will create a new middleware from a http.Handler.
func New(h http.Handler) seed.MiddlewareFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request, next seed.MiddleWareQueue) bool {
		h.ServeHTTP(w, req)
		return next.Next(ctx, w, req)
	}
}
