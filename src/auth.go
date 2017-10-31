package main

// GenerateApp generates new pair of appid and client secret
import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	hashids "github.com/speps/go-hashids"
)

// - appid is to be recorded in DNS TXT record openddns_appid=appid
// - secret is to authorize IP changes
func GenerateApp(userID string) (appid string, secret string, ok bool) {
	fmt.Println("Generating appid and secret...")

	var err error
	var db *sql.DB
	var stmt *sql.Stmt

	appid = hex.EncodeToString(uuid.NewV4().Bytes())
	secret, ok = GenerateSecret(appid)

	// Assign to a persisted user identified by userID
	db, err = sql.Open("sqlite3", "../auth.db")
	if err != nil {
		ok = false
		db.Close()

		return
	}
	stmt, err = db.Prepare("INSERT INTO apps (appid, secret, user_id) VALUES (?, ?, ?)")

	if err != nil {
		ok = false
		db.Close()

		return
	}

	// TODO: Has secret prior to saving to DB
	_, err = stmt.Exec(appid, secret, userID)

	if err != nil {
		fmt.Errorf(err.Error())

		ok = false
		db.Close()

		return
	}

	ok = true
	db.Close()
	return
}

// GenerateSecret generates a random secret for every invocation
func GenerateSecret(appid string) (secret string, ok bool) {
	var err interface{}
	var db *sql.DB
	var stmt *sql.Stmt

	fmt.Printf("Generating secret for appid %s...\n", appid)

	hd := hashids.NewData()
	hd.Salt = appid
	hd.MinLength = 16
	h, _ := hashids.NewWithData(hd)

	randomList := make([]int, 4)
	for i := 0; i < 4; i++ {
		randomList[i] = rand.Int()
	}

	secret, _ = h.Encode(randomList)

	db, err = sql.Open("sqlite3", "../auth.db")
	if err != nil {
		ok = false
		db.Close()

		return
	}

	stmt, err = db.Prepare("UPDATE apps SET secret = ? WHERE appid = ?)")
	if err != nil {
		ok = false
		db.Close()

		return
	}

	_, err = stmt.Exec(secret, appid)

	ok = true
	return
}

// Authenticate authenticate appid to be modified using secret
func Authenticate(appid string, secret string) (string, bool) {
	var userID string = ""
	var err error
	var db *sql.DB
	var row *sql.Row

	db, err = sql.Open("sqlite3", "../auth.db")
	if err != nil {
		return "", false
	}
	defer db.Close()

	row = db.QueryRow("SELECT user_id FROM apps WHERE appid = ? AND secret = ?", appid, secret)

	var scanned string
	if row != nil {
		err = row.Scan(&scanned)
		if err != nil {
			return "", false
		}

		userID = scanned
	}

	return userID, true
}
