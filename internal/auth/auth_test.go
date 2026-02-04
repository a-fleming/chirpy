package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func TestGetBearerToken_Valid(t *testing.T) {
	expected := "access_granted"
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+expected)

	received, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("FAIL: GetBearerToken returned error: %v", err)
	}
	if received != expected {
		t.Fatalf("FAIL: got bearerToken %v, expected %v", received, expected)
	}
}

func TestBearerToken_Invalid(t *testing.T) {
	cases := []struct {
		name string
		auth string
	}{
		{"empty header value", ""},
		{"empty token", "Bearer "},
		{"missing space after Bearer", "BearerSUPER_SECRET"},
		{"wrong scheme", "Basic abc123"},
		{"extra segments", "Bearer token extra"},
		{"leading space", " Bearer token"},
		{"missing authorization", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			headers := http.Header{}
			if tc.name != "missing authorization" {
				headers.Set("Authorization", tc.auth)
			}
			_, err := GetBearerToken(headers)
			if err == nil {
				t.Fatalf("expected error, got nil (Authorization=%q)", tc.auth)
			}
		})
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

func TestValidateJWTRejectsWrongIssuer(t *testing.T) {
	userID := uuid.New()
	secret := "super_secret_key"
	expiresIn := time.Hour

	// Manually craft a token with a bad issuer
	badIssuer := "other-issuer"
	claims := jwt.RegisteredClaims{
		Issuer:    badIssuer,
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatalf("expected error for wrong issuer, got nil")
	}
}
