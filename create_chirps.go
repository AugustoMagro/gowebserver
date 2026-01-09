package main

import (
	"encoding/json"
	"net/http"

	"github.com/AugustoMagro/gowebserver/internal/database"
)

func (cfg *apiConfig) createChirpy(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		IserId string `json:"user_id"`
	}

	type response struct {
		User []database.User
	}

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	if err := dec.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	users, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, 201, users)
}
