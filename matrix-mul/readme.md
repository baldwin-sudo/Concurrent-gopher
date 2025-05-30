# Concurrent Matrix Multiplication in Go

This project implements matrix multiplication in Go using two approaches:

1. **Sequential multiplication**
2. **Concurrent multiplication** using a worker pool and chunked task division

The concurrent version assigns rectangular submatrices (chunks) to workers instead of computing each element individually.

## Structure

- `Task` struct: defines a rectangular block of matrix entries to compute.
- `worker`: processes a block and sends results.
- `collector`: gathers computed blocks into the result matrix.
- `SequentialMatrixMultiplication`: baseline reference implementation.
- `ConcurrentMatrixMultiplication`: concurrent version using goroutines and channels.

## Usage

Run the program:

```bash
go run main.go
