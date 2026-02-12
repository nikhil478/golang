package main

import (
	"context"
	"fmt"
	"time"
)

type data struct {
	result string
}

// func main() {

// 	// TODO: set deadline for goroutine to return computational result

// 	compute := func() <-chan data {

// 		ch := make(chan data)

// 		go func() {
// 			defer close(ch)

// 			time.Sleep(50* time.Millisecond)
			
// 			ch <- data{"123"}
// 		}()

// 		return ch
// 	}

// 	// wait for the work to finfish, it is takes too long move on

// 	ch := compute()

// 	d := <-ch

// 	fmt.Printf("work is complete: %s\n", d)
// }


func main() {

	// TODO: set deadline for goroutine to return computational result

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel() 

	compute := func() <-chan data {

		ch := make(chan data)

		go func() {
			defer close(ch)
			deadline, ok := ctx.Deadline()
			if ok {
				if deadline.Sub(time.Now().Add(50 * time.Millisecond)) < 0 {
					fmt.Println("not sufficient time given ...  terminating")
				}
			}

			time.Sleep(50* time.Millisecond)

			select {
			case ch <- data{"123"}:
			case <- ctx.Done():
				return
			}
			
			
		}()

		return ch
	}

	// wait for the work to finfish, it is takes too long move on

	ch := compute()

	d := <-ch

	fmt.Printf("work is complete: %s\n", d)
}