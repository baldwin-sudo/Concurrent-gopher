
# Go Concurrency Patterns — My Learning Journey

Exploration of **concurrency patterns in Go**! 

This repository contains projects I've built (or i'am currently  building) while learning how to leverage Go's powerful concurrency features, including goroutines, channels, sync primitives, and more. The projects here demonstrate real-world use cases where concurrency helps improve performance and efficiency.

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

