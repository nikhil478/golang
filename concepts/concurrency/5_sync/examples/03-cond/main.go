package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedRsc = make(map[string]interface{})

// func main() {
// 	var wg sync.WaitGroup

// 	wg.Add(1)

// 	go func ()  {

// 		defer wg.Done()

// 		// TODO: suspend goroutine until sharedRsc is populated


// 		for len(sharedRsc) == 0 {
// 			time.Sleep(1 * time.Millisecond)
// 		}

// 		fmt.Println(sharedRsc["rsc1"])
// 	}()



// }/

func SingleWaitCond() {
	var wg sync.WaitGroup

	mu := sync.Mutex{}
	c := sync.NewCond(&mu)


	wg.Add(1)

	go func ()  {

		defer wg.Done()

		// TODO: suspend goroutine until sharedRsc is populated

		c.L.Lock()
		for len(sharedRsc) == 0 {
			c.Wait()
			time.Sleep(1 * time.Millisecond)
		}

		fmt.Println(sharedRsc["rsc1"])
		c.L.Unlock()
	}()


	c.L.Lock()

	// writes changes to sharedRsc
	sharedRsc["rsc1"] = "foo"
	c.Signal()
	c.L.Unlock()

	wg.Wait()
}


// func BroadcastExample(){

// 	var wg sync.WaitGroup

// 	wg.Add(1)

// 	go func ()  {
		
// 		defer wg.Done()

// 		// TODO: suspend goroutine until sharedRsc is populated.

// 		for len(sharedRsc) < 1 {
// 			time.Sleep(1 * time.Millisecond)
// 		}

// 		fmt.Println(sharedRsc["rsc1"])
// 	}()

// 	wg.Add(1)
// 	go func ()  {
// 		defer wg.Done()

// 		for len(sharedRsc) < 2 {
// 			time.Sleep(1 * time.Millisecond)
// 		}

// 		fmt.Println(sharedRsc["rsc2"])
// 	}()


// 	sharedRsc["rsc1"] = "foo"
// 	sharedRsc["rsc2"] = "bar"

// 	wg.Wait()
// }


func BroadcastExample(){

	var wg sync.WaitGroup

	mu := sync.Mutex{}
	c := sync.NewCond(&mu)


	wg.Add(1)

	go func ()  {
		
		defer wg.Done()

		// TODO: suspend goroutine until sharedRsc is populated.

		c.L.Lock()
		for len(sharedRsc) < 1 {
			c.Wait()
		}

		fmt.Println(sharedRsc["rsc1"])
		c.L.Unlock()
	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()

		c.L.Lock()
		for len(sharedRsc) < 2 {
			c.Wait()
		}

		fmt.Println(sharedRsc["rsc2"])
		c.L.Unlock()
	}()

	c.L.Lock()
	sharedRsc["rsc1"] = "foo"
	sharedRsc["rsc2"] = "bar"
	c.Broadcast()
	c.L.Unlock()

	wg.Wait()
}

func main(){
	BroadcastExample()
}