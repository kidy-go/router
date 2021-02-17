// router_test.go kee > 2021/02/08

package router

import (
	"fmt"
	"net/http"
	"time"
	//"regexp"
	//"strings"
	"encoding/json"
	"encoding/xml"
	//yaml "gopkg.in/yaml.v2"
	"testing"
)

type J struct {
	Name   string `json:"name"`
	Sex    int    `json:"sex"`
	Age    int    `json:"age"`
	Status int    `json:"status"`
}

type Index struct {
}

func (c *Index) Index(ctx Context) {
	resp := ctx.Writer(&J{
		Name:   "笑傲江湖",
		Status: 200,
		Age:    29,
		Sex:    1,
	})
	//resp := ctx.Writer("Hello World")
	resp.Cookie("x-ssr", "ertyuioiyt567uijhbv==", map[string]interface{}{
		"Expires": time.Time{},
		"MaxAge":  3600,
	})
}

func (c *Index) PostHeader(ctx Context) http.Header {
	ctx.Writer().Cookie("x-sword", "S.O.A.", CookieOptions{
		"MaxAge": 3600,
	})
	return ctx.Request().Header
}

func (c *Index) PutMIMeType() {
}

func (c *Index) Post() {

}

func (c *Index) GetMy() (int, string) {
	return 200, `My info`
}

func (c *Index) GetBy(id int) (string, []string) {
	return ContentTypeJSON, []string{
		"A", "B", "C",
	}
}

func (c *Index) GetByInfo(id bool) (bool, string) {
	return id, fmt.Sprintf("Status: %v", id)
}

func (c *Index) GetByInfoBy(id, vid int) (string, J) {
	return ContentTypeXML, J{
		Name:   "笑傲江湖曲",
		Sex:    1,
		Age:    29,
		Status: 2021,
	}
}

func (c *Index) GetInBy(id int) (error, string) {
	if id%2 > 0 {
		return fmt.Errorf(`{"msg": "Has Error"}`), ContentTypeJSON
	}
	return nil, ContentTypeText
}

func (c *Index) GetCodeBy(ctx Context, code int) (int, string) {
	ctx.ContentType(ContentTypeXML)

	xmlByte, _ := xml.Marshal(&J{Name: "Kee", Status: code})
	xmlByte = append([]byte(xml.Header), xmlByte...)

	jsonByte, _ := json.Marshal([]string{
		"Hello", "World",
	})

	A := []interface{}{
		"Hello World",
		[]byte("keesely.net"),
		&J{Name: "Kee", Status: code},
		string(xmlByte),
		string(jsonByte),
		[]string{"Hello", "Kee"},
	}

	for _, s := range A {
		switch v := s.(type) {
		case string, []byte:
			fmt.Println("SB", v)
		case interface{}:
			fmt.Println("Interface", v)
		}
	}

	return code, string(xmlByte)
	//return code, fmt.Sprintf(`{"Status Code": %d}`, code)
}

func TestRouter(t *testing.T) {
	s := &http.Server{Addr: ":89"}

	router := NewRouter()

	router.Handler("/test", func(ctx Context) {
		ctx.Writer(map[string]string{
			"code": "200",
			"msg":  "successful",
		})
	})

	s.Handler = router

	router.Handler("/", new(Index))

	// paths := []string{
	// 	"/w/{id:int}",
	// 	"/w/{name:string}",
	// 	"/w/{id:int}/info/{format:string}",
	// 	"/w/{id:int}/info/{format:string}/xx",
	// 	"/w/{id:int}/{format:string}/xx",
	// }

	// cpath := []string{
	// 	"/w/1101/info/中==对/xx",
	// 	"/w/1211",
	// 	"/w/2011/json/xx",
	// }

	// for _, path := range paths {
	// 	regx := regexp.MustCompile(`\{(.*?)(:.*?)?(\{[0-9,]+\})?\}`)
	// 	r := regx.ReplaceAllStringFunc(path, func(match string) string {
	// 		match = match[1 : len(match)-1]
	// 		ss := strings.Split(match, ":")
	// 		switch len(ss) {
	// 		case 3:
	// 		case 2:
	// 			switch ss[1] {
	// 			case "string":
	// 				return `([^/]+)`
	// 			case "int":
	// 				return `(\d+)`
	// 			}
	// 		case 1:
	// 			return `([^/]+)`
	// 		}
	// 		return match
	// 	})
	// 	r = `^(?U)` + r + `\z`
	// 	regx = regexp.MustCompile(r)
	// 	for _, p := range cpath {
	// 		if regx.MatchString(p) {
	// 			fmt.Println(p, path)
	// 		}
	// 	}
	// }

	s.ListenAndServe()
}
