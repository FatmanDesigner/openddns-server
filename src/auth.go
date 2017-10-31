package main

// GenerateApp generates new pair of appid and client secret
import (
	"encoding/hex"
	"fmt"
	"math/rand"

	uuid "github.com/satori/go.uuid"
	hashids "github.com/speps/go-hashids"
)

// - appid is to be recorded in DNS TXT record openddns_appid=appid
// - secret is to authorize IP changes
func GenerateApp(userID string) (appid string, secret string) {
	fmt.Println("Generating appid and secret...")

	appid = hex.EncodeToString(uuid.NewV4().Bytes())
	secret = GenerateSecret(appid)

	// Assign to a persisted user identified by userID

	return
}

// GenerateSecret generates a random secret for every invocation
func GenerateSecret(appid string) string {
	fmt.Printf("Generating secret for appid %s...\n", appid)

	hd := hashids.NewData()
	hd.Salt = appid
	hd.MinLength = 16
	h, _ := hashids.NewWithData(hd)

	randomList := make([]int, 4)
	for i := 0; i < 4; i++ {
		randomList[i] = rand.Int()
	}

	secret, _ := h.Encode(randomList)

	return secret
}
