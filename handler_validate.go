package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		msg := "Chirp is too long"
		respondWithError(w, http.StatusBadRequest, msg, nil)
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
	})
}
