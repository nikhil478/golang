## Race condition


One of the complexity to manage in concurrency is race condition, the race condition occurs when two or more operations must be executed in the correct order to produce the desired result  but the program has not been written so that order is guaranteed to be maintained 

- race condition occurs when order of execution is not guaranteed.
- concurrent programs does not execute in the order they are coded


goroutines are executed async from the main routine , the order in which the main routine and go routine execute is undeterminisitic

what could be the possible outcomes of this program?

func main(){

    var data int

    go func(){
        data++
    }()

    if data == 0 {
        fmt.Printf("the value is %v \n", data)
    }
}

Output.  |.  Execution seq

nothing is printed |  data++ line get executed till main reaches data == 0 ocndition

0 is printed | data++ not reached and programs prints

1 is printed | data == 0 execute then data++ execute and print statement executes

Go flows logical concurrency model called fork and joins , if main does not wait for the goroutine, then it is very much possible that the program will finish before goroutine get chance to run


In order to create join point we uses sync wait group to deterministically block the main routine 