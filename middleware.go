package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/stonoy/PriceTracker/internal/database"
	"github.com/stonoy/PriceTracker/utils"
)

// define a func type, that need a user based on the apikey in the request header
type funcNeedJwtToken func(http.ResponseWriter, *http.Request, database.User)

// a middleware that accept the above mentioned func type and provide it the apikey from request header also it returns http.HandleFunc signature to be valid in router.method
func (cfg *apiConfig) authMiddleware(anySuitableFunc funcNeedJwtToken) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the jwt token
		token, err := utils.GetTokenFromHeader(r)
		if err != nil {
			respWithError(w, 403, fmt.Sprintf("error getting jwt token: %v", err))
			return
		}

		// get user id from token
		userId, _, err := utils.GetDataFromToken(token, cfg.Jwt_Secret)
		if err != nil {
			respWithError(w, 403, fmt.Sprintf("error parsing jwt token: %v", err))
			return
		}

		// fmt.Println(userId)

		// Parse the UUID string
		userIdUUID, err := uuid.Parse(userId)
		if err != nil {
			respWithError(w, 403, fmt.Sprintf("Error parsing UUID: %v\n", err))
			return
		}

		// get user from apikey using func generated by sqlc
		user, err := cfg.DB.GetUserByJwtToken(r.Context(), userIdUUID)
		if err != nil {
			respWithError(w, 403, fmt.Sprintf("error getting user by jwt token: %v", err))
			return
		}

		anySuitableFunc(w, r, user)

	}
}