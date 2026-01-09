package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AugustoMagro/gowebserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		User
	}

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	if err := dec.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, 201, response{
		User: User{
			ID:         user.ID,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
			Email:      user.Email,
		},
	})

}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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

type User struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}
