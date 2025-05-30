# 🧠 Simple In-Memory DB Server

A lightweight implementation of an in-memory database server written in Python. It supports basic data storage, retrieval, and manipulation operations using different data structures like strings, lists, and hashes — similar in spirit to Redis, but simplified for learning and experimentation.

---

## 🚀 Features

- ✅ In-memory key-value storage
- 🧵 Supports:
  - Strings (`SET`, `GET`)
  - Lists (`LPUSH`, `RPUSH`, `LPOP`, `RPOP`)
  - Hashes (`HSET`, `HGET`, `HDEL`)
- 💬 Simple TCP-based command protocol
- 🛠 Modular design for easy extension
- ⚡ Fast execution with no disk persistence

---

## 📦 Commands Supported

| Command              | Description                        |
|----------------------|------------------------------------|
| `SET key val`        | Set a string value                 |
| `GET key`            | Get a string value                 |
| `LPUSH key val`      | Push value to the head of a list   |
| `RPUSH key val`      | Push value to the tail of a list   |
| `LPOP key`           | Pop value from the head of a list  |
| `RPOP key`           | Pop value from the tail of a list  |
| `HSET key field val` | Set a hash field                   |
| `HGET key field`     | Get a hash field                   |
| `HDEL key field`     | Delete a hash field                |
| `DEL key`            | Delete a key                       |
| `KEYS`               | List all stored keys               |

---

## 🛠 Architecture

```
+-----------------------+
|  Client (e.g., netcat)|  <---> TCP socket
+-----------------------+
             |
             v
+-----------------------+
|  Command Parser       |
+-----------------------+
             |
             v
+-----------------------+
|  Dispatcher           | --> routes to correct handler
+-----------------------+
     |       |       |
     v       v       v
  Strings   Lists   Hashes
   Dict      List    Dict of Dicts
```

---

## 🧪 Example Usage

You can test the server using telnet or netcat:

```sh
$ nc localhost 6379
SET name ChatDB
OK
GET name
ChatDB
LPUSH users alice
OK
RPUSH users bob
OK
LPOP users
alice
```

---

## 📁 Project Structure

```
.
├── server.py         # Main server loop
├── parser.py         # Parses incoming commands
├── db.py             # In-memory data storage and operations
├── commands/         # Handlers for each command type
│   ├── string.py
│   ├── list.py
│   └── hash.py
└── README.md
```

---

## 🧑‍💻 Getting Started

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/simple-inmem-db.git
   cd simple-inmem-db
   ```

2. Run the server:

   ```sh
   python3 server.py
   ```

3. Connect using netcat or write a custom client.

---

## 🧩 Future Improvements

- Add support for expiration (TTL)
- Add persistence (snapshot or AOF)
- Add transactions or pipelining
- Implement pub/sub model

---

## 📝 License

MIT License. Use it freely for learning or extending.

---
