## Go Scheduler

-> Go scheduler is part of go runtime(application binary). it is known as M:N scheduler

-> Go scheduler runs in user space

-> Go scheduler uses OS threads to schedule goroutines for execution

-> go routines run in the context of OS threads

-> go runtime create no of worker os threads = GOMAXPROCS (default value is no of processor on machine)

-> go scheduler distributes runnable goroutimes over multiple worker OS threads

-> At any time, N go routines coculd be scheduled on M Os threads that runs on at most GOMAXPROCS numbers of processors.

As of go 1.14 go schdeluar implements async preemption, 

earlier its is cooperative schdular but then question comes what happen if long running go routine process hogs into the cpu, others gorutine would just get blocked

thats why async preemption implemented , what actually happen is , a goroutine os given a time slice of 10ms when that time slice is over , go scheduler will try to preempt it this provides other goroutines the opportunity to run even when there are long running cpu bound goroutines 

when its is created, it will be in runnable state, waiting in the run queue , it moves to the executing state once the goroutine is scheduled on the os thread, if the goroutines run through its time twice, then it is preempted and placed back into the run queue

if the goroutines gets blocked on any condition, like blocked on channel, blocked on syscall or waiting for the mutex lock, then they moved to waiting state, once the I/O operation is complete, they are moved back to the runnable state


## Elements 

For a cpu core, go runtime creates a OS thread, which is represented by the letter M, OS thread works pretty much like POSIX thread, Go runitme also create a logical processor P and associate that with the OS thread M, the logical processor hold the context for scheduling, which can be seen as a local scheduler running on a thread , G represents a goroutine

each logical processor P has a local run queue when runnable goroutines are queued

there is a global run queue, once the local queue is exhausted, the logical processor will pull goroutines from global run queue when the new goroutines created , it added at end of global run queue

There is no change as far as OS is concerned while context switching in goroutines , it is still scheduling same OS thread

-> Context switching between goroutines is managed my logical processor 

Major components are 
    -> M(os thread)
    -> P(logical processor which manages scheduling of goroutines)
    -> G (is the goroutine, which also includes scheduling information like stack and instruction pointer)
    -> LRQ (local runn queue where runnable goroutines are arranged)
    -> GRQ (global run queue when a goroutine is created they are placed inside global run queue)