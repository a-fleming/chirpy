package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "password1234"
	passwordHash := "$argon2id$v=19$m=65536,t=1,p=32$PP0DkHE23AGNcFUIhGy5fw$UyW02lhBazNZL/T1R06hQPSkMQsS6prRIIHUyJzLpPQ"

	same, err := CheckPasswordHash(password, passwordHash)
	if err != nil {
		t.Fatalf("FAIL: CheckPasswordHash returned error: %v", err)
	}
	if !same {
		t.Fatalf("FAIL: expected password to match hash")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "super_secret_key"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("FAIL: MakeJWT returned error: %v", err)
	}
	if token == "" {
		t.Fatalf("FAIL: MakeJWT returned empty token")
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("FAIL: ValidateJWT returned error: %v", err)
	}
	if gotID != userID {
		t.Fatalf("FAIL: got userID %v, expected %v", gotID, userID)
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	userID := uuid.New()
	rightSecret := "right_secret"
	wrongSecret := "wrong_secret"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, rightSecret, expiresIn)
	if err != nil {
		t.Fatalf("FAIL: MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatalf("FAIL: expected error for wrong secret, got nil")
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "super_secret_key"

	// Negative duration means token is already expired at creation
	expiresIn := -1 * time.Minute
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("FAIL: MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("FAIL: expected error for expired token, got nil")
	}
}
