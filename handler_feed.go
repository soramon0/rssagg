package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/soramon0/rssagg/internal/database"
)

func (apiCfg *apiConfig) handleListFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.ListFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't fetch feeds:", err))
		return
	}
	respondWithJSON(w, 200, feeds)
}

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
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
	if p.Url == "" {
		respondWithError(w, 400, "Url is required")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      p.Name,
		Url:       p.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed:", err))
		return
	}
	respondWithJSON(w, 201, feed)
}
