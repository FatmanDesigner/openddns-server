package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var dnsPort, httpPort, dbString, staticRoot string
	dnsPort = "53"
	httpPort = "9000"
	staticRoot = "./web-ui/dist"

	configs := map[string]*string{
		"DNS_PORT":         &dnsPort,
		"HTTP_PORT":        &httpPort,
		"DB_STRING":        &dbString,
		"STATIC_ROOT":      &staticRoot,
		"GH_CLIENT_ID":     nil,
		"GH_CLIENT_SECRET": nil}
	if !ensureEnvParams(configs) {
		log.Fatal("Could not ensure env params required. Exiting...")
	}

	var db *sql.DB
	if absoluteFilePath, err := filepath.Abs(dbString); err == nil {
		db = InitDB(absoluteFilePath)
		if db == nil {
			log.Fatalf("Could not open DB %s", dbString)
			return
		}
		defer db.Close()
	} else {
		log.Fatal(err.Error())
		return
	}

	log.Printf("Starting OpenDDNS Server...\n"+
		"- DNS Server at port %s\n"+
		"- HTTP Server at port %s\n", dnsPort, httpPort)

	// Execute HTTP
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		server := &HttpServer{DB: db}
		if port, err := strconv.Atoi(httpPort); err == nil {
			server.HttpServe(port)
		}
		wg.Done()
	}()
	go func() {
		if port, err := strconv.Atoi(dnsPort); err == nil {
			DnsServe(port)
		}
		wg.Done()
	}()

	wg.Wait()
}

func ensureEnvParams(configs map[string]*string) (ok bool) {
	ok = true
	for key := range configs {
		if value, found := os.LookupEnv(key); found && len(value) != 0 {
			if configs[key] != nil {
				*configs[key] = value
				log.Printf("Env param %s applied %s", key, value)
				ok = ok && true
			}
		} else if len(*configs[key]) == 0 {
			log.Printf("Env param %s not found", key)
			ok = false
		} else {
			os.Setenv(key, *configs[key])
			log.Printf("Env param %s has been defaulted to %s", key, *configs[key])
			ok = ok && true
		}
	}

	return
}
