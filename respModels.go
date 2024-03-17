package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/stonoy/PriceTracker/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type Product struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Url          string    `json:"url"`
	UserID       uuid.UUID `json:"user_id"`
	BasePrice    int       `json:"base_price"`
	CurrentPrice int       `json:"current_price"`
}

func userDbtoJson(DbUser database.User, token string) User {
	return User{
		ID:        DbUser.ID,
		CreatedAt: DbUser.CreatedAt,
		UpdatedAt: DbUser.UpdatedAt,
		Name:      DbUser.Name,
		Email:     DbUser.Email,
		Token:     token,
	}
}

func productDbtoJson(DbProduct database.Product) Product {
	return Product{
		ID:           DbProduct.ID,
		CreatedAt:    DbProduct.CreatedAt,
		UpdatedAt:    DbProduct.UpdatedAt,
		Name:         DbProduct.Name,
		Url:          DbProduct.Url,
		UserID:       DbProduct.UserID,
		BasePrice:    int(DbProduct.BasePrice.Int32),
		CurrentPrice: int(DbProduct.CurrentPrice.Int32),
	}
}

func allProductsDbToJson(DbProducts []database.Product) []Product {
	products := []Product{}

	for _, item := range DbProducts {
		products = append(products, productDbtoJson(item))
	}

	return products
}
