package main

import (
	"net/http"
)

// type router interface {
// 	http.Handler
//
// 	get(res http.ResponseWriter, req *http.Request)
// 	put(res http.ResponseWriter, req *http.Request)
// 	post(res http.ResponseWriter, req *http.Request)
// 	delete(res http.ResponseWriter, req *http.Request)
// 	patch(res http.ResponseWriter, req *http.Request)
// 	head(res http.ResponseWriter, req *http.Request)
// 	options(res http.ResponseWriter, req *http.Request)
// }

type router struct {
	route
	children map[string]*router
}

var _ http.Handler = (*router)(nil)

func (r *router) get(e endpoint)     { r.route.get = e }
func (r *router) put(e endpoint)     { r.route.put = e }
func (r *router) post(e endpoint)    { r.route.post = e }
func (r *router) delete(e endpoint)  { r.route.delete = e }
func (r *router) patch(e endpoint)   { r.route.patch = e }
func (r *router) head(e endpoint)    { r.route.head = e }
func (r *router) options(e endpoint) { r.route.options = e }

func (r *router) subroute(path string) *router {
	child := &router{}
	r.children[path] = child
	return child
}

// ............................... //

type route struct {
	get, put, post, delete, patch, head, options endpoint
}

type endpoint func(http.ResponseWriter, *http.Request)

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.RequestURI()
	if path != "/" {
		subrouter := r.children[path]
		if subrouter == nil {
			// Endpoint does not Exist; Return
			return
		}
		subrouter.ServeHTTP(res, req)
		return
	}

	// Request for Root Endpoint
	var e endpoint
	switch req.Method {
	case "GET":
		e = r.route.get
	case "PUT":
		e = r.route.put
	case "POST":
		e = r.route.post
	case "DELETE":
		e = r.route.delete
	case "PATCH":
		e = r.route.patch
	case "HEAD":
		e = r.route.head
	case "OPTIONS":
		e = r.route.options
	default:
		e = nil
	}

	if e != nil {
		e(res, req)
	}
}

func routerMain() {
	r := &router{}

	r.get(func(res http.ResponseWriter, req *http.Request) {})
	r.post(func(res http.ResponseWriter, req *http.Request) {})

	r.subroute("/users").get(func(res http.ResponseWriter, req *http.Request) {})

	http.ListenAndServe(":8080", r)
}
