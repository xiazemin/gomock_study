package foo

import "fmt"

type Foo interface {
	Bar(x int) int
	Bar1(x int) int
}

func SUT(f Foo) {
	// ...
	if 99 == f.Bar(88) {
		fmt.Print("ok")
	}
}
