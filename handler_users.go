package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, req *http.Request) {
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

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		msg := "401 Unauthorized"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		msg := "error validating JWT"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	}
	newEmail, newHashedPassword, err := DecodeEmailAndHashedPassword(req)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}
	params := database.UpdateUserEmailAndPasswordParams{
		Email:          newEmail,
		HashedPassword: newHashedPassword,
		ID:             userID,
	}
	user, err := cfg.db.UpdateUserEmailAndPassword(req.Context(), params)
	if err != nil {
		msg := "Unable to update email and password"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}
	respondWithJSON(w, http.StatusOK, SanitizeUser(user))
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}
