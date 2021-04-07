// route.go kee > 2021/03/23

package router

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"regexp"
	"strings"
)

type Route struct {
	methods []string
	uri     string
	typ     reflect.Type
	val     reflect.Value
	handler interface{}
	ctx     Context
	uses    []interface{}
	partNum int
	params  []string
}

type Handler func(Context)

func (r *Route) Match(parttern string) bool {
	var params []string

	mIndex := strings.Index(parttern, "/")
	method := parttern[:mIndex]
	if !r.HasMethod(method) {
		return false
	}

	keys, rp := regexMatch(method + r.uri)
	regx := regexp.MustCompile(rp)
	if regx.MatchString(parttern) {
		matchs := regx.FindAllStringSubmatch(parttern, -1)
		for i, val := range matchs[0][1:] {
			val = keys[i] + ":" + val
			params = append(params, val)
		}
		r.params = params
		return true
	}
	return false
}

func (r *Route) HasMethod(method string) bool {
	for _, m := range r.methods {
		if m == method {
			return true
		}
	}
	return false
}

func (r *Route) dispatch() []reflect.Value {
	typ := r.typ
	cType := reflect.TypeOf((*Context)(nil)).Elem()
	var params []reflect.Value
	for i, j := 0, 0; i < typ.NumIn(); i++ {
		if typ.In(i) == reflect.TypeOf(r.handler) {
			continue
		}
		if typ.In(i).Implements(cType) {
			params = append(params, reflect.ValueOf(r.ctx))
		} else {
			v := r.params[j]
			pv := strings.Split(r.params[j], ":")
			if len(pv) > 2 {
				v = strings.Join(pv[2:], ":")
			}

			var val interface{}
			switch typ.In(i).Kind() {
			case reflect.String:
				val = v
			case reflect.Int:
				val = cast.ToInt(v)
			case reflect.Int64:
				val = cast.ToInt64(v)
			case reflect.Int32:
				val = cast.ToInt32(v)
			case reflect.Uint8:
				val = cast.ToUint8(v)
			case reflect.Uint32:
				val = cast.ToUint32(v)
			case reflect.Uint64:
				val = cast.ToUint64(v)
			case reflect.Bool:
				b := strings.ToLower(v)
				val = false
				if b == "1" || b == "t" || b == "true" {
					val = true
				}
			default:
				val = v
			}
			if val != nil {
				params = append(params, reflect.ValueOf(v))
				j++
			}
		}
	}
	//ctx.params = RequestParams{params}
	result := r.val.Call(params)
	return result
}

// 参数形式:
// {param}
// {param:string}
// {param:regexp([a-zA-Z0-9]+)}
// {:string}
// {:regexp(\d+)}
//
// 预定义参数格式:
// [string]: 字符串类型, 遇到/停止
// [int]:    数值类型 (取值范围: -2147483648 ~ 2147483647)
// [long]:   长整型数值 (取值范围: -9223372036854775808 ~ 9223372036854775807)
// [bool]:	 布尔值类型(取值: 1 、 t 、 T 、 TRUE 、 true 、 True 、 0 、 f 、 F 、 FALSE 、 false 、 False)
// [number]: 数字类型, 不限取值范围, 将会以字符串类型给出
// [path]:	 路径地址, 例: /home/www or /home/www/1.txt
func regexMatch(part string) (keys []string, rp string) {
	regx := regexp.MustCompile(`\{(.*?)\}`)
	part = regx.ReplaceAllStringFunc(part, func(match string) string {
		match = match[1 : len(match)-1]
		key, typ := match, "string"
		if i := strings.Index(match, ":"); i > -1 {
			key, typ = match[:i], match[i+1:]
		}
		keys = append(keys, key+":"+typ)

		switch typ {
		case "string":
			return `([^/]+)`
		case "uint", "ulong":
			return `(\d+)`
		case "int":
			return `([-]?\d{1,10})`
		case "long":
			return `([-]?\d{0,19})`
		case "number":
			return `([-]?(\d+))`
		case "bool":
			return `(1|t|T|TRUE|true|True|0|f|F|FALSE|false|False)?`
		case "path":
			return `((\w+\/?)+)`
		default:
			if len(typ) > 6 && typ[:6] == "regexp" {
				return typ[6:]
			}
			return `([^/]+)`
		}
		return match
	})
	return keys, `^(?U)` + part + `$`
}

func VarDump(v ...interface{}) {
	fmt.Println(v...)
}
