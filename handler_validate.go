package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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
		CleanedBody: basicCleanChirp(params.Body),
	})
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
