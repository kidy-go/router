package router

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"net/http"
	"reflect"
	"strings"
)

type (
	HttpError interface {
		Error() string
	}
	response struct {
		statusCode  int
		contentType string
		charset     string
		body        []byte
		raw         interface{}
		w           http.ResponseWriter
	}
)

// 解析路由函数返回值, 构造response(statusCode, contentType, content)
// 支持多返回结果:
// int = StatusCode
// string = ContentType | response body
// error = bad request
// []byte = response body
// struct | map | array = Result
func dispatchResult(result []reflect.Value, ctx Context) (err error) {
	var (
		statusCode  = 0
		contentType string
		content     interface{}
	)

	if len(result) > 0 {
		for _, val := range result {
			if !val.IsValid() {
				continue
			}
			switch v := val.Interface().(type) {
			case bool:
				if !v {
					statusCode = StatusCodeNotFound
					continue
				}
			case int:
				if statusCode == 0 {
					statusCode = v
				}
			case string:
				if contentType == "" && strings.Index(v, "/") > 0 {
					contentType = v
				} else {
					content = []byte(v)
				}
			case []byte:
				content = v
			case HttpError:
				if v == nil || val.IsNil() {
					continue
				}
				if statusCode == 0 || statusCode < 400 {
					statusCode = StatusCodeBadRequest
				}
				if content == nil {
					content = []byte(v.Error())
				}
				err = v
				break
			default:
				// if val.Kind() == reflect.Ptr && content == nil && v != nil {
				//	content = v
				// }
				content = v
			}
		}
		w := ctx.Writer(content, statusCode)
		if contentType != "" {
			w.SetContentType(contentType)
		}
	}
	_, err = ctx.Response().ResponseWrite()
	return err
}

func NewResponse(v interface{}, code int, w http.ResponseWriter) *response {
	r := &response{
		statusCode: code,
		raw:        v,
		w:          w,
	}
	switch s := v.(type) {
	case string:
		r.body = []byte(s)
	case []byte:
		r.body = s
	default:
		r.raw = s
	}
	return r
}

func (r *response) SetStatusCode(code int) {
	r.statusCode = code
}
func (r *response) SetCharset(charset string) {
	r.charset = charset
}
func (r *response) SetContentType(ctype string, charset ...string) {
	if len(charset) > 0 {
		r.charset = charset[0]
	}
	r.contentType = ctype
}

func (r *response) Header(key, val string) *response {
	if v := r.w.Header().Get(key); v != "" {
		r.w.Header().Add(key, val)
		return r
	}
	r.w.Header().Set(key, val)
	return r
}

func (r *response) WithHeaders(headers http.Header) *response {
	for k, h := range headers {
		for _, v := range h {
			r.Header(k, v)
		}
	}
	return r
}

func (r *response) Cookie(name, value string, options CookieOptions) *response {
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		Domain: "",
	}

	if v, ok := options.Path(); ok {
		cookie.Path = v
	}
	if v, ok := options.Domain(); ok {
		cookie.Domain = v
	}
	if v, ok := options.Expires(); ok {
		cookie.Expires = v
	}
	if v, ok := options.MaxAge(); ok {
		cookie.MaxAge = v
	}
	if v, ok := options.Secure(); ok {
		cookie.Secure = v
	}
	if v, ok := options.HttpOnly(); ok {
		cookie.HttpOnly = v
	}

	r.WithCookie(cookie)
	return r
}

func (r *response) WithCookie(cookie *http.Cookie) *response {
	http.SetCookie(r.w, cookie)
	return r
}

func (r *response) JSON(v interface{}) *response {
	r.SetContentType(ContentTypeJSON)
	switch b := v.(type) {
	case string:
		r.body = []byte(b)
	case []byte:
		r.body = b
	default:
		var err error
		r.body, err = json.Marshal(b)
		if err != nil {
			panic(err)
		}
	}
	return r
}

func (r *response) JSONP(callback string, v interface{}) *response {
	r.SetContentType(ContentTypeJS)
	switch b := v.(type) {
	case string:
		r.body = []byte(b)
	case []byte:
		r.body = b
	default:
		var err error
		r.body, err = json.Marshal(b)

		if err != nil {
			panic(err)
		}
		r.body = append([]byte(callback+`(`), r.body...)
		r.body = append(r.body, []byte(`)`)...)
	}
	return r
}

func (r *response) YAML(v interface{}) *response {
	r.SetContentType(ContentTypeYAML)
	switch b := v.(type) {
	case string:
		r.body = []byte(b)
	case []byte:
		r.body = b
	default:
		var err error
		r.body, err = yaml.Marshal(b)

		if err != nil {
			panic(err)
		}
	}
	return r
}

func (r *response) YAMLText(v interface{}) *response {
	r.SetContentType(ContentTypeYAMLText)
	switch b := v.(type) {
	case string:
		r.body = []byte(b)
	case []byte:
		r.body = b
	default:
		var err error
		r.body, err = yaml.Marshal(b)

		if err != nil {
			panic(err)
		}
	}
	return r
}

func (r *response) XML(v interface{}) *response {
	r.SetContentType(ContentTypeXML)
	switch b := v.(type) {
	case string:
		r.body = []byte(b)
	case []byte:
		r.body = b
	default:
		var err error
		r.body, err = xml.Marshal(b)
		if err != nil {
			panic(err)
		}
		r.body = append([]byte(xml.Header), r.body...)
	}
	return r
}

func (r *response) Stream(v []byte) (int, error) {
	r.SetContentType(ContentTypeStream)
	r.body = v
	return r.ResponseWrite()
}

func (r *response) Write(v []byte) (int, error) {
	r.body = v
	return r.ResponseWrite()
}

func (r *response) WriteString(v string) (int, error) {
	r.body = []byte(v)
	return r.ResponseWrite()
}

func (r *response) Writef(format string, args ...interface{}) (int, error) {
	r.body = []byte(fmt.Sprintf(format, args...))
	return r.ResponseWrite()
}

func (r *response) WriteResponse(v []byte, cType string) (int, error) {
	r.SetContentType(cType)
	r.body = v
	return r.ResponseWrite()
}

func (r *response) ResponseWrite() (int, error) {
	if r.statusCode == 0 {
		r.SetStatusCode(StatusCodeOK)
	}

	if len(r.body) == 0 && r.raw != nil {
		switch r.contentType {
		case ContentTypeJSON:
			r.JSON(r.raw)
		case ContentTypeYAML:
			r.YAML(r.raw)
		case ContentTypeYAMLText:
			r.YAMLText(r.raw)
		case ContentTypeXML:
			r.XML(r.raw)
		case "":
			r.JSON(r.raw)
		default:
		}
	}

	if r.contentType == "" {
		r.SetContentType(ContentTypeText)
	}

	cType := r.contentType + "; charset=" + r.charset
	r.Header("Content-Type", cType)
	r.w.WriteHeader(r.statusCode)
	return r.w.Write(r.body)
}
