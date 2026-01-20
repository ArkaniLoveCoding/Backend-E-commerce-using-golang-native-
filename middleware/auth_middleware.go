package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/utils"
)

func AuthenticateProfile (next http.Handler) http.Handler {
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

func GetValueTokenID (w http.ResponseWriter, r *http.Request) (uuid.UUID, error)  {

	user := r.Context().Value("id")
	if user == "" {
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

// later,  i will fix it

func GetValueTokenRole (w http.ResponseWriter, r *http.Request) (string, error)  {

	role := r.Context().Value("role")
	if role == "" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get role of the user!", false)
		return "", nil
	}

	return "", nil

}