package seed

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"slices"
	"strings"

	HRouter "github.com/julienschmidt/httprouter"
)

var MethodSep = ","

var allowedMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

// Router 路由器
type Router interface {
	http.Handler

	// Use 注册过滤器(中间件)
	//
	// 	可用于给 server 增加切片的功能
	// 	如可以创建注册一个用于校验是否登录的 中间件
	// 	注册的中间件会对此后的 handler 生效，在此之前注册的 handler 方法不会生效
	Use(ms ...MiddlewareFunc) Router

	// HandleStd 以http.Handler方式注册业务handler
	//
	// 	method  是http方法，如GET、POST,也可以使用逗号来连接同时传入多个，如 "GET,POST"
	// 	也可以用特殊的 ANY,会自动注册所有( ANY 的取值详见 MethodAny )
	// 	handler 是业务的逻辑
	// 	ms 是该接口特有的中间件函数
	HandleStd(methods string, path string, handler http.Handler, ms ...MiddlewareFunc)

	// HandleFunc 以HandlerFunc方式注册业务handler
	//
	// 	method  是http方法，如GET、POST,也可以使用逗号来连接同时传入多个，如 "GET,POST"
	// 	也可以用特殊的 ANY,会自动注册所有( ANY 的取值详见 MethodAny )
	// 	handler 是业务的逻辑
	// 	ms 是该接口特有的中间件函数
	HandleFunc(methods string, path string, handlerFunc HandlerFunc, ms ...MiddlewareFunc)

	// Group 路由分组
	//
	// 	如 可以将 /user/xxx 系列分成一个分组
	// 	prefix路由前缀，如 "/user"
	// 	ms 是该分组的中间件函数
	Group(prefix string, f func(r Router), ms ...MiddlewareFunc)

	// notFound 设置全局404状态处理器
	notFound(http.Handler)

	// 静态资源
	static(path string, root http.FileSystem)
}

type router struct {
	*HRouter.Router

	// prefix 路由前缀
	//
	// 	用于新建路由组等情况暂存前缀信息
	prefix string

	// middlewareFuncs 路由中间件
	//
	//用于新建路由组等情况暂存中间件信息，并一起并入最终的handler
	middlewareFuncs MiddlewareFuncs
}

func (r *router) Group(prefix string, f func(r Router), ms ...MiddlewareFunc) {
	var mws = make([]MiddlewareFunc, len(r.middlewareFuncs))

	//copy middlewares
	_ = copy(mws, r.middlewareFuncs)
	mws = append(mws, ms...)

	//make new router prefix
	var router = &router{Router: r.Router, middlewareFuncs: mws, prefix: r.prefix + prefix}
	f(router)
}

func (r *router) HandleFunc(methods string, path string, handlerFunc HandlerFunc, ms ...MiddlewareFunc) {
	r.HandleStd(methods, path, handlerFunc.Handler(), ms...)
}

func (r *router) HandleStd(methods string, mpath string, handler http.Handler, ms ...MiddlewareFunc) {
	var h = r.Trans2Handle(handler, ms...)
	var mss = strings.Split(methods, MethodSep)
	var apath = path.Clean(fmt.Sprintf("%s%s", r.prefix, mpath))
	for _, v := range mss {
		if slices.Index(allowedMethods, v) == -1 {
			panic(fmt.Sprintf("invalid router method '%s' for path '%s'", v, apath))
		}
		r.Handle(v, apath, h)
	}
}

func (r *router) Use(ms ...MiddlewareFunc) Router {
	if len(ms) > 0 {
		r.middlewareFuncs = append(r.middlewareFuncs, ms...)
	}
	return r
}

func (r *router) Trans2Handle(h http.Handler, ms ...MiddlewareFunc) HRouter.Handle {
	ms = append(r.middlewareFuncs, ms...)
	var f = func(w http.ResponseWriter, r *http.Request, pr HRouter.Params) {
		var mw MiddlewareFunc = func(ctx context.Context, ww http.ResponseWriter, rr *http.Request, next MiddleWareQueue) bool {
			h.ServeHTTP(ww, rr)
			return false
		}

		var mws MiddlewareFuncs = append(ms, mw)
		mws.Next(r.Context(), w, r)
	}
	return f
}

func (r *router) notFound(h http.Handler) {
	r.NotFound = h
}

func (r *router) static(path string, root http.FileSystem) {
	r.ServeFiles(path, root)
}

func NewRouter() Router {
	var r = &HRouter.Router{
		RedirectTrailingSlash:  false,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		HandleOPTIONS:          false,
		NotFound:               notFound,
	}
	return &router{Router: r, prefix: "", middlewareFuncs: []MiddlewareFunc{}}
}
