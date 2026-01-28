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

type mockStore struct {}

func (m *mockStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}

func (m *mockStore) GetUserById(id uuid.UUID) (*types.User, error) {
	return nil, nil
}

func (m *mockStore) CreateUser(context.Context, *types.User) error {
	return nil
}

func (m *mockStore) GetAllUser() ([]types.User, error) {
	return nil, nil
}

func (m *mockStore) UpdateToken(
	id uuid.UUID, token string, token_refresh string,
) error {
	return nil
}

func (m *mockStore) GetUsersRole(email string) (*types.User, error) {
	return nil, nil
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
	) error {
	return nil
}

func (m *mockStore) DeleteUsersOnlyAdmin(id uuid.UUID, ctx context.Context) error {
	return nil
}