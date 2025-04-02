package seed

import (
	"context"
	"net/http"
)

// HandlerFunc 标准的HandlerFunc
type HandlerFunc func(ctx context.Context, req Request) Response

// Handler HandlerFunc自身转换为http.Handler
func (h HandlerFunc) Handler() http.Handler {
	var f http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		var request = NewRequest(r)
		var response = h(r.Context(), request)
		if response != nil {
			_ = response.WriteTo(w)
		}
	}
	return f
}

var notFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Preflight request for OPTIONS request method
	if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.NotFoundHandler().ServeHTTP(w, r)
})
