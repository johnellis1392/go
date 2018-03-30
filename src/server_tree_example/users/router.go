package users

import (
	"github.com/gorilla/mux"
)

// UserRouter - Router for Users
// type UserRouter struct{}

// NewRouter - Router
func NewRouter() *mux.Router {
	// return &UserRouter{}
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").
		Path("/users").
		Name("ListUsers").
		HandlerFunc(listUsers)

	return router
}
