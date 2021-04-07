// trie.go kee > 2021/03/20

package router

import (
	"fmt"
	"regexp"
	"strings"
)

type node struct {
	parttern string
	part     string
	children []*node
	dynamic  bool
}

func (t *node) Insert(parts []string, height int) *node {
	if len(parts) == height {
		t.parttern = strings.Join(parts, "/")
		return t
	}

	part := parts[height]
	child := t.matchChild(part)
	if child == nil {
		child = &node{
			part:    part,
			dynamic: strings.Index(part, "{") > -1,
		}
		t.children = append(t.children, child)
	}
	return child.Insert(parts, height+1)
}

func (t *node) matchChild(part string) *node {
	for _, child := range t.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

func (t *node) matchChildren(part string) []*node {
	var children []*node

	for _, child := range t.children {
		if child.part == part || (child.dynamic && child.dynamicMatch(part)) {
			children = append(children, child)
		}
	}
	return children
}

func (t *node) Search(parts []string, height int) *node {
	if len(parts) == height || t.dynamic {
		if t.parttern == "" {
			return nil
		}
		return t
	}

	part := parts[height]
	children := t.matchChildren(part)
	for _, child := range children {
		if result := child.Search(parts, height+1); result != nil {
			return result
		}
	}
	return nil
}

func (t *node) dynamicMatch(part string) bool {
	if t.dynamic {
		match := t.part[strings.Index(t.part, "{"):strings.Index(t.part, "}")]

		typ := "sting"
		if i := strings.Index(match, ":"); i > -1 {
			typ = match[i+1:]
		}

		rp := `([^/]+)`
		switch typ {
		case "string":
			rp = `([^/]+)`
		case "uint", "ulong":
			rp = `(\d)+`
		case "int":
			rp = `[-]?(\d){1,10}`
		case "long":
			rp = `[-]?(\d){0,19}`
		case "number":
			rp = `[-]?(\d+)`
		case "bool":
			rp = `(1|t|T|TRUE|true|True|0|f|F|FALSE|false|False)?`
		case "path":
			rp = `((\w+\/?)+)`
		default:
			if len(typ) > 6 && typ[:6] == "regexp" {
				rp = typ[6:]
			}
		}
		return (regexp.MustCompile(`^` + rp)).MatchString(part)
	}

	return false
}

func echo(v ...interface{}) {
	fmt.Println(v...)
}
