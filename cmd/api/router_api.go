package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/middleware"
	serviceProduct "github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/service/products"
	serviceUser "github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/service/users"
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

	userStore := serviceUser.NewStore(s.db)
	userService := serviceUser.NewHandlerUser(userStore)

	userStores := serviceUser.NewStore(s.db)
	userServices := serviceUser.NewHandlerUserForAuthenticate(userStores)

	productStore := serviceProduct.NewStoreProduct(s.db)
	productServices := serviceProduct.NewHandlerProduct(productStore)

	// register for user

	subRouter.Handle(
		"/registration",
		http.HandlerFunc(
			userService.RegistrationFunc,
		),
	).Methods("POST")

	// login for user

	subRouter.Handle(
		"/login",
		http.HandlerFunc(
			userService.LoginFunc,
		),
	).Methods("POST")

	// update data user 
	subRouter.Handle(
		"/users/{id}",
		middleware.AuthenticateForIdUser(http.HandlerFunc(
			userService.UpdateUser,
		)),
	).Methods("PUT")
	// is authenticate

	subRouter.Handle("/profile", middleware.AuthenticateForIdUser(http.HandlerFunc(
		userServices.GetProfileUser,
	))).Methods("GET")

	// delete users 
	subRouter.Handle(
		"/users/{id}",
		middleware.AuthenticateForRole(http.HandlerFunc(
			userService.DeleteUser,
		)),
	).Methods("DELETE")

	// get all users
	subRouter.Handle(
		"/users",
		middleware.AuthenticateForRole(http.HandlerFunc(
			userService.GetAllUser,
		)),
	).Methods("GET")

	// get one users
	subRouter.Handle(
		"/users/{id}",
		http.HandlerFunc(
			userService.GetOneUsersById,
		),
	).Methods("GET")

	// products router (this router is compeletly authenticate for a several methods)

	// create new product
	subRouter.Handle(
		"/products",
		middleware.AuthenticateForRole(http.HandlerFunc(
			productServices.CreateProductHandler,
		)),
	).Methods("POST")
	
	// create products for testing the function of router (this is not authenticate!)
	subRouter.Handle(
		"/products/test",
		http.HandlerFunc(
			productServices.CreateNewProductTesting,
		),
	).Methods("POST")

	// get all product
	subRouter.Handle(
		"/products",
		http.HandlerFunc(
			productServices.GetAllProduct,
		),
	).Methods("GET")

	// get one product by id
	subRouter.Handle(
		"/products/{id}",
		http.HandlerFunc(
			productServices.GetProductByIDHandler,
		),
	).Methods("GET")

	// delete product 
	subRouter.Handle(
		"/products/{id}",
		middleware.AuthenticateForRole(http.HandlerFunc(
			productServices.DeleteProduct,
		)),
	).Methods("DELETE")

	// update product (admin only)
	subRouter.Handle(
		"/products/{id}",
		middleware.AuthenticateForRole(http.HandlerFunc(
			productServices.UpdateProductsOnlyAdmin,
		)),
	).Methods("PUT")

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
