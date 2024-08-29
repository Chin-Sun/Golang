package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chin-Sun/Golang/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreat(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		NAME string
	}
	decoder := json.NewDecoder(r.Body)
	parame := parameters{}
	err := decoder.Decode(&parame)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode parameters: %v", err))
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      parame.NAME,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}
