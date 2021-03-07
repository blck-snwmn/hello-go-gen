//go:generate go run ./cmd/hello-go-gen -source=./$GOFILE -destination=./out/output.go
package main

// Foo is samle
type Foo struct {
	// inner comment
	a int // a is sample
}

type Bar struct {
	a int
	b string
}

// hoge
type (
	// Baz1 is sample
	Baz1 struct {
		aa int
	}
	Baz2 struct {
		b string
	}
)
