package chirpyapi

import "net/http"

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		RespondWithError(w, http.StatusForbidden, "Reset is only allowed in dev environment.", nil)
		return
	}

	type response struct {
		Body string `json:"body"`
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusForbidden, "Failed to reset the database", err)
	}

	RespondWithJSON(w, http.StatusOK, response{Body: "Database reseted"})
}
