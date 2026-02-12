## Concurrency in go 

-> based on the paper written by tony hoare, communicating sequential processes (CSP)

    1. Each process is built for sequential execution. 
       (Every process has a local state and the process operates on local state)

    2. Data is communicated between processes
        No shared memory (we send copy of data to other resources)
        As there is no sharing of memory, there would be no race condition and deadlock

    3. Scale by adding more of the same
        we can scale easily as each process can run independently 
        if the computation is taking more time, we can add more processes of the same type and run the computation faster

## Gos concurrency tool set

1. Go routines: are concurrently executing functions
2. Channels: are use to communicate data between go routines
3. Select: is use to multiplex the channels
4. Sync package: provides classical syncronization tools like the mutex, conditional variables and others


-> We can think goroutines as user space threads managed by go runtime, go runtime is part of executable, it is built into the executable of the application
-> Go routines extremely lightweight, goroutines starts with 2kb of stack, which grows and shrink as required
-> Low CPU overhead - three instructions per function call (amount of CPU instructions required to create go routine is very less)
    this enables us to create hundreds and thousands of go routines in the same address space 
-> channels are used for communication of data between goroutines. Sharing of data should be avoided 
-> Context switching between goroutines is much cheaper than thread context switching 
-> Go runtime can be more selective in what is persisted for retrievel, how it is persisted and when the persisting needs to occur
-> Go runtime creates OS threads, go routines run in the context of the OS thread
    Many goroutines can execute in the context of the single OS thread, The operating system schedules, the OS threads and the go runtime schedules, multiple goroutines on the OS thread 
-> Go runtime manages the scheduling of the goroutines on the OS threads