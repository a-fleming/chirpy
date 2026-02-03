package main

import (
	"chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	email, password, err := DecodeEmailAndPassword(req)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), email)
	isMatch, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if !isMatch || err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, SanitizeUser(user))
}
