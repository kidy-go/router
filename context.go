// context.go kee > 2021/03/16

package router

import "net/http"

type Context interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

type context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
}

func (ctx *context) Request() *http.Request {
	return ctx.request
}

func (ctx *context) ResponseWriter() http.ResponseWriter {
	return ctx.responseWriter
}

func NewContext(r *http.Request, w http.ResponseWriter) *context {
	return &context{
		request:        r,
		responseWriter: w,
	}
}
