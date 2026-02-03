package main

import (
	"chirpy/internal/database"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	email, hashedPassword, err := DecodeEmailAndHashedPassword(req)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	if !isValidEmail(email) {
		msg := "Invalid email address"
		respondWithError(w, http.StatusBadRequest, msg, err)
		return
	}

	paramsWithHashedPassword := database.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(req.Context(), paramsWithHashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, SanitizeUser(user))
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}
