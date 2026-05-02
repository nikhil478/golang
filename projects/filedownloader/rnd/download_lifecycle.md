Download (Total ~6.8 min)
│
├── 🧵 Go Runtime Initialization (happens instantly at start)
│   │
│   ├── Goroutine created (G)
│   │   ├── Stack allocated (~2 KB)
│   │   └── Scheduled onto M (OS thread) via P (processor)
│   │
│   ├── Scheduler (G-M-P model)
│   │   ├── G = goroutine
│   │   ├── M = OS thread
│   │   └── P = logical processor
│   │
│   └── Netpoller initialized
│       └── Uses epoll/kqueue/IOCP under the hood
│
├── 1. Setup Phase (~0–300 ms)
│   │
│   ├── DNS Resolution (~10–50 ms)
│   │   │
│   │   ├── Go:
│   │   │   ├── net/http → net.Resolver
│   │   │   ├── Goroutine blocks → parked by scheduler
│   │   │   └── M reused for other goroutines
│   │   │
│   │   ├── OS:
│   │   │   └── DNS query via UDP
│   │   │
│   │   └── Resume:
│   │       └── Netpoller wakes goroutine when response arrives
│   │
│   ├── TCP Handshake (~50–150 ms)
│   │   │
│   │   ├── Go:
│   │   │   ├── Dial() → non-blocking connect
│   │   │   ├── Goroutine parked
│   │   │   └── Registered in netpoller
│   │   │
│   │   ├── OS:
│   │   │   ├── SYN/SYN-ACK/ACK
│   │   │   └── Socket buffers allocated
│   │   │
│   │   └── Resume:
│   │       └── epoll/kqueue signals "writable" → goroutine resumes
│   │
│   ├── TLS Handshake (~100–200 ms)
│   │   │
│   │   ├── Go:
│   │   │   ├── crypto/tls runs in user space
│   │   │   ├── CPU-heavy (encryption, key exchange)
│   │   │   └── Goroutine actively running (not parked)
│   │   │
│   │   └── OS:
│   │       └── network round trips
│   │
│   └── HTTP Request Send (~1–5 ms)
│       │
│       ├── Go:
│       │   ├── Serialize request
│       │   └── write() syscall
│       │
│       └── OS:
│           └── send packet
│
├── 2. TCP Slow Start (~0–3 sec)
│   │
│   ├── Go Runtime Behavior:
│   │   │
│   │   ├── read() called
│   │   ├── If no data:
│   │   │   ├── Goroutine parked
│   │   │   └── Registered in netpoller
│   │   │
│   │   └── When data arrives:
│   │       └── Netpoller wakes goroutine
│   │
│   ├── Scheduler:
│   │   ├── switches between goroutines
│   │   └── keeps CPU busy while waiting for network
│   │
│   └── OS:
│       └── TCP congestion window grows
│
├── 3. Steady-State Transfer (~3 sec → ~6.5 min)
│   │
│   ├── Go: io.Copy loop (core execution)
│   │   │
│   │   ├── LOOP:
│   │   │   │
│   │   │   ├── (1) Read from network
│   │   │   │   │
│   │   │   │   ├── Go:
│   │   │   │   │   ├── read() syscall
│   │   │   │   │   ├── If buffer empty → goroutine parked
│   │   │   │   │   └── If data ready → continues
│   │   │   │   │
│   │   │   │   └── OS:
│   │   │   │       └── kernel → user memcpy
│   │   │   │
│   │   │   ├── (2) Write to file
│   │   │   │   │
│   │   │   │   ├── Go:
│   │   │   │   │   └── write() syscall
│   │   │   │   │
│   │   │   │   └── OS:
│   │   │   │       └── user → page cache memcpy
│   │   │   │
│   │   │   └── (3) Loop repeats (~30K times total)
│   │   │
│   │   ├── Goroutine State:
│   │   │   ├── RUNNING (copying data)
│   │   │   ├── WAITING (network)
│   │   │   └── WAITING (disk if slow)
│   │   │
│   │   └── Scheduler:
│   │       ├── preempts long-running goroutines
│   │       └── balances CPU usage
│   │
│   ├── Netpoller Role:
│   │   │
│   │   ├── monitors socket readiness
│   │   ├── integrates with epoll/kqueue/IOCP
│   │   └── wakes goroutine when data available
│   │
│   ├── Memory (Go + OS combined):
│   │   │
│   │   ├── Go heap:
│   │   │   └── ~32 KB buffer reused
│   │   │
│   │   ├── Kernel:
│   │   │   ├── socket buffer
│   │   │   └── page cache
│   │   │
│   │   └── GC:
│   │       ├── minimal pressure (buffer reused)
│   │       └── runs occasionally in background
│   │
│   └── CPU:
│       ├── memcpy (kernel + user)
│       ├── TLS decrypt (if HTTPS)
│       └── scheduler overhead (tiny)
│
├── 4. Page Cache Growth (parallel)
│   │
│   ├── Go:
│   │   └── unaware (handled by OS)
│   │
│   └── OS:
│       ├── accumulates writes in RAM
│       └── flushes asynchronously
│
├── 5. Backpressure (continuous)
│   │
│   ├── If disk slow:
│   │   │
│   │   ├── Go:
│   │   │   ├── write() blocks
│   │   │   └── goroutine parked
│   │   │
│   │   └── Scheduler:
│   │       └── runs other goroutines
│   │
│   ├── If network slow:
│   │   │
│   │   ├── Go:
│   │   │   └── read() blocks → parked
│   │   │
│   │   └── Netpoller:
│   │       └── wakes when data arrives
│   │
│   └── OS:
│       └── TCP flow control adjusts speed
│
├── 6. Final Phase (~last seconds)
│   │
│   ├── Go:
│   │   ├── final read() returns EOF
│   │   ├── loop exits
│   │   └── file.Close()
│   │
│   ├── OS:
│   │   ├── remaining page cache flushed
│   │   └── disk write completes
│   │
│   └── CPU:
│       └── small spike (flush)
│
└── 7. Completion
    │
    ├── Go:
    │   ├── goroutine exits
    │   ├── stack freed
    │   └── GC may reclaim memory later
    │
    ├── OS:
    │   ├── socket closed (FIN)
    │   ├── buffers released
    │   └── file descriptor closed
    │
    └── Final State:
        └── file fully written + durable