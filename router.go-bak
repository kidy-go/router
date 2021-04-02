// router.go kee > 2021/02/08

package router

import (
	"fmt"
	"github.com/spf13/cast"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

const (
	ANY = iota
	GET
	POST
	HEAD
	OPTIONS
	PUT
	PATCH
	DELETE
	TRACE
	CONNECT
)

func MethodCode(method string) int {
	switch strings.ToUpper(method) {
	case "ANY":
		return ANY
	case "GET":
		return GET
	case "POST":
		return POST
	case "HEAD":
		return HEAD
	case "OPTIONS":
		return OPTIONS
	case "PUT":
		return PUT
	case "PATCH":
		return PATCH
	case "DELETE":
		return DELETE
	case "TRACE":
		return TRACE
	case "CONNECT":
		return CONNECT
	default:
		return ANY
	}
}

type Route struct {
	handler  interface{}
	value    reflect.Value
	typeof   reflect.Type
	context  *context
	path     string
	methods  []string
	children map[int][]Route
	params   []string
}

type Router struct {
	routes  map[int][]Route
	context *context
	host    string
}

func NewRouter() *Router {
	return &Router{
		routes: map[int][]Route{},
	}
}

func InitRoutes() map[int][]Route {
	return map[int][]Route{
		ANY:     []Route{},
		GET:     []Route{},
		POST:    []Route{},
		HEAD:    []Route{},
		OPTIONS: []Route{},
		PUT:     []Route{},
		PATCH:   []Route{},
		DELETE:  []Route{},
		TRACE:   []Route{},
		CONNECT: []Route{},
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, w)
	cType := reflect.TypeOf((*Context)(nil)).Elem()

	method := MethodCode(req.Method)
	routes, ok := r.routes[method]
	if ok {
		if any, ok := r.routes[ANY]; ok {
			routes = append(routes, any...)
		}
	}
	ok, route := r.pathMatch(method, req.URL.Path, routes)
	if ok {
		typ := route.typeof
		var params []reflect.Value
		for i, j := 0, 0; i < typ.NumIn(); i++ {
			if typ.In(i).Implements(cType) {
				params = append(params, reflect.ValueOf(ctx))
			} else if typ.In(i) == reflect.TypeOf(route.handler) {
				//params = append(params, reflect.ValueOf(route.handler))
			} else {
				var v interface{}
				switch typ.In(i).Kind() {
				case reflect.String:
					v = route.params[j]
				case reflect.Int:
					v = cast.ToInt(route.params[j])
				case reflect.Int64:
					v = cast.ToInt64(route.params[j])
				case reflect.Int32:
					v = cast.ToInt32(route.params[j])
				case reflect.Bool:
					b := strings.ToLower(route.params[j])
					v = false
					if b == "1" || b == "t" || b == "true" {
						v = true
					}
				default:
					v = route.params[j]
				}
				if v != nil {
					params = append(params, reflect.ValueOf(v))
					j++
				}
			}
		}
		//ctx.params = RequestParams{params}
		result := route.value.Call(params)
		dispatchResult(result, ctx)
	} else {
		ctx.NotFound()
	}
}

func (r *Router) pathMatch(method int, path string, routes []Route) (ok bool, route Route) {
	for _, rot := range routes {
		ok, route = rot.PathMatch(method, path)
		if ok {
			break
		}
	}
	return
}

func (r *Router) Handler(path string, handler interface{}) Route {
	typ := reflect.TypeOf(handler)
	val := reflect.ValueOf(handler)
	route := Route{
		handler:  handler,
		value:    val,
		typeof:   typ,
		path:     path,
		children: map[int][]Route{},
	}
	return r.handler(ANY, route)
}

func (r *Router) handler(method int, route Route) Route {
	typ := route.typeof
	val := route.value
	path := route.path

	switch typ.Kind() {
	case reflect.Func:
	case reflect.Ptr:
		for i := 0; i < typ.NumMethod(); i++ {
			tM := typ.Method(i)
			m, spath := parsePtr(tM.Name, typ.Method(i).Type)
			spath = path + spath
			if len(spath) > 1 && spath[len(spath)-1:] == "/" {
				spath = spath[:len(spath)-1]
			}
			spath = strings.Replace(spath, "//", "/", -1)
			sroute := Route{
				path:     spath,
				handler:  route.handler,
				typeof:   typ.Method(i).Type,
				value:    val.Method(i),
				methods:  []string{m},
				children: map[int][]Route{},
			}
			mcode := MethodCode(m)
			if _, ok := route.children[mcode]; !ok {
				route.children[mcode] = []Route{}
			}
			route.children[mcode] = append(
				route.children[mcode],
				r.handler(mcode, sroute),
			)
		}
	case reflect.String:
	}
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = []Route{}
	}
	r.routes[method] = append(r.routes[method], route)
	return route
}

