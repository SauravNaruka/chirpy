package main

import (
	"encoding/json"
	"net/http"

	"github.com/SauravNaruka/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersRedUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to get api key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "api key is invalid", err)
		return
	}

	param := parameter{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&param)
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
