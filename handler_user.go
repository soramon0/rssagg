package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/soramon0/rssagg/auth"
	"github.com/soramon0/rssagg/internal/database"
)

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "user not found")
			return
		}
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 200, user)
}

func (apiCfg *apiConfig) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.ListUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't fetch user:", err))
		return
	}
	respondWithJSON(w, 200, users)
}

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string `json:"name"`
	}
	p := params{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON:", err))
		return
	}
	if p.Name == "" {
		respondWithError(w, 400, "Name is required")
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      p.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user:", err))
		return
	}
	respondWithJSON(w, 201, user)
}
