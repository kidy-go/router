// test/router_test.go kee > 2021/04/20

package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type request struct {
	r *http.Request
	n string
}

func newRequest(method, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	return req
}

func TestRouter(t *testing.T) {
	var tests = []request{
		request{
			n: "get:/prefix/1/content",
			r: newRequest("GET", "/prefix/1/content"),
		},
		request{
			n: "get:/power/kee",
			r: newRequest("GET", "/power/kee"),
		},
	}

	for _, r := range tests {
		w := httptest.NewRecorder()
		ru.ServeHTTP(w, r.r)
		fmt.Println(r.n)
		fmt.Println(w.Body)
		fmt.Println("--------------------------------------------------------------")
	}

}
