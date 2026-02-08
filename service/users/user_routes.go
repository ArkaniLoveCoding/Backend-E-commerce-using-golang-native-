package service

import (
	"context"
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

func (h *HandleRequest) RegistrationFunc(w http.ResponseWriter, r *http.Request) {

	var validate *validator.Validate
	var request types.Register
	if err := utils.DecodeData(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot decode data payload!", false)
		return
	}

	isValid := utils.IsValidEmail(request.Email)
	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, "The format of your gmail is use @!", false)
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

	isValid := utils.IsValidEmail(request.Email)
	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, "The format of your gmail is @!", false)
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

	ctx_token, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()

	if err := h.db.UpdateToken(ctx_token, user.Id, token, refresh_token, user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update token!", err.Error())
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

	isValid := utils.IsValidEmail(user.Email)
	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get the profile because your email is not an valid email!", false)
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

	isValid := utils.IsValidEmail(request_update.Email)
	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update an email because the email not valid !", false)
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
	
	hashed_password, err := utils.HashPassword(request_update.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to load hashed password!", err.Error())
		return
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()

	var user = &types.User{
		Id: uuid_parse_id,
		Firstname: request_update.Firstname,
		Lastname: request_update.Lastname,
		Password: hashed_password,
		Email: request_update.Email,
		Country: request_update.Country,
		Address: request_update.Address,
	}


	if err := h.db.UpdateDataUser(
		uuid_parse_id, 
		ctx,
		request_update.Firstname,
		request_update.Lastname,
		request_update.Password,
		request_update.Email,
		request_update.Country,
		request_update.Address, 
		user,
		); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update token!", err.Error())
		return 
	}


	user_update_response := types.UserUpdateResponse{
		Id: user.Id,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Password: user.Password,
		Email: user.Email,
		Country: user.Country,
		Address: user.Address,
		Role: users.Role,
		Token: users.Token,
		Rerfresh_token: users.Rerfresh_token,
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

func (h *HandleRequest) GetAllUser(w http.ResponseWriter, r *http.Request) {

	middleware_checking_role, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get role from middleware!", err.Error())
		return 
	}

	if middleware_checking_role != "ADMIN" {
		utils.WriteError(w, http.StatusBadRequest, "Only admin can access this method!", false)
	}

	users, err := h.db.GetAllUser()
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "Failed to get all data of user!", err.Error())
		return 
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to get all the data from users db!", users)

}

func (h *HandleRequest) GetOneUsersById(w http.ResponseWriter, r *http.Request) {

	vars_id := mux.Vars(r)
	id := vars_id["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to load params because the params is nill!", false)
		return 
	}

	uuid_parse_id, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Faield to convert data into uuid type!", err.Error())
		return 
	}

	users, err := h.db.GetUserById(uuid_parse_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Faile dot get one id from db!", err.Error())
		return 
	}

	if users == nil {
		utils.WriteError(w, http.StatusBadRequest, "Faield to get users because users is nil!", false)
	}

	time_created_format := users.Created_at.Format("2006-01-02")
	time_updated_format := users.Updated_at.Format("2006-01-02")

	user_response := types.UserResponse{
		Id: users.Id,
		Firstname: users.Firstname,
		Lastname: users.Lastname,
		Password: users.Password,
		Email: users.Email,
		Country: users.Country,
		Address: users.Address,
		Role: users.Role,
		Token: users.Token,
		Rerfresh_token: users.Rerfresh_token,
		Created_at: time_created_format,
		Updated_at: time_updated_format,
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to get one data from users db!", user_response)

}
