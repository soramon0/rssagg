package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/soramon0/rssagg/internal/database"
)

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollow, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows", err))
		return
	}
	respondWithJSON(w, 200, feedFollow)
}

func (apiCfg *apiConfig) handleDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, 400, "invalid id")
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	p := params{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON:", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    p.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed:", err))
		return
	}
	respondWithJSON(w, 201, feed)
}
