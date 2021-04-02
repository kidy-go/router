// response.go kee > 2021/02/17

package context

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"net/http"
	"time"
)

type response struct {
	statusCode  int
	contentType string
	body        []byte
	raw         interface{}
	w           http.ResponseWriter
}

func NewResponse(v interface{}, code int, w http.ResponseWriter) *response {
	return &response{
		statusCode: code,
		raw:        v,
		w:          w,
	}
}

func (r *response) SetStatusCode(code int) {
	r.statusCode = code
}
func (r *response) SetContentType(ctype string) {
	r.contentType = ctype
	r.Header("Content-Type", r.contentType)
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

func (r *response) Cookie(name, value string, options map[string]interface{}) *response {
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		Domain: "",
	}

	if v, ok := options["Path"]; ok {
		cookie.Path = v.(string)
	}
	if v, ok := options["Domain"]; ok {
		cookie.Domain = v.(string)
	}
	if v, ok := options["Expires"]; ok {
		cookie.Expires = v.(time.Time)
	}
	if v, ok := options["MaxAge"]; ok {
		cookie.MaxAge = v.(int)
	}
	if v, ok := options["Secure"]; ok {
		cookie.Secure = v.(bool)
	}
	if v, ok := options["HttpOnly"]; ok {
		cookie.HttpOnly = v.(bool)
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
		r.body, _ = json.Marshal(b)
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
		r.body, _ = yaml.Marshal(b)
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
		r.body, _ = yaml.Marshal(b)
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
		r.body, _ = xml.Marshal(b)
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
	if r.contentType == "" {
		r.SetContentType(ContentTypeText)
	}
	r.w.WriteHeader(r.statusCode)
	return r.w.Write(r.body)
}
