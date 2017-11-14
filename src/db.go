package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DomainEntry is a domain entry
type DomainEntry struct {
	DomainName string `json:"domainName"`
	IP         string `json:"ip"`
	UpdatedAt  int    `json:"updatedAt"`
}

type AppInfo struct {
	AppID  string `json:"appid"`
	Secret string `json:"secret,omitempty"`
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
	if err = createTable(db, "CREATE TABLE if not exists `domains` ( `domain_name` TEXT NOT NULL, `ip` TEXT NOT NULL, `updated_at` INTEGER NOT NULL, `owner_id` TEXT NOT NULL, PRIMARY KEY(`domain_name`) )"); err != nil {
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

func QueryAppInfosUserID(db *sql.DB, userID string) ([]AppInfo, error) {
	log.Printf("Querying apps by userID=%s", userID)

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT appid FROM apps WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appInfos []AppInfo
	for rows.Next() {
		appInfo := AppInfo{}
		if err := rows.Scan(&appInfo.AppID); err != nil {
			return nil, err
		}

		appInfos = append(appInfos, appInfo)
	}

	return appInfos, nil
}
