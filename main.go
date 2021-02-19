package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"

	"golang.org/x/xerrors"
)

func main() {
	in := "target.go"
	out := "./out/output.go"
	if err := execute(in, out); err != nil {
		fmt.Printf("failed to execute: %+v\n", err)
	}
}

func execute(in, out string) error {
	result, err := generate("target.go")
	if err != nil {
		return xerrors.Errorf("failed to generate: %w", err)
	}
	if err := write(result, "./out/output.go"); err != nil {
		return xerrors.Errorf("failed to write: %w", err)
	}
	return nil
}

func generate(inputPath string) ([]byte, error) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, inputPath, nil, parser.Mode(0))

	var (
		buf bytes.Buffer
		err error
	)
	fmt.Fprintln(&buf, "package generate")

	fmt.Fprintln(&buf, `import "fmt"`)

	ast.Inspect(f, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				fmt.Fprint(&buf, "type ")

				errIn := format.Node(&buf, fset, typeSpec)
				if errIn != nil {
					err = errIn
					return true
				}
				fmt.Fprint(&buf, "\n")

				fmt.Fprintf(&buf, templateFunc, typeSpec.Name.Name)
			}
		}
		return true
	})
	return buf.Bytes(), err
}

const tempalteFile = `
package gen

import "fmt"

type(
	%s
)
`

const templateFunc = `
func (s *%s)Hello(){
	fmt.Println("hello world")
}
`

func write(srcBase []byte, outputPath string) error {
	// fmt.Println(string(srcBase))
	src, err := format.Source(srcBase)
	if err != nil {
		return xerrors.Errorf("failed to format.Source: %w", err)
	}
	// 書き込み
	err = ioutil.WriteFile(outputPath, src, 0664)
	if err != nil {
		fmt.Printf("failed to ioutil.WriteFile: %+v\n", err)
		return xerrors.Errorf("failed to format.Source: %w", err)
	}
	return nil
}
