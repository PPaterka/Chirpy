package main

import (
	"encoding/json"
	"net/http"
)

func chirpVerification(w http.ResponseWriter, r *http.Request) {
	type chirpStruct struct {
		Body string `json:"body"`
	}
	type respBody struct {
		Valid bool `json:"valid"`
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
	resp := respBody{Valid: true}
	respondWithJSON(w, 200, resp)

}
