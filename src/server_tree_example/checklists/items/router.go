package items

import "github.com/gorilla/mux"

// ChecklistItemRouter - Router
// type ChecklistItemRouter struct{}

// NewRouter - Router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	return router
}
