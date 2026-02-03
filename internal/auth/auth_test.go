package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password := "password1234"
	passwordHash := "$argon2id$v=19$m=65536,t=1,p=32$PP0DkHE23AGNcFUIhGy5fw$UyW02lhBazNZL/T1R06hQPSkMQsS6prRIIHUyJzLpPQ"

	same, err := CheckPasswordHash(password, passwordHash)
	if err != nil {
		t.Errorf("FAIL: Error checking password.")
	}
	if !same {
		t.Errorf("FAIL: Received false, Expected true")
	}
}
