## sync.Cond

-> Condition variable is one of the sync mechanisms

-> a condition variable is basically a container of goroutines that are waiting for a certain condition

How to make a goroutine wait till some event(condition) occur?

-> One way - Wait in a loop for the condition

```

func WaitInALoop(){

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	var sharedRsc = make(map[string]string)

	go func ()  {
		
		defer wg.Done()
		mu.Lock()

		for len(sharedRsc) == 0 {
			mu.Unlock()
			time.Sleep(100*time.Millisecond)
			mu.Lock()
		}

		// Do processing...
		fmt.Println(sharedRsc["rsc"])
		mu.Unlock()
	}()
}

func main(){
	WaitInALoop()
}

```


    [this is very inefficient]

-> We need some way to make goroutine suspend while waiting

-> We need some way to signal the suspended goroutine when particular event has occured 

Channels ?

-> We can use channels to block a goroutine on receive

-> Sender goroutine to indicate occurance of event


-> What us there are multiple goroutines waiting on multiple conditions/events ?


thats where conditional variable comes into the picture

conditional variable are type 

```
var c *sync.Cond
```

-> we use constructor method sync.NewCond() to create a conditional variable, it takes sync.Locker interface as inpupt, which is usually sync.Mutex


```
m := sync.Mutex{}
c := sync.NewCond(&m)
```

sync.Cond has 3 methods

c.Wait()
c.Signal()
c.Broadcast()

-> suspends execution of the calling goroutine

->automatically unlcoks c.L

-> we cannot return unless awoken by broadcast or signal

-> wait locks c.l before returning

-> Because c.L is noit locked when wait first resumes, the caller typically cannot assume that the condition is truye when wait returns, Instead the called should wait in a loop


c.Signal()

-> wakes one goroutine waiting on c, if there is any

-> finds and notify longest waiting goroutine

-> it is allowed but not required for the caller to hold c.L during the call


c.Broadcast()

-> Wakes all goroutines waiting on c

-> it is allowed but not required for the caller to hold c.L during the call


```

func CondExample(){

	
	mu := &sync.Mutex{}
    wg := &sync.NewCond(mu)

	var sharedRsc = make(map[string]string)

	go func ()  {
		
		defer wg.Done()
		c.L.Lock()

		for len(sharedRsc) == 0 {
			c.Wait()
		}

		// Do processing...
		fmt.Println(sharedRsc["rsc"])
		c.L.Unlock()
	}()
}

func main(){
	CondExample()
}

```


Sumamry :

-> conditional var is used to sync exec of goroutines

-> wait suspends the execution of goroutine 

-> Signal wakes one goroutine waiting on c

-> Broadcast wakes all goroutine waiting on c