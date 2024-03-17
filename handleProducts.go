package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/stonoy/PriceTracker/internal/database"
)

func (cfg *apiConfig) createProduct(w http.ResponseWriter, r *http.Request, user database.User) {
	type reqStruct struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	reqObj := reqStruct{}
	err := decoder.Decode(&reqObj)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("error in parsing incoming json : %v", err))
		return
	}

	product, err := cfg.DB.CreateProduct(r.Context(), database.CreateProductParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqObj.Name,
		Url:       reqObj.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("error in creating a new product : %v", err))
		return
	}

	respWithJson(w, 201, productDbtoJson(product))
}

func (cfg *apiConfig) productByUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	products, err := cfg.DB.FindProductsByUser(r.Context(), user.ID)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("error in getting user's products : %v", err))
		return
	}

	type respStruct struct {
		Products []Product `json:"products"`
	}

	respWithJson(w, 200, respStruct{
		Products: allProductsDbToJson(products),
	})
}
