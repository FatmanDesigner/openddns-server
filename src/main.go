package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Starting OpenDDNS Server...\n")

	Serve(9000)
}
