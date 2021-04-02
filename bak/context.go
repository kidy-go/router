package router

import (
	"net/http"
	"time"
)

type (
	Context interface {
		StatusCode(int)
		ContentType(string)
		Header() http.Header
		Request() *http.Request
		ResponseWriter() http.ResponseWriter
		Writer(...interface{}) *response
		Response() *response
		Get(string) interface{}
		GetStatusCode() int
		NotFound()

		// 获取客户端上传的文件
		// TODO:
		// 需要返回一个封装后的文件系统对象
		// File(string) FileStore
		// Cookie(string) string
		// Header(string) string
		// Params()
	}

	RequestParams map[string]interface{}
	CookieOptions map[string]interface{}

	context struct {
		request        *http.Request
		responseWriter http.ResponseWriter
		params         RequestParams
		response       *response
		charset        string
	}
)

const (
	ContentTypeStream       = "application/octet-stream"
	ContentTypeWasm         = "application/wasm"
	ContentTypeJSON         = "application/json"
	ContentTypeJSONP        = "application/problem+json"
	ContentTypeXML          = "application/xml"
	ContentTypeXMLP         = "application/problem+xml"
	ContentTypeYAML         = "application/x-yaml"
	ContentTypeProtobuf     = "application/x-protobuf"
	ContentTypeMsgpack      = "application/msgpack"
	ContentTypeXMsgpack     = "application/x-msgpack"
	ContentTypeGRPC         = "application/grpc"
	ContentTypeRPC          = "application/rpc"
	ContentTypeForm         = "application/x-www-form-urlencoded"
	ContentTypeFormMutipart = "multipart/form-data"
	ContentTypeHTML         = "text/html"
	ContentTypeJS           = "text/javascript"
	ContentTypeText         = "text/plain"
	ContentTypeMarkdown     = "text/markdown"
	ContentTypeXMLText      = "text/xml"
	ContentTypeYAMLText     = "text/yaml"
)

func NewContext(req *http.Request, w http.ResponseWriter) *context {
	return &context{
		request:        req,
		responseWriter: w,
		charset:        "utf-8",
	}
}

func (c *context) StatusCode(code int) {
	c.responseWriter.WriteHeader(code)
}

func (c *context) ContentType(cType string) {
	c.responseWriter.Header().Set("Content-Type", cType)
}

func (c *context) Header() http.Header {
	return c.request.Header
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *context) Writer(options ...interface{}) *response {
	var (
		body   interface{}
		status = 0
	)
	if len(options) > 0 {
		for _, o := range options {
			switch v := o.(type) {
			case int:
				if status > 0 {
					body = v
					break
				}
				status = v
			default:
				body = v
			}
		}
	}
	c.response = NewResponse(body, status, c.responseWriter)
	c.response.SetCharset(c.charset)
	return c.response
}

func (c *context) Response() *response {
	if c.response == nil {
		return c.Writer()
	}
	return c.response
}

func (c *context) Get(key string) interface{} {
	return key
}

func (c *context) GetStatusCode() int {
	return 200
}

func (c *context) NotFound() {
	c.responseWriter.WriteHeader(404)
	c.responseWriter.Write([]byte(`Not Found!`))
}

func (co CookieOptions) Path() (string, bool) {
	if v, ok := co["Path"]; ok {
		return v.(string), ok
	}
	return "", false
}
func (co CookieOptions) Domain() (string, bool) {
	if v, ok := co["Domain"]; ok {
		return v.(string), ok
	}
	return "", false
}
func (co CookieOptions) Expires() (time.Time, bool) {
	if v, ok := co["Expires"]; ok {
		return v.(time.Time), ok
	}
	return time.Time{}, false
}
func (co CookieOptions) MaxAge() (int, bool) {
	if v, ok := co["MaxAge"]; ok {
		return v.(int), ok
	}
	return 0, false
}
func (co CookieOptions) Secure() (bool, bool) {
	if v, ok := co["Secure"]; ok {
		return v.(bool), ok
	}
	return false, false
}
func (co CookieOptions) HttpOnly() (bool, bool) {
	if v, ok := co["HttpOnly"]; ok {
		return v.(bool), ok
	}
	return false, false
}
