package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/stonoy/PriceTracker/internal/database"
	"github.com/stonoy/PriceTracker/utils"
)

func (cfg *apiConfig) registerUser(w http.ResponseWriter, r *http.Request) {
	type reqStruct struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	reqObj := reqStruct{}
	err := decoder.Decode(&reqObj)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("error in parsing incoming json : %v", err))
		return
	}

	if reqObj.Name == "" {
		respWithError(w, 400, "provide a name ")
		return
	}

	if reqObj.Email == "" {
		respWithError(w, 400, "provide a email ")
		return
	}

	if reqObj.Password == "" || len(reqObj.Password) < 6 {
		respWithError(w, 400, "provide a valid password ")
		return
	}

	hashedPassword, err := utils.HashPassword(reqObj.Password)
	if err != nil {
		respWithError(w, 500, fmt.Sprintf("error in hashing password : %v", err))
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqObj.Name,
		Email:     reqObj.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("error in creating a new user : %v", err))
		return
	}

	token, err := utils.CreateJwtAccessToken(cfg.Jwt_Secret, user)
	if err != nil {
		respWithError(w, 500, fmt.Sprintf("error in creating a new token : %v", err))
		return
	}

	respWithJson(w, 201, userDbtoJson(user, token))
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type reqStruct struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	reqObj := reqStruct{}
	err := decoder.Decode(&reqObj)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("error in parsing incoming json : %v", err))
		return
	}

	if reqObj.Email == "" {
		respWithError(w, 400, "provide a email ")
		return
	}

	if reqObj.Password == "" || len(reqObj.Password) < 6 {
		respWithError(w, 400, "provide a valid password ")
		return
	}

	user, err := cfg.DB.FindUserByEmail(r.Context(), reqObj.Email)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("No such user exists: %v", err))
		return
	}

	hasPasswordMatched := utils.IsPasswordMatches(reqObj.Password, user.Password)
	if !hasPasswordMatched {
		respWithError(w, 403, fmt.Sprintf("password not matches: %v", err))
		return
	}

	token, err := utils.CreateJwtAccessToken(cfg.Jwt_Secret, user)
	if err != nil {
		respWithError(w, 500, fmt.Sprintf("error in creating a new token : %v", err))
		return
	}

	respWithJson(w, 200, userDbtoJson(user, token))
}
