package types

import "time"

type ProductStore interface {
	GetAllProduct([]User) (*[]User, error)
}

type Products struct {
	ID         		int    		`db:"id"`
	Name      		string 		`db:"name"`
	Stock     		int    		`db:"stock"`
	Price      		string 		`db:"price"`
	Expired    		string 		`db:"expired"`
	Category   		string 		`db:"category"`
	Created_at 		time.Time	`db:"created_at"`
	Updated_at 		time.Time 	`db:"updated_at"`
}

type ProductResponse struct {
	Name 			string 		`json:"name" validate:"required,min=2,max=100"`
	Stock 			int 		`json:"stock" validate:"required,min=2,max=100"`
	Expired 		string		`json:"expired" validate:"required,min=2,max=100"`
	Category 		string 		`json:"category" validate:"required"`
}