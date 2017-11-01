package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB returns a pointer to DB with tables fully structured
func InitDB(filepath string) *sql.DB {
	var db *sql.DB
	var stmt *sql.Stmt
	var err error

	if db, err = sql.Open("sqlite3", filepath); err != nil {
		log.Printf("LINE 17: Could not open DB. %s\n", err.Error())
		return nil
	}

	if stmt, err = db.Prepare("CREATE TABLE if not exists `apps` ( `appid` TEXT, `secret` TEXT, `user_id` TEXT, PRIMARY KEY(`appid`) )"); err != nil {
		defer db.Close()
		return nil
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		defer db.Close()
		return nil
	}

	return db
}
