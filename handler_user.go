package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/soramon0/rssagg/internal/database"
)

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
	respondWithJSON(w, 200, user)
}
