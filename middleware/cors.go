package middleware

import (
	"context"
	"net/http"

	"github.com/ninthsoft/seed"
)

func CORS(ctx context.Context, w http.ResponseWriter, req *http.Request, next seed.MiddleWareQueue) bool {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	return next.Next(ctx, w, req)
}
