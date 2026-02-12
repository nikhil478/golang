## Work Stealing 

1. Work stealing helps to balance the goroutines across all logical processors

2. Work get better distributed and get done more efficiently

Work Stealing Rule

-> If there is no goroutine in local run queue
    -> try to steal from other logical processors
    -> if not found, check the global runnable queue for a G
    -> if not found check netpoller