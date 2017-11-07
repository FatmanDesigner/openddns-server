package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DomainEntry is a domain entry
type DomainEntry struct {
	DomainName string `json:"domainName"`
	IP         string `json:"ip"`
	UpdatedAt  int    `json:"updatedAt"`
}

// InitDB returns a pointer to DB with tables fully structured
func InitDB(filepath string) *sql.DB {
	log.Printf("Initializing DB %s", filepath)

	var db *sql.DB
	var err error

	if db, err = sql.Open("sqlite3", filepath); err != nil {
		log.Printf("Could not open DB. %s\n", err.Error())
		return nil
	}

	// CREATE TABLE `apps`
	if err = createTable(db, "CREATE TABLE if not exists `apps` ( `appid` TEXT, `secret` TEXT, `user_id` TEXT, PRIMARY KEY(`appid`) )"); err != nil {
		defer db.Close()
		return nil
	}

	// CREATE TABLE `domains`
	if err = createTable(db, "CREATE TABLE if not exists `domains` ( `domain_name` TEXT NOT NULL, `ip` TEXT NOT NULL, `updated_at` INTEGER NOT NULL, `owner_id` TEXT NOT NULL, `appid` TEXT NOT NULL, PRIMARY KEY(`domain_name`) )"); err != nil {
		defer db.Close()
		return nil
	}

	return db
}

func createTable(db *sql.DB, createStatement string) error {
	var stmt *sql.Stmt
	var err error

	if stmt, err = db.Prepare(createStatement); err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func QueryDomainOwnerByDomainName(db *sql.DB, domainName string) (string, error) {
	log.Printf("Querying domain owner by domainName=%s", domainName)

	var row *sql.Row
	var err error
	var ownerID string

	row = db.QueryRow("SELECT owner_id FROM domains WHERE domain_name = ?", domainName)
	if err = row.Scan(&ownerID); err != nil {
		return "", err
	}

	return ownerID, nil
}

// QueryDomainEntriesByUserID gets all domain entries by a userID
func QueryDomainEntriesByUserID(db *sql.DB, userID string) ([]DomainEntry, error) {
	log.Printf("Querying domain entries by userID=%s", userID)

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT domain_name, ip, updated_at FROM domains WHERE owner_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domainEntries []DomainEntry

	for rows.Next() {
		domainEntry := DomainEntry{}
		if err := rows.Scan(&domainEntry.DomainName, &domainEntry.IP, &domainEntry.UpdatedAt); err != nil {
			return nil, err
		}

		domainEntries = append(domainEntries, domainEntry)
	}

	return domainEntries, nil
}

func UpsertDomainEntry(db *sql.DB, ownerID string, appid string, domainName string, domainIP string) error {
	var row *sql.Row
	var stmt *sql.Stmt
	var err error
	var count int

	row = db.QueryRow("SELECT COUNT(domain_name) FROM domains WHERE domain_name = ?", domainName)
	if err = row.Scan(&count); err != nil {
		return err
	}

	if count > 1 {
		panic("Database integrity")
	} else if count == 0 {
		if stmt, err = db.Prepare("INSERT domains (domain_name, ip, updated_at, owner_id, appid) VALUES (?, ?, ?, ?, ?)"); err != nil {
			return err
		}

		updatedAt := time.Now().UnixNano() / 1000000
		log.Printf("Inserting a new domain entry: domainName=%s, domainIP=%s", domainName, domainIP)
		if _, err = stmt.Exec(domainName, domainIP, updatedAt, ownerID, appid); err != nil {
			return err
		}
	} else {
		if stmt, err = db.Prepare("UPDATE domains SET ip=?, updated_at=?, owner_id=?, appid=? WHERE domain_name = ?"); err != nil {
			return err
		}

		updatedAt := time.Now().UnixNano() / 1000000
		log.Printf("Updating existing domain entry: domainName=%s, domainIP=%s", domainName, domainIP)
		if _, err = stmt.Exec(domainIP, updatedAt, ownerID, appid, domainName); err != nil {
			return err
		}
	}

	return nil
}
