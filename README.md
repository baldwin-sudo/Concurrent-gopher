# Go Concurrency Patterns — My Learning Journey

This repository contains projects I've built (or am currently building) while exploring **concurrency patterns in Go**! These projects showcase practical use cases leveraging goroutines, channels, sync primitives, and more to improve performance and efficiency.

---

## Projects

### 1. Concurrent File Counter

Counts **lines, words, and characters** in all files within a specified directory.

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
```bash
go run main.go -host=target_hostname_or_ip
```

⸻

### 3. Concurrent Load Balancer Simulation (Work In Progress)

An exploration of load balancing techniques using Go concurrency features.

Simple Version:
- Implements a worker pool pattern where incoming requests are distributed evenly among workers.
- Uses channels to send requests and receive responses.
- Demonstrates basic load distribution with goroutines and WaitGroups.

Advanced Version (Planned):
- A more sophisticated load balancer that assigns requests to workers based on a heap-based priority queue or load metrics.
- Features:
- Dynamic worker selection
- Backpressure handling
- Timeouts
- Aims to achieve better balancing under varying workloads.

⸻

### 4. Concurrent Matrix Multiplication

Demonstrates matrix multiplication with two approaches:
- Sequential matrix multiplication for baseline performance.
- Concurrent matrix multiplication that splits the task into rectangular chunks and distributes them across a worker pool.

Key Features:
- Utilizes goroutines, channels, and sync.WaitGroup.
- Divides multiplication work into blocks rather than individual elements for better performance balance.
- Dynamically adapts the number of workers based on runtime.NumCPU.

Usage:
```bash
go run main.go
```
Generates two 2000x2000 matrices, multiplies them sequentially and concurrently, and prints the time taken by each approach.

⸻

### 5. Mini Redis — Concurrent In-Memory Key-Value Store

A lightweight Redis-like server written from scratch using Go’s concurrency primitives.

Core Features:
- TCP server that handles multiple client connections.
- Built-in command parser supporting GET, SET, DEL, HELP, and QUIT.
- Worker pool that processes commands concurrently using channels.
- Shared key-value store protected by a sync.RWMutex.
- Per-client Output channel and isDone flag for clean connection management.
- Graceful shutdown via context.Context and os.Signal.

Concurrency Model:
- Main server accepts clients and spawns a goroutine for each.
- Jobs are sent to workers via a buffered channel, results are routed back to clients.
- The dispatcher monitors the results and delivers output to each client’s writer goroutine.
- context.WithCancel is used to cancel all goroutines when the server receives a shutdown signal (e.g. Ctrl+C).

Usage:
```bash
go run .
```
Then connect via:
```bash
nc localhost 9999
```
Try commands like:

- SET key value
- GET key
- DEL key
- HELP
- QUIT


⸻

What I’ve Learned
- Structuring concurrent workers using channels and goroutines.
- Using sync.WaitGroup for graceful goroutine management.
- Coordinating producer-consumer patterns with buffered channels.
- Comparing sequential vs concurrent execution performance.
- Managing concurrent-safe maps with sync.Map and sync.RWMutex.
- Basic network programming and file system traversal in Go.
- Designing and implementing concurrent load balancers using Go primitives.
- Implementing chunk-based parallelism for matrix computations.
- Building a concurrent TCP server with job dispatching and graceful shutdown using context.

