package main

import (
	"encoding/json"
	"net/http"

	"github.com/SauravNaruka/chirpy/internal/auth"
	"github.com/SauravNaruka/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {
	type parameter struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization header", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	param := parameter{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	hashPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	dbUser, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          param.Email,
		HashedPassword: hashPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while updating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirypRed: dbUser.IsChirpyRed.Bool,
		},
	})
}
