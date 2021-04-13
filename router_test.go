// router_test.go kee > 2021/03/20

package router

import (
	"net/http"
	"testing"
)

func setRouter() *Router {
	router := NewRouter()
	router.AddRoute("get", "/i8u", demoFunc)
	router.AddRoute("Get", "/i8u/power", (&rHandler{}).Index)
	router.Group(GroupStack{
		prefix: "/prefix",
	}, func(router *Router) {
		router.Handle("/", new(rHandler))
		router.AddRoute("GET", "/rhome", "rHandler@Index")
	})
	return router
}

func TestRouter(t *testing.T) {
	router := setRouter()
	s := &http.Server{Addr: ":89"}
	s.Handler = router

	s.ListenAndServe()
}
