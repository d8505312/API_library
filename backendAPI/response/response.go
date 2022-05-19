package response

//response/response.go

import "github.com/gin-gonic/gin"

type Response struct {
	data interface{}
	code int
}

func Resp() *Response {

	return &Response{}
}

func (r *Response) Json(data gin.H, code int) *Response {
	r.data = data
	r.code = code
	return r
}

func (r *Response) String(data string, code int) *Response {
	r.data = data
	r.code = code
	return r
}

func (r *Response) Byte(data []byte, code int) *Response {
	r.data = data
	r.code = code
	return r
}

func (r *Response) GetData() interface{} {

	return r.data
}

func (r *Response) GetCode() int {

	return r.code
}
