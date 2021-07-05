package main

import (
	"fmt"
	"time"

	"github.com/nhatvu148/go-lang/async"
)

func DoneAsync() int {
	fmt.Println("Warming up ...")
	time.Sleep(3 * time.Second)
	fmt.Println("Done ...")
	return 1
}

func main() {
	fmt.Println("Let's start ...")
	future := async.Exec(func() interface{} {
		return DoneAsync()
	})
	fmt.Println("Done is running ...")
	val := future.Await()
	fmt.Println(val)
}
