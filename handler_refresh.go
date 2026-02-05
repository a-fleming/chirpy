package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type refreshResponse struct {
		JWTToken string `json:"token"`
	}
	now := time.Now().UTC()

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		msg := "401 Unauthorized"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	}
	refreshTokenDetails, err := cfg.db.GetRefreshTokenDetails(req.Context(), refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "401 Unauthorized"
			respondWithError(w, http.StatusUnauthorized, msg, err)
		} else {
			msg := "Database error"
			respondWithError(w, http.StatusInternalServerError, msg, err)
		}
		return
	} else if refreshTokenDetails.RevokedAt.Valid {
		msg := "401 Unauthorized"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	} else if refreshTokenDetails.ExpiresAt.Compare(now) <= 0 {
		msg := "401 Unauthorized"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	}
	accessToken, err := auth.MakeJWT(refreshTokenDetails.UserID, cfg.jwtSecret, defaultAccessTokenExpiration)
	if err != nil {
		msg := "Unable to generate JWT"
		respondWithError(w, http.StatusInternalServerError, msg, err)
	}

	response := refreshResponse{
		JWTToken: accessToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}
