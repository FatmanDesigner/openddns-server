package main

import "testing"

// Register an A Record entry
// WARNING: transactional
func TestRegister(t *testing.T) {
	result := Register("google.com", "8.8.4.4")
	if !result {
		t.Error("Expecting google.com to be registered at 8.8.4.4")
	}
}

// Lookup returns IP string for a giving domainName
// @returns tuple (bool string)
func TestLookup(t *testing.T) {
	Register("google.com", "8.8.4.4")
	found, ip := Lookup("google.com")

	if !found {
		t.Error("google.com must be resolved to 8.8.4.4")
	}

	if ip != "8.8.4.4" {
		t.Error("google.com must be resolved to 8.8.4.4")
	}
}
