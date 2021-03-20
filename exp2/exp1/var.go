package main

import (
	"fmt"

	"github.com/prashantv/gostub"
)

var counter = 100

func stubGlobalVariable() {
	stubs := gostub.Stub(&counter, 200)
	defer stubs.Reset()
	fmt.Println("Counter:", counter)
	// 可以多次打桩
	stubs.Stub(&counter, 10000)
	fmt.Println("Counter:", counter)
}

func main() {
	stubGlobalVariable()
}
