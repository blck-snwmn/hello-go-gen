package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "target.go", nil, parser.Mode(0))

	ast.Inspect(f, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			fmt.Println(typeSpec.Name.Name)
			if st, ok := typeSpec.Type.(*ast.StructType); ok {
				ast.Print(fset, st.Fields)
			}
		}
		return true
	})
}
