package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
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

func printHeader(writer io.Writer) {
	fmt.Fprintln(writer, "// Code generated sample DO NOT EDIT.")
	fmt.Fprintln(writer, "package generate")
	fmt.Fprintln(writer, `import "fmt"`)
}

func printBySpec(writer io.Writer, fset *token.FileSet, gd *ast.GenDecl, spec ast.Spec) error {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil
	}
	_, ok = typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}
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
			fmt.Fprintln(writer, d.Text)
		}
	}
	fmt.Fprint(writer, "type ")

	err := format.Node(writer, fset, typeSpec)
	if err != nil {
		return err
	}
	fmt.Fprint(writer, "\n")

	fmt.Fprintf(writer, templateFunc, typeSpec.Name.Name)
	return nil
}

func printByGenDel(writer io.Writer, fset *token.FileSet, gd *ast.GenDecl) error {
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
						fmt.Fprintln(writer, d.Text)
					}
				}
				fmt.Fprint(writer, "type ")

				err := format.Node(writer, fset, typeSpec)
				if err != nil {
					return err
				}
				fmt.Fprint(writer, "\n")

				fmt.Fprintf(writer, templateFunc, typeSpec.Name.Name)
			}
		}
	}
	return nil
}

func printBody(writer io.Writer, f *ast.File, fset *token.FileSet) error {
	var err error
	ast.Inspect(f, func(n ast.Node) bool {
		if gd, ok := n.(*ast.GenDecl); ok {
			if errIn := printByGenDel(writer, fset, gd); errIn != nil {
				err = errIn
				return false
			}
		}
		return true
	})
	return err
}

func generate(inputPath string) ([]byte, error) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)

	var (
		buf bytes.Buffer
		err error
	)

	printHeader(&buf)

	printBody(&buf, f, fset)

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
