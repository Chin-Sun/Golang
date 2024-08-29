package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chin-Sun/Golang/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerFeedFollowsCreat(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	parame := parameters{}
	err := decoder.Decode(&parame)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode parameters: %v", err))
		return
	}

	feed_follow, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		FeedID:    parame.FeedID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowToFeedFollow(feed_follow))
}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get feed follows")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollows(feed_follows))
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
