### What Happen when we send on recv on channel

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