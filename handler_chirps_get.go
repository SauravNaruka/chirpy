package main

import "net/http"

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while getting chirps", err)
		return
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

	respondWithJSON(w, http.StatusOK, chirps)
}
