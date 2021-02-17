package gen

import "fmt"

type (
	Foo struct {
		a int
	}
	Bar struct {
		a int
		b string
	}
	Baz1 struct {
		aa int
	}
	Baz2 struct {
		b string
	}
)

func (s *Foo) Hello() {
	fmt.Println("hello world")
}

func (s *Bar) Hello() {
	fmt.Println("hello world")
}

func (s *Baz1) Hello() {
	fmt.Println("hello world")
}

func (s *Baz2) Hello() {
	fmt.Println("hello world")
}
