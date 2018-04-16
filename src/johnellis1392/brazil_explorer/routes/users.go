package routes

import (
	"github.com/gorilla/mux"
)

func newUserRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}
