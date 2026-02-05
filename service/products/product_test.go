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

	productStore := &mockStore{
		CreateProductFn: func(ctx context.Context, product *types.Products) error {
			return nil
		},
	}
	Handler := NewHandlerProduct(productStore)

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

func TestProductUpdate(t *testing.T) {

	productStore := &mockStore{
		UpdateProductFn: func(id uuid.UUID, ctx context.Context, name, price string, stock int, category, expired string) error {
			return nil
		},
	}
	handler := NewHandlerProduct(productStore)
	
	t.Run("TESTING UPDATE PRODUCT", func(t *testing.T) {

		id_fake := uuid.New()

		payload := types.Products{
			Id: id_fake,
			Name: "Indomie aceh",
			Price: "Rp10.0000",
			Stock: 10,
			Category: "makanan instan",
			Expired: "10 Januari 2030",
		}

		encode, err := json.Marshal(&payload)
		if err != nil {
			t.Errorf("Failed to marshall json payload!")
		}

		req, err := http.NewRequest(
			http.MethodPut,
			"/product/"+id_fake.String(),
			bytes.NewBuffer(encode),
		)
		if err != nil {
			t.Errorf("Erorr cannot be request for http new request!")
		}

		router := mux.NewRouter()
		rr := httptest.NewRecorder()
		router.HandleFunc("/product/{id}", handler.UpdateProductsOnlyAdmin)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("the response mus be return bad request")
		}

		t.Log(rr.Body.String())

	})

}

func TestProductDelete (t *testing.T) {

	productStore := &mockStore{
		DeleteProductFn: func(id uuid.UUID, ctx context.Context) error {
			return nil
		},
	}

	handler := NewHandlerProduct(productStore)

	t.Run("this testing must be return bad request!", func(t *testing.T) {

		id_fake := uuid.New()

		req, err := http.NewRequest(
			http.MethodDelete,
			"/products/"+id_fake.String(),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to make a new request for delete!")
		}

		router := mux.NewRouter()
		rr := httptest.NewRecorder()

		router.HandleFunc("/products/{id}", handler.DeleteProduct)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("The response should be bad request response!")
		}

		t.Log(rr.Body.String())

	})

}

func TestGetOneProduct (t *testing.T) {


	productStore := &mockStore{
		GetOneProductFn: func(id uuid.UUID) (*types.Products, error) {
			return nil, nil
		},
	}

	handler := NewHandlerProduct(productStore)

	t.Run("The response should be return a bad request response", func(t *testing.T) {

		id_fake := uuid.New()

		req, err := http.NewRequest(
			http.MethodGet,
			"/products/"+id_fake.String(),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to make new http request")
		}

		router := mux.NewRouter()
		rr := httptest.NewRecorder()

		router.HandleFunc("/products/{id}", handler.GetProductByIDHandler)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("The response must be a bad request response!")
		}

		t.Log(rr.Code)
		t.Log(rr.Body.String())

	})

}

func TestGetAllProduct (t *testing.T) {

	productStore := &mockStore{
		GetAllProductFn: func() ([]types.Products, error) {
			return nil, nil
		},
	}

	handler := NewHandlerProduct(productStore)

	t.Run("The returnof this response must be a bad request response!", func(t *testing.T) {

		req, err := http.NewRequest(
			http.MethodGet,
			"/products",
			nil,
		)
		if err != nil {
			t.Errorf("Failed to make new http request for this testing!")
		}

		router := mux.NewRouter()
		rr := httptest.NewRecorder()

		router.HandleFunc("/products", handler.GetAllProduct)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("The respone should be bad request response!")
		}

		t.Log(rr.Body.String())

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
	SearchManyProductsFn func(ctx context.Context, keyword string, offset int, limit int) ([]types.Products, error)
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

func (m *mockStore) SearchManyProducts(ctx context.Context, keyword string, offset int, limit int) ([]types.Products, error) {

	return m.SearchManyProductsFn(ctx, keyword, offset, limit)

}
