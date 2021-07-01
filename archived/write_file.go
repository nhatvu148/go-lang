package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// hello := "Hello world"
	// myStringArr := []string{"hello", "my", "world"}

	// fmt.Println(strings.Join(myStringArr, ", "))
	// fmt.Println([]byte(hello))
	// ioutil.WriteFile("hello.txt", []byte(hello), 0644)

	result, err := ioutil.ReadFile("hello.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
}
