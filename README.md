

# Go Concurrency Patterns — My Learning Journey

Welcome to my exploration of **concurrency patterns in Go**! 🚀

This repository contains projects I've built while learning how to leverage Go's powerful concurrency features, including goroutines, channels, sync primitives, and more. The projects here demonstrate real-world use cases where concurrency helps improve performance and efficiency.

---

## Projects

### 1. Concurrent File Counter

This project counts **lines, words, and characters** in all files within a specified directory.

* Uses **goroutines** and **channels** to process files concurrently.
* Employs a **worker pool pattern** with buffered channels to balance workload.
* Includes both **concurrent** and **sequential** implementations to compare performance.
* Measures and displays elapsed time for both approaches.

*Usage:*

```bash
go run main.go -dir=your_directory_path
```

---

### 2. Concurrent TCP Port Scanner

A simple port scanner that:

* Scans a set of **common TCP ports** on a target host.
* Implements both **sequential** and **concurrent** scanning using goroutines.
* Uses **sync.Map** and **WaitGroup** for safe concurrent state management.
* Compares timing results between sequential and concurrent scans.

*Usage:*

```bash
go run main.go -host=target_hostname_or_ip
```

---

## What I've Learned

* How to structure concurrent workers using channels and goroutines.
* Using **sync.WaitGroup** for graceful goroutine management.
* Coordinating producer-consumer patterns with buffered channels.
* Comparing sequential vs concurrent execution performance.
* Managing concurrent safe maps with `sync.Map`.
* Basic network programming and file system traversal in Go.

