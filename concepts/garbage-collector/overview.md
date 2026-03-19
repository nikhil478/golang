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


## Tracing Garbage Collection

Garbage collection may refer to many different methods of automatically recycling memory

for example, reference counting.

In the context of this doc, gc refers to tracing gc, which identifies in use, so called live, objects by following pointers transitevly 

- Object: An object is a dynamically allocated piece of memory that contains one or more Go values

- Pointer: A memory address that references any value within an object. This naturally includes Go values of the form *T, but also includes parts of built in Go values,
             but also includes parts of built-in Go values. Strings, slices, channels, maps, and interface values all contain memory addresses that the GC must trace

Together, objects and pointers to other objects form the object graph. To identify live memory, the GC walks the object graph starting at the program's roots, pointers that identify objects that are definitely in-use by the program. Two examples of roots are local variables and global variables. The process of walking the object graph is referred to as scanning. Another phrase you might see in the Go documentation is whether an object is reachable, which just means that the object can be discovered by the scanning process. Note also that, with one exception, once memory becomes unreachable, it stays unreachable.

This basic algorithm is common to all tracing GCs. Where tracing GCs differ is what they do once they discover memory is live. Go's GC uses the mark-sweep technique, which means that in order to keep track of its progress, the GC also marks the values it encounters as live. Once tracing is complete, the GC then walks over all memory in the heap and makes all memory that is not marked available for allocation. This process is called sweeping.

One alternative technique you may be familiar with is to actually move the objects to a new part of memory and leave behind a forwarding pointer that is later used to update all the application's pointers. We call a GC that moves objects in this way a moving GC; Go has a non-moving GC.

## The GC cycle

Because the Go GC is a mark-sweep GC, it broadly operates in two phases: the mark phase, and the sweep phase.
While this statement might seem tautological, it contains an important insight: it's not possible to release memory back to be allocated until all memory has been traced, because there may still be an un-scanned pointer keeping an object alive
As a result, the act of sweeping must be entirely separated from the act of marking
Furthermore, the GC may also not be active at all, when there's no GC-related work to do. The GC continuously rotates through these three phases of sweeping, off, and marking in what's known as the GC cycle. For the purposes of this document, consider the GC cycle starting with sweeping, turning off, then marking.

The next few sections will focus on building intuition for the costs of the GC to aid users in tweaking GC parameters for their own benefit

# Understanding costs

The GC is inherently a complex piece of software built on even more complex systems. It's easy to become mired in detail when trying to understand the GC and tweak its behavior. This section is intended to provide a framework for reasoning about the cost of the Go GC and its tuning parameters.

To begin with, consider this model of GC cost based on three simple axioms.

1. The GC involves only two resources: physical memory, and CPU time

2. The GC's memory costs consist of live heap memory, new heap memory allocated before the mark phase, and space for metadata that, even if proportional to the previous costs, are small in comparison.

GC memory cost for cycle N = live heap from cycle N-1 + new heap

Live heap memory is memory that was determined to be live by the previous GC cycle, while new heap memory is any memory allocated in the current cycle, which may or may not be live by the end. How much memory is live at any given point in time is a property of the program, and not something the GC can directly control.

The GC's CPU costs are modeled as a fixed cost per cycle, and a marginal cost that scales proportionally with the size of the live heap.

GC CPU time for cycle N = Fixed CPU time cost per cycle + average CPU time cost per byte * live heap memory found in cycle N

The fixed CPU time cost per cycle includes things that happen a constant number of times each cycle, like initializing data structures for the next GC cycle. This cost is typically small, and is included just for completeness.

Most of the CPU cost of the GC is marking and scanning, which is captured by the marginal cost. The average cost of marking and scanning depends on the GC implementation, but also on the behavior of the program. For example, more pointers means more GC work, because at minimum the GC needs to visit all the pointers in the program. Structures like linked lists and trees are also more difficult for the GC to walk in parallel, increasing the average cost per byte.

This model ignores sweeping costs, which are proportional to total heap memory, including memory that is dead (it must be made available for allocation). For Go's current GC implementation, sweeping is so much faster than marking and scanning that the cost is negligible in comparison.

To see why, let's explore a constrained but useful scenario: the steady state. The steady state of an application, from the GC's perspective, is defined by the following properties:

The rate at which the application allocates new memory (in bytes per second) is constant.

This means that, from the GC's perspective, the application's workload looks approximately the same over time. For example, for a web service, this would be a constant request rate with, on average, the same kinds of requests being made, with the average lifetime of each request staying roughly constant.

