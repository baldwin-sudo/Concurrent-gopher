package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	var host = flag.String("host", "localhost", "host is the ip address or the hostname of the machine you wish to scan . ")
	flag.Parse()
	var commonPorts = []int{20, 21, 22, 23, 25, 53, 80, 110, 123, 143, 161, 443, 465, 587, 993, 995, 3306, 3389, 5432, 6379, 8080, 8443}
	var addresses []string = make([]string, len(commonPorts))
	// populate addresses slice
	for idx, p := range commonPorts {
		addresses[idx] = fmt.Sprintf("%s:%d", *host, p)

	}
	// Measure sequential scan time
	startSeq := time.Now()
	seq := scanPortsSequentially(addresses)
	durationSeq := time.Since(startSeq)

	// Measure concurrent scan time
	startConc := time.Now()
	conc := scanPortsConccurently1(addresses)
	durationConc := time.Since(startConc)

	fmt.Println("Sequential scan results:")
	fmt.Println(syncMapToString(seq))
	fmt.Printf("Sequential scan took: %v\n\n", durationSeq)

	fmt.Println("Concurrent scan results:")
	fmt.Println(syncMapToString(conc))
	fmt.Printf("Concurrent scan took: %v\n", durationConc)

}
func scanPortsSequentially(addresses []string) *sync.Map {
	portsStatus := &sync.Map{}
	for _, address := range addresses {
		open := scanPort(address)
		portsStatus.Store(address, open)
	}
	return portsStatus
}
func scanPortsConccurently1(addresses []string) *sync.Map {
	portsStatus := &sync.Map{}
	var wg sync.WaitGroup
	for _, address := range addresses {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			open := scanPort(addr)
			portsStatus.Store(addr, open)
		}(address)

	}
	wg.Wait()
	return portsStatus

}
func scanPortsConccurently2(addresses []string) *sync.Map {
	portsStatus := &sync.Map{}
	var wg sync.WaitGroup
	for _, address := range addresses {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			open := scanPort(addr)
			portsStatus.Store(addr, open)
		}(address)

	}
	wg.Wait()
	return portsStatus

}
func scanPort(address string) (open bool) {
	_, err := net.Dial("tcp", address)

	if err != nil {
		open = false
	} else {
		open = true

	}
	return
}
func syncMapToString(m *sync.Map) string {
	result := ""
	m.Range(func(key, value any) bool {
		result += fmt.Sprintf("%s: %v\n", key, value)
		return true
	})
	return result
}
