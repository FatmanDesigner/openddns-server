package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	var dnsPort, httpPort int

	if _, ok := os.LookupEnv("DNS_PORT"); ok {
		dnsPort, _ = strconv.Atoi(os.Getenv("DNS_PORT"))
	} else {
		dnsPort = 53
	}

	if _, ok := os.LookupEnv("HTTP_PORT"); ok {
		httpPort, _ = strconv.Atoi(os.Getenv("HTTP_PORT"))
	} else {
		httpPort = 9000
	}

	db := InitDB("../auth.db")
	defer db.Close()
	fmt.Printf("Starting OpenDDNS Server...\n- DNS Server at port %d\n- HTTP Server at port %d\n", dnsPort, httpPort)

	// Execute HTTP
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		HttpServe(httpPort)
		wg.Done()
	}()
	go func() {
		DnsServe(dnsPort)
		wg.Done()
	}()

	wg.Wait()
}
