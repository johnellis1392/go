package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func listUsers(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("listUsers: Hello, World!"))
}

func describeUesr(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("describeUser: Hello, World!"))
}

func createUser(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("createUser: Hello, World!"))
}

func updateUser(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("updateUser: Hello, World!"))
}

func deleteUser(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("deleteUser: Hello, World!"))
}

func newUserRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("GET").HandlerFunc(listUsers)
	r.Methods("GET").Path("/:userId").HandlerFunc(describeUesr)
	r.Methods("POST").HandlerFunc(createUser)
	r.Methods("PATCH", "PUT").Path("/:userId").HandlerFunc(updateUser)
	r.Methods("DELETE").Path("/:userId").HandlerFunc(deleteUser)
	return r
}
