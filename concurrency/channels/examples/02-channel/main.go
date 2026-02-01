package main

import "fmt"

// func main(){

// 	go func() {
// 		for i := 0; i < 6; i++ {
// 			// TODO: send iterator over channel
// 		}
// 	}()

// 	// TODO: range over channel to recv values
// }


func main(){

	ch := make(chan int)
	go func() {
		for i := 0; i < 6; i++ {
			ch <- i
		}
		close(ch)
	}()

	for value := range ch {
		fmt.Println(value)
	}
}