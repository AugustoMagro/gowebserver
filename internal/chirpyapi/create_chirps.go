package chirpyapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AugustoMagro/gowebserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {

	chirpys, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirpy", err)
		return
	}

	var newChirps []Chirp
	for _, i := range chirpys {
		newChirps = append(newChirps, Chirp{
			ID:         i.ID,
			Created_at: i.CreatedAt,
			Updated_at: i.UpdatedAt,
			Body:       i.Body,
			User_id:    i.UserID,
		})
	}

	RespondWithJSON(w, 200, newChirps)
}

func (cfg *ApiConfig) GetChirpsByID(w http.ResponseWriter, r *http.Request) {

	idChirpy, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid chirp ID", err)
		return
	}

	chirpy, err := cfg.db.GetChirpsID(r.Context(), idChirpy)
	if err != nil {
		RespondWithError(w, 404, "Couldn't retrieve chirpy", err)
		return
	}

	newChirps := Chirp{
		ID:         chirpy.ID,
		Created_at: chirpy.CreatedAt,
		Updated_at: chirpy.UpdatedAt,
		Body:       chirpy.Body,
		User_id:    chirpy.UserID,
	}

	RespondWithJSON(w, 200, newChirps)
}

func (cfg *ApiConfig) CreateChirpy(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp Chirp
	}

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	params := parameters{}
	if err := dec.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	chirpy, err := cfg.db.CreateChirpy(r.Context(), database.CreateChirpyParams{Body: params.Body, UserID: params.UserId})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirpy", err)
		return
	}

	RespondWithJSON(w, 201, Chirp{
		ID:         chirpy.ID,
		Created_at: chirpy.CreatedAt,
		Updated_at: chirpy.UpdatedAt,
		Body:       chirpy.Body,
		User_id:    chirpy.UserID,
	})
}

type Chirp struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_id    uuid.UUID `json:"user_id"`
}
