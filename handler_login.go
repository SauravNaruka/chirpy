package main

import (
	"encoding/json"
	"net/http"

	"github.com/SauravNaruka/chirpy/internal/auth"
)

func (c *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	param := parameter{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := c.db.GetUserByEmail(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(param.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		},
	})
}
