package middleware

import (
	"context"
	"net/http"

	"github.com/ninthsoft/seed"
)

func CORS(ctx context.Context, w http.ResponseWriter, req *http.Request, next seed.MiddleWareQueue) bool {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	//preflight request
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
		return false
	}
	return next.Next(ctx, w, req)
}
