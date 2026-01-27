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

