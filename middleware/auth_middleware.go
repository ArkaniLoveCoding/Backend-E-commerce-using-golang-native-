package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/utils"
)

// middleware func to get id from user token as a user token

func AuthenticateForIdUser (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token_header := r.Header.Get("Authorization")
		if token_header == "" {
			utils.WriteError(w, http.StatusBadRequest, "Cant find the auth!", false)
			return 
		}

		token_key := strings.TrimPrefix(token_header, "Bearer ")
		if token_key == "" {
			utils.WriteError(w, http.StatusBadRequest, "Failed to authencticate the bearer method!", false)
			return
		}

		validate, err := utils.ValidateToken(token_key)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Failed to verify jwt!", err.Error())
			return
		}

		uuid_user_id, err := uuid.Parse(validate.Id)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Failed to convert uuid!", err.Error())
			return 
		}

		ctx_id := context.WithValue(r.Context(), "id", uuid_user_id)
		r = r.WithContext(ctx_id)

		next.ServeHTTP(w, r)


	})
}

// func that wants to return id from user token as a token id

func GetValueTokenID (w http.ResponseWriter, r *http.Request) (uuid.UUID, error)  {

	user := r.Context().Value("id")
	if user == "" || user == nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get id of the user !", false)
		return uuid.Nil, nil
	}
	convert_into_uuid, ok := user.(uuid.UUID)
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "Failed!", false)
		return uuid.Nil, nil
	}

	return convert_into_uuid, nil

}

// middleware func for authenticate role user 

func AuthenticateForRole (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header_token := r.Header.Get("Authorization")
		if header_token == "" {
			utils.WriteError(w, http.StatusBadRequest, "Failed to get token!", false)
			return 
		}

		token_bearer := strings.TrimPrefix(header_token, "Bearer ")
		if token_bearer == "" {
			utils.WriteError(w, http.StatusBadRequest, "Failed to convert into bearer token as a token!", false)
			return
		}

		token_key, err  := utils.ValidateToken(token_bearer)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Failed to validate token as a token from user!", err.Error())
			return 
		}

		context_save := context.WithValue(r.Context(), "role", token_key.Role)
		r = r.WithContext(context_save)

		next.ServeHTTP(w, r)

	})
}

// func that wants to take the role from token as a user role default 

func GetValueTokenRole (w http.ResponseWriter, r *http.Request) (string, error) {

	role_user := r.Context().Value("role")
	if role_user == "" || role_user == nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get user role, because the data is nothing!", false)
		return "", nil
	}

	role_string, ok := role_user.(string)
	if !ok {
		return "", errors.New("Failed to convert from type any to string!")
	}

	return role_string, nil

}