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

// this is for router that token is verified in their function!

type HandleRequestForAuthenticate struct {
	db types.UserStore
}

func NewHandlerUserForAuthenticate (db types.UserStore) *HandleRequestForAuthenticate {
	return &HandleRequestForAuthenticate{db: db}
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
	time_updated := time.Now().UTC()
	time_format_created := time_created.Format("2006-01-02")
	time_format_updated := time_updated.Format("2006-01-02")

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
		Created_at: time_created,
		Updated_at: time_updated,
	}

	user_response := types.UserResponse{
		Id: user.Id,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Password: user.Password,
		Email: user.Email,
		Country: user.Country,
		Address: user.Address,
		Role: user.Role,
		Created_at: time_format_created,
		Updated_at: time_format_updated,
	}

	if err := h.db.CreateUser(ctx, user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), false)
		return      
	}

	utils.WriteSuccess(w, http.StatusCreated, "Success to create new user !", user_response)
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

	if user == nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot find the email !", false)
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

	time_created := time.Now().UTC()
	time_updated := time.Now().UTC()
	time_format_created := time_created.Format("2006-01-02")
	time_format_updated := time_updated.Format("2006-01-02")

	user_response := types.UserResponse{
		Id: user.Id,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Password: user.Password,
		Email: user.Email,
		Country: user.Country,
		Address: user.Address,
		Role: user.Role,
		Token: token,
		Rerfresh_token: refresh_token,
		Created_at: time_format_created,
		Updated_at: time_format_updated,
	}

	utils.WriteSuccess(w, http.StatusOK, "Sucessfully!", user_response)

}

func (h *HandleRequestForAuthenticate) GetProfileUser(w http.ResponseWriter, r *http.Request) {

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

	user_response := types.UserResponse{
		Id: user.Id,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Password: user.Password,
		Email: user.Email,
		Country: user.Country,
		Address: user.Address,
		Role: user.Role,
	}

	utils.WriteSuccess(w, http.StatusAccepted, "Successfully to get profile!", user_response)

}

func (h *HandleRequest) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var validate *validator.Validate
	var request_update types.UserUpdate

	if err := utils.DecodeData(r, &request_update); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to decode data!", err.Error())
		return
	}

	validate = validator.New()
	if err := validate.Struct(&request_update); err != nil {
		var errorValidate []string
		for _, errors := range err.(validator.ValidationErrors) {
			errorValidate = append(errorValidate, fmt.Sprintf("Fatal Error ! : %v, %v", errors.Field(), errors.Tag()))
			return
		}
	}

	vars_id := mux.Vars(r)
	id := vars_id["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to detect id !", false)
		return 
	}

	uuid_parse_id, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert data as a uuid!", err.Error())
		return 
	}

	middleware_check_account_id, err := middleware.GetValueTokenID(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get token id from jwt token!", err.Error())
		return 
	}

	if middleware_check_account_id != uuid_parse_id {
		utils.WriteError(w, http.StatusInternalServerError, "Cannot change other profile!", false)
		return
	}


	users, err := h.db.GetUserById(uuid_parse_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get id!", err.Error())
		return 
	}

	if users == nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot find the data from id!", false)
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()


	if err := h.db.UpdateDataUser(
		uuid_parse_id, 
		ctx,
		request_update.Firstname,
		request_update.Lastname,
		request_update.Password,
		request_update.Email,
		request_update.Country,
		request_update.Address, 
		); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update token!", err.Error())
		return 
	}
	
	var user = &types.User{
		Id: uuid_parse_id,
		Firstname: request_update.Firstname,
		Lastname: request_update.Lastname,
		Password: request_update.Password,
		Email: request_update.Email,
		Country: request_update.Country,
		Address: request_update.Address,
	}


	user_update_response := types.UserUpdateResponse{
		Id: user.Id,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Password: user.Password,
		Email: user.Email,
		Country: user.Country,
		Address: user.Address,
		Created_at: users.Created_at.Format("2006-01-02"),
		Updated_at: time.Now().UTC().Format("2006-01-02"),
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to update data!", user_update_response)
}

func (h *HandleRequest) DeleteUser(w http.ResponseWriter, r *http.Request) {

	middleware_checking_role, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Something went wrong!", err.Error())
		return 
	}

	if middleware_checking_role != "ADMIN" {
		utils.WriteError(w, http.StatusBadGateway, "Only admin can delete the other users account!", false)
		return
	}

	vars_id := mux.Vars(r)
	id := vars_id["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "the id that you want to delete is nil!", false)
		return 
	}

	uuid_parse_id, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert into uuid type!", err.Error())
		return 
	}

	users, err := h.db.GetUserById(uuid_parse_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to checking data from db!", err.Error())
		return 
	}

	if users == nil {
		utils.WriteError(w, http.StatusBadRequest, "The users is nil from db!", false)
		return 
	}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	if err := h.db.DeleteUsersOnlyAdmin(uuid_parse_id, ctx); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to delete the users from db only admin!", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to delete users account as a admin!", true)

}
