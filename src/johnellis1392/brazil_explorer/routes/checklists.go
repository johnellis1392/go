package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func listChecklists(res http.ResponseWriter, req *http.Request) {}

func describeChecklist(res http.ResponseWriter, req *http.Request) {}

func createChecklist(res http.ResponseWriter, req *http.Request) {}

func updateChecklist(res http.ResponseWriter, req *http.Request) {}

func deleteChecklist(res http.ResponseWriter, req *http.Request) {}

func newChecklistRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}