The marginal costs of the GC are constant.

This means that statistics of the object graph, such as the distribution of object sizes, the number of pointers, and the average depth of data structures, remain the same from cycle to cycle.

Let's work through an example. Assume some application is operating in a steady-state, allocating 10 MiB/s, while the GC can scan memory at a rate of 100 MiB/cpu-second (this is made up). The steady state makes no assumptions about the size of the live heap, but for simplicity, let's say this application's live heap is always 10 MiB. Let's also assume, again, for simplicity, that the fixed GC costs are zero. Let's play around with the GC cycle period.

Suppose each GC cycle happens after exactly 1 cpu-second. Then, by the end of each GC cycle our example application will have allocated 10 MiB of additional memory, resulting in a 20 MiB total heap size. And with every GC cycle, the GC will spend 0.1 cpu-seconds scanning the 10 MiB live heap, resulting in a 10% CPU overhead. Recall that the GC only needs to walk the live heap, not the whole heap. (Note: a constant live heap does not mean that all newly allocated memory is dead. It means that, after the GC runs, some mix of old and new heap memory dies, and only that the end result is 10 MiB found live each cycle.)

Now suppose each GC cycle happens less often, once every 2 cpu-seconds. Then, our example application, in the steady state, will have a 30 MiB total heap size on each GC cycle, since it'll allocate 20 MiB in that time. But with every GC cycle, the GC will still only need 0.1 cpu-seconds to scan the 10 MiB of live memory. Again, we're assuming the live heap size stays the same, regardless of how much memory is allocated. So this means that our GC overhead just went down, from 10% to 5%, at the cost of 50% more memory being used.

This change in overheads is the fundamental time/space trade-off mentioned earlier. And GC frequency is at the center of this trade-off: if we execute the GC more frequently, then we use less memory, and vice versa. But how often does the GC actually execute? In Go, deciding when the GC should start is the main parameter which the user has control over.

# GOGC

At a high level, GOGC determines the trade-off between GC CPU and memory.

It works by determining the target heap size after each GC cycle, a target value for the total heap size in the next cycle. The GC's goal is to finish a collection cycle before the total heap size exceeds the target heap size. Total heap size is defined as the live heap size at the end of the previous cycle, plus any new heap memory allocated by the application since the previous cycle. Meanwhile, target heap memory is defined as:

Target heap memory = Live heap + (Live heap + GC roots) * GOGC / 100

As an example, consider a Go program with a live heap size of 8 MiB, 1 MiB of goroutine stacks, and 1 MiB of pointers in global variables. Then, with a GOGC value of 100, the amount of new memory that will be allocated before the next GC runs will be 10 MiB, or 100% of the 10 MiB of work, for a total heap footprint of 18 MiB. With a GOGC value of 50, then it'll be 50%, or 5 MiB. With a GOGC value of 200, it'll be 200%, or 20 MiB.

Note: GOGC includes the root set only as of Go 1.18. Previously, it would only count the live heap. Often, the amount of memory in goroutine stacks is quite small and the live heap size dominates all other sources of GC work, but in cases where programs had hundreds of thousands of goroutines, the GC was making poor judgements.

The heap target controls GC frequency: the bigger the target, the longer the GC can wait to start another mark phase and vice versa. While the precise formula is useful for making estimates, it's best to think of GOGC in terms of its fundamental purpose: a parameter that picks a point in the GC CPU and memory trade-off. The key takeaway is that doubling GOGC will double heap memory overheads and roughly halve GC CPU cost, and vice versa. (To see a full explanation as to why, see the appendix.)

Note: the target heap size is just a target, and there are several reasons why the GC cycle might not finish right at that target. For one, a large enough heap allocation can simply exceed the target. However, other reasons appear in GC implementations that go beyond the GC model this guide has been using thus far. For some more detail, see the latency section, but the complete details may be found in the additional resources.

GOGC may be configured through either the GOGC environment variable (which all Go programs recognize), or through the SetGCPercent API in the runtime/debug package.

Note that GOGC may also be used to turn off the GC entirely (provided the memory limit does not apply) by setting GOGC=off or calling SetGCPercent(-1). Conceptually, this setting is equivalent to setting GOGC to a value of infinity, as the amount of new memory before a GC is triggered is unbounded.

