package main

import (
	"database/sql"
	"net/http"

	"github.com/soramon0/rssagg/auth"
	"github.com/soramon0/rssagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, 404, "user not found")
				return
			}
			respondWithError(w, 400, err.Error())
			return
		}

		handler(w, r, user)
	}
}
