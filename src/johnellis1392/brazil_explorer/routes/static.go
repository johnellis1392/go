package routes

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	assetDir = "web"
)

func sanitize(path string) string {
	ps := strings.Split(path, "/")
	if len(ps) == 0 {
		return ""
	}

	var res string
	res += ps[0]
	for _, p := range ps[1:] {
		switch p {
		case "", ".", "..":
			continue
		default:
			res += p
		}
	}

	return res
}

func contentTypeFor(asset string) string {
	args := strings.Split(asset, ".")
	if len(args) == 0 {
		return ""
	}

	switch args[len(args)-1] {
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	case "html":
		return "text/html"
	default:
		return "text/plain"
	}
}

// TODO: Hnadle Errors
func getStaticAsset(res http.ResponseWriter, req *http.Request) {
	vs := mux.Vars(req)
	asset := assetDir + "/" + sanitize(vs["asset"])
	// f, err := os.Open(asset)
	// if err != nil {
	// 	panic(err)
	// }

	// res.WriteHeader(http.StatusOK)
	// res.Header().Add("Content-Type", contentTypeFor(asset))

	http.ServeFile(res, req, asset)

	// http.ServeFile(res, req, r.URL.Path)
}

// NewStaticAssetRouter returns a new Router for Static Assets
func NewStaticAssetRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/:asset", getStaticAsset)
	return r
}
