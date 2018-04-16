package routes

import (
	"github.com/gorilla/mux"
)

func newChecklistRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}
