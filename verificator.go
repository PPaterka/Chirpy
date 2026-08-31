package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func chirpVerification(w http.ResponseWriter, r *http.Request) {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	type chirpStruct struct {
		Body string `json:"body"`
	}
	type respBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpStruct{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp too long")
		return
	}
	words := strings.Split(chirp.Body, " ")
	cleanedWords := []string{}
	for _, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	respondWithJSON(w, 200, respBody{CleanedBody: cleanMessage(chirp.Body, badWords)})
}

func cleanMessage(body string, badWords []string) string {
	words := strings.Split(body, " ")
	cleanedWords := []string{}
	for _, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	return strings.Join(cleanedWords, " ")
}
