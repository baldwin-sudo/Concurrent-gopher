# Mini Redis — Concurrent In-Memory Key-Value Store in Go

## Overview

This project is a minimalist Redis-like key-value store implemented in Go. It supports concurrent clients over TCP and basic Redis commands such as `GET`, `SET`, `DEL`, `QUIT`, and `HELP`. The focus is on practicing Go concurrency patterns, channel-based worker pools, and graceful server shutdown using contexts.

---
  
## Features

- **TCP Server:** Listens on a configurable port and accepts multiple concurrent client connections.
- **Command Parsing:** Supports a simple parser to handle commands with argument validation.
- **Load Balancer:**Uses a pool of workers and balances the load between them , by assigning the next job to the least loaded worker  using a heap datastructure . Workers are concurrent ofc !
- **Shared Store:** Maintains a global thread-safe map protected by mutexes for storing keys and values.
- **Client Management:** Each client has dedicated input reading and output writing goroutines communicating via channels.
- **Graceful Shutdown:** Uses Go's `context` package and OS signal handling to allow clean server shutdown on Ctrl+C.
- **Command Handling:** Supports basic commands with proper responses and error handling.
- **Per-Client Done Signaling:** Clients have a done flag or channel for clean connection teardown and avoiding panics.

---

## Architecture & Components

- **`server.go`:**  
  - Initializes the listener, context, and channels.  
  - Handles incoming TCP connections, spawning goroutines for each client.  
  - Manages worker pool and dispatch loop with cancellation support.  

- **`client.go`:**  
  - Defines the `Client` struct with connection, output channel, and done signaling.  
  - Implements client writer goroutine to serialize writes to TCP connection.

- **`store.go`:**  
  - Implements the shared key-value store with concurrency-safe methods.  

- **`parser.go`:**  
  - Parses raw input lines into typed command requests with argument validation.  

- **`handler.go`:**  
  - Maps parsed requests to store operations and returns formatted responses.  

---

## Concurrency Model

- **Worker Pool:** A fixed number of worker goroutines consume jobs from a buffered `jobs` channel, process them by interacting with the store, and send results on a `results` channel.
- **Dispatcher:** A goroutine listens on `results` and forwards responses to appropriate clients via their output channels.
- **Per-client Goroutines:** Each client connection runs:
  - A reader goroutine parsing input and pushing jobs.
  - A writer goroutine serializing output from the client’s output channel.
- **Context Cancellation:**  
  - A root `context.Context` is used to propagate shutdown signals to all goroutines.  
  - OS signals (`SIGINT`) trigger cancellation and graceful shutdown of the server and workers.

---

## How Shutdown Works

- Ctrl+C triggers `os.Interrupt` signal, captured by a dedicated goroutine.
- This cancels the root context and closes the listener socket, unblocking accept.
- Workers and dispatcher detect context cancellation and exit gracefully.
- Client connections detect closed output channels or done flags and close cleanly.

---

## Usage Recap

- Run server with `go run .`
- Connect with `nc localhost 9999` or `telnet localhost 9999`
- Use commands:  
  - `SET key value`  
  - `GET key`  
  - `DEL key`  
  - `HELP`  
  - `QUIT` to disconnect  
- Server logs connections and shutdown events.

---

## Key Learnings / Reminders

- Use **Go contexts** for clean cancellation across multiple goroutines.
- Structure concurrency with **channels** for safe communication.
- Separate concerns: parsing, handling, storage, client IO.
- Avoid closing channels from multiple places to prevent panics.
- Start channel-consuming goroutines before sending data to avoid deadlocks.
- Use **worker pools** to scale request processing without spawning unlimited goroutines.
- Handle client disconnects and errors gracefully to avoid leaks.
- Use buffered channels sized to your CPU count for efficient throughput.

