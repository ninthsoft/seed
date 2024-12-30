package seed

import (
	"net/http"
)

type MSeed interface {
	Router

	// HTTPServer 返回实例的 *http.Server
	HTTPServer() *http.Server

	// 设置默认的http server
	SetHTTPServer(srv *http.Server) MSeed

	// Run 启动HTTPServer
	//
	// 	addrs 监听的端口,非必选，但如果没有通过HTTPServer pointer修改且值不给可能导致服务启动失败
	Run(addrs ...string) error

	// Run 启动HTTPServer
	//
	// 	certFile https certFile
	// 	keyFile https keyFile
	// 	addrs 监听的端口,非必选，但如果没有通过HTTPServer pointer修改且值不给可能导致服务启动失败
	RunTLS(certFile, keyFile string, addrs ...string) error
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

// start http server
func (c *mseed) Run(addrs ...string) error {
	if len(addrs) > 0 {
		c.server.Addr = addrs[0]
	}
	// set handler as itself
	c.server.Handler = c

	if c.enableTLS {
		return c.server.ListenAndServeTLS(c.certFile, c.keyFile)
	}
	return c.server.ListenAndServe()
}

func (c *mseed) RunTLS(certFile, keyFile string, addrs ...string) error {
	c.certFile = certFile
	c.keyFile = keyFile
	c.enableTLS = true
	return c.Run(addrs...)
}

// New return *mseed
func New() MSeed {
	return &mseed{
		Router: NewRouter(),
		server: &http.Server{},
	}
}
