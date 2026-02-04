package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
	"time"
)

type loginParams struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ExpiresInSec *int32 `json:"expires_in_seconds,omitempty"`
}

func DecodeLoginParams(req *http.Request) (loginParams, error) {
	const defaultExpirationSec int32 = 3600 // 1 Hour
	params := loginParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return loginParams{}, err
	}

	// add expiration time if not provided
	if params.ExpiresInSec == nil {
		v := defaultExpirationSec
		params.ExpiresInSec = &v
	}

	// validate expiration time is less than the default
	if *params.ExpiresInSec > defaultExpirationSec {
		*params.ExpiresInSec = defaultExpirationSec
	}
	return params, nil
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	params, err := DecodeLoginParams(req)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	isMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !isMatch || err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}

	expiresInSec := time.Duration(*params.ExpiresInSec) * time.Second
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresInSec)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate JWT.", err)
	}
	userWithToken := SanitizeUser(user)
	userWithToken.JWTToken = token
	respondWithJSON(w, http.StatusOK, userWithToken)
}
