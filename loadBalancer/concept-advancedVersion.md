# 🌐 Concurrent Server with Load Balancer (Go)

This server architecture handles multiple client connections concurrently and uses a load balancer to distribute requests to workers.

---

## 🧭 Workflow

1. **Server Listens**
   - Listens for incoming TCP connections on a port.

2. **Client Connection Handler**
   - For each connection, spawn a goroutine.
   - Read the request from the connection (`conn` implements `io.Reader`).
   - Create a `Request` struct with:
     - The task data (e.g., parsed command).
     - A `resultChan` (channel for response).
   - Send the request to the load balancer.

3. **Load Balancer**
   - Receives requests and forwards them to workers (e.g., via round-robin or other strategy).

4. **Worker**
   - Processes the request.
   - Sends the result back via `resultChan`.

5. **Client Handler**
   - Waits on `resultChan`.
   - Writes the response back to the client (`conn` also implements `io.Writer`).

---

## ✅ Why It Works

- `conn` is both `Reader` and `Writer` compatible.
- Each handler goroutine owns the connection and its response channel.
- No shared state or mutexes required.
- Enables high concurrency and clean separation of responsibilities.

---

## 🧠 Go Concepts Used

- Goroutines
- Channels (including channels of channels)
- TCP networking (`net.Conn`)
- Load balancing strategies (round-robin, etc.)
