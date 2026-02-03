package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		msg := "403 Forbidden"
		respondWithError(w, http.StatusForbidden, msg, nil)
		return
	}
	err := cfg.db.Reset(req.Context())
	if err != nil {
		msg := "Unable to reset users"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	msg := "Hits reset to 0"
	w.Write([]byte(msg))

}
