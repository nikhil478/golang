package main

import (
	"fmt"
	"runtime"
	"sync"
)

// func main() {
// 	runtime.GOMAXPROCS(4)

// 	var balance int
// 	var wg sync.WaitGroup

// 	deposit := func(amount int) {
// 		balance += amount
// 	}

// 	withdraw := func(amount int) {
// 		balance -= amount
// 	}

// 	//make 100 deposits of $1
// 	// and 100 windrows of $1 concurrently
// 	// run the program and check result

// 	// TODO: fix the issue for consistent output.

// 	wg.Add(100)
// 	for i := 0; i < 100; i++ {
// 		go func(){
// 			defer wg.Done()
// 			deposit(1)
// 		}()
// 	}

// 	wg.Add(100)
// 	for i := 0; i < 100; i++ {
// 		go func(){
// 			defer wg.Done()
// 			withdraw(1)
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Println(balance)
// }


func main() {
	runtime.GOMAXPROCS(4)

	var balance int
	var wg sync.WaitGroup
	var mu sync.Mutex

	deposit := func(amount int) {
		mu.Lock()
		balance += amount
		mu.Unlock()
	}

	withdraw := func(amount int) {
		mu.Lock()
		balance -= amount
		mu.Unlock()
	}

	//make 100 deposits of $1
	// and 100 windrows of $1 concurrently
	// run the program and check result

	// TODO: fix the issue for consistent output.

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(){
			defer wg.Done()
			deposit(1)
		}()
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(){
			defer wg.Done()
			withdraw(1)
		}()
	}

	wg.Wait()
	fmt.Println(balance)
}