package main

import (
	"chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		msg := "401 Unauthorized"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		msg := "Unable to revoke refresh token"
		respondWithError(w, http.StatusInternalServerError, msg, err)
	}
	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
