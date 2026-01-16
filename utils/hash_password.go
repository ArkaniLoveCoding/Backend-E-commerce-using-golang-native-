package utils

import "golang.org/x/crypto/bcrypt"


func HashPassword (password string) ([]byte, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return []byte{}, nil 
	}

	return hash, nil

}