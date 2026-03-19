### Virtual Memory 

Resource: https://go.dev/doc/gc-guide#A_note_about_virtual_memory

we are focused on the physical memory use of the GC, but a question that comes up regularly is what exactly that means and how it compares to virtual memory (typically presented in programs like top as "VSS")

Physical memory is memory housed in the actual physical RAM chip in most computers. Virtual memory is an abstraction over physical memory provided by the operating system to isolate programs from one another. It's also typically acceptable for programs to reserve virtual address space that doesn't map to any physical addresses at all

Because virtual memory is just a mapping maintained by the operating system, it is typically very cheap to make large virtual memory reservations that don't map to physical memory.

The Go runtime generally relies upon this view of the cost of virtual memory in a few ways:

The Go runtime never deletes virtual memory that it maps. Instead, it uses special operations that most operating systems provide to explicitly release any physical memory resources associated with some virtual memory range.

This technique is used explicitly to manage the memory limit and return memory to the operating system that the Go runtime no longer needs. The Go runtime also releases memory it no longer needs continuously in the background. See the additional resources for more information.

On 32-bit platforms, the Go runtime reserves between 128 MiB and 512 MiB of address space up-front for the heap to limit fragmentation issues.

The Go runtime uses large virtual memory address space reservations in the implementation of several internal data structures. On 64-bit platforms, these typically have a minimum virtual memory footprint of about 700 MiB. On 32-bit platforms, their footprint is negligible.

As a result, virtual memory metrics such as "VSS" in top are typically not very useful in understanding a Go program's memory footprint. Instead, focus on "RSS" and similar measurements, which more directly reflect physical memory usage.