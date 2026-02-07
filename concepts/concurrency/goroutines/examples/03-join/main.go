package main

import (
	"fmt"
	"sync"
)


func main(){

	var data int

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		data++
	}()

	wg.Wait()
	fmt.Printf("the value of data is %v\n", data)

	fmt.Println("Done...")
}