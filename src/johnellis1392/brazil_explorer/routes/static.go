package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getStaticAsset(res http.ResponseWriter, req *http.Request) {
	// ...
}

// NewStaticAssetRouter returns a new Router for Static Assets
func NewStaticAssetRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", getStaticAsset)
	return r
}