# Memory limit

 GOGC was the sole parameter that could be used to modify the GC's behavior. While it works great as a way to set a trade-off, it doesn't take into account that available memory is finite. Consider what happens when there's a transient spike in the live heap size: because the GC will pick a total heap size proportional to that live heap size, GOGC must be configured such for the peak live heap size, even if in the usual case a higher GOGC value provides a better trade-off.

 Go added support for setting a runtime memory limit. The memory limit may be configured either via the GOMEMLIMIT environment variable which all Go programs recognize, or through the SetMemoryLimit function available in the runtime/debug package.

This memory limit sets a maximum on the total amount of memory that the Go runtime can use. The specific set of memory included is defined in terms of runtime.MemStats as the expression

Sys - HeapReleased

or equivalently in terms of the runtime/metrics package,

/memory/classes/total:bytes - /memory/classes/heap/released:bytes

Because the Go GC has explicit control over how much heap memory it uses, it sets the total heap size based on this memory limit and how much other memory the Go runtime uses.

The visualization below depicts the same single-phase steady state workload from the GOGC section, but this time with an extra 10 MiB of overhead from the Go runtime and with an adjustable memory limit. Try shifting around both GOGC and the memory limit and see what happens.

Notice that when the memory limit is lowered below the peak memory that's determined by GOGC (42 MiB for a GOGC of 100), the GC runs more frequently to keep the peak memory within the limit.

Returning to our previous example of the transient heap spike, by setting a memory limit and turning up GOGC, we can get the best of both worlds: no memory limit breach, and better resource economy. 

Notice that with some values of GOGC and the memory limit, peak memory use stops at whatever the memory limit is, but that the rest of the program's execution still obeys the total heap size rule set by GOGC.

This observation leads to another interesting detail: even when GOGC is set to off, the memory limit is still respected! In fact, this particular configuration represents a maximization of resource economy because it sets the minimum GC frequency required to maintain some memory limit. In this case, all of the program's execution has the heap size rise to meet the memory limit.

Now, while the memory limit is clearly a powerful tool, the use of a memory limit does not come without a cost, and certainly doesn't invalidate the utility of GOGC.

Consider what happens when the live heap grows large enough to bring total memory use close to the memory limit. In the steady state visualization above, try turning GOGC off and then slowly lowering the memory limit further and further to see what happens. Notice that the total time the application takes will start to grow in an unbounded manner as the GC is constantly executing to maintain an impossible memory limit.

This situation, where the program fails to make reasonable progress due to constant GC cycles, is called thrashing. It's particularly dangerous because it effectively stalls the program. Even worse, it can happen for exactly the same situation we were trying to avoid with GOGC: a large enough transient heap spike can cause a program to stall indefinitely! Try reducing the memory limit (around 30 MiB or lower) in the transient heap spike visualization and notice how the worst behavior specifically starts with the heap spike.

In many cases, an indefinite stall is worse than an out-of-memory condition, which tends to result in a much faster failure.

For this reason, the memory limit is defined to be soft. The Go runtime makes no guarantees that it will maintain this memory limit under all circumstances; it only promises some reasonable amount of effort. This relaxation of the memory limit is critical to avoiding thrashing behavior, because it gives the GC a way out: let memory use surpass the limit to avoid spending too much time in the GC.

How this works internally is the GC sets an upper limit on the amount of CPU time it can use over some time window (with some hysteresis for very short transient spikes in CPU use). This limit is currently set at roughly 50%, with a 2 * GOMAXPROCS CPU-second window. The consequence of limiting GC CPU time is that the GC's work is delayed, meanwhile the Go program may continue allocating new heap memory, even beyond the memory limit.

The intuition behind the 50% GC CPU limit is based on the worst-case impact on a program with ample available memory. In the case of a misconfiguration of the memory limit, where it is set too low mistakenly, the program will slow down at most by 2x, because the GC can't take more than 50% of its CPU time away.

Note: the visualizations on this page do not simulate the GC CPU limit.

# Suggested uses

While the memory limit is a powerful tool, and the Go runtime takes steps to mitigate the worst behaviors from misuse, it's still important to use it thoughtfully. Below is a collection of tidbits of advice about where the memory limit is most useful and applicable, and where it might cause more harm than good.

- Do take advantage of the memory limit when the execution environment of your Go program is entirely within your control, and the Go program is the only program with access to some set of resources (i.e. some kind of memory reservation, like a container memory limit).

A good example is the deployment of a web service into containers with a fixed amount of available memory.

In this case, a good rule of thumb is to leave an additional 5-10% of headroom to account for memory sources the Go runtime is unaware of.

- Do feel free to adjust the memory limit in real time to adapt to changing conditions.

