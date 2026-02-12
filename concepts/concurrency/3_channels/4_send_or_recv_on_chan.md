### What Happen when we send on recv on channel


### 1. Buffered Channel 

```

ch := make(chan int, 3)

func G1(ch chan<- int){
    for _, v := range []{1,2,3,4} {
        ch <- v
    }
}

func G2(ch <-chan int){
    for v := range ch {
        fmt.Println(v)
    }
}

```


When we create a channel, this will be the representation

hchan 

-
buf. -> empty [|||]
-
lock
-
sendx
-
recvx
-
....
-

scenario when the G1 executes first, G1 is trying to send a value on the channel, which has empty buffer
    -> 1. goroutine has to aquire a lock on hchan struct
    -> 2. enqueues the element in circular ring buffer (this is memory copy which means element is copied into buffer)
    -> 3. increase value of sendx to 1
    -> 4. then it releases lock on the channel and proceed with further computation

now G2 comes along and tried to recieve value from the buffer 
    -> 1. g2 will aquire a lock on hchan struct
    -> 2. dequeue the element from the circular buffer and copies the value into variable v
    -> 3. increase recvx by 1 
    -> 4. releases lock on the channel struct and proceeds with further computation 

This is simple send and recv on buffered channel

-> There is no memory share between goroutines 
-> gororutine copy elements into and from hchan
-> hchan is protected by mutex lock

"Do not communicate by sharing memory instead, share memory by communicating"

## Full Buffer Scenario

lets say G1 
    1. enqueues the value 1,2,3 
    2. buffer full
    3. ch <- 4
    4. it will get blocked and wait to recv on the channel right ?
        how does this happen?
            -> G1 create sudog g struct and g ele will hold ref to goroutine and value to be sent will be saved in the elem field
            -> this is structue {G1, elem= 4} is include into sendq list
            -> g1 calls gopark()
            -> scheduler will move g1 out of the execution from the os thread and other local goroutines get scheduled on OS thread
            -> Now G2 comes along 
                -> tries to receive the value from the chan
                1. aquire lock
                2. dequeue value copies into variable g
                3. pops the waiting g1 on the sendq
                4. enqueue vale in of the elem field into the buffer
                5. it is g2 which will include value to the buffer on which g1 is blocked
                6. when g1 is enqueue is done g2 set state of g1 to runnable ( G2 calls goready(G1) )
                7. then g1 state changed to runnable it gets added to LRQ

                
Summary:
    1. when channel buffer is full and a goroutine tries to send value 
    2. Sender goroutine gets blocked, it is parked on SendQ 
    3. Data will be saved in the elem field of the sudog structure
    4. when receiver comes along, it dequeues the value from buffer
    5. enqueues the data from elem field to the buffer
    6. pops the goroutine in sendq, and puts it into runnable state



## Empty Buffer Scenario 

Lets say G2 get executed first and try to receive value from the empty channel


    -> G2 will create sudog struct for itself {G, elem = v (variable)}
    -> enqueue into recvq 
    -> G2 calls gopark()

    -> context switching with next goroutine in LRQ
    -> G1 comes and tries to send value on empty channel
        -> it checks if there are any goroutines in recvq of the channel
            -> now g2 will directly copies the value into v variable in g2 stack
                (g1 can access stack memory of g2 and can write directly to stack)
                [Only scenario where one go routine can access stack of another gorouitne]
                --> This is done for optimization reasons
    -> pops G2 from recvq
    -> G1 calls goready(G2)
    -> now g2 get scheduled in LRQ 


Summary:

    -> when goroutine calls receive on empty buffer
    -> goroutine is blocked, it is parked into recvq
    -> elem field of the sudog structure holds the ref to the stack variable of receiver goroutine
    -> when sender comes along, sender finds the goroutine in recvq
    -> sender copies the data, into the stack variable, on the receiver goroutines directly
    -> pops the goroutine in recvq, and puts it into runnable state


### Unbuffered channel


## 1. Send 

-> when sender goroutine wants to send values
-> if there is corresponding receiver waiting in recvq
-> sender will write the value directly into receiver goroutine stack variable
-> sender goroutine puts the receiver goroutine back to runnable state
-> if there is no receiver goroutine in recvq
-> sender gets parked into sendq
-> data is saved in elem field in sudog struct
-> receiver comes and copies the data
-> puts the sender to runnable state again

## 2. Receive

-> receiver goroutine wants to receive value
-> if it find a goroutine in waiting in sendq
-> receiver copies the value in elem field to its variable 
-> puts the sender goroutine to the runnable state
-> if there are no sender goroutine in sendq
-> receiver gets parked into recvq
-> reference to variable is saved in elem field in sudog struct
-> sender comes along and copies the data directly to receiver stack variable
-> puts the receiver back to runnable state



Summary :

-> hchan struct represents channel
-> it contains circular ring buffer and mutex lock
-> goroutines that gets blocked on send or recv are parked in sendq or recvq
-> go scheduler moves the blocked goroutines, out of OS thread
-> once channel operation is complete, goroutine is moved back to LRQ