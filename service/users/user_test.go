package service

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


func TestUser(t *testing.T) {

	newStore := &mockStore{}
	handler := NewHandlerUser(newStore)

	t.Run("this testing should be return false", func(t *testing.T) {

	payload := types.Register{
		Firstname: "Lalu Ahmad Arkani",
		Lastname: "Arkani",
		Password: "arkan123",
		Email: "laluahmadarkani@gmail.com",
		Country: "United States Of America",
		Address: "New York City",
		Role: "USER",
	}

	encoding, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("Failed to encoding the json payload!, => %s", err.Error())
		return
	}

    req, err := http.NewRequest(

		http.MethodPost, 
		"/registration",
		bytes.NewBuffer(encoding),

	)
	if err != nil {
		t.Errorf("Failed to make new request for http payload json!, => %s", err.Error())
		return
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()

	router.HandleFunc("/registration", handler.RegistrationFunc)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("The status must be a bad request!")
		return
	}
	})

}

type mockStore struct {
	GetUserByEmailFn func(email string) (*types.User, error)
	GetUserByIdFn func(id uuid.UUID) (*types.User, error)
	CreateUserFn func(ctx context.Context, user *types.User) error
	GetUserByRoleFn func(role string) (*types.User, error)
	UpdateUserFn func(
		id uuid.UUID,
		ctx context.Context,
		firstname string,
		lastname string,
		password string,
		email string,
		country string,
		address string,
	) error
	DeleteUsersOnlyAdminFn func(id uuid.UUID, ctx context.Context) error
	UpdateTokenFn func(
		ctx context.Context, id uuid.UUID, token string, refresh_token string, user *types.User) error
	GetAllUserFn func() ([]types.User, error)
}

func (m *mockStore) GetUserByEmail(email string) (*types.User, error) {
	
	return m.GetUserByEmailFn(email)

}

func (m *mockStore) GetUserById(id uuid.UUID) (*types.User, error) {
	
	return m.GetUserByIdFn(id)

}

func (m *mockStore) CreateUser(ctx context.Context, user *types.User) error {
	
	return m.CreateUserFn(ctx, user)

}

func (m *mockStore) GetAllUser() ([]types.User, error) {
	
	return m.GetAllUserFn()

}

func (m *mockStore) UpdateToken(
	ctx context.Context, id uuid.UUID, token string, token_refresh string, user *types.User,
) error {
	
	return m.UpdateTokenFn(ctx, id, token, token_refresh, user)

}

func (m *mockStore) GetUsersRole(role string) (*types.User, error) {
	
	return m.GetUserByRoleFn(role)

}

func (m *mockStore) UpdateDataUser(
	id uuid.UUID,
	ctx context.Context,
	firstname string,
	lastname string,
	password string,
	email string,
	country string,
	address string,
	user *types.User, 
	) error {
	
		return m.UpdateUserFn(id, ctx, firstname, lastname, password, email, country, address)

}

func (m *mockStore) DeleteUsersOnlyAdmin(id uuid.UUID, ctx context.Context) error {
	
	return m.DeleteUsersOnlyAdminFn(id, ctx)

}