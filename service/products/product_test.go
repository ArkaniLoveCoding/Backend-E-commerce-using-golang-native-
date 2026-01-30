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


func TestProductsCreate(t *testing.T) {

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
		if err != nil {
			t.Errorf("Failed to make new http request!")
		}
		
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

type mockStore struct {
	DeleteProductFn func(id uuid.UUID, ctx context.Context) error 
	CreateProductFn func(ctx context.Context, product *types.Products) error
	GetOneProductFn func(id uuid.UUID) (*types.Products, error)
	GetAllProductFn func() ([]types.Products, error)
	UpdateProductFn func(
		id uuid.UUID,
		ctx context.Context,
		name string,
		price string,
		stock int, 
		category string,
		expired string, 
	) error 
}

func (m *mockStore) GetAllProduct() ([]types.Products, error) {
	
	return m.GetAllProductFn()

}

func (m *mockStore) GetProductByID(id uuid.UUID) (*types.Products, error) {

	return m.GetOneProductFn(id)

}

func (m *mockStore) CreateNewProduct(ctx context.Context, products *types.Products) error {

	return m.CreateProductFn(ctx, products)

}

func (m *mockStore) DeleteProductsOnlyAdmin(id uuid.UUID, ctx context.Context) error {

	return m.DeleteProductFn(id, ctx)

}

func (m *mockStore) UpdateProductsOnlyAdmin(
	id uuid.UUID,
	name string,
	stock int,
	category string,
	price string,
	expired string,
	ctx_update context.Context,
) error {

	return m.UpdateProductFn(id, ctx_update, name, price, stock, category, expired)

}
