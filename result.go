// result.go kee > 2021/04/07

package router

import (
	"github.com/kidy-go/utils"
	yaml "gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

type HttpError interface {
	Error() string
}

func display(result []reflect.Value, ctx Context) error {
	var (
		statusCode  = 0
		contentType string
		content     []byte
		//err         error
	)

	if len(result) > 0 {
		for _, res := range result {
			if !res.IsValid() {
				continue
			}

			switch v := res.Interface().(type) {
			case bool:
				if !v {
					statusCode = StatusCodeNotFound
					continue
				}
			case int:
				statusCode = utils.HasOr(statusCode == 0, v, statusCode).(int)
			case string:
				if contentType == "" && strings.Index(v, "/") > 0 {
					contentType = v
				} else {
					content = []byte(v)
				}
			case []byte:
				content = v
			case HttpError:
				if v == nil || res.IsNil() {
					continue
				}

				if statusCode == 0 || statusCode < 400 {
					statusCode = StatusCodeBadRequest
				}

				if content == nil {
					content = []byte(v.Error())
				}

				//err = v
				break
			default:
				content, _ = yaml.Marshal(v)
				//content = []byte(v)
			}
		}
	}

	resp := NewResponse(ctx)
	resp.WithStatus(403)
	resp.WithBody(content)
	_, e := resp.Write()
	return e

	ctx.ResponseWriter().Header().Set("Content-Type", contentType)
	if statusCode > 0 {
		ctx.ResponseWriter().WriteHeader(statusCode)
	}
	_, err := ctx.ResponseWriter().Write(content)
	return err
}
