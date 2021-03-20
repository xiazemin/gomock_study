package main

import (
	"fmt"
	"time"

	"github.com/prashantv/gostub"
)

func main() {
	var XNow = time.Now
	fmt.Println("Before stub:", XNow())
	stubs := gostub.Stub(&XNow, func() time.Time {
		d, _ := time.ParseDuration("+1h")
		return time.Now().Add(d)
	})
	defer stubs.Reset()
	fmt.Println("After stub:", XNow())
	// Output:
	// Before stub: 2020-06-02 18:03:25.791691 +0800 CST m=+0.000065259
	// After stub: 2020-06-02 19:03:25.791843 +0800 CST m=+3600.000217139
}
