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
