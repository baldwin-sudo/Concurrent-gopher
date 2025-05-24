package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Job struct {
	path string
}

type Result struct {
	Lines int
	Words int
	Char  int
}

func main() {
	rootDir := flag.String("dir", "test", "Will count total Words, Lines and Characters in this Directory.")
	flag.Parse()
	fmt.Println("root directory:", *rootDir)

	fileSystem := os.DirFS(*rootDir)

	// ---------------------------
	// Concurrent Processing
	// ---------------------------
	fmt.Println("Starting concurrent processing...")
	startConcurrent := time.Now()
	numCPU := runtime.NumCPU()
	workersNum := numCPU * 4 // or 6–8× depending on I/O capacity
	fmt.Println("workers : ", workersNum)
	jobs := make(chan Job, workersNum*2)
	results := make(chan Result, workersNum*2)

	var wg sync.WaitGroup
	for i := 0; i < workersNum; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	var aggWg sync.WaitGroup
	aggWg.Add(1)

	go func() {
		defer aggWg.Done()
		totalLines := 0
		totalWords := 0
		totalChars := 0
		for result := range results {
			totalLines += result.Lines
			totalWords += result.Words
			totalChars += result.Char
		}
		fmt.Println("\n[Concurrent] TOTALS:")
		fmt.Println("Lines:", totalLines)
		fmt.Println("Words:", totalWords)
		fmt.Println("Characters:", totalChars)
	}()

	// Ensure walk finishes before closing jobs
	walkErr := func() error {
		return fs.WalkDir(fileSystem, ".", func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				// fmt.Println("Skipping:", path, "->", err)
				return nil
			}
			if !entry.IsDir() {
				jobs <- Job{filepath.Join(*rootDir, path)}
			}
			return nil
		})
	}()

	close(jobs)
	wg.Wait()
	close(results)
	aggWg.Wait()

	if walkErr != nil {
		fmt.Println("Error walking directory:", walkErr)
	}

	elapsedConcurrent := time.Since(startConcurrent)
	fmt.Printf("[Concurrent] Processing took: %s\n", elapsedConcurrent)

	// ---------------------------
	// Sequential Processing
	// ---------------------------
	fmt.Println("\nStarting sequential processing...")
	startSequential := time.Now()

	totalLines := 0
	totalWords := 0
	totalChars := 0

	_ = fs.WalkDir(fileSystem, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			// fmt.Println("Skipping:", path, "->", err)
			return nil
		}
		if !entry.IsDir() {
			result, err := ProcessFile(filepath.Join(*rootDir, path))
			if err == nil {
				totalLines += result.Lines
				totalWords += result.Words
				totalChars += result.Char
			}
		}
		return nil
	})

	fmt.Println("\n[Sequential] TOTALS:")
	fmt.Println("Lines:", totalLines)
	fmt.Println("Words:", totalWords)
	fmt.Println("Characters:", totalChars)

	elapsedSequential := time.Since(startSequential)
	fmt.Printf("[Sequential] Processing took: %s\n", elapsedSequential)
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		ProcessFileConcurrent(job.path, results)
	}
}

func ProcessFileConcurrent(path string, results chan<- Result) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineCounter := 0
	wordCounter := 0
	charCounter := 0

	for scanner.Scan() {
		lineCounter++
		line := scanner.Text()
		words := strings.Fields(line)
		wordCounter += len(words)
		for _, word := range words {
			charCounter += len(word)
		}
	}

	if err := scanner.Err(); err != nil {
		return
	}

	results <- Result{
		Lines: lineCounter,
		Words: wordCounter,
		Char:  charCounter,
	}
}

func ProcessFile(path string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, errors.New("failed to open file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineCounter := 0
	wordCounter := 0
	charCounter := 0

	for scanner.Scan() {
		lineCounter++
		line := scanner.Text()
		words := strings.Fields(line)
		wordCounter += len(words)
		for _, word := range words {
			charCounter += len(word)
		}
	}

	if err := scanner.Err(); err != nil {
		return Result{}, err
	}

	return Result{
		Lines: lineCounter,
		Words: wordCounter,
		Char:  charCounter,
	}, nil
}
