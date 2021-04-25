// init.go kee > 2021/04/20

package test

import (
	"github.com/kidy-go/router"
)

var ru *router.Router

func init() {
	ru = router.NewRouter()
	ru.AddRoute("get", "/i8u", demoFunc)
	ru.AddRoute("Get", "/power/{name}", (&rHandler{}).PutPowerBy)

	ru.Group(router.GroupStack{
		Prefix: "/prefix",
	}, func(r *router.Router) {
		r.Handle("/", new(rHandler))
		r.AddRoute("GET", "/rhome", "rHandler@Index")
	})
}
