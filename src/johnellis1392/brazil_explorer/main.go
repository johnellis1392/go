package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/api", NewAPIRouter())
	r.Handle("/assets", NewStaticAssetRouter())
	return r
}

func main() {
	c := envConfig()
	addr := c.AddressString()
	r := newRouter()

	fmt.Printf("Listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
