package product

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


type HandleRequest struct {
	db types.ProductStore
}

func NewHandlerProduct(db types.ProductStore) *HandleRequest {
	return &HandleRequest{db: db}
}

// testing 

func (h *HandleRequest) CreateNewProductTesting(w http.ResponseWriter, r *http.Request) {

	var validate *validator.Validate
	var request types.ProductResponse
	if err := utils.DecodeData(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to decode the response!", err.Error())
		return 
	}

	validate = validator.New()
	if err := validate.Struct(&request); err != nil {
		var errors []string
		for _, errorValidate := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Fatal Error ! : %v, %v", errorValidate.Field(), errorValidate.Tag()))
		}
	}
	time_created := time.Now().UTC()
	time_updated := time.Now().UTC()

	products := &types.Products{
		Id: request.Id,
		Name: request.Name,
		Stock: request.Stock,
		Price: request.Price,
		Expired: request.Expired,
		Category: request.Category,
		Created_at: time_created,
		Updated_at: time_updated,
	}
	
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	if err := h.db.CreateNewProduct(ctx, products); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create new data of products!", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusAccepted, "Successfully to create new products!", products)


}

// controllers

func (h *HandleRequest) CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	role_user, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get role for the user from token!", err.Error())
		return 
	}

	if role_user != "ADMIN" {
		utils.WriteError(w, http.StatusBadRequest, "Cannot access this method because your role is not admin!", false)
		return
	}

	var validate *validator.Validate
	var request types.ProductResponse
	if err := utils.DecodeData(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to decode the response!", err.Error())
		return 
	}

	validate = validator.New()
	if err := validate.Struct(&request); err != nil {
		var errors []string
		for _, errorValidate := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Fatal Error ! : %v, %v", errorValidate.Field(), errorValidate.Tag()))
		}
	}
	time_created := time.Now().UTC()
	time_updated := time.Now().UTC()

	time_format_created := time_created.Local().Format("2006-01-02")
	time_format_updated := time_updated.Format("2006-01-02")

	products := &types.Products{
		Id: uuid.New(),
		Name: request.Name,
		Stock: request.Stock,
		Price: request.Price,
		Expired: request.Expired,
		Category: request.Category,
		Created_at: time_created,
		Updated_at: time_updated,
	}

	product_response := types.ProductResponse{
		Id: products.Id,
		Name: products.Name,
		Stock: products.Stock,
		Price: products.Price,
		Expired: products.Expired,
		Category: products.Category,
		Created_at: time_format_created,
		Updated_at: time_format_updated,
	}
	
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	if err := h.db.CreateNewProduct(ctx, products); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create new data of products!", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusAccepted, "Successfully to create new products!", product_response)

}

func (h *HandleRequest) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	uuid_parse_id, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert into uuid!", err.Error())
		return
	}

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Cannot find the params!", false)
		return 
	}

	products, err := h.db.GetProductByID(uuid_parse_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot find the data from id!", err.Error())
		return 
	}

	if products == nil {
		utils.WriteError(w, http.StatusBadRequest, "The id that you want to find is nil!", false)
		return 
	}

	product_response := &types.Products{
		Id: products.Id,
		Name: products.Name,
		Stock: products.Stock,
		Price: products.Price,
		Expired: products.Expired,
		Category: products.Category,
		Created_at: products.Created_at,
		Updated_at: products.Updated_at,
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully!", product_response)

}

func (h *HandleRequest) GetAllProduct(w http.ResponseWriter, r *http.Request) {

	products, err := h.db.GetAllProduct()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to find all data of product!", err.Error())
		return 
	}

	utils.WriteSuccess(w, http.StatusAccepted, "Successfully to get all data of product", products)

}

func (h *HandleRequest) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	middleware_checking_role, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get role from jwt token!", err.Error())
		return 
	}

	if middleware_checking_role != "ADMIN" {
		utils.WriteError(w, http.StatusBadRequest, "Only admin can delete one of the data of the product db!", false)
		return
	}

	vars_id := mux.Vars(r)
	id := vars_id["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "the id that you want to be a params is nil!", false)
		return 
	}

	uuid_parse_id, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert the data to uuid!", err.Error())
		return
	}

	users, err := h.db.GetProductByID(uuid_parse_id)

	if users == nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get users from db use id at param !", false)
		return 
	}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	if err := h.db.DeleteProductsOnlyAdmin(uuid_parse_id, ctx); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to delete the product!", err.Error())
		return 
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to delete the one of data of the product db!", true)
	
}

func (h *HandleRequest) UpdateProductsOnlyAdmin(w http.ResponseWriter, r *http.Request) {

	var validate *validator.Validate
	var payload_update types.PayloadUpdate

	if err := utils.DecodeData(r, &payload_update); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update the product data!", err.Error())
		return
	}

	validate = validator.New()
	if err := validate.Struct(&payload_update); err != nil {
		var errors []string
		for _, errorValidate := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Fatal Error ! : %v, %v", errorValidate.Field(), errorValidate.Tag()))
		}
	}

	middleware_get_role, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get the role from token!", err.Error())
		return
	}

	if middleware_get_role != "ADMIN" {
		utils.WriteError(w, http.StatusBadRequest, "the role cant be access here is admin role!", false)
		return
	}

	vars_id := mux.Vars(r)
	id := vars_id["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get params from id postman!", false)
	}

	uuid_parse, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert from string into an uuid type !", err.Error())
		return
	}

	products, err := h.db.GetProductByID(uuid_parse)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get data products from id!", err.Error())
		return 
	}

	if products == nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get products because is nill!", false)
		return
	}

	var product = &types.Products{
		Id: payload_update.Id,
		Name: payload_update.Name,
		Price: payload_update.Price,
		Stock: payload_update.Stock,
		Category: payload_update.Category,
		Expired: payload_update.Expired,
	}

	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()

	if err := h.db.UpdateProductsOnlyAdmin(
		uuid_parse,
		payload_update.Name,
		payload_update.Stock,
		payload_update.Category,
		payload_update.Price,
		payload_update.Expired,
		ctx,
	); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to update the data of products!", err.Error())
		return
	}

	product_response := types.ProductResponse{
		Id: product.Id,
		Name: product.Name,
		Stock: product.Stock,
		Price: product.Price,
		Category: product.Category,
		Expired: product.Category,
		Created_at: products.Created_at.Format("2006-01-02"),
		Updated_at: products.Updated_at.Format("2006-01-02"),
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to update the products!", product_response)

}

