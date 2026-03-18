## Garbage Collector

-> Go language takes responsibility for arranging the storage of Go values
-> Go Developer need not care about where these values are stored
-> In practice, these values often to be store in computer physical memory and physical memory is a finite resource
-> As it is finite, memory must be managed carefully and recycled in order to avoid running out of it, while executing a go program
-> Its the job of a Go implementation to allocate and recycle memory as needed

-> Another term for automatically recycling memory is Garbage Collection

At a high level , a gc (garbage collector) is a system that recycles memory on behalf of the application by identifying which parts of memory are no longer needed. The go standard toolchain provides a runtime library that ships with every application, and the runtime library includes a gc

Note: Existence of a gc is not guaranteed by GO specification, only that the underlying storage for GO values is managed by the language, This omission is intentional and enables the use of radically different management memory techniques

1. Go compliers options (gc(default), Gccgo and Gollvm)
2. Manual Handling
3. Ownership-like patterns (Rust-style thinking in Go)
    - You mimic ownership by clearly defining who controls data
    - Avoid shared mutable state
    - Prevent memory leaks via references
4. Arena-style allocation (simulate in Go) : Go doesn’t officially support arenas (except experimental), but you can mimic it.
    - Group objects and discard them together.



## memory that doesn't need to be managed by the GC

For instance, non-pointer Go values stored in local variables will likely not be managed by the Go GC at all, and Go will instead arrange for memory to be allocated that's tied to the lexical scope in which it's created. In general, this is more efficient than relying on the GC, because the Go compiler is able to predetermine when that memory may be freed and emit machine instructions that clean up. Typically, we refer to allocating memory for Go values this way as "stack allocation," because the space is stored on the goroutine stack

Go values whose memory cannot be allocated this way, because the Go compiler cannot determine its lifetime, are said to escape to the heap. "The heap" can be thought of as a catch-all for memory allocation, for when Go values need to be placed somewhere. The act of allocating memory on the heap is typically referred to as "dynamic memory allocation" because both the compiler and the runtime can make very few assumptions as to how this memory is used and when it can be cleaned up. That's where a GC comes in: it's a system that specifically identifies and cleans up dynamic memory allocations.

There are many reasons why a Go value might need to escape to the heap. One reason could be that its size is dynamically determined. Consider for instance the backing array of a slice whose initial size is determined by a variable, rather than a constant. Note that escaping to the heap must also be transitive: if a reference to a Go value is written into another Go value that has already been determined to escape, that value must also escape.

Whether a Go value escapes or not is a function of the context in which it is used and the Go compiler's escape analysis algorithm. It would be fragile and difficult to try to enumerate precisely when values escape: the algorithm itself is fairly sophisticated and changes between Go releases. For more details on how to identify which values escape and which do not, see the section on eliminating heap allocations (https://go.dev/doc/gc-guide#Eliminating_heap_allocations).