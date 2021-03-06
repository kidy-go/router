// demo_handler.go kee > 2021/03/18

package router

import (
	"fmt"
)

type ModResult struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Sex  string `json:"sex"`
	Age  int    `json:"age"`
}

type rHandler struct{}

func (c *rHandler) Index() (int, string) {
	return 403, "Hi index"
}

func (c *rHandler) Get(ctx Context) map[string]string {
	return map[string]string{
		"path":   ctx.Request().URL.Path,
		"method": ctx.Request().Method,
	}
}

func (c *rHandler) GetByContent(id int) *ModResult {
	return &ModResult{
		Name: "Kee",
		Sex:  "Men",
		Id:   id,
		Age:  29,
	}
}

func (c *rHandler) Post() string {
	return "POST Successful."
}

func (c *rHandler) GetBy(b bool) (bool, string) {
	return b, "get successful."
}

func (c *rHandler) PutBy(id int) {}

func (c *rHandler) PatchBy(id int) {}

func (c *rHandler) DeleteBy(id int) {}

func (c *rHandler) PutPowerBy(ctx Context, name string) string {
	return "hello " + name
}

func demoFunc(i8 uint8) string {
	return fmt.Sprintf("%d", i8)
}
