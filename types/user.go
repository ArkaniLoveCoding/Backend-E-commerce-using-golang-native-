package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)


type UserStore interface {
	UpdateToken(id uuid.UUID, token string, token_refresh string) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*[]User, error)
	CreateUser(context.Context, *User) error 
	GetAllUser([]User) (*[]User, error)
}

type User struct {
	ID 				uuid.UUID	`db:"id"`
	Firstname 		string  	`db:"firstname"`
	Lastname 		string 		`db:"lastname"`
	Password 		string 		`db:"password"`
	Email 			string 		`db:"email"`
	Country 		string 		`db:"country"`
	Address 		string 		`db:"address"`
	Role 			string 		`db:"role"`
	Token 			string  	`db:"token"`
	Rerfresh_token 	string 		`db:"refresh_token"`
	CreatedAt 		time.Time 	`db:"created_at"`
}

type Register struct {
	ID 				uuid.UUID	`json:"id"`
	Firstname 		string 		`json:"firstname" validate:"required,min=2,max=100"`
	Lastname 		string  	`json:"lastname" validate:"required,min=2,max=100"`
	Password 		string  	`json:"password" validate:"required,min=2,max=100"`
	Email 	 		string 		`json:"email" validate:"required,email"`
	Country 		string 		`json:"country" validate:"required,min=2,max=100"`
	Address 		string 		`json:"address" validate:"required,min=2,max=100"`
	Role  			string 		`json:"role" validate:"required,oneof=USER ADMIN"`
	Token 			string 		`json:"token"`
	Rerfresh_token	string 		`json:"refresh_token"`
	Created_at 		time.Time 	`json:"created_at"`
}

type Login struct {
	Email 			string	`json:"email" validate:"required,email"`
	Password 		string	`json:"pasword" validate:"required,min=2,max=100"`
}
