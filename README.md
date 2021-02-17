# hello-go-gen

input
```go
type Foo struct {
	a int
}

type (
	Baz1 struct {
		aa int
	}
	Baz2 struct {
		b string
	}
)
```

output
```go
func (s *Foo) Hello() {
	fmt.Println("hello world")
}

func (s *Baz1) Hello() {
	fmt.Println("hello world")
}

func (s *Baz2) Hello() {
	fmt.Println("hello world")
}
```
