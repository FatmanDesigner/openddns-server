package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Starting OpenDDNS Server...\n")

	go Serve(9000)
	DnsServe()
}
