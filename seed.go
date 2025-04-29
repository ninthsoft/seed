package seed

import (
	"context"
	"net/http"
)

type MSeed interface {
	Router

	// HTTPServer 返回实例的 *http.Server
	HTTPServer() *http.Server

	// 设置默认的http server
	SetHTTPServer(srv *http.Server) MSeed

	// 静态资源服务
	Static(path string, root http.FileSystem) MSeed

	// Run 启动HTTPServer
	//
	// 	如果需要手动设定启动占用的端口请调用SetHTTPServer
	Run() error

	// Run 启动HTTPServer
	//
	// 	certFile https certFile
	// 	keyFile https keyFile
	// 	如果需要手动设定启动占用的端口请调用SetHTTPServer
	RunTLS(certFile, keyFile string) error

	// Shutdown gracefully shuts down the server
	Shutdown(ctx context.Context) error

	// NotFound 注册全局的404处理器
	NotFound(h http.Handler)
}

// mseed is driven by Router
type mseed struct {
	Router

	certFile  string
	keyFile   string
	enableTLS bool

	server *http.Server
}

func (c *mseed) HTTPServer() *http.Server {
	return c.server
}

func (c *mseed) SetHTTPServer(srv *http.Server) MSeed {
	c.server = srv
	return c
}

func (c *mseed) Static(path string, root http.FileSystem) MSeed {
	c.static(path, root)
	return c
}

func (c *mseed) Run() error {
	// set handler as itself
	c.server.Handler = c

	if c.enableTLS {
		return c.server.ListenAndServeTLS(c.certFile, c.keyFile)
	}
	return c.server.ListenAndServe()
}

func (c *mseed) RunTLS(certFile, keyFile string) error {
	c.certFile = certFile
	c.keyFile = keyFile
	c.enableTLS = true
	return c.Run()
}

func (c *mseed) Shutdown(ctx context.Context) error {
	return c.HTTPServer().Shutdown(ctx)
}

func (c *mseed) NotFound(h http.Handler) {
	c.notFound(h)
}

// New return *mseed
func New() MSeed {
	return &mseed{
		Router: NewRouter(),
		server: &http.Server{},
	}
}
