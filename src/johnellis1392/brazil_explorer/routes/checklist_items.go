package routes

import (
	"github.com/gorilla/mux"
)

func newChecklistItemRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}
