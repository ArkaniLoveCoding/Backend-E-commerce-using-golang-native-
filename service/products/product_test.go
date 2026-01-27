package product

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
)

type mockStore struct {}

func TestProducts(t *testing.T) {

	userStore := &mockStore{}
	Handler := NewHandlerProduct(userStore)

	t.Run("This is the products create test!", func(t *testing.T) {

		payload := types.Products{
			Name: "Indomie",
			Stock: 10,
			Price: "4.000 Rp",
			Expired: "10 Januari 2028",
			Category: "makanan instan",
		}

		encoding, err := json.Marshal(&payload)
		if err != nil {
			t.Errorf("Error for payloading json! %v", err.Error())
			return
		}

		req, err := http.NewRequest(
			http.MethodPost,
			"/products",
			bytes.NewBuffer(encoding),
		)
		
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products/test", Handler.CreateProductHandler)
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusBadRequest {
			t.Errorf("Something went wrong! %v", err.Error())
			return
		}

	})


}

func (m *mockStore) GetAllProduct() ([]types.Products, error) {
	
	return nil, nil

}

func (m *mockStore) GetProductByID(id uuid.UUID) (*types.Products, error) {

	return nil, nil

}

func (m *mockStore) CreateNewProduct(ctx context.Context, products *types.Products) error {

	return nil

}
