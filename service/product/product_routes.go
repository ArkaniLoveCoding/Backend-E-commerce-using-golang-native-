package product

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
)


type HandleRequest struct {
	db types.ProductStore
}

func NewHandlerProduct (db types.ProductStore) *HandleRequest {
	return &HandleRequest{db: db}
}

func (h *HandleRequest) GetALlHandlerFunc(router *mux.Router) {
	router.HandleFunc("/products", h.GetAllProducts).Methods("GET")
}

func (h *HandleRequest) GetAllProducts(w http.ResponseWriter, r *http.Request) {

}