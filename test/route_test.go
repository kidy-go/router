// route_test.go kee > 2021/03/16

package router

import (
	"fmt"
	//"github.com/kidy-go/utils/test"
	//"net"
	"testing"
)

func TestRoute(t *testing.T) {

	return
	context := &context{}
	routes := []Route{}

	routes = append(routes, NewRoute().Handle("/", new(rHandler)))
	routes = append(routes, NewRoute().Handle("/i8/{:uint}", demoFunc))

	type T struct {
		method string
		path   string
		expect string
	}

	tRoutes := []T{
		T{"get", "/index", "Hi index"},
		T{"post", "/index", "Hi index"},
		T{"get", "/", "Map"},
		T{"post", "/", "POST Successful."},
		T{"get", "/F", "get successful."},
		T{"put", "/power/kee", "hello kee"},
		T{"get", "/i8/123", "123"},
	}

	for _, ts := range tRoutes {
		fmt.Println(ts)
		for _, route := range routes {
			if ok, r := route.Search(ts.path); ok && r.HasMethod(ts.method) {
				fmt.Println("> CHECKOUT >>", r.methods, r.path, r.params)
				r.dispatch(context)
				fmt.Println(`>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>`)
			}
		}
	}
}
