package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/PPaterka/Chirpy/internal/database"
	"github.com/google/uuid"
)

type respBody struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	resp := []respBody{}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, chirp := range chirps {
		jsonChirp := respBody{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		resp = append(resp, jsonChirp)
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) newChirp(w http.ResponseWriter, r *http.Request) {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp too long")
		return
	}

	body, err := validateMessage(chirp.Body, badWords)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	resp, err := cfg.db.AddChirp(r.Context(), database.AddChirpParams{Body: body, UserID: chirp.UserID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, respBody{
		ID:        resp.ID,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
		Body:      resp.Body,
		UserID:    resp.UserID,
	})
}

func validateMessage(body string, badWords []string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	words := strings.Split(body, " ")
	cleanedWords := []string{}
	for _, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	return strings.Join(cleanedWords, " "), nil
}
