package main

import (
	"fmt"
	"johnellis1392/brazil_explorer/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/api", routes.NewAPIRouter())
	r.Handle("/assets", routes.NewStaticAssetRouter())

	// r.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
	// 	res.Write([]byte("base: Hello, World!"))
	// })

	return r
}

func main() {
	c := envConfig()
	addr := c.AddressString()
	r := newRouter()

	fmt.Printf("Listening on \"%s\"...\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
