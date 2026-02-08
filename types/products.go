package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProductStore interface {
	GetAllProduct() ([]Products, error)
	GetProductByID(id uuid.UUID) (*Products, error)
	CreateNewProduct(ctx context.Context, products *Products) error
	DeleteProductsOnlyAdmin(id uuid.UUID, ctx context.Context) error 
	UpdateProductsOnlyAdmin(
		id uuid.UUID,
		name string,
		stock int,
		image string,
		category string,
		price string,
		expired string,
		ctx_update context.Context,
	) error
	SearchManyProducts(ctx context.Context, keyword string, offset int, limit int) ([]Products, error)
}

type Products struct {
	Id        		uuid.UUID    `db:"id"`
	Name      		string 		`db:"name"`
	Stock     		int    		`db:"stock"`
	Image 			string 		`db:"image"`
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
	Image 			string		`json:"image" validate:"required"`
	Price 			string 		`json:"price" validate:"required,min=2,max=100"`
	Expired 		string		`json:"expired" validate:"required,min=2,max=100"`
	Category 		string 		`json:"category" validate:"required,min=2,max=100"`
	Created_at 		string  	`json:"created_at"`
	Updated_at 		string  	`json:"updated_at"`
}

type PayloadUpdateAndCreate struct {
	Id				uuid.UUID 	`json:"id"`
	Name 			string 		`json:"name" validate:"required,min=2,max=100"`
	Stock 			int 		`json:"stock" validate:"required,min=2,max=100"`
	Image 			string 		`json:"image" validate:"required"`
	Price 			string 		`json:"price" validate:"required,min=2,max=100"`
	Expired 		string		`json:"expired" validate:"required,min=2,max=100"`
	Category 		string 		`json:"category" validate:"required,min=2,max=100"`
	Created_at 		string  	`json:"created_at"`
	Updated_at 		string  	`json:"updated_at"`
}