package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	values := make(chan string, 2)
	defer close(values)

	go func() {
		defer wg.Done()
		sendValue(values) // spin up a goroutine
	}()

	go func() {
		defer wg.Done()
		sendValue(values) // spin up a goroutine
	}()

	value := <-values

	fmt.Println(value)
}

func sendValue(c chan string) {
	fmt.Println("Executing Goroutine")
	time.Sleep(1 * time.Second)
	c <- "Hello world" // block when received
	fmt.Println("Finished Executing Goroutine")
}
