package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	if err := dec.Decode(&params); err != nil {
		log.Printf("Error deconding parameters %s", err)
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	params.Body = cleanWord(params.Body)
	respondWithJSON(w, 200, responseBody{CleanedBody: params.Body})
}

func cleanWord(s string) string {

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanWords := []string{}
	for word := range strings.FieldsSeq(s) {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				word = "****"
			}
		}
		cleanWords = append(cleanWords, word)
	}

	return strings.Join(cleanWords, " ")
}
