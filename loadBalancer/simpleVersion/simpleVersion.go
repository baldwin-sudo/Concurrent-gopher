package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Start listening on TCP port 1999 on localhost
	ln, err := net.Listen("tcp", "localhost:1999")
	if err != nil {
		log.Fatal("server couldn't listen:", err)
	}
	// workers pool size
	numWorkers := runtime.NumCPU() * 2
	// Create a root context and a cancel function to control goroutine lifecycle on shutdown
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Buffered channel to queue incoming client requests to workers
	requests := make(chan Request, numWorkers*2)

	// Start a pool of workers equal to twice the CPU cores
	var workers_wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		workers_wg.Add(1)
		go func() {
			defer workers_wg.Done()
			worker(requests, ctx)
		}()
	}

	// Channel to receive OS signals (SIGINT, SIGTERM) for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup to track active client connections
	var conns_wg sync.WaitGroup

	// Goroutine to accept incoming connections continuously
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				// On error (likely due to ln.Close() on shutdown), cancel context and exit loop
				cancel()
				return
			}
			fmt.Println("user connected:", conn.LocalAddr())

			// Track active connection goroutine
			conns_wg.Add(1)
			go func() {
				defer conns_wg.Done()
				handleConn(conn, requests, ctx)
			}()
		}
	}()

	// Block main goroutine until shutdown signal received ,
	//  so the main doesnt exit directly it is waiting for the signal
	<-sigs

	// Close listener to stop accepting new connections
	ln.Close()

	// Wait for all active connections to finish cleanly
	conns_wg.Wait()

	// Close requests channel to signal workers no more jobs will come
	close(requests)

	// Wait for all workers to finish processing
	workers_wg.Wait()

	log.Println("Server stopped gracefully ...")
}

// Request represents a client command and a channel to send the response back
type Request struct {
	Command       int
	ClientChannel chan string
}

// parseCommand parses the client's input string into an integer command.
// Returns an error if the input is empty, contains multiple tokens, or isn't a valid integer.
func parseCommand(cmd string) (int, error) {
	c := strings.Fields(cmd)
	if len(c) == 0 {
		return 0, errors.New("enter something")
	} else if len(c) > 1 {
		return 0, errors.New("incorrect command, please type help to understand")
	}
	num, err := strconv.Atoi(c[0])
	if err != nil {
		return 0, err
	}
	return num, nil
}

// handleConn manages a single client connection.
// It reads commands line-by-line from the client, sends requests to the worker pool,
// and writes back the results.
// It respects context cancellation for graceful shutdown.
func handleConn(conn net.Conn, requests chan<- Request, ctx context.Context) {
	defer conn.Close()

	// Buffered reader and writer for efficient I/O with client
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// Print initial prompt for client input
	rw.Writer.WriteString("> ")
	rw.Flush()

	for {
		// Check if server is shutting down; if so, stop reading further commands
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Set a read deadline to prevent hanging indefinitely if client is idle
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))

		// Create a fresh channel per request so responses don't get mixed up between commands
		client_channel := make(chan string)

		// Read a full line (command) from the client
		request_str, err := rw.ReadString('\n')
		if err != nil {
			// If read times out, just retry reading without closing connection
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			// Other errors (e.g. client closed connection) break the loop
			break
		}

		request_str = strings.TrimSpace(request_str)
		if request_str == "exit" {
			// Client requested to close connection gracefully
			break
		}

		// Parse the command from client input
		command, err := parseCommand(request_str)
		if err != nil {
			rw.Writer.WriteString(err.Error() + "\n")
			rw.Flush()
			continue
		}

		// Package the command and the response channel into a Request struct
		request := Request{
			Command:       command,
			ClientChannel: client_channel,
		}

		// Send request to worker pool for processing
		requests <- request

		// Wait for worker response on client-specific channel
		response := <-client_channel

		// Write the response back to client, followed by a prompt for next input
		rw.Writer.WriteString(response + "\n")
		rw.Writer.WriteString("> ")
		rw.Flush()
	}
}

// worker processes incoming requests from the requests channel.
// It listens for server shutdown through the context and exits cleanly when cancelled.
func worker(requests chan Request, ctx context.Context) {
	for request := range requests {
		select {
		case <-ctx.Done():
			// Server is shutting down, exit worker goroutine
			return
		default:
		}

		command, client_chan := request.Command, request.ClientChannel
		// For now, just convert command integer to string and send back as response
		client_chan <- fmt.Sprint(command)
	}
}
