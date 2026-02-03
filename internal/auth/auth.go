package auth

import (
	"github.com/alexedwards/argon2id"
)

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func HashPassword(password string) (string, error) {
	var defaultParams = argon2id.DefaultParams
	return argon2id.CreateHash(password, defaultParams)
}