A good example is a cgo program where C libraries temporarily need to use substantially more memory.

- Don't set GOGC to off with a memory limit if the Go program might share some of its limited memory with other programs, and those programs are generally decoupled from the Go program. Instead, keep the memory limit since it may help to curb undesirable transient behavior, but set GOGC to some smaller, reasonable value for the average case.

While it may be tempting to try and "reserve" memory for co-tenant programs, unless the programs are fully synchronized (e.g. the Go program calls some subprocess and blocks while its callee executes), the result will be less reliable as inevitably both programs will need more memory. Letting the Go program use less memory when it doesn't need it will generate a more reliable result overall. This advice also applies to overcommit situations, where the sum of memory limits of containers running on one machine may exceed the actual physical memory available to the machine.

- Don't use the memory limit when deploying to an execution environment you don't control, especially when your program's memory use is proportional to its inputs.

A good example is a CLI tool or a desktop application. Baking a memory limit into the program when it's unclear what kind of inputs it might be fed, or how much memory might be available on the system can lead to confusing crashes and poor performance. Plus, an advanced end-user can always set a memory limit if they wish.

- Don't set a memory limit to avoid out-of-memory conditions when a program is already close to its environment's memory limits.

This effectively replaces an out-of-memory risk with a risk of severe application slowdown, which is often not a favorable trade, even with the efforts Go makes to mitigate thrashing. In such a case, it would be much more effective to either increase the environment's memory limits (and then potentially set a memory limit) or decrease GOGC (which provides a much cleaner trade-off than thrashing-mitigation does).

# Latency

The Go GC, however, is not fully stop-the-world and does most of its work concurrently with the application. This is primarily to reduce application latencies. Specifically, the end-to-end duration of a single unit of computation (e.g. a web request). Thus far, this document mainly considered application throughput (e.g. web requests handled per second). Note that each example in the GC cycle section focused on the total CPU duration of an executing program. However, such a duration is far less meaningful for say, a web service. While throughput is still important for a web service (i.e. queries per second), often the latency of each individual request matters even more.

In terms of latency, a stop-the-world GC may require a considerable length of time to execute both its mark and sweep phases, during which the application, and in the context of a web service, any in-flight request, is unable to make further progress. Instead, the Go GC avoids making the length of any global application pauses proportional to the size of the heap, and that the core tracing algorithm is performed while the application is actively executing. (The pauses are more strongly proportional to GOMAXPROCS algorithmically, but most commonly are dominated by the time it takes to stop running goroutines.) Collecting concurrently is not without cost: in practice it often leads to a design with lower throughput than an equivalent stop-the-world garbage collector. However, it's important to note that lower latency does not inherently mean lower throughput, and the performance of the Go garbage collector has steadily improved over time, in both latency and throughput.

The concurrent nature of Go's current GC does not invalidate anything discussed in this document so far: none of the statements relied on this design choice. GC frequency is still the primary way the GC trades off between CPU time and memory for throughput, and in fact, it also takes on this role for latency. This is because most of the costs for the GC are incurred while the mark phase is active.

The key takeaway then, is that reducing GC frequency may also lead to latency improvements. This applies not only to reductions in GC frequency from modifying tuning parameters, like increasing GOGC and/or the memory limit, but also applies to the optimizations described in the optimization guide.

However, latency is often more complex to understand than throughput, because it is a product of the moment-to-moment execution of the program and not just an aggregation of costs. As a result, the connection between latency and GC frequency is less direct. Below is a list of possible sources of latency for those inclined to dig deeper.

Brief stop-the-world pauses when the GC transitions between the mark and sweep phases,
Scheduling delays because the GC takes 25% of CPU resources when in the mark phase,
User goroutines assisting the GC in response to a high allocation rate,
Pointer writes requiring additional work while the GC is in the mark phase, and
Running goroutines must be suspended for their roots to be scanned.
These latency sources are visible in execution traces, except for pointer writes requiring additional work.

# Finalizers, cleanups, and weak pointers

Garbage collection provides the illusion of infinite memory using only finite memory. Memory is allocated but never explicitly freed, which enables simpler APIs and concurrent algorithms compared to bare-bones manual memory management. (Some languages with manually managed memory use alternative approaches such as "smart pointers" and compile-time ownership tracking to ensure that objects are freed, but these features are deeply embedded into the API design conventions in these languages.)

