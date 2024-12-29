package bind

import (
	"github.com/gookit/validate"
	"github.com/gorilla/schema"
	"github.com/hexthink/seed"
	"github.com/hexthink/seed/render"
)

type BindType int

const (
	JSON BindType = iota
	Query
	Form
	PostForm
	MultipartForm
)

var decoder = schema.NewDecoder()

func Should(r seed.Request, dst interface{}, bts ...BindType) (err error) {
	if dst == nil {
		return render.NewError("dst object cannot be nil", 4000)
	}
	var bt = JSON
	if len(bts) > 0 {
		bt = bts[0]
	}
	var request = r.HTTPRequest()
	_ = request.ParseForm()

	var values = request.URL.Query()
	switch bt {
	case JSON:
		if err = r.JsonUnmarshal(dst); err != nil {
			return render.NewError(err.Error(), 4000)
		}
		if v := validate.Struct(dst); !v.Validate() {
			return render.NewError(v.Errors.String(), 4000)
		}
	case Form:
		values = request.Form
	case PostForm:
		values = request.PostForm
	case MultipartForm:
		values = request.MultipartForm.Value
	default:
	}
	if err = decoder.Decode(dst, values); err == nil {
		if v := validate.Struct(dst); !v.Validate() {
			return render.NewError(v.Errors.String(), 4000)
		}
		return
	}
	return render.NewError(err.Error(), 4000)
}
