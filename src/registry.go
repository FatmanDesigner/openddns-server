package main

import "fmt"

// DNSEntry Consists of domainName and ip for doamin lookup
type DNSEntry struct {
	domainName string
	ip         string
}

// In-memory entries
// TODO: Externalize
var entries []DNSEntry

// Register an A Record entry
// WARNING: transactional
func Register(domainName string, ip string) bool {
	fmt.Printf("Registering...\n")

	for index := 0; index < len(entries); index++ {
		if domainName == entries[index].domainName {
			entries[index].ip = ip
			return true
		}
	}

	entries = append(entries, DNSEntry{domainName, ip})

	return true
}

// Lookup returns IP string for a giving domainName
// @returns tuple (bool string)
func Lookup(domainName string) (bool, string) {
	for index := 0; index < len(entries); index++ {
		if domainName == entries[index].domainName {
			return true, entries[index].ip
		}
	}

	return false, ""
}
