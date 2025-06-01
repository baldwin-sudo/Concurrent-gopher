package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
)

type Job struct {
	req    *Request
	client *Client
}
type Result struct {
	result string
	client *Client
}
type Worker struct {
	id   string
	Load int32 // number of jobs its currently handling
	jobs chan Job
}
type LoadBalancer struct {
	Pool WorkersPool
}

func InitLoadBalancer(num_workers int, ctx context.Context, store *Store, results chan Result) *LoadBalancer {
	pool := make(WorkersPool, num_workers)
	LoadBalancer := &LoadBalancer{Pool: pool}
	for i := 0; i < num_workers; i++ {
		worker := initWorker(i)
		LoadBalancer.Push(worker)
		go worker.work(ctx, store, results)
	}
	return LoadBalancer
}
func (lb *LoadBalancer) Push(worker *Worker) {
	lb.Pool.Push(worker)
}
func (lb *LoadBalancer) AssignJob(job Job) {
	workerAny := lb.Pool.Pop()
	worker, ok := workerAny.(*Worker)
	if !ok {
		panic("Pop did not return a *Worker")
	}
	worker.AddWork(job)
	lb.Pool.Push(worker)
}
func (lb *LoadBalancer) Shutdown() {

}

func initWorker(id int) *Worker {
	return &Worker{id: "worker_id_" + strconv.Itoa(id), Load: 0, jobs: make(chan Job, 2)}
}
func (worker *Worker) AddWork(job Job) {
	worker.Load++
	worker.jobs <- job
}

func (worker *Worker) work(ctx context.Context, store *Store, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s shutting down due to context cancellation\n", worker.id)
			return
		case job, ok := <-worker.jobs:
			if !ok {
				return
			}
			result, err := Handle(*job.req, store)
			worker.Load--
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
func handleConnWithLB(conn net.Conn, id int, lb *LoadBalancer, results chan Result) {
	client := initClient(conn, id)
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
		job := Job{client: client, req: req}
		lb.AssignJob(job)

		if req.Command == CMD_QUIT {
			client.isDone = true
			break
		}
	}
	defer conn.Close()
	fmt.Printf("Client %s disconnected\n", client.ID)
}
