package main

import (
	"context"
	"fmt"
)

// func main(){

// 	// TODO: generator generates integers in a seperate goroutine and
// 	// sends them to the returned channel
// 	// callers of gen need to cancel the goroutine once
// 	// they consume 5 integers
// 	// so that internal goroutine
// 	// started by gen is not leaked

// 	generator := func () <-chan int  {

// 	}

// 	// create a context that is cancellable
// }


func main(){
	

	// TODO: generator generates integers in a seperate goroutine and
	// sends them to the returned channel
	// callers of gen need to cancel the goroutine once
	// they consume 5 integers
	// so that internal goroutine
	// started by gen is not leaked


	generator := func (ctx context.Context) <-chan int  {
		
		ch := make(chan int)
		n := 1
		go func ()  {
			defer close(ch)
			for {
				select {
				case ch <- n:
				case <- ctx.Done():
					return
				}
				n++
			}
		}()
		return ch
	}

	// create a context that is cancellable

	ctx, cancel := context.WithCancel(context.Background())

	ch := generator(ctx)

	for n := range ch {
		fmt.Println(n)
		if n == 5 {
			cancel()
		}
	}
}