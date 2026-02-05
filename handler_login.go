package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const defaultAccessTokenExpiration = time.Hour

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DecodeLoginParams(req *http.Request) (loginParams, error) {
	params := loginParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return loginParams{}, err
	}
	return params, nil
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type loginResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		JWTToken     string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}
	params, err := DecodeLoginParams(req)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}
	isMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !isMatch || err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", nil)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, defaultAccessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate JWT", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate refresh token", err)
		return
	}
	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), createRefreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to store refresh token", err)
		return
	}

	userWithTokens := loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		JWTToken:     token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, userWithTokens)
}
