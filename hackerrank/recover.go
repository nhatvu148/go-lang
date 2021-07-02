package main

import (
	"fmt"
	"log"
)

func test() {
	fmt.Println("Panic Start")
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error:", err)
		}
	}()
	panic("Golang Panic Occoured!")
	fmt.Println("This will not be executed.")
}

func main() {

	fmt.Println("Start")
	test()
	fmt.Println("end")

}
