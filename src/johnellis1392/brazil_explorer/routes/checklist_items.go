package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func listChecklistItems(res http.ResponseWriter, req *http.Request) {}

func describeChecklistItem(res http.ResponseWriter, req *http.Request) {}

func createChecklistItem(res http.ResponseWriter, req *http.Request) {}

func updateChecklistItem(res http.ResponseWriter, req *http.Request) {}

func deleteChecklistItem(res http.ResponseWriter, req *http.Request) {}

func newChecklistItemRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}
