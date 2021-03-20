package main

import (
	"fmt"

	"github.com/prashantv/gostub"
)

// Greet return "hello,xxx"
func Greet(name string) string {
	return "hello," + name
}

func main() {
	var GreetFunc = Greet
	fmt.Println("Before stub:", GreetFunc("axiaoxin"))
	stubs := gostub.Stub(&GreetFunc, func(name string) string {
		return "fuck u," + name
	})
	defer stubs.Reset()
	fmt.Println("After stub:", GreetFunc("axiaoxin"))
	// Output:
	// Before stub: hello,axiaoxin
	// After stub: fuck u,axiaoxin
}
