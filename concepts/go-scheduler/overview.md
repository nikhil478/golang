source : https://github.com/golang/go/blob/master/src/runtime/HACKING.md

The scheduler manages three types of resources that pervade the runtime: Gs, Ms, and Ps. It's important to understand these even if you're not working on the scheduler.

Gs, Ms, Ps

A "G" is simply a goroutine. It's represented by type g. When a goroutine exits, its g object is returned to a pool of free gs and can later be reused for some other goroutine.

An "M" is an OS thread that can be executing user Go code, runtime code, a system call, or be idle. It's represented by type m. There can be any number of Ms at a time since any number of threads may be blocked in system calls.

Finally, a "P" represents the resources required to execute user Go code, such as scheduler and memory allocator state. It's represented by type p. There are exactly GOMAXPROCS Ps. A P can be thought of like a CPU in the OS scheduler and the contents of the p type like per-CPU state. This is a good place to put state that needs to be sharded for efficiency, but doesn't need to be per-thread or per-goroutine.

The scheduler's job is to match up a G (the code to execute), an M (where to execute it), and a P (the rights and resources to execute it). When an M stops executing user Go code, for example by entering a system call, it returns its P to the idle P pool. In order to resume executing user Go code, for example on return from a system call, it must acquire a P from the idle pool.

All g, m, and p objects are heap allocated, but are never freed, so their memory remains type stable. As a result, the runtime can avoid write barriers in the depths of the scheduler.

Stacks

Every non-dead G has a user stack associated with it, which is what user Go code executes on. User stacks start small (e.g., 2K) and grow or shrink dynamically.

When a goroutine exits, its stack memory may be freed immediately or retained for reuse by another goroutine. If the stack has the default starting size, it is kept with the g object for reuse. If the stack has grown beyond the starting size, it is freed and a new stack will be allocated when the g is reused. Note that the g object itself is never freed (as described above), but its associated stack memory is managed separately and can be reclaimed.

Every M has a system stack associated with it (also known as the M's "g0" stack because it's implemented as a stub G) and, on Unix platforms, a signal stack (also known as the M's "gsignal" stack). System and signal stacks cannot grow, but are large enough to execute runtime and cgo code (8K in a pure Go binary; system-allocated in a cgo binary).

Runtime code often temporarily switches to the system stack using systemstack, mcall, or asmcgocall to perform tasks that must not be preempted, that must not grow the user stack, or that switch user goroutines. Code running on the system stack is implicitly non-preemptible and the garbage collector does not scan system stacks. While running on the system stack, the current user stack is not used for execution.