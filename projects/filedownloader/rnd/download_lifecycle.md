# Download Lifecycle (~6.8 min)

```mermaid
sequenceDiagram
    participant G as Goroutine (Go)
    participant S as Scheduler (G-M-P)
    participant N as Netpoller
    participant OS as OS / Kernel
    participant NET as Network
    participant DISK as Disk

    %% Initialization
    Note over G,S: Go runtime initializes<br/>Goroutine created (~2KB stack)
    Note over N: Netpoller initialized (epoll/kqueue/IOCP)

    %% DNS
    G->>OS: DNS query (via Resolver)
    G->>S: Park goroutine
    OS->>NET: UDP DNS request
    NET-->>OS: DNS response
    OS-->>N: Socket ready
    N-->>G: Wake goroutine

    %% TCP Handshake
    G->>OS: Dial() (non-blocking connect)
    G->>S: Park goroutine
    OS->>NET: SYN
    NET-->>OS: SYN-ACK
    OS->>NET: ACK
    OS-->>N: Writable socket
    N-->>G: Wake goroutine

    %% TLS Handshake
    G->>NET: TLS handshake (multiple RTTs)
    Note over G: CPU-heavy (crypto/tls)

    %% HTTP Request
    G->>OS: write() request
    OS->>NET: Send HTTP request

    %% Slow Start
    loop TCP Slow Start (~0–3s)
        G->>OS: read()
        alt No data
            G->>S: Park goroutine
            OS-->>N: Waiting for data
            N-->>G: Wake on data arrival
        else Data available
            OS-->>G: Data (memcpy)
        end
    end

    %% Steady State Transfer
    loop io.Copy (~30K iterations)
        G->>OS: read()
        alt No data
            G->>S: Park
            N-->>G: Wake on readiness
        else Data ready
            OS-->>G: Data (kernel → user)
        end

        G->>OS: write()
        OS->>DISK: Page cache write (user → kernel)
    end

    %% Backpressure
    alt Disk slow
        OS-->>G: write() blocks
        G->>S: Park
    else Network slow
        G->>S: Park on read()
        N-->>G: Wake when data arrives
    end

    %% Final Phase
    G->>OS: read()
    OS-->>G: EOF
    G->>OS: file.Close()

    OS->>DISK: Flush page cache
    Note over DISK: Final disk write

    %% Completion
    G-->>S: Goroutine exits
    OS->>NET: FIN (close socket)
    OS-->>G: Resources released

    Note over G,DISK: File fully written and durable
````