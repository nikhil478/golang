In this module we will look into select

* G1 wants to receive result of computation from G2 and G3

        G1  <- G2 (Task1)
            <- G3 (Task2)

In what order are we going to receive results?

g1 <- g2
g1 <- g3

OR

g1 <- g3
g1 <- g2

What if g3 was much faster in one instance and g2 is faster than g3 in another ?

can we do operation on channel which ever is ready and dont worry about the order ?

* Select statement is like a switch 

```
select {
    case <- ch1:
        // block of statements
    case <- ch2:
        // block of statements
    case ch3 <- struct{}{}:
        // block of statements
}

```

-> Each cases specifies communication
-> All channel operation are considered simulatenously to see any of them is ready
-> select waits until some case is ready to proceed
-> of none of the channel is ready entire select statement is going to be blocked until some case is ready for communication
-> when one the channels is ready, that operation will proceed and execute associated block of statements
-> if multiple channel are ready it will pick one at random

-> Select is very helpful in implementing
    -> Timeouts
    -> Non blocking communication 
        (default allows you to exit a select block without blocking)


[Empty select statement will block forever] 

```
Select{}

```

[Select on nil channel will block forever]

```

var ch chan string

select {

    case v := <-ch:
    case ch <- v:
}


```

Summary 

-> Select is like switch statement with each statement specifing channel operation
-> select will block until any of the case statement is ready
-> with select we can implement timeout and non blocking communication
-> select on nil channel or empty channel will block forever 