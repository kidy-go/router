// response.go kee > 2021/04/08

package router

import (
	"fmt"
	// "github.com/kidy-go/utils"
	//"io"
	"net/http"
)

type Writer interface {
	SetCharset(string) Writer
	Header(string, string) Writer
	WithHeaders(map[string]string) Writer
	WithStatus(int) Writer
	WithBody([]byte) Writer
	WithCookie(*http.Cookie) Writer
	WithCookies([]*http.Cookie) Writer

	Write([]byte) (int, error)
	WriteString(string) (int, error)
	Writef(string, ...interface{}) (int, error)
	Writer() http.ResponseWriter
}

type response struct {
	rw         http.ResponseWriter
	header     http.Header
	body       []byte
	statusCode int
	statusText string
	charset    string
}

const (
	IsInvalid = iota
	IsInformational
	IsSuccessful
	IsRedirection
	IsClientError
	IsServerError
	IsOk
	IsForbidden
	IsNotFound
)

func NewResponse(ctx Context) *response {
	rw := ctx.ResponseWriter()
	return &response{
		rw:         rw,
		header:     rw.Header(),
		statusCode: 200,
	}
}

func (resp *response) WithBody(body []byte) *response {
	resp.body = body
	return resp
}

func (resp *response) WithHeaders(header map[string]string) *response {
	for k, v := range header {
		resp.WithHeader(k, v)
	}
	return resp
}

func (resp *response) WithHeader(key string, value string) *response {
	resp.header.Set(key, value)
	return resp
}

func (resp *response) WithStatus(status int) *response {
	resp.statusCode = status
	resp.rw.WriteHeader(resp.statusCode)
	return resp
}

func (resp *response) Write(body []byte) (int, error) {
	return resp.rw.Write(body)
}

func (resp *response) Writef(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(resp.rw, format, args...)
}

func (resp *response) prepare() {
}

func (resp *response) statusCodeCheck(status int) bool {
	switch status {
	// is response invalid?
	// @see https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
	case IsInvalid:
		return resp.statusCode < 100 || resp.statusCode >= 600
	// is response informative?
	case IsInformational:
		return resp.statusCode >= 100 && resp.statusCode < 200
	// is response successful?
	case IsSuccessful:
		return resp.statusCode >= 200 && resp.statusCode < 300
	// is the response a redirect?
	case IsRedirection:
		return resp.statusCode >= 300 && resp.statusCode < 400
	// is there a client error?
	case IsClientError:
		return resp.statusCode >= 400 && resp.statusCode < 500
	// was there a server side error?
	case IsServerError:
		return resp.statusCode >= 500 && resp.statusCode < 600
	// is the response OK?
	case IsOk:
		return 200 == resp.statusCode
	// is the response forbidden?
	case IsForbidden:
		return 403 == resp.statusCode
	// is the response a not found error?
	case IsNotFound:
		return 404 == resp.statusCode
	}
	return false
}

// is the response a redirect of some form?
func (resp *response) IsRedirect(location string) bool {
	var flag bool
	for _, code := range []int{
		http.StatusCreated,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	} {
		if code == resp.statusCode {
			flag = true
			break
		}
	}
	return flag && location == resp.header.Get("Location")
}
