package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/utils"
)

type HandleRequest struct {
	db types.UserStore
}

func NewHandlerUser(db types.UserStore) *HandleRequest {
	return &HandleRequest{db: db}
}

func (h *HandleRequest) RegistrationUserHandler(router *mux.Router) {
	router.HandleFunc("/registration", h.RegistrationFunc).Methods("POST")
}

func (h *HandleRequest) LoginUserHandler(router *mux.Router) {
	router.HandleFunc("/login", h.LoginFunc).Methods("POST")
}

func (h *HandleRequest) UpdateTokenFunc (id int, token string, token_refresh string) error {

	_, err := h.db.UpdateToken(id, token, token_refresh)
	if err != nil {
		return errors.New("Failed to update the token!")
	}

	return nil

}

func (h *HandleRequest) RegistrationFunc(w http.ResponseWriter, r *http.Request) {

	var validate *validator.Validate
	var request types.Register
	if err := utils.DecodeData(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot decode data payload!", false)
		return
	}

	user, err := 
	h.db.GetUserByEmail(request.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to load the email", err.Error())
		return
	}
	if user != nil {
		utils.WriteError(w, http.StatusBadRequest, "Email has been already exist!", nil)
		return
	}

	hash, err := utils.HashPassword(request.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot hash the password of the user!", false)
		return
	}

	validate = validator.New()
	if err := validate.Struct(&request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), false)
	}
	context := context.Background()

	if err := h.db.CreateUser(context, types.User{
		Firstname: request.Firstname,
		Lastname: request.Lastname,
		Password: string(hash),
		Email: request.Email,
		Country: request.Country,
		Address: request.Address,
		Role: request.Role,
		Token: request.Token,
		Rerfresh_token: request.Rerfresh_token,
	}); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Cant create new user !", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, "Success to create new user!", request)

}

func (h *HandleRequest) LoginFunc(w http.ResponseWriter, r *http.Request) {

	

}

