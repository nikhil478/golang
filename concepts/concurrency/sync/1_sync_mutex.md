### Sync


When to use channels and when to use mutex ?


Channels:
-> Passing copy of data
-> Distributing units of work
-> Communicating async results


1. Mutex (data is big or we want the access to these data to be concurrent safe)

-> Caches
-> State


Mutex

-> used for protect shared resources
-> sync.Mutex - provide exclusive access to a shared resource
[if one goroutine other goroutine will be going to block]

Region between the lock and unlock is called the critical section


sync.RWMutex - Allows multiple readers. Writers get exclusive lock.

Summary 

-> Mutex is used guards access to shared resources

-> It is developers convention to call Lock() to access shared memory and call Unlock() when done


-> Critical section represents the bottleneck between the goroutines