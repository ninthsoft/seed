package render

type Error struct {
	code int
	msg  string
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) Code() int {
	return e.code
}

func NewError(msg string, codes ...int) Error {
	var code = 1
	if len(codes) > 0 {
		code = codes[0]
	}
	return Error{code: code, msg: msg}
}
