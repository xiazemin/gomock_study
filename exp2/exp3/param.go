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
	// StubFunc 第一个参数必须是一个函数变量的指针，该指针指向的必须是一个函数变量，第二个参数为函数 mock 的返回值
	stubs := gostub.StubFunc(&GreetFunc, "fuck u,axiaoxin")
	defer stubs.Reset()
	fmt.Println("After stub:", GreetFunc("axiaoxin"))
	// Output:
	// Before stub: hello,axiaoxin
	// After stub: fuck u,axiaoxin
}
