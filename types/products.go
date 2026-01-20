package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProductStore interface {
	GetAllProduct([]Products) (*Products, error)
	GetProductByID(id string) (*Products, error)
	CreateNewProduct(ctx context.Context, products *Products) error
}

type Products struct {
	Id        		uuid.UUID    `db:"id"`
	Name      		string 		`db:"name"`
	Stock     		int    		`db:"stock"`
	Price      		string 		`db:"price"`
	Expired    		string 		`db:"expired"`
	Category   		string 		`db:"category"`
	Created_at 		time.Time	`db:"created_at"`
	Updated_at 		time.Time 	`db:"updated_at"`
}

type ProductResponse struct {
	Id				uuid.UUID 	`json:"id"`
	Name 			string 		`json:"name" validate:"required,min=2,max=100"`
	Stock 			int 		`json:"stock" validate:"required,min=2,max=100"`
	Price 			string 		`json:"price" validate:"required,min=2,max=100"`
	Expired 		string		`json:"expired" validate:"required,min=2,max=100"`
	Category 		string 		`json:"category" validate:"required"`
	Created_at 		time.Time 	`json:"created_at"`
	Updated_at 		time.Time 	`json:"updated_at"`
}