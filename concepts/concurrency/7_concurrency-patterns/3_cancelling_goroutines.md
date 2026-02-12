### Cancelling Goroutines


-> Upstream stages close their outbound channels when all the send operations are done
example : generator function

-> Downstream stage keep receiving values from inbound channel until the channel is closed
example : square function

-> All goroutines exit once all the values have been successfully sent downstream
example: merge function

But in real pipelines

-> Receiever stages may only need a subset of values to make progress

-> a stage can exit early because an inbound value represents an error in an earlier stage

-> receiver should not have to wait for the remaining values to arrive

-> we want earlier stages to stop producing values that later stages dont need


* Main goroutuines just receive one value
* Abondones the inbound channel from merge
* Merge goroutines will be blocked on channel send operation
* square and generate goroutine will be blocked on send
*  This leads to GOROUTINE LEAK

How can we signal to goroutine to abandon what ther are doing and terminate ?

Cancellation of goroutines
    -> pass a read only done channel to goroutine
    -> close the channel, to send broadcast signal to all gorotuine

    -> on receiving the signal on done channel, goroutines need to abandon their work and terminate

    -> we use select to make send/receive operation on channel pre-emptible