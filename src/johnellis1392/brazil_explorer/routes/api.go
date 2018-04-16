package routes

import (
	"github.com/gorilla/mux"
)

// NewAPIRouter returns a new API Router
func NewAPIRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/users", newUserRouter())
	r.Handle("/checklists", newChecklistRouter())
	r.Handle("/checklists/items", newChecklistItemRouter())
	return r
}
