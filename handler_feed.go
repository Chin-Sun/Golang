package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chin-Sun/Golang/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerFeedsCreat(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		NAME string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	parame := parameters{}
	err := decoder.Decode(&parame)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode parameters: %v", err))
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      parame.NAME,
		Url:       parame.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedToFeed(feed))
}

func (cfg *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
