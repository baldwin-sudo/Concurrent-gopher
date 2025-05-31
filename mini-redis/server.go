package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var (
	NUM_WORKERS int = runtime.NumCPU() * 2
	BUFFER_SIZE int = NUM_WORKERS * 2
)

func server() {
	ln, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server started successfully ... ")

	// Create cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal catching: listen for SIGINT or SIGTERM (Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// When signal received, cancel context and close listener to unblock Accept
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt, shutting down (freeing ressources)...")
		cancel()
		ln.Close()
	}()

	// Channels for jobs and results
	jobs := make(chan Job, BUFFER_SIZE)
	results := make(chan Result, BUFFER_SIZE)

	// Initialize store
	store := New()

	// Launch workers with context
	launchWorkers(ctx, store, jobs, results)

	// Start dispatcher goroutine with context
	go dispatch(ctx, results)

	id := 0
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, shutdown
			fmt.Println("Server shutdown signal received. Closing job and results channels.")
			close(jobs)
			close(results)
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return // Exit accept loop if shutting down
				default:
					log.Printf("Accept error: %v\n", err)
					continue
				}
			}
			fmt.Fprintf(os.Stdout, "user : %s connected\n", conn.RemoteAddr().String())
			go handleConn(conn, id, jobs, results)
			id++
		}
	}
}

func launchWorkers(ctx context.Context, store *Store, jobs <-chan Job, results chan<- Result) {
	for i := 0; i < NUM_WORKERS; i++ {
		worker := Worker{id: "worker-" + strconv.Itoa(i)}
		go worker.work(ctx, store, jobs, results)
	}
}

func (worker *Worker) work(ctx context.Context, store *Store, jobs <-chan Job, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s shutting down due to context cancellation\n", worker.id)
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			result, err := Handle(*job.req, store)
			if err != nil {
				results <- Result{err.Error(), job.client}
				continue
			}
			results <- Result{result, job.client}
		}
	}
}

func dispatch(ctx context.Context, results <-chan Result) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Dispatcher shutting down due to context cancellation")
			return
		case result, ok := <-results:
			if !ok {
				return
			}
			if result.client.isDone {
				continue
			}
			result.client.Output <- result.result + "\n> "
		}
	}
}

// do a separte error channel to keep this function clean
func handleConn(conn net.Conn, id int, jobs chan Job, results chan Result) {
	client := initClient(conn, id)
	// start the individual client writer :

	go client.ClientWrite()
	client.Output <- "Welcome to mini-redis by baldwin!\n> "

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		input := scanner.Text()
		req, err := ParseCommand(input)
		if err != nil {
			results <- Result{client: client, result: err.Error()}
			continue
		}

		jobs <- Job{client: client, req: req}
		if req.Command == CMD_QUIT {
			client.isDone = true

			break
		}
	}
	defer conn.Close()
	fmt.Printf("Client %s disconnected\n", client.ID)
}

type Job struct {
	req    *Request
	client *Client
}
type Result struct {
	result string
	client *Client
}
type Worker struct {
	id string
}
