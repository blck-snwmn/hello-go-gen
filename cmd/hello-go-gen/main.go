package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"

	"golang.org/x/xerrors"
)

var (
	sourceArg      = flag.String("source", "", "")
	destinationArg = flag.String("destination", "", "")
)

func main() {
	flag.Parse()

	source, destination := *sourceArg, *destinationArg

	if source == "" {
		fmt.Println("source args is empty")
		return
	}
	if destination == "" {
		fmt.Println("destination args is empty")
		return
	}
	if err := execute(source, destination); err != nil {
		fmt.Printf("failed to execute: %+v\n", err)
	}
}

func execute(in, out string) error {
	result, err := generate(in)
	if err != nil {
		return xerrors.Errorf("failed to generate: %w", err)
	}
	if err := write(result, out); err != nil {
		return xerrors.Errorf("failed to write: %w", err)
	}
	return nil
}

func generate(inputPath string) ([]byte, error) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)

	var (
		buf bytes.Buffer
		err error
	)

	fmt.Fprintln(&buf, "// Code generated sample DO NOT EDIT.")

	fmt.Fprintln(&buf, "package generate")

	fmt.Fprintln(&buf, `import "fmt"`)

	ast.Inspect(f, func(n ast.Node) bool {
		if gd, ok := n.(*ast.GenDecl); ok {
			for _, spec := range gd.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						if len(gd.Specs) == 1 && typeSpec.Doc == nil {
							// 以下のようなケースでコメントを取得するための処理
							//  `gd.Lparen == token.NoPos` で判定もできそうだが、以下の場合,
							// TypeSpec.Doc == nil となることが多そうなので、それで判定している
							//
							// // Foo is ...
							// type Foo {}
							typeSpec.Doc = gd.Doc
						}
						doc := typeSpec.Doc
						typeSpec.Doc = nil
						if doc != nil {
							for _, d := range doc.List {
								fmt.Fprintln(&buf, d.Text)
							}
						}
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
			}
		}
		return true
	})
	return buf.Bytes(), err
}

const templateFunc = `
func (s *%[1]s)Hello(){
	fmt.Println("hello world")
}
func (s *%[1]s)Hello2(){
	fmt.Println("hello world2")
}
func (s *%[1]s)Hello3(){
	fmt.Println("hello world2")
}
func (s *%[1]s)Hello4(){
	fmt.Println("hello world2")
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
