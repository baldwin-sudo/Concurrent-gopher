
# Go Concurrency Patterns : Learning Journey

Exploration of **concurrency patterns in Go**!

This repository contains projects I've built (or am currently building) while learning how to leverage Go's powerful concurrency features, including goroutines, channels, sync primitives, and more. 

---

## Projects

### 1. Concurrent File Counter

This project counts **lines, words, and characters** in all files within a specified directory.

- Uses **goroutines** and **channels** to process files concurrently.
- Employs a **worker pool pattern** with buffered channels to balance workload.
- Includes both **concurrent** and **sequential** implementations to compare performance.
- Measures and displays elapsed time for both approaches.

**Usage:**

```bash
go run main.go -dir=your_directory_path
```

⸻

### 2. Concurrent TCP Port Scanner

A simple port scanner that:
- Scans a set of common TCP ports on a target host.
- Implements both sequential and concurrent scanning using goroutines.
- Uses sync.Map and WaitGroup for safe concurrent state management.
- Compares timing results between sequential and concurrent scans.

Usage:
=======
- Scans a set of common TCP ports on a target host.
- Implements both sequential and concurrent scanning using goroutines.
- Uses `sync.Map` and `WaitGroup` for safe concurrent state management.
- Compares timing results between sequential and concurrent scans.

**Usage:**

```bash
go run main.go -host=target_hostname_or_ip
```

⸻

3. Concurrent Load Balancer Simulation (Work In Progress)

---

### 3. Concurrent Load Balancer Simulation *(Work In Progress)*

An exploration of load balancing techniques using Go concurrency features.

#### Simple Version

- **Implements a worker pool pattern** where incoming requests are distributed evenly among workers.
- Uses **channels** to send requests and receive responses.
- Demonstrates basic load distribution with **goroutines** and **WaitGroups**.

#### Advanced Version *(Planned)*

- A more sophisticated load balancer that assigns requests to workers based on a **heap-based priority queue** or **load metrics**.
- Features:
    - **Dynamic worker selection**
    - **Backpressure handling**
    - **Timeouts**
- Aims to achieve better balancing under varying workloads.



## What I’ve Learned
- How to structure concurrent workers using channels and goroutines.
- Using sync.WaitGroup for graceful goroutine management.
- Coordinating producer-consumer patterns with buffered channels.
- Comparing sequential vs concurrent execution performance.
- Managing concurrent safe maps with sync.Map.
- Basic network programming and file system traversal in Go.
- Designing and implementing concurrent load balancers using Go primitives.
=======

### 4. Concurrent Matrix Multiplication

This project demonstrates matrix multiplication using two approaches:

- **Sequential matrix multiplication** for baseline performance.
- **Concurrent matrix multiplication** that splits the task into rectangular chunks and distributes them across a worker pool.

**Key Features:**

- Utilizes goroutines, channels, and `sync.WaitGroup`.
- Divides the multiplication work into blocks rather than individual elements for better performance balance.
- Dynamically adapts the number of workers based on `runtime.NumCPU`.

**Usage:**

```bash
go run main.go
```

Generates two 2000x2000 matrices, multiplies them both sequentially and concurrently, and prints the time taken by each approach.

---

## What I’ve Learned

- How to structure concurrent workers using channels and goroutines.
- Using `sync.WaitGroup` for graceful goroutine management.
- Coordinating producer-consumer patterns with buffered channels.
- Comparing sequential vs concurrent execution performance.
- Managing concurrent-safe maps with `sync.Map`.
- Basic network programming and file system traversal in Go.
- Designing and implementing concurrent load balancers using Go primitives.
- Chunk-based parallelism for matrix computations using concurrency.
### 5. Mini Redis — Concurrent In-Memory Key-Value Store in Go
#### Overview :
This project is a minimalist Redis-like key-value store implemented in Go. It supports concurrent clients over TCP and basic Redis commands such as GET, SET, DEL, QUIT, and HELP. The focus is on practicing Go concurrency patterns, channel-based worker pools, and graceful server shutdown using contexts.
#### Features :

- TCP Server: Listens on a configurable port and accepts multiple concurrent client connections.
- Command Parsing: Supports a simple parser to handle commands with argument validation.
- Load Balancer: Uses a pool of workers and balances the load between them , by assigning the next job to the least loaded worker -using a heap datastructure . Workers are concurrent ofc !
- Shared Store: Maintains a global thread-safe map protected by mutexes for storing keys and values.
- Client Management: Each client has dedicated input reading and output writing goroutines communicating via channels.
- Graceful Shutdown: Uses Go's context package and OS signal handling to allow clean server shutdown on Ctrl+C.
- Command Handling: Supports basic commands with proper responses and error handling.
- Per-Client Done Signaling: Clients have a done flag or channel for clean connection teardown and avoiding panics.
