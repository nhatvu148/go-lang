package main

import (
	"flag"
	"fmt"
	"strconv"
)

type in int32

var serverChan = make(chan in)

func main() {
	// var num int64

	// num = 126

	// oct_num := strconv.FormatInt(num, 2)

	// fmt.Println("binary num: ", oct_num)
	// fmt.Println(bits.OnesCount(uint(num)))
	flag.Parse()
	x := flag.Args()
	y := x[1:]

	oddChan := make(chan int32)
	evenChan := make(chan int32)

	odd := []int32{}
	even := []int32{}
	for _, n := range y {
		num, err := strconv.ParseInt(n, 10, 32)
		num32 := int32(num)
		if err == nil {
			if num32%2 == 0 {
				evenChan <- num32
			} else {
				oddChan <- num32
			}
		}
	}

	go func() {
		for {
			select {
			case x1 := <-oddChan:
				odd = append(odd, x1)
			case x1 := <-evenChan:
				even = append(even, x1)
			default:
				fmt.Println("Default case..!")
			}
		}
	}()

	// process(ch)
	for _, n := range odd {
		fmt.Println(n)
	}

	for _, n := range even {
		fmt.Println(n)
	}

	fmt.Println()
}

func process(ch <-chan int) {
	s := <-ch
	fmt.Println(s)
	//ch <- 2
}

func Server() {

}
