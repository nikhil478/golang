package main

import "fmt"

// func main(){
// 	go func (a, b int)  {
// 		c := a + b
// 	}(1, 2)

// 	// TODO: get the value computed from goroutine
// 	// fmt.Printf("computed value %v \n",)
// }

func main() {

	ch := make(chan int)
	go func(a, b int) {
		c := a + b
		ch <- c
	}(1, 2)
	r := <-ch

	fmt.Printf("computed value %v \n", r)
}
