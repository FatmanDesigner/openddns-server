package main

import (
	"log"
	"os"
	"path"
	"strconv"
	"sync"
)

func main() {
	var dnsPort, httpPort int
	var dbString string

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

	if _, ok := os.LookupEnv("DB_STRING"); ok {
		dbString = os.Getenv("DB_STRING")
	} else {
		pwd, _ := os.Getwd()
		dbString = path.Join(pwd, "auth.db")
	}

	db := InitDB(dbString)
	if db == nil {
		log.Printf("Could not open DB %s", dbString)
		return
	}
	defer db.Close()
	log.Printf("Starting OpenDDNS Server...\n- DNS Server at port %d\n- HTTP Server at port %d\n", dnsPort, httpPort)

	// Execute HTTP
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {

		server := &HttpServer{DB: db}
		server.HttpServe(httpPort)
		wg.Done()
	}()
	go func() {
		DnsServe(dnsPort)
		wg.Done()
	}()

	wg.Wait()
}
