package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}
	params.Body = basicCleanChirp(params.Body)

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		msg := "Chirp is too long"
		respondWithError(w, http.StatusBadRequest, msg, nil)
		return
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams(params))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp(chirp))
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}
	jsonFormattedChirps := []Chirp{}
	for _, chirp := range chirps {
		jsonFormattedChirps = append(jsonFormattedChirps, Chirp(chirp))
	}
	respondWithJSON(w, http.StatusOK, jsonFormattedChirps)
}

func basicCleanChirp(text string) string {
	const replacementStr = "****"
	profaneWords := []string{
		"fornax",
		"kerfuffle",
		"sharbert",
	}
	words := strings.Split(text, " ")
	for idx, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[idx] = replacementStr
		}
	}
	return strings.Join(words, " ")
}

func advancedCleanChirp(text string) string {
	const replacementStr = "****"
	profaneWords := []string{
		"fornax",
		"kerfuffle",
		"sharbert",
	}
	cleaned := text

	for _, profane := range profaneWords {
		lowerCase := strings.ToLower(cleaned)
		for idx := strings.Index(lowerCase, profane); idx > -1; {
			fmt.Printf("found '%s' at idx: %d\n", profane, idx)
			endIdx := idx + len(profane)
			cleaned = cleaned[:idx] + replacementStr + cleaned[endIdx:]

			lowerCase = strings.ToLower(cleaned)
			idx = strings.Index(lowerCase, profane)
		}
	}
	return cleaned
}
