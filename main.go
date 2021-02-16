package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

const tempalteFile = `
package gen

import "fmt"

`

const templateFunc = `
func (s *%s)Hello(){
	fmt.Println("hello world")
}
`

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "target.go", nil, parser.Mode(0))

	var structs []string

	// structの取得
	ast.Inspect(f, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				structName := typeSpec.Name.Name
				structs = append(structs, structName)
			}
		}
		return true
	})

	// コード生成
	var builder strings.Builder
	builder.WriteString(tempalteFile)
	for _, s := range structs {
		builder.WriteString(fmt.Sprintf(templateFunc, s))
	}
	// ソースコードとして整形等
	src, err := format.Source([]byte(builder.String()))
	if err != nil {
		fmt.Println(err)
		return
	}
	// 書き込み
	err = ioutil.WriteFile("./out/output.go", src, 0664)
	if err != nil {
		fmt.Println(err)
	}
}
