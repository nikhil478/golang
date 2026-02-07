package main

import "fmt"

// func main() {

// 	// TODO: create channel owner, which creates channel,
// 	// return receive only channel to caller
// 	// spins  a go routine, which
// 	// writes data into channel and
// 	// closes the channel when done

// 	consumer := func(ch <- chan int) {
// 		// read values from channel

// 		for v := range ch {
// 			fmt.Printf("Received: %d\n", v)
// 		}

// 		fmt.Println("Done receiving !")
// 	}

// 	ch := owner()
// 	consumer(ch)
// }

func main() {

	// TODO: create channel owner, which creates channel,
	// return receive only channel to caller
	// spins  a go routine, which
	// writes data into channel and
	// closes the channel when done

	owner := func() <-chan int {
		ch := make(chan int, 10)
		go func() {
			defer close(ch)
			for i := range 10 {
				ch <- i
			}
		}()
		return ch
	}

	consumer := func(ch <-chan int) {
		// read values from channel

		for v := range ch {
			fmt.Printf("Received: %d\n", v)
		}

		fmt.Println("Done receiving !")
	}

	ch := owner()
	consumer(ch)
}
