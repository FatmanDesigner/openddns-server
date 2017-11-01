package main

import "testing"

// GenerateApp generates new pair of appid and client secret

func TestGenerateApp(t *testing.T) {
	userID := "khanh"
	appid, secret, ok := GenerateApp(userID)

	if !ok {
		t.Error("Generate App should be OK")
	}

	if len(appid) == 0 {
		t.Error("appid should be non-empty")
	}

	if len(secret) == 0 {
		t.Error("secret should be non-empty")
	}
}

func TestGenerateSecret(t *testing.T) {
	secretA, _ := GenerateSecret("itsalongappid")

	if len(secretA) == 0 {
		t.Error("secret should be non-empty")
	}
	t.Logf("Secret generated for \"itsalongappid\": %s", secretA)

	secretB, _ := GenerateSecret("itsalongappid")

	if len(secretB) == 0 {
		t.Error("secret should be non-empty")
	}
	t.Logf("Secret generated for \"itsalongappid\": %s", secretB)

	if secretA == secretB {
		t.Error("secrets should not be the same")
	}
}

func TestAuthenticate(t *testing.T) {
	userID := "khanh"
	appid, secret, ok := GenerateApp(userID)

	if !ok {
		t.Error("Generate App should be OK")
	}

	t.Logf("appid: %s, secret: %s", appid, secret)
	storedUserID, ok := Authenticate(appid, secret)

	t.Logf("UserID: %s", storedUserID)
	if !ok {
		t.Error("Authorize should be OK")
	}
	if storedUserID != "khanh" {
		t.Error("UserID should be \"khanh\"")
	}
}
