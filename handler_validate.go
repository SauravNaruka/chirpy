package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json: "body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	const maxChirpLength = 140

	deconder := json.NewDecoder(r.Body)
	params := parameters{}
	err := deconder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanChrip := getCleanedBody(params.Body, badWords)
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanChrip,
	})
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		_, exists := badWords[strings.ToLower(word)]
		if !exists {
			continue
		}
		words[i] = "****"
	}

	return strings.Join(words, " ")
}
