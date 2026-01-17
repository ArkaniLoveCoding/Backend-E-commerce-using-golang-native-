package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (h *HandleRequest) UpdateTokenFunc (id uuid.UUID, token string, token_refresh string) error {

	err := h.db.UpdateToken(id, token, token_refresh)
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

	validate = validator.New()
	if err := validate.Struct(&request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), false)
	}

	users, err := 
	h.db.GetUserByEmail(request.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to load the email", err.Error())
		return
	}
	if users != nil {
		utils.WriteError(w, http.StatusBadRequest, "Email has been already exist!", nil)
		return
	}

	hash, err := utils.HashPassword(request.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot hash the password of the user!", false)
		return
	}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()
	time_created := time.Now().UTC()

	var user = &types.User{
		ID: uuid.New(),
		Firstname: request.Firstname,
		Lastname: request.Lastname,
		Password: hash,
		Email: request.Email,
		Country: request.Country,
		Address: request.Address,
		Role: request.Role,
		Token: request.Token,
		Rerfresh_token: request.Rerfresh_token,
		CreatedAt: time_created,
	}

	if err := h.db.CreateUser(ctx, user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), false)
		return      
	}

	utils.WriteSuccess(w, http.StatusCreated, "Success to create new user !", user)
}

func (h *HandleRequest) LoginFunc(w http.ResponseWriter, r *http.Request) {

	var validate *validator.Validate
	var request types.Login
	if err := utils.DecodeData(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to decode the request", err.Error())
		return
	}

	validate = validator.New()
	if err := validate.Struct(request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create new user, because something missing!", err.Error())
		return
	}

	_, err := h.db.GetUserByEmail(request.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to find the email taht mactheds", err.Error())
		return 
	}

	var u types.User
	if err := utils.ComparePassword(u.Password, request.Password); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to compare the password!", err.Error())
		return
	}

	token, refresh_token, err := utils.GenerateJwt(
		u.ID, 
		u.Firstname, u.Lastname, 
		u.Password,
		u.Email,
		u.Role,
	)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to genereate jwt!", err.Error())
	}

	if err := h.UpdateTokenFunc(u.ID, token, refresh_token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update the token!", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Sucessfully!", u)

}

