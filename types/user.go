package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)


type UserStore interface {
	UpdateToken(id uuid.UUID, token string, token_refresh string) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(id uuid.UUID) (*User, error)
	CreateUser(context.Context, *User) error 
	GetAllUser() ([]User, error)
	GetUsersRole(role string) (*User, error)
	UpdateDataUser(
		id uuid.UUID,
		ctx context.Context,
		firstname string,
		lastname string,
		password string,
		email string,
		country string,
		address string,
		user *User,
		) error
	DeleteUsersOnlyAdmin(id uuid.UUID, ctx context.Context) error
}

type User struct {
	Id				uuid.UUID	`db:"id"`
	Firstname 		string  	`db:"firstname"`
	Lastname 		string 		`db:"lastname"`
	Password 		string 		`db:"password"`
	Email 			string 		`db:"email"`
	Country 		string 		`db:"country"`
	Address 		string 		`db:"address"`
	Role 			string 		`db:"role"`
	Token 			string  	`db:"token"`
	Rerfresh_token 	string 		`db:"refresh_token"`
	Created_at 		time.Time 	`db:"created_at"`
	Updated_at		time.Time 	`db:"updated_at"`
}

type Register struct {
	Id 				uuid.UUID	`json:"id"`
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
	Updated_at 		time.Time 	`json:"updated_at"`
}

type Login struct {
	Email 			string	`json:"email" validate:"required,email"`
	Password 		string	`json:"password" validate:"required,min=2,max=100"`
}

type UserUpdate struct {
	Id 				uuid.UUID	`json:"id"`
	Firstname 		string 		`json:"firstname" validate:"required,min=2,max=100"`
	Lastname 		string  	`json:"lastname" validate:"required,min=2,max=100"`
	Password 		string  	`json:"password" validate:"required,min=2,max=100"`
	Email 	 		string 		`json:"email" validate:"required,email"`
	Country 		string 		`json:"country" validate:"required,min=2,max=100"`
	Address 		string 		`json:"address" validate:"required,min=2,max=100"`
}

type UserUpdateResponse struct {
	Id 				uuid.UUID	`json:"id"`
	Firstname 		string 		`json:"firstname" validate:"required,min=2,max=100"`
	Lastname 		string  	`json:"lastname" validate:"required,min=2,max=100"`
	Password 		string  	`json:"password" validate:"required,min=2,max=100"`
	Email 	 		string 		`json:"email" validate:"required,email"`
	Country 		string 		`json:"country" validate:"required,min=2,max=100"`
	Address 		string 		`json:"address" validate:"required,min=2,max=100"`
	Role  			string 		`json:"role" validate:"required,oneof=USER ADMIN"`
	Token 			string 		`json:"token"`
	Rerfresh_token	string 		`json:"refresh_token"`
	Created_at 		string 		`json:"created_at"`
	Updated_at 		string 		`json:"updated_at"`
}

type UserResponse struct {
	Id				uuid.UUID	`json:"id"`
	Firstname 		string 		`json:"firstname"`
	Lastname 		string  	`json:"lastname"`
	Password 		string  	`json:"password"`
	Email 	 		string 		`json:"email"`
	Country 		string 		`json:"country"`
	Address 		string 		`json:"address"`
	Role  			string 		`json:"role"`
	Token 			string 		`json:"token"`
	Rerfresh_token	string 		`json:"refresh_token"`
	Created_at 		string  	`json:"created_at"`
	Updated_at 		string 		`json:"updated_at"`
}
