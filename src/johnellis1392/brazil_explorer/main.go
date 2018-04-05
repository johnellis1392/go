package main

import (
  "fmt"
  "net/http"
  "log"
  "github.com/gorilla/mux"
)


func init() {}


func BaseEndpoint(res http.ResponseWriter, req *http.Request) {
  fmt.Println("Received Request in BaseEndpoint")
  fmt.Fprintf(res, "Hello, World! From BaseEndpoint\n")
}

func NewRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)

  // Workspace Router
  _ = NewWorkspaceRouter(router.PathPrefix("/workspaces"))

  // Base Handlers
  router.HandleFunc("/", BaseEndpoint)

  return router
}


func main() {
  config := NewEnvConfig()
  addr := config.AddressString()
  fmt.Printf("Listening on '%s'...\n", addr)

  router := NewRouter()
  log.Fatal(http.ListenAndServe(addr, router))
}