func parsePtr(name string, typ reflect.Type) (string, string) {
	// method: GET / POST / HEAD / OPTIONS / PUT / PATCH / DELETE / TRACE / CONNECT
	regM := regexp.MustCompile(`^(Get|Post|Head|Options|Put|Patch|Delete|Trace|Connect|Any)?`)
	method := strings.ToUpper(regM.FindString(name))
	name = regM.ReplaceAllString(name, "")

	// Has By
	regBy := regexp.MustCompile("(By)+")

	i := 0
	name = regBy.ReplaceAllStringFunc(name, func(s string) string {
		if i <= typ.NumIn() {
			i += 1
			var v string
			switch typ.In(i).Kind() {
			case reflect.String:
				v = "string"
			case reflect.Int:
				v = "int"
			case reflect.Int64:
				v = "int64"
			case reflect.Int32:
				v = "int32"
			case reflect.Bool:
				v = "bool"
			default:
				v = "string"
			}
			return fmt.Sprintf("/{key%d:%s}", i, v)
		}
		return "/{key:string}"
	})

	regPart := regexp.MustCompile("([A-Z]+)")
	name = regPart.ReplaceAllStringFunc(name, func(s string) string {
		return "/" + strings.ToLower(s)
	})
	if len(name) <= 0 || name[0:1] != "/" {
		name = "/" + name
	}
	if len(method) <= 0 {
		method = "GET"
	}
	return method, name
}

func (r *Router) Group(option map[string]interface{}, handler func(*Router)) {

}

func (r *Router) dispatch(path string, route *Route) {

}

func (r Route) PathMatch(method int, path string) (bool, Route) {
	regx := regexp.MustCompile(`\{(.*?)(:.*?)?(\{[0-9,]+\})?\}`)
	rp := regx.ReplaceAllStringFunc(r.path, func(match string) string {
		match = match[1 : len(match)-1]
		ss := strings.Split(match, ":")
		switch len(ss) {
		case 3:
		case 2:
			switch ss[1] {
			case "string":
				return `([^/]+)`
			case "int", "int64", "int32":
				return `(\d+)`
			case "bool":
				return `(1|t|T|TRUE|true|True|0|f|F|FALSE|false|False)?`
			case "path":
				return `((\w+\/?)+)`
			default:
				return `([^/]+)`
			}
		case 1:
			return `([^/]+)`
		}
		return match
	})
	rp = `^(?U)` + rp + `\z`
	regx = regexp.MustCompile(rp)

	ok := regx.MatchString(path)
	if ok {
		params := regx.FindStringSubmatch(path)
		if len(params) > 1 {
			r.params = params[1:]
		}
		if children, _ := r.children[method]; len(children) > 0 {
			for _, crp := range children {
				crp.params = append(crp.params, r.params...)
				ok, rr := crp.PathMatch(method, path)
				if ok {
					return ok, rr
				}
			}
		} else {
			return ok, r
		}
	}
	return false, Route{}
}
