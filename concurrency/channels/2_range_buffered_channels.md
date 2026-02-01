## Range


Range over the channel, the receiver goroutine can use range to receive a sequence of values from the channel

1. Iterate over values received from a channel

2. Loop automatically breaks, when a channel is closed

3. range does not return the second boolean value

        for value := range ch {}

##  UnBuffered Channels 

Channels which we are using till now is unbuffered channels that mean there is no buffer in between 

since there is no buffer, sender goroutine will block until there is a receiver, to receive the value and the receiver until there is a sender

## Buffered Channels

In a buffered channels, there is a buffer between the sender and the receiver goroutine, and we can specify the capacity, that is the buffer size, which indicates the number of elements that can be sent without the receiver being ready to receive the values, the sender can keep sending the values without blocking, till the buffer gets full, when the buffer get full, the sender will block

* channels are given capacity
* in memory, fifo queue
* async 

ch := make(chan Type, capacity)

the receiver can keep receiving the values without blocking till the buffer gets empty, when the buffer get empty receiver will block 

the buffered channels are in memory FIFO queues, so the element that is sent first, will be element that will be read first 