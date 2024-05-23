package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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

	// id, _ := uuid.Parse("98adf320-6789-4cec-b1bd-0fa5ffe7db04")

	// product, _ := cfg.DB.GetProductById(r.Context(), id)
	// log.Println(product.Priority)

	type respStruct struct {
		Products []Product `json:"products"`
	}

	respWithJson(w, 200, respStruct{
		Products: allProductsDbToJson(products),
	})
}

func (cfg *apiConfig) updateProductPriority(w http.ResponseWriter, r *http.Request, user database.User) {

	idString := chi.URLParam(r, "productId")
	id, err := uuid.Parse(idString)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("can not parse feed id: %v", err))
		return
	}

	// check user has created the product
	product, err := cfg.DB.GetProductById(r.Context(), id)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("can get the product of the id - %v: %v", id, err))
		return
	}

	if product.UserID != user.ID {
		respWithError(w, 403, "Not authrised to update the product")
		return
	}

	_, err = cfg.DB.UpdateProductPriority(r.Context(), id)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("can not update product priority : %v", err))
		return
	}

	type respStruct struct {
		Updated string `json:"updated"`
	}

	respWithJson(w, 200, respStruct{Updated: "ok"})

}

func (cfg *apiConfig) deleteProduct(w http.ResponseWriter, r *http.Request, user database.User) {
	idString := chi.URLParam(r, "productId")
	id, err := uuid.Parse(idString)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("can not parse feed id: %v", err))
		return
	}

	// check user has created the product
	product, err := cfg.DB.GetProductById(r.Context(), id)
	if err != nil {
		respWithError(w, 400, fmt.Sprintf("can get the product of the id - %v: %v", id, err))
		return
	}

	if product.UserID != user.ID {
		respWithError(w, 403, "Not authrised to delete the product")
		return
	}

	err = cfg.DB.DeleteProduct(r.Context(), id)
	if err != nil {
		respWithError(w, 200, fmt.Sprintf("can not delete the product : %v", err))
		return
	}

	type respStruct struct {
		Deleted string `json:"deleted"`
	}

	respWithJson(w, 200, respStruct{Deleted: "ok"})
}
