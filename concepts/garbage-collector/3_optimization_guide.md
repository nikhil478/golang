## Optimization guide

Resource : https://go.dev/doc/gc-guide#Optimization_guide

# Identifying costs

Before trying to optimize how your Go application interacts with the GC, it's important to first identify that the GC is a major cost in the first place.

The Go ecosystem provides a number of tools for identifying costs and optimizing Go applications. For a brief overview of these tools, see the guide on diagnostics(https://go.dev/doc/diagnostics). Here, we'll focus on a subset of these tools and a reasonable order to apply them in in order to understand GC impact and behavior.

1. CPU Profiles [https://pkg.go.dev/runtime/pprof#hdr-Profiling_a_Go_program]

A good place to start is with CPU profiling. CPU profiling provides an overview of where CPU time is spent, though to the untrained eye it may be difficult to identify the magnitude of the role the GC plays in a particular application. Luckily, understanding how the GC fits in mostly boils down to knowing what different functions in the `runtime` package mean. Below is a useful subset of these functions for interpreting CPU profiles.

Note that the functions listed below are not leaf functions, so they may not show up in the default the pprof tool provides with the top command. Instead, use the top -cum command or use the list command on these functions directly and focus on the cumulative percent column.

runtime.gcBgMarkWorker: Entrypoint to the background mark worker goroutines. Time spent here scales with GC frequency and the complexity and size of the object graph. It represents a baseline for how much time the application spends marking and scanning.

Note that within these goroutines, you will find calls to runtime.gcDrainMarkWorkerDedicated, runtime.gcDrainMarkWorkerFractional, and runtime.gcDrainMarkWorkerIdle, which indicate worker type. In a largely idle Go application, the Go GC is going to use up additional (idle) CPU resources to get its job done faster, which is indicated with the runtime.gcDrainMarkWorkerIdle symbol. As a result, time here may represent a large fraction of CPU samples, which the Go GC believes are free. If the application becomes more active, CPU time in idle workers will drop. One common reason this can happen is if an application runs entirely in one goroutine but GOMAXPROCS is >1.

runtime.mallocgc: Entrypoint to the memory allocator for heap memory. A large amount of cumulative time spent here (>15%) typically indicates a lot of memory being allocated.

runtime.gcAssistAlloc: Function goroutines enter to yield some of their time to assist the GC with scanning and marking. A large amount of cumulative time spent here (>5%) indicates that the application is likely out-pacing the GC with respect to how fast it's allocating. It indicates a particularly high degree of impact from the GC, and also represents time the application spend marking and scanning. Note that this is included in the runtime.mallocgc call tree, so it will inflate that as well.


2. Execution traces

While CPU profiles are great for identifying where time is spent in aggregate, they're less useful for indicating performance costs that are more subtle, rare, or related to latency specifically. Execution traces on the other hand provide a rich and deep view into a short window of a Go program's execution. They contain a variety of events related to the Go GC and specific execution paths can be directly observed, along with how the application might interact with the Go GC. All the GC events tracked are conveniently labeled as such in the trace viewer.

See the documentation for the runtime/trace (https://pkg.go.dev/runtime/trace) package for how to get started with execution traces.

3. GC traces

When all else fails, the Go GC provides a few different specific traces that provide much deeper insights into GC behavior. These traces are always printed directly to STDERR, one line per GC cycle, and are configured through the GODEBUG environment variable that all Go programs recognize. They're mostly useful for debugging the Go GC itself since they require some familiarity with the specifics of the GC's implementation, but nonetheless can occasionally be useful to gain a better understanding of GC behavior.

The core GC trace is enabled by setting GODEBUG=gctrace=1. The output produced by this trace is documented in the environment variables section in the documentation (https://pkg.go.dev/runtime#hdr-Environment_Variables) for the runtime package.

A supplementary GC trace called the "pacer trace" provides even deeper insights and is enabled by setting GODEBUG=gcpacertrace=1. Interpreting this output requires an understanding of the GC's "pacer" (see additional resources), which is outside the scope of this guide(https://go.dev/doc/gc-guide#Additional_resources).


# Heap profiling

Heap profiling¶

After identifying that the GC is a source of significant costs, the next step in eliminating heap allocations is to find out where most of them are coming from. For this purpose, memory profiles (really, heap memory profiles) are very useful. Check out the documentation for how to get started with them.

Memory profiles describe where in the program heap allocations come from, identifying them by the stack trace at the point they were allocated. Each memory profile can break down memory in four ways.

inuse_objects—Breaks down the number of objects that are live.
inuse_space—Breaks down live objects by how much memory they use in bytes.
alloc_objects—Breaks down the number of objects that have been allocated since the Go program began executing.
alloc_space—Breaks down the total amount of memory allocated since the Go program began executing.
Switching between these different views of heap memory may be done with either the -sample_index flag to the pprof tool, or via the sample_index option when the tool is used interactively.

Note: memory profiles by default only sample a subset of heap objects so they will not contain information about every single heap allocation. However, this is sufficient to find hot-spots. To change the sampling rate, see runtime.MemProfileRate.

For the purposes of reducing GC costs, alloc_space is typically the most useful view as it directly corresponds to the allocation rate. This view will indicate allocation hot spots that would provide the most benefit.

# Escape analysis

Once candidate heap allocation sites have been identified with the help of heap profiles(https://go.dev/doc/gc-guide#Heap_profiling), how can they be eliminated? The key is to leverage the Go compiler's escape analysis to have the Go compiler find alternative, and more efficient storage for this memory, for example in the goroutine stack. Luckily, the Go compiler has the ability to describe why it decides to escape a Go value to the heap. With that knowledge, it becomes a matter of reorganizing your source code to change the outcome of the analysis (which is often the hardest part, but outside the scope of this guide).

As for how to access the information from the Go compiler's escape analysis, the simplest way is through a debug flag supported by the Go compiler that describes all optimizations it applied or did not apply to some package in a text format. This includes whether or not values escape. Try the following command, where [package] is some Go package path.

$ go build -gcflags=-m=3 [package]
This information can also be visualized as an overlay in an LSP-capable editor; it is exposed as a code action. For example, in VS Code, invoke the "Source Action... > Show compiler optimization details" command to enable diagnostics for the current package. (You can also run the "Go: Toggle compiler optimization details" command.) Use this configuration setting to control which annotations are displayed:

Enable the overlay for escape analysis by setting ui.diagnostic.annotations to include escape .
Finally, the Go compiler provides this information in a machine-readable (JSON) format that may be used to build additional custom tooling. For more information on that, see the documentation in the source Go code.[https://github.com/golang/vscode-go/wiki/settings#uidiagnosticannotations]


# Implementation-specific optimizations

The Go GC is sensitive to the demographics of live memory, because a complex graph of objects and pointers both limits parallelism and generates more work for the GC. As a result, the GC contains a few optimizations for specific common structures. The most directly useful ones for performance optimization are listed below.

Note: Applying the optimizations below may reduce the readability of your code by obscuring intent, and may fail to hold up across Go releases. Prefer to apply these optimizations only in the places they matter most. Such places may be identified by using the tools listed in the section on identifying costs.

Pointer-free values are segregated from other values.

As a result, it may be advantageous to eliminate pointers from data structures that do not strictly need them, as this reduces the cache pressure the GC exerts on the program. As a result, data structures that rely on indices over pointer values, while less well-typed, may perform better. This is only worth doing if it's clear that the object graph is complex and the GC is spending a lot of time marking and scanning.

The GC will stop scanning values at the last pointer in the value.

As a result, it may be advantageous to group pointer fields in struct-typed values at the beginning of the value. This is only worth doing if it's clear the application spends a lot of its time marking and scanning. (In theory the compiler can do this automatically, but it is not yet implemented, and struct fields are arranged as written in the source code.)

Furthermore, the GC must interact with nearly every pointer it sees, so using indices into an slice, for example, instead of pointers, can aid in reducing GC costs.

# Linux transparent huge pages (THP)

Resource: https://go.dev/doc/gc-guide#Linux_transparent_huge_pages

When a program accesses memory, the CPU needs to translate the virtual memory addresses it uses into physical memory addresses that refer to the data it was trying to access. To do this, the CPU consults the "page table," a data structure that represents the mapping from virtual to physical memory, managed by the operating system. Each entry in the page table represents an indivisible block of physical memory called a page, hence the name.

Transparent huge pages (THP) is a Linux feature that transparently replaces pages of physical memory backing contiguous virtual memory regions with bigger blocks of memory called huge pages. By using bigger blocks, fewer page table entries are needed to represent the same memory region, improving page table lookup times. However, bigger blocks mean more waste if only a small part of the huge page is used by the system.

When running Go programs in production, enabling transparent huge pages on Linux can improve throughput and latency at the cost of additional memory use. Applications with small heaps tend not to benefit from THP and may end up using a substantial amount of additional memory (as high as 50%). However, applications with big heaps (1 GiB or more) tend to benefit quite a bit (up to 10% throughput) without very much additional memory overhead (1-2% or less). Being aware of your THP settings in either case can be helpful, and experimentation is always recommended.

One can enable or disable transparent huge pages in a Linux environment by modifying /sys/kernel/mm/transparent_hugepage/enabled. See the official Linux admin guide for more details. If you choose to have your Linux production environment enable transparent huge pages, we recommend the following additional settings for Go programs.

Set /sys/kernel/mm/transparent_hugepage/defrag to defer or defer+madvise.

This setting controls how aggressively a Linux kernel coalesces regular pages into huge pages. defer tells the kernel to coalesce huge pages lazily and in the background. A more aggressive setting can induce stalls in memory constrained systems and can often hurt application latencies. defer+madvise is like defer, but is friendlier to other applications on the system that request huge pages explicitly and require them for performance.

Set /sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none to 0.

This setting controls how many additional pages the Linux kernel daemon can allocate when trying to allocate a huge page. The default setting is maximally aggressive, and can often undo work the Go runtime does to return memory to the OS. Before Go 1.21, the Go runtime tried to mitigate the negative effects of the default setting, but it came with a CPU cost. With Go 1.21+ and Linux 6.2+, the Go runtime no longer mutates huge page state.

If you experience an increase in memory usage when upgrading to Go 1.21.1 or later, try applying this setting; it will likely resolve your issue. As an additional workaround, you can call the Prctl function with PR_SET_THP_DISABLE to disable huge pages at the process level, or you can set GODEBUG=disablethp=1 (to be added in Go 1.21.6 and Go 1.22) to disable huge pages for heap memory. Note that the GODEBUG setting may be removed in a future release.