package main

import (
	"fmt"
	"sync"
)

// func main() {
// 	var wg sync.WaitGroup

// 	load := func ()  {
// 		fmt.Println("Run only once init function")
// 	}

// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func ()  {
// 			defer wg.Done()

// 			// TODO: modify so that load function gets called only once
// 			load()
// 		}()
// 	}
// 	wg.Wait()
// }


/* we can introduce concurrent safe flag b/w go routine and check the flag */
/* we will use once struct which have done and mutex flag internally*/
func main() {
	var wg sync.WaitGroup
	var once sync.Once

	load := func ()  {
		fmt.Println("Run only once init function")
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func ()  {
			defer wg.Done()

			// TODO: modify so that load function gets called only once
			once.Do(load)
		}()
	}
	wg.Wait()
}