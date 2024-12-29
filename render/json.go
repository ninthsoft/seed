package render

import (
	"context"
	"net/http"

	"github.com/hexthink/seed"
)

type jsonTmpl struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var DefaultJsonRender JsonRender = defaultJsonRender{}

type JsonRender interface {
	Format(ctx context.Context, data interface{}, err error) (int, interface{})
}

type defaultJsonRender struct {
}

func (r defaultJsonRender) Format(ctx context.Context, data interface{}, err error) (int, interface{}) {
	var resp = jsonTmpl{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	if e, ok := err.(Error); ok {
		resp.Code = e.Code()
		resp.Msg = e.Error()
	}
	return http.StatusOK, resp
}

var JSON = func(ctx context.Context, data interface{}, err error) (r seed.Response) {
	return seed.JsonResponse(DefaultJsonRender.Format(ctx, data, err))
}
