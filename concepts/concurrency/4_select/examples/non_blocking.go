package examples

import (
	"fmt"
	"time"
)

// func ExNonBlockingCommunication(){
// 	ch := make(chan string)

// 	go func(){
// 		for i := 0; i < 3; i++ {
// 			time.Sleep(1 * time.Second)
// 			ch <- "message"
// 		}
// 	}()

// 	// TODO: if there is no value on channel, do not block.

// 	for i := 0; i < 2; i++ {
// 		m := <-ch
// 		fmt.Println(m)


// 		// Do some processing
// 		fmt.Println("processing...")
// 		time.Sleep(1500* time.Millisecond)
// 	}
// }


func ExNonBlockingCommunication(){
	ch := make(chan string)

	go func(){
		for range 3 {
			time.Sleep(2 * time.Second)
			ch <- "message"
		}
	}()

	for i := 0; i < 2; i++ {
		select {
			case m := <-ch:
				fmt.Println(m)
				// Do some processing
			default:
				fmt.Println("no messages received")
		}
		fmt.Println("processing...")
		time.Sleep(3* time.Second)
	}
}