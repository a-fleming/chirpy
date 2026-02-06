package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type PolkaEvent struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, req *http.Request) {
	polkaEvent := PolkaEvent{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&polkaEvent)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	if polkaEvent.Event == "user.upgraded" {
		userID, err := uuid.Parse(polkaEvent.Data.UserID)
		if err != nil {
			msg := "404 Not Found"
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}
		_, err = cfg.db.UpgradeUserToChirpyRed(req.Context(), userID)
		if err != nil {
			if err == sql.ErrNoRows {
				msg := "404 Not Found"
				respondWithError(w, http.StatusNotFound, msg, err)
				return
			} else {
				msg := "Database error"
				respondWithError(w, http.StatusInternalServerError, msg, err)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		// ignore other events
		w.WriteHeader(http.StatusNoContent)
		return
	}

}
