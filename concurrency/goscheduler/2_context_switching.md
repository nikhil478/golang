## Context Switch due to synchronous system call

-> What happens in general when sync system ccall are made (like reading or writing to a file with sync flag set)
    1. there will be a disc I/O to be performed, so sync system call will block for I/O operation to be completed
    2. Os thread is moved out of the CPU to waiting queue for I/O to complete 

    (so we will not able to schedule any other goroutine on that thread)
    -> Sync system call can reduce parallelism

    how does goroutine handles this scenario ?

    Go scheduer identifies that G1 has caused OS thread M1 to block, so it brings in a new OS thread, either from the thread pool cache or it creates a new OS thread if a thread is not available in the thread pool
    then go schedular will detach the logical processor P form the OS thread M1 and movesit to the new OS thread M2

    G1 still attached to old OS thread M1, the logical processor P can now schedule other goroutines in its local run queue for execution on the OS thread M2

    Once the sync system call that was made by G1 is complete, then it is moved back to the end of the local run queue on the logical processor P and M1 is put to sleep and placed in the thread pool cache so it can be utilized in the future when similar scenario happens

## Context switching due to async system calls or http api call

What happens in general when async system call are made?

-> file broker is set to non blocking mode
-> if the file descripton is not ready, for I/O operation, system call does not block, but returns an error
    (for example socket buffer is empty and we are trying to read from it or if the socket buffer is full and we are trying to write to it, then the read or the write does not block, but returns an error)
    and the application will have to retry the operation again at a later point in time so this is good but it increases the application complexity
-> Asnyc IO increases the application complexity
    (applications will have to create any event loop and setup callbacks or it has to maintain a table mapping the file descriptor and the function pointer, it has to maintain a state to keep track of how much data was read last time or how much data was written last time) all these thing does add up to the complexity of the application and if it not implemented properly then it does make the application a bit inefficient

So how does handle this scenario ?
    Go scheduler uses netpoller to convert async system call to blocking system call

    -> when a goroutine makes a async system call and file descriptor is not ready, goroutine is parked at netpoller os thread

    -> netpoller uses interface provided by OS to do polling on file descriptor
        - kqueue(MacOS)
        - epoll(Linux)
        - iocp(Windows)
    to poll on the file descriptor
    once the netpoller gets a notif from the os, it in turn notifies the goroutine to retry I/O operation

    -> so complexity of managing async system call is moved from application to go runtime, which manages it efficiently
        (so the application need not have to make a call to select or poll and wait for the file descriptor to be ready, but instead it will be done by the netpoller in an efficient manner)

        lets say G1 is running on OS thread M1 and opens an network connection with net.Dial the file descriptor used for the conenction is set to non-blocking mode. when goroutines tries to read and write to that connection, the networking call will do the operation until it receives an error EAGAIN , then it calls into the netpoller, then the scheduler will ove the goroutine G1 out of the OS thread M1 to the netpoller Thread and another goroutine in the LRQ, in this case G2 gets scheduled to run on the OS thread M1 

        the netpoller uses the interface provided by the os to poll on the file descriptor
        when the netpoller receives the notif from the operating system that it can perform and I/O operation on the file descriptor, then it will look through its internal data structure to see if there are any goroutines that are blocked on that file descriptior

        then it notifies that go routine, then that goroutine can retry the I/O operation. Once the I/O operation is complete, goroutine is moved back to the LRQ of M1 , In this way to process an async system call, no extra OS thread is used, instead the netpoller OS thread is used to process the goroutines
