package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	if !isValidEmail(params.Email) {
		msg := "Invalid email address"
		respondWithError(w, http.StatusBadRequest, msg, err)
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, User(user))
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}
