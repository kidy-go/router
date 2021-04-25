// router.go kee > 2021/03/20

package router

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type Router struct {
	routes     map[string][]*Route
	groupStack GroupStack
}

type GroupStack struct {
	Prefix string
	Suffix string
	Domain string
	Uses   []interface{}
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string][]*Route),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, w)

	var (
		result []reflect.Value
		ok     bool
	)
	if ok, result = r.dispatch(ctx); !ok {
		result = []reflect.Value{
			reflect.ValueOf(http.StatusNotFound),
		}
	}
	ctx.response.WithResult(result)
	ctx.response.Sended()
	//display(result, ctx)
}

func (r *Router) Handle(uri string, handler interface{}) {
	typ, val := reflect.TypeOf(handler), reflect.ValueOf(handler)

	if r.groupStack.Prefix != "" {
		uri = r.groupStack.Prefix + uri
	}
	if typ.Kind() == reflect.Ptr {
		for i := 0; i < typ.NumMethod(); i++ {
			tM := typ.Method(i)
			m, path := parsePtr(tM.Name, typ.Method(i).Type)
			path = uri + path
			if len(path) > 1 && path[len(path)-1:] == "/" {
				path = path[:len(path)-1]
			}
			path = strings.Replace(path, "//", "/", -1)
			for _, method := range m {
				if nil == r.routes[method] {
					r.routes[method] = []*Route{}
				}
				r.routes[method] = append(r.routes[method], &Route{
					uri:     path,
					handler: handler,
					typ:     typ.Method(i).Type,
					val:     val.Method(i),
					methods: m,
				})
			}
		}
	}
}

func (r *Router) Get(uri string, handler Handler) *Route {
	return r.AddRoute("GET", uri, handler)
}

func (r *Router) Post(uri string, handler Handler) *Route {
	return r.AddRoute("POST", uri, handler)
}

func (r *Router) Put(uri string, handler Handler) *Route {
	return r.AddRoute("PUT", uri, handler)
}

func (r *Router) Patch(uri string, handler Handler) *Route {
	return r.AddRoute("PATCH", uri, handler)
}

func (r *Router) Delete(uri string, handler Handler) *Route {
	return r.AddRoute("DELETE", uri, handler)
}

func (r *Router) Options(uri string, handler Handler) *Route {
	return r.AddRoute("OPTIONS", uri, handler)
}

func (r *Router) Head(uri string, handler Handler) *Route {
	return r.AddRoute("HEAD", uri, handler)
}

func (r *Router) Trace(uri string, handler Handler) *Route {
	return r.AddRoute("HEAD", uri, handler)
}

func (r *Router) Connect(uri string, handler Handler) *Route {
	return r.AddRoute("CONNECT", uri, handler)
}

func (r *Router) Group(stack GroupStack, gHandler func(*Router)) {
	router := &Router{
		groupStack: stack,
		routes:     make(map[string][]*Route),
	}
	gHandler(router)
	// merge routes
	for m, routes := range router.routes {
		if nil == r.routes[m] {
			r.routes[m] = []*Route{}
		}
		r.routes[m] = append(r.routes[m], routes...)
	}
}

func (r *Router) AddRoute(method, uri string, handler interface{}) *Route {
	var route *Route
	typ, val := reflect.TypeOf(handler), reflect.ValueOf(handler)
	if r.groupStack.Prefix != "" {
		uri = r.groupStack.Prefix + uri
	}
	method = strings.ToUpper(method)
	if typ.Kind() == reflect.Func {
		route = &Route{
			uri:     uri,
			handler: handler,
			typ:     typ,
			val:     val,
			methods: []string{method},
		}
		if nil == r.routes[method] {
			r.routes[method] = []*Route{}
		}
		r.routes[method] = append(r.routes[method], route)
	}
	return route
}

func (r *Router) dispatch(ctx Context) (ok bool, result []reflect.Value) {
	method, uri := ctx.Request().Method, ctx.Request().URL.Path
	for _, route := range r.routes[method] {
		if route.Match(method + uri) {
			route.ctx = ctx
			return true, route.dispatch()
		}
	}
	return false, nil
}

func parsePtr(name string, typ reflect.Type) ([]string, string) {
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
			return "/{:" + v + "}"
		}
		return "/{:string}"
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
	var methods []string
	if method == "ANY" {
		methods = append(methods, "GET", "POST", "HEAD", "OPTIONS", "PUT", "PATCH", "DELETE", "TRACE", "CONNECT")
	} else {
		methods = append(methods, method)
	}

	return methods, name
}
