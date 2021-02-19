package generate

import "fmt"

type Foo struct {
	a int
}

func (s *Foo) Hello() {
	fmt.Println("hello world")
}

type Bar struct {
	a int
	b string
}

func (s *Bar) Hello() {
	fmt.Println("hello world")
}

type Baz1 struct {
	aa int
}

func (s *Baz1) Hello() {
	fmt.Println("hello world")
}

type Baz2 struct {
	b string
}

func (s *Baz2) Hello() {
	fmt.Println("hello world")
}