Only the live objects—those reachable from a global variable or a computation in some goroutine—can affect the behavior of the program. Any time after an object becomes unreachable ("dead"), it may be safely recycled by the GC. This allows for a wide variety of GC designs, such as the tracing design used by Go today. The death of an object is not an observable event at the language level.

However, Go's runtime library provides three features that break that illusion: cleanups, weak pointers, and finalizers. Each of these features provides some way to observe and react to object death, and in the case of finalizers, even reverse it. This of course complicates Go programs and adds an additional burden to the GC implementation. Nonetheless, these features exist because they are useful in a variety of circumstances, and Go programs use them and benefit from them all the time.

For the details of each feature, refer to its package documentation (runtime.AddCleanup, weak.Pointer, runtime.SetFinalizer). Below is some general advice for using these features, outlines of common issues you can run into with each feature, and advice for testing uses of these features

-> General Advice

    * Write unit tests: The exact timing of cleanups, weak pointers, and finalizers can be difficult to predict, and it's easy to convince yourself that everything works, even after many consecutive executions. But it's also easy to make subtle mistakes. Writing tests for them can be tricky, but given that they're so subtle to use, testing is even more important usual.

    * Avoid using these features directly in typical Go code : These are low-level features with subtle restrictions and behaviors. For instance, there's no guarantee cleanups or finalizers will be run at program exit, or at all for that matter. The long comments in their API documentation should be seen as a warning. The vast majority of Go code does not benefit from using these features directly, only indirectly.

    * Encapsulate the use of these mechanisms within a package : Where possible, do not allow the use of these mechanisms to leak into the public API of your package; provide interfaces that make it hard or impossible for users to misuse them. For example, instead of asking the user to set up a cleanup on some C-allocated memory to free it, write a wrapper package and hide that detail inside.

    * Restrict access to objects that have finalizers, cleanups, and weak pointers to the package that created and applied them. : This is related to the previous point, but is worth calling out explicitly, since it's a very powerful pattern for using these features in a less error-prone way. For example, the unique package uses weak pointers under the hood, but completely encasulates the objects that are weakly pointed-to. Those values can never be mutated by the rest of the application, it can only be copied through the Value method, preserving the illusion of infinite memory for package users.

    * Prefer cleaning up non-memory resources deterministically when possible, with finalizers and cleanups as a fallback : Cleanups and finalizers are a good fit for memory resources such as memory allocated externally, like from C, or references to an mmap mapping. Memory allocated by C's malloc must eventually be freed by C's free. A finalizer that calls free, attached to a wrapper object for the C memory, is a reasonable way to ensure that C memory is eventually reclaimed as a consequence of garbage collection.

However, non-memory resources, like file descriptors, tend to be subject to system limits that the Go runtime is generally unaware of. In addition, the timing of the garbage collector in a given Go program is usually something a package author has little control over (for instance, how often the GC runs is controlled by GOGC, which can be set by operators to a variety of different values in practice). These two facts conspire to make cleanups and finalizers a bad fit to use as the only mechanism for releasing non-memory resources.

If you're a package author exposing an API that wraps some non-memory resource, consider providing an explicit API for releasing the resource deterministically (through a Close method, or something similar), rather than relying on the garbage collector through cleanups or finalizers. Instead, prefer to use cleanups and finalizers as a best-effort handler for programmer mistakes, either by cleaning up the resource anyway like os.File does, or by reporting the failure to deterministically clean up back to the user.

    * Prefer cleanups to finalizers: Historically, finalizers were added to simplify the interface between Go code and C code and to clean up non-memory resources. The intended use was to apply them to wrapper objects that owned C memory or some other non-memory resource, so that the resource could be released once Go code was done using it. These reasons at least partially explain why finalizers are narrowly scoped, why any given object can only have one finalizer, and why that finalizer must be attached to the first byte of the object only. This limitation already stifles some use-cases. For example, any package that wishes to internally cache some information about an object passed to it cannot clean up that information once the object is gone.

But worse than that, finalizers are inefficient and error-prone due to the fact that they resurrect the object they're attached to, so that it can be passed to the finalizer function (and even continue to live beyond that, too). This simple fact means that if the object is part of a reference cycle it can never be freed, and the memory backing the object cannot be reused until at least until the following garbage collection cycle.

Because finalizers resurrect objects, though, they do have a better-defined execution order than cleanups. For this reason, finalizers are still potentially (but rarely) useful for cleaning up structures that have complex destruction ordering requirements.

But for all other uses in Go 1.24 and beyond, we recommend you use cleanups because they are more flexible, less error-prone, and more efficient than finalizers.

