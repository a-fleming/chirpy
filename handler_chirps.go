package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"
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

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, req *http.Request) {
	type createChirpParams struct {
		Body string `json:"body"`
	}

	params := createChirpParams{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		msg := "Unable to decode parameters"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		msg := "401 Unauthorized"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		msg := "error validating JWT"
		respondWithError(w, http.StatusUnauthorized, msg, err)
		return
	}

	// Verify user exists in database
	_, err = cfg.db.GetUserByID(req.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "401 Unauthorized"
			respondWithError(w, http.StatusUnauthorized, msg, err)
		} else {
			msg := "Database error"
			respondWithError(w, http.StatusInternalServerError, msg, err)
		}
		return
	}

	params.Body = basicCleanChirp(params.Body)

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		msg := "Chirp is too long"
		respondWithError(w, http.StatusBadRequest, msg, nil)
		return
	}
	dbParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp(chirp))
}

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, req *http.Request) {
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

	chirp_id, err := uuid.Parse(req.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "404 Not Found", nil)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "404 Not Found", nil)
		return
	}
	if userID != chirp.UserID {
		msg := "403 Forbidden"
		respondWithError(w, http.StatusForbidden, msg, nil)
		return
	}
	err = cfg.db.DeleteChirpByID(req.Context(), chirp_id)
	if err != nil {
		msg := "Database error"
		respondWithError(w, http.StatusInternalServerError, msg, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerChirpsGetChirpByID(w http.ResponseWriter, req *http.Request) {
	chirp_id, err := uuid.Parse(req.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "404 Not Found", nil)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "404 Not Found", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp(chirp))
}

func (cfg *apiConfig) handlerChirpsGetChirps(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	authorIDStr := query.Get("author_id")
	sortOrder := query.Get("sort")
	if authorIDStr != "" {
		authorID, err := uuid.Parse(authorIDStr)
		if err != nil {
			msg := "404 Not Found"
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}

		dbChirps, err := cfg.db.GetChirpsByUserID(req.Context(), authorID)
		if err != nil && err != sql.ErrNoRows {
			msg := "404 Not Found"
			respondWithError(w, http.StatusNotFound, msg, err)
			return
		}
		chirps := convertDBChirpsToChirps(dbChirps)
		chirps = sortChirps(chirps, sortOrder)
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}
	chirps := convertDBChirpsToChirps(dbChirps)
	chirps = sortChirps(chirps, sortOrder)
	respondWithJSON(w, http.StatusOK, chirps)
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

func convertDBChirpsToChirps(dbChirps []database.Chirp) []Chirp {
	convertedChirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		convertedChirps[i] = Chirp(dbChirp)
	}
	return convertedChirps
}

func sortChirps(chirps []Chirp, order string) []Chirp {
	if strings.ToLower(order) == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}
	if strings.ToLower(order) == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}
	return chirps
}
