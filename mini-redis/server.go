package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
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

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT/SIGTERM for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt, shutting down (freeing resources)...")
		cancel()
		ln.Close()
	}()

	results := make(chan Result, BUFFER_SIZE)
	store := New()

	// Initialize load balancer and workers
	lb := InitLoadBalancer(NUM_WORKERS, ctx, store, results)

	// Start dispatcher
	go dispatch(ctx, results)

	id := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server shutdown signal received.")
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					log.Printf("Accept error: %v\n", err)
					continue
				}
			}
			fmt.Fprintf(os.Stdout, "user : %s connected\n", conn.RemoteAddr().String())
			go handleConnWithLB(conn, id, lb, results)
			id++
		}
	}
}
