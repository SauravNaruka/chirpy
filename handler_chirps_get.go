package main

import (
	"net/http"
	"sort"

	"github.com/SauravNaruka/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")

	var dChirps []database.Chirp
	var err error
	if authorID != "" {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author_id format", err)
			return
		}
		dChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
	} else {
		dChirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while getting chirps", err)
		return
	}

	sortDirection := "asc"
	sortDirectionParam := r.URL.Query().Get("sort")
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}

	chirps := make([]Chirp, len(dChirps))

	for i, c := range dChirps {
		chirp := Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserId:    c.UserID,
		}
		chirps[i] = chirp
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	})
}
