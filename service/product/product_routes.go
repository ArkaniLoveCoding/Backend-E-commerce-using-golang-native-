package product

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/middleware"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/utils"
)


type HandleRequest struct {
	db types.ProductStore
}

func NewHandlerProduct (db types.ProductStore) *HandleRequest {
	return &HandleRequest{db: db}
}

func (h *HandleRequest) GetAllHandlerFunc(router *mux.Router) {
	router.Use(middleware.AuthenticateProfile)
	router.HandleFunc("/products", h.GetAllProductHandler).Methods("GET")
}

func (h *HandleRequest) CreateNewProductFunc(router *mux.Router) {
	router.Use(middleware.AuthenticateProfile)
	router.HandleFunc("/products", h.CreateProductHandler).Methods("POST")
}

func (h *HandleRequest) GetProductByIDFunc(router *mux.Router) {
	router.Use(middleware.AuthenticateProfile)
	router.HandleFunc("/products/:id", h.GetProductByIDHandler).Methods("POST")
}

func (h *HandleRequest) GetAllProductHandler(w http.ResponseWriter, r *http.Request) {
	
}

func (h *HandleRequest) CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	role_user, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get role for the user from token!", err.Error())
		return 
	}

	if role_user != "ADMIN" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to access this method, because your role is not as a admin!", false)
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

func (h *HandleRequest) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {

}