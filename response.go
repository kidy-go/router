// response.go kee > 2021/04/08

package router

import (
	"fmt"
	"github.com/kidy-go/utils"
	//"io"
	"encoding/json"
	"encoding/xml"
	yaml "gopkg.in/yaml.v2"
	"net/http"
	"reflect"
	"strings"
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
	req        *http.Request
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

type HttpError interface {
	Error() string
}

func NewResponse(req *http.Request, rw http.ResponseWriter) *response {
	return &response{
		rw:         rw,
		req:        req,
		statusCode: 200,
	}
}

func (resp *response) WithBody(body []byte) *response {
	resp.body = body
	return resp
}

func (resp *response) WithResult(result []reflect.Value) *response {
	if len(result) > 0 {
		for _, res := range result {
			if !res.IsValid() {
				continue
			}

			cType := resp.Header().Get("Content-Type")
			switch v := res.Interface().(type) {
			case bool:
				if !v {
					resp.WithStatus(http.StatusNotFound)
					continue
				}
			case int:
				resp.WithStatus(v)
			case string:
				if cType == "" && strings.Index(v, "/") > 0 {
					resp.WithHeader("Content-Type", cType)
				} else {
					resp.WithBody([]byte(v))
				}
			case []byte:
				resp.WithBody(v)
			case HttpError:
				if v == nil || res.IsNil() {
					continue
				}
				if resp.statusCode == 0 || resp.statusCode < 400 {
					resp.WithStatus(http.StatusBadRequest)
				}

				if len(resp.body) == 0 {
					resp.WithBody([]byte(v.Error()))
				}
				break
			default:
				switch reflect.TypeOf(res).Kind() {
				case reflect.Struct, reflect.Array, reflect.Map:
					resp.JSON(v)
				}
			}
		}
	}
	return resp
}

func (resp *response) WithHeaders(header map[string]string) *response {
	for k, v := range header {
		resp.WithHeader(k, v)
	}
	return resp
}

func (resp *response) WithHeader(key string, value string) *response {
	resp.rw.Header().Set(key, value)
	return resp
}

func (resp *response) Header() http.Header {
	return resp.rw.Header()
}

func (resp *response) WithStatus(status int) *response {
	resp.statusCode = status
	resp.rw.WriteHeader(resp.statusCode)
	return resp
}

func (resp *response) SetContentType(contentType string) *response {
	resp.WithHeader("Content-Type", contentType)
	return resp
}

func (resp *response) JSON(v interface{}) *response {
	resp.WithHeader("Content-Type", ContentTypeJSON)
	switch b := v.(type) {
	case string:
		resp.WithBody([]byte(b))
	case []byte:
		resp.WithBody(b)
	default:
		var (
			err error
			jbt []byte
		)
		jbt, err = json.Marshal(b)
		if err != nil {
			resp.WithStatus(http.StatusInternalServerError)
			resp.WithBody([]byte(err.Error()))
			break
		}
		resp.WithBody(jbt)
	}
	return resp
}

func (resp *response) JSONP(callback string, v interface{}) *response {
	resp.WithHeader("Content-Type", ContentTypeJS)
	switch b := v.(type) {
	case string:
		resp.WithBody([]byte(b))
	case []byte:
		resp.WithBody(b)
	default:
		var (
			err error
			jbt []byte
		)
		jbt, err = json.Marshal(b)
		if err != nil {
			resp.WithStatus(http.StatusInternalServerError)
			resp.WithBody([]byte(err.Error()))
			break
		}
		s := callback + `(` + string(jbt) + `)`
		resp.WithBody([]byte(s))
	}
	return resp
}

func (resp *response) YAML(v interface{}) *response {
	resp.WithHeader("Content-Type", ContentTypeYAML)
	switch b := v.(type) {
	case string:
		resp.WithBody([]byte(b))
	case []byte:
		resp.WithBody(b)
	default:
		var (
			err error
			jbt []byte
		)
		jbt, err = yaml.Marshal(b)
		if err != nil {
			resp.WithStatus(http.StatusInternalServerError)
			resp.WithBody([]byte(err.Error()))
			break
		}
		resp.WithBody(jbt)
	}
	return resp
}

func (resp *response) YAMLText(v interface{}) *response {
	resp.WithHeader("Content-Type", ContentTypeYAMLText)
	switch b := v.(type) {
	case string:
		resp.WithBody([]byte(b))
	case []byte:
		resp.WithBody(b)
	default:
		var (
			err error
			jbt []byte
		)
		jbt, err = yaml.Marshal(b)
		if err != nil {
			resp.WithStatus(http.StatusInternalServerError)
			resp.WithBody([]byte(err.Error()))
			break
		}
		resp.WithBody(jbt)
	}
	return resp
}

func (resp *response) XML(v interface{}) *response {
	resp.WithHeader("Content-Type", ContentTypeXML)
	switch b := v.(type) {
	case string:
		resp.WithBody([]byte(b))
	case []byte:
		resp.WithBody(b)
	default:
		var (
			err error
			jbt []byte
		)
		jbt, err = xml.Marshal(b)
		if err != nil {
			resp.WithStatus(http.StatusInternalServerError)
			resp.WithBody([]byte(err.Error()))
			break
		}
		resp.WithBody(jbt)
	}
	return resp
}

func (resp *response) Stream(v []byte) (int, error) {
	resp.WithHeader("Content-Type", ContentTypeStream)
	return resp.Write(v)
}

func (resp *response) Write(body []byte) (int, error) {
	return resp.rw.Write(body)
}

func (resp *response) Writef(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(resp.rw, format, args...)
}

func (resp *response) Writer() http.ResponseWriter {
	return resp.rw
}

func (resp *response) Sended() (int, error) {
	//VarDump(resp.statusCode, resp.body)
	return resp.Write(resp.body)
}

func (resp *response) prepare() {
	charset := utils.HasOr(resp.charset, "GBK").(string)
	cType := resp.Header().Get("Content-Type")

	if "" == cType {
		resp.WithHeader("Content-Type", ContentTypeHTML+"; charset="+charset)
	} else if 0 == strings.Index(cType, "text/") && -1 == strings.Index(cType, "charset") {
		resp.WithHeader("Content-Type", cType+"; charset="+charset)
	}

	if "" == resp.Header().Get("Transfer-Encoding") {
		resp.Header().Del("Content-Length")
	}
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
	return flag && location == resp.Header().Get("Location")
}
