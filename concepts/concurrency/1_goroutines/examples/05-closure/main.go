package main

import (
	"fmt"
	"sync"
)

// usually variable inside the function go out of scope when function returns
// but here goroutine is smart , as it sees there is ref to variable i so it pins it , it moves it from the stack to heap, so that goroutine still has the access to the variable
// evens after the enclosing function returns

func main(){

	var wg sync.WaitGroup

	incr := func(wg *sync.WaitGroup) {
		var i int
		wg.Add(1)

		go func(){
			defer wg.Done()
			i++
			fmt.Println("value of i: %v \n", i)
		}()

		fmt.Println("return from function")
		return
	}

	incr(&wg)
	wg.Wait()
	fmt.Println("done...")
}


