package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Task struct {
	iStart int
	iEnd   int
	jStart int
	jEnd   int
}

type Result struct {
	// Not needed anymore, results written directly by worker
	// Could use a shared matrix with synchronization or pre-allocated matrix
}

// worker now computes a block of the result matrix
func worker(tasks <-chan Task, A, B [][]float64, C [][]float64, wg *sync.WaitGroup) {
	defer wg.Done()
	n := len(A[0]) // common dimension

	for task := range tasks {
		for i := task.iStart; i < task.iEnd; i++ {
			for j := task.jStart; j < task.jEnd; j++ {
				sum := 0.0
				for k := 0; k < n; k++ {
					sum += A[i][k] * B[k][j]
				}
				C[i][j] = sum
			}
		}
	}
}

// generator creates chunked tasks covering the whole matrix C (m x p)
func createTasks(m, p, chunkSize int) []Task {
	var tasks []Task
	for i := 0; i < m; i += chunkSize {
		iEnd := i + chunkSize
		if iEnd > m {
			iEnd = m
		}
		for j := 0; j < p; j += chunkSize {
			jEnd := j + chunkSize
			if jEnd > p {
				jEnd = p
			}
			tasks = append(tasks, Task{iStart: i, iEnd: iEnd, jStart: j, jEnd: jEnd})
		}
	}
	return tasks
}

// concurrent multiplication using chunked tasks
func ConcurrentMatrixMultiplication(A, B [][]float64) [][]float64 {
	m := len(A)
	p := len(B[0])
	C := make([][]float64, m)
	for i := range C {
		C[i] = make([]float64, p)
	}

	numWorkers := runtime.NumCPU() * 2
	chunkSize := 50 // adjust chunk size for better performance (experiment)

	tasks := make(chan Task, numWorkers*2)

	var wg sync.WaitGroup

	// Start workers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(tasks, A, B, C, &wg)
	}

	// Send chunked tasks
	go func() {
		for _, task := range createTasks(m, p, chunkSize) {
			tasks <- task
		}
		close(tasks)
	}()

	wg.Wait()
	return C
}

// Sequential matrix multiplication (unchanged)
func SequentialMatrixMultiplication(A, B [][]float64) [][]float64 {
	m := len(A)
	n := len(A[0])
	p := len(B[0])
	C := make([][]float64, m)
	for i := range C {
		C[i] = make([]float64, p)
		for j := 0; j < p; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += A[i][k] * B[k][j]
			}
			C[i][j] = sum
		}
	}
	return C
}

// generate matrix (unchanged)
func generateMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64() * 10
		}
	}
	return matrix
}

func main() {
	rand.Seed(time.Now().UnixNano())
	A := generateMatrix(2000, 1500)
	B := generateMatrix(1500, 3000)

	start := time.Now()
	C2 := ConcurrentMatrixMultiplication(A, B)
	fmt.Printf("Concurrent: %v\n", time.Since(start))

	start = time.Now()
	C1 := SequentialMatrixMultiplication(A, B)
	fmt.Printf("Sequential: %v\n", time.Since(start))

	// Optional correctness check (not printing, just verifying)
	for i := 0; i < len(C1); i++ {
		for j := 0; j < len(C1[0]); j++ {
			if (C1[i][j] - C2[i][j]) > 1e-6 {
				log.Fatalf("Mismatch at %d,%d: seq=%f conc=%f", i, j, C1[i][j], C2[i][j])
			}
		}
	}
}
