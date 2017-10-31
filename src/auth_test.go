package main

import "testing"

// GenerateApp generates new pair of appid and client secret

func TestGenerateApp(t *testing.T) {
	userID := "khanh"
	appid, secret := GenerateApp(userID)

	if len(appid) == 0 {
		t.Error("appid should be non-empty")
	}

	if len(secret) == 0 {
		t.Error("secret should be non-empty")
	}
}

func TestGenerateSecret(t *testing.T) {
	secretA := GenerateSecret("itsalongappid")

	if len(secretA) == 0 {
		t.Error("secret should be non-empty")
	}

	secretB := GenerateSecret("itsalongappid")

	if len(secretB) == 0 {
		t.Error("secret should be non-empty")
	}

	if secretA == secretB {
		t.Error("secrets should not be the same")
	}
}
