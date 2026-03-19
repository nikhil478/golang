# Additional notes on GOGC

The GOGC section claimed that doubling GOGC doubles heap memory overheads and halves GC CPU costs. To see why, let's break it down mathematically.

Firstly, the heap target sets a target for the total heap size. This target, however, mainly influences the new heap memory, because the live heap is fundamental to the application.

Target heap memory = Live heap + (Live heap + GC roots) * GOGC / 100

Total heap memory = Live heap + New heap memory

⇒
New heap memory = (Live heap + GC roots) * GOGC / 100

From this we can see that doubling GOGC would also double the amount of new heap memory that application will allocate each cycle, which captures heap memory overheads. Note that Live heap + GC roots is an approximation of the amount of memory the GC needs to scan.

Next, let's look at GC CPU cost. Total cost can be broken down as the cost per cycle, times GC frequency over some time period T.

Total GC CPU cost = (GC CPU cost per cycle) * (GC frequency) * T

GC CPU cost per cycle can be derived from the GC model:

GC CPU cost per cycle = (Live heap + GC roots) * (Cost per byte) + Fixed cost

Note that sweep phase costs are ignored here as mark and scan costs dominate.

The steady state is defined by a constant allocation rate and a constant cost per byte, so in the steady state we can derive a GC frequency from this new heap memory:

GC frequency = (Allocation rate) / (New heap memory) = (Allocation rate) / ((Live heap + GC roots) * GOGC / 100)

Putting this together, we get the full equation for the total cost:

Total GC CPU cost = (Allocation rate) / ((Live heap + GC roots) * GOGC / 100) * ((Live heap + GC roots) * (Cost per byte) + Fixed cost) * T

For a sufficiently large heap (which represents most cases), the marginal costs of a GC cycle dominate the fixed costs. This allows for a significant simplification of the total GC CPU cost formula.

Total GC CPU cost = (Allocation rate) / (GOGC / 100) * (Cost per byte) * T

From this simplified formula, we can see that if we double GOGC, we halve total GC CPU cost. (Note that the visualizations in this guide do simulate fixed costs, so the GC CPU overheads reported by them will not exactly halve when GOGC doubles.) Furthermore, GC CPU costs are largely determined by allocation rate and the cost per byte to scan memory. For more information on how to reduce these costs specifically, see the optimization guide.

Note: there exists a discrepancy between the size of the live heap, and the amount of that memory the GC actually needs to scan: the same size live heap but with a different structure will result in a different CPU cost, but the same memory cost, resulting a different trade-off. This is why the structure of the heap is part of the definition of the steady state. The heap target should arguably only include the scannable live heap as a closer approximation of memory the GC needs to scan, but this leads to degenerate behavior when there's a very small amount of scannable live heap but the live heap is otherwise large.