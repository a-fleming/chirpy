package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password,omitempty"`
	JWTToken       string    `json:"token,omitempty"`
}

func DecodeEmailAndPassword(req *http.Request) (string, string, error) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return "", "", err
	}
	return params.Email, params.Password, nil
}

func DecodeEmailAndHashedPassword(req *http.Request) (string, string, error) {
	email, password, err := DecodeEmailAndPassword(req)
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", "", err
	}
	return email, hashedPassword, nil
}

func SanitizeUser(user any) User {
	switch u := user.(type) {
	case User:
		return User{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Email:     u.Email,
			JWTToken:  u.JWTToken,
		}
	case database.User:
		return User{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Email:     u.Email,
			// TODO: probably need to change database to include JWTToken
		}
	default:
		return User{}
	}
}
