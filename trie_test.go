// trie_test.go kee > 2021/03/20

package router

import (
	"testing"
)

// func TestTrie(t *testing.T) {
// 	paths := []string{
// 		"GET/",
// 		"GET/home/{id}",
// 		"GET/home",
// 		"GET/auth/login",
// 		"POST/auth/login",
// 		"PATCH/ticket/{id}",
// 		"DELETE/ticket/{id}",
// 		"POST/ticket",
// 		"GET/auth/login/ticket",
// 		`DELETE/auth/login/{id:regexp(\d+)}`,
// 		"GET/hello/admin",
// 		"GET/hello/{id:int}",
// 		"GET/ticket",
// 		"GET/ticket/{id}/{name}",
// 		"PUT/ticket/{id}",
// 	}
//
// 	var routes []*Route
// 	for _, path := range paths {
// 		index := strings.Index(path, "/")
// 		method, uri := path[:index], path[index:]
// 		routes = append(routes, &Route{
// 			methods: []string{method},
// 			uri:     uri,
// 		})
// 	}
//
// 	matchPaths := []string{
// 		"GET/home/121",
// 		"PUT/ticket/222",
// 		"GET/home",
// 		"GET/ticket",
// 		"GET/",
// 		"GET/hello/admin",
// 		"GET/hello/2011",
// 		"DELETE/auth/login/10110",
// 		"GET/ticket/112345678/kee",
// 	}
//
// 	for _, path := range matchPaths {
// 		fmt.Println("Define> ", path)
// 		for _, route := range routes {
// 			if route.Match(path) {
// 				VarDump(">>> ", route.params)
// 				return
// 			}
// 		}
// 		fmt.Println("--------------------------------------------------------------")
// 	}
// }

func TestRouter(t *testing.T) {
	//h := new(rHandler)
	var methods []string
	methods = append(methods, "GET", "POST", "HEAD", "OPTIONS", "PUT", "PATCH", "DELETE", "TRACE", "CONNECT")
	router := NewRouter()
	router.AddRoute("get", "/i8u", demoFunc)
	router.AddRoute("Get", "/i8u/power", (&rHandler{}).Index)
	router.Group(GroupStack{
		prefix: "/prefix",
	}, func(router *Router) {
		router.Handle("/home", new(rHandler))
		router.AddRoute("GET", "/rhome", "rHandler@Index")
	})

	for _, m := range methods {
		for _, r := range router.routes[m] {
			VarDump(r.methods, r.uri, r.typ)
		}
	}
}
