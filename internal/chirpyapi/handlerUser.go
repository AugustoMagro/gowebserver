package chirpyapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AugustoMagro/gowebserver/internal/auth"
	"github.com/AugustoMagro/gowebserver/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *ApiConfig) HandleUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	if err := dec.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashed_password})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	RespondWithJSON(w, 201, response{
		User: User{
			ID:         user.ID,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
			Email:      user.Email,
		},
	})

}

func (cfg *ApiConfig) GetUsers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User []database.User
	}

	defer r.Body.Close()

	users, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	RespondWithJSON(w, 201, users)
}

func (cfg *ApiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email              string `json:"email"`
		Passwrod           string `json:"password"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	if err := dec.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}

	valid, err := auth.CheckPasswordHash(params.Passwrod, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn´t check password", err)
		return
	}

	if valid {

		expires_time := time.Hour

		if params.Expires_in_seconds > 0 && params.Expires_in_seconds < 3600 {
			expires_time = time.Duration(params.Expires_in_seconds) * time.Second
		}

		token, err := auth.MakeJWT(user.ID, cfg.secret_key, expires_time)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Couldn´t create token", err)
			return
		}

		RespondWithJSON(w, http.StatusOK, response{User: User{
			ID:         user.ID,
			Created_at: user.CreatedAt,
			Updated_at: user.UpdatedAt,
			Email:      user.Email,
		},
			Token: token,
		})
	} else {
		RespondWithError(w, 401, "401 Unauthorized", nil)
	}

}
