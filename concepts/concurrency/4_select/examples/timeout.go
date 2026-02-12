package examples

import (
	"fmt"
	"time"
)


// func ExTimeout(){

// 	ch := make(chan string, 1)

// 	go func(){
// 		time.Sleep(2 * time.Second)
// 		ch <- "one"
// 	}()

// 	// TODO: implement timeout for recv on channel ch

// 	m := <- ch
// 	fmt.Println(m)
// }


func ExTimeout(){

	ch := make(chan string, 1)

	go func(){
		time.Sleep(2 * time.Second)
		ch <- "one"
	}()

	select {
		case m := <- ch:
			fmt.Println(m)
		case <- time.After(1 * time.Second):
			fmt.Println("timeout")
	}
}