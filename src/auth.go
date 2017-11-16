package main

// GenerateApp generates new pair of appid and client secret
import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"log"
	"math/rand"

	uuid "github.com/satori/go.uuid"
	hashids "github.com/speps/go-hashids"
)

type Auth struct {
	DB *sql.DB
}

// - appid is to be recorded in DNS TXT record openddns_appid=appid
// - secret is to authorize IP changes
func (instance *Auth) GenerateApp(userID string) (appid string, secret string, ok bool) {
	var err error
	var stmt *sql.Stmt
	var row *sql.Row

	log.Printf("Looking for user %s", userID)
	row = instance.DB.QueryRow("SELECT user_id FROM apps WHERE user_id = ?", userID)
	var scanned string

	if err = row.Scan(&scanned); err == nil {
		ok = true
		log.Printf("No op. UserID already has appid: %s", scanned)

		return
	}

	log.Println("Generating appid and secret...")

	appid = hex.EncodeToString(uuid.NewV4().Bytes())
	secret, ok = internalGenerateSecret(appid)

	if stmt, err = instance.DB.Prepare("INSERT INTO apps (appid, secret, user_id) VALUES (?, ?, ?)"); err != nil {
		ok = false
		log.Printf("Could not prepare insert statement. %s", err.Error())

		return
	}

	secretHashed := hex.EncodeToString(sha1.New().Sum([]byte(secret)))
	_, err = stmt.Exec(appid, secretHashed, userID)

	if err != nil {
		ok = false
		log.Println("Could not insert into app")

		return
	}

	ok = true
	return
}

// GenerateSecret generates a random secret for every invocation
func (instance *Auth) GenerateSecret(userID string, appid string) (secret string, ok bool) {
	secret, ok = internalGenerateSecret(appid)

	if !ok {
		return
	}

	var err error
	var db *sql.DB = instance.DB
	var stmt *sql.Stmt

	if stmt, err = db.Prepare("UPDATE apps SET secret = ? WHERE user_id = ? and appid = ?"); err != nil {
		log.Println("Could not prepare UPDATE statement. " + err.Error())
		ok = false
		return
	}

	secretHashed := hex.EncodeToString(sha1.New().Sum([]byte(secret)))
	if _, err = stmt.Exec(secretHashed, userID, appid); err != nil {
		log.Println("Could not execute UPDATE statement. " + err.Error())
		ok = false
		return
	}

	return
}

// Authenticate authenticate appid to be modified using secret
func (instance *Auth) Authenticate(appid string, secret string) (string, bool) {
	var userID string = ""
	var err error
	var db *sql.DB = instance.DB
	var row *sql.Row

	secretHashed := hex.EncodeToString(sha1.New().Sum([]byte(secret)))
	row = db.QueryRow("SELECT user_id FROM apps WHERE appid = ? AND secret = ?", appid, secretHashed)

	var scanned string
	if row != nil {
		err = row.Scan(&scanned)
		if err != nil {
			log.Printf("Could not find user_id for appid = %s, secret = %s. %s", appid, secret, err.Error())
			return "", false
		}

		userID = scanned
	}

	return userID, true
}

func internalGenerateSecret(appid string) (secret string, ok bool) {
	log.Printf("Generating secret for appid %s...\n", appid)

	hd := hashids.NewData()
	hd.Salt = appid
	hd.MinLength = 16
	h, _ := hashids.NewWithData(hd)

	randomList := make([]int, 4)
	for i := 0; i < 4; i++ {
		randomList[i] = rand.Int()
	}

	secret, _ = h.Encode(randomList)

	ok = true
	return
}
