package main

import (
	"fmt"
)

func main() {
	lowercase := "abcdefghijklmnopqrstunwxyz"
	for _, c := range lowercase {
		fmt.Println(c)
	}
	fmt.Println(int('A'))
	fmt.Println(string(rune(65)))
}
