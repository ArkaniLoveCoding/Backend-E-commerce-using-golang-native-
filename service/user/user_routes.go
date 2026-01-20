package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/middleware"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/utils"
)

// this is for router that token is not verified in their function!

type HandleRequest struct {
	db types.UserStore
}

func NewHandlerUser(db types.UserStore) *HandleRequest {
	return &HandleRequest{db: db}
}

//

// this is for router that token is verified in their function!

type HandleRequestForAuthenticate struct {
	db types.UserStore
}

func NewHandlerUserForAuthenticate (db types.UserStore) *HandleRequestForAuthenticate {
	return &HandleRequestForAuthenticate{db: db}
}

//


func (h *HandleRequest) RegistrationUserHandler(router *mux.Router) {
	router.HandleFunc("/registration", h.RegistrationFunc).Methods("POST")
}

func (h *HandleRequest) LoginUserHandler(router *mux.Router) {
	router.HandleFunc("/login", h.LoginFunc).Methods("POST")
}

func (h *HandleRequestForAuthenticate) GetProfileHandler(router *mux.Router) {
	router.HandleFunc("/profile", h.GetProfileUser).Methods("GET")
}

func (h *HandleRequest) UpdateTokenFunc (id uuid.UUID, token string, token_refresh string) error {

	err := h.db.UpdateToken(id, token, token_refresh)
	if err != nil {
		return errors.New("Failed to update the token!" + err.Error())
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
		var errorValidate []string
		for _, errors := range err.(validator.ValidationErrors) {
			errorValidate = append(errorValidate, fmt.Sprintf("Fatal Error ! : %v, %v", errors.Field(), errors.Tag()))
			return
		}
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
		Id: uuid.New(),
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
		var errorValidate []string 
		for _, errors := range err.(validator.ValidationErrors) {
			errorValidate = append(errorValidate, fmt.Sprintf("Fatal Erorr ! : %v, %v", errors.Field(), errors.Tag()))
			return
		}
	}

	user, err := h.db.GetUserByEmail(request.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to find the email taht mactheds", err.Error())
		return 
	}

	if err := utils.ComparePassword(user.Password, request.Password); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to compare the password!", err.Error())
		return
	}

	token, refresh_token, err := utils.GenerateJwt(
		user.Id, 
		user.Firstname, 
		user.Lastname, 
		user.Password,
		user.Email,
		user.Role,
	)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to genereate jwt!", err.Error())
	}

	if err := h.UpdateTokenFunc(user.Id, token, refresh_token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update the token!", err.Error())
		return
	}

	// place to take the token after the token has been created
	fmt.Println(token)
	fmt.Println(refresh_token)

	user_response := types.UserResponse{
		Id: user.Id,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Password: user.Password,
		Email: user.Email,
		Country: user.Country,
		Address: user.Address,
		Role: user.Role,
		Token: user.Token,
		Rerfresh_token: user.Rerfresh_token,
		Created_at: user.CreatedAt,
	}

	utils.WriteSuccess(w, http.StatusOK, "Sucessfully!", user_response)

}

func (h *HandleRequestForAuthenticate) GetProfileUser (w http.ResponseWriter, r *http.Request) {

	user_id, err := middleware.GetValueTokenID(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to find the id from token !", err.Error())
		return 
	}

	user, err := h.db.GetUserById(user_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "Failed to get Id user from db!", err.Error())
		return 
	}

	utils.WriteSuccess(w, http.StatusAccepted, "Successfully to get profile!", user)

}
