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

