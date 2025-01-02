package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersRedUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	param := parameter{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid json", err)
		return
	}

	if param.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), param.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
