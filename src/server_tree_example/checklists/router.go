package checklists

import (
	"server_tree_example/checklists/items"

	"github.com/gorilla/mux"
)

// ChecklistRouter - Router for Checklists
// type ChecklistRouter struct{}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/checklists/items", items.NewRouter())

	return router
}
