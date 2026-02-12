package main

import (
	"fmt"
	"sync"
)

func main(){

	var wg sync.WaitGroup

	// what is the output
	// TODO: fix the issue

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		// go func ()  {
		// 	defer wg.Done()
		// 	fmt.Println(i)     this will print 4 // go routine operate on the current value of the variable at time of their execution
		// }(). // if you want to operate at specific value then u need to pass that to goroutine inside params
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}
