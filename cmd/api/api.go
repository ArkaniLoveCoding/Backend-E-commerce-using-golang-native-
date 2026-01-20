package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/middleware"
	service "github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/service/user"
)

type ApiServer struct {
	Addr   string
	db     *sqlx.DB
	server *http.Server
}

func ApiServerAddr(addr string, db *sqlx.DB) *ApiServer {
	return &ApiServer{
		Addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// testing if the server is working!
	subRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "Successfully to testing the web server, now the web server is working!",
			"data": "Hello world!"
		}`))
	})

	// not authenticate

	userStore := service.NewStore(s.db)
	userService := service.NewHandlerUser(userStore)
	userService.RegistrationUserHandler(subRouter)
	userService.LoginUserHandler(subRouter)

	// products router

	// is authenticate 

	userStores := service.NewStore(s.db)
	userServices := service.NewHandlerUserForAuthenticate(userStores)
	subRouter.Handle("/profile", middleware.AuthenticateProfile(http.HandlerFunc(
		userServices.GetProfileUser,
	))).Methods("GET")

	// Create HTTP server
	s.server = &http.Server{
		Addr:   s.Addr,
		Handler: router,
	}

	log.Printf("Server starting on %s", s.Addr)

	if err := s.server.ListenAndServe(); err != nil {
		return errors.New(err.Error())
	}
	
	return nil
}

// Shutdown gracefully shuts down the server
func (s *ApiServer) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
