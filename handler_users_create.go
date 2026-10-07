package main

import (
	"encoding/json"
	"net/http"

	"chirpy/internal/auth"
	"chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	createParams := database.CreateUserParams{Email: params.Email, HashedPassword: hash}
	user, err := cfg.db.CreateUser(r.Context(), createParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	usr := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, usr)
}
