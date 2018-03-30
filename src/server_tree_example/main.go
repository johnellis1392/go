package main

import (
	"fmt"
	"log"
	"net/http"
	"server_tree_example/checklists"
	"server_tree_example/users"

	"github.com/gorilla/mux"
)

const (
	addressString = "localhost:8080"
)

// Result - Result
type Result struct {
	Message string
}

// The following is an experiment to determine some of the finer points
// of how function arguments propagate in golang.
//
// func errorProducer() (*Result, error) {
// 	return nil, fmt.Errorf("Example Error")
// }
//
// func errorHandler(result *Result, err error) {
// 	if result != nil {
// 		fmt.Println("Result is defined")
// 	} else {
// 		fmt.Println("Result is NOT defined")
// 	}
//
// 	if err != nil {
// 		fmt.Println("Received Error:", err.Error())
// 	} else {
// 		fmt.Println("No Error to process")
// 	}
// }
//
// func main() {
// 	fmt.Println("Running experiment...")
// 	errorHandler(errorProducer())
// }

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// Setup Endpoints
	router.Handle("/users", users.NewRouter())
	router.Handle("/checklists", checklists.NewRouter())

	fmt.Println("Listening on address", addressString, "...")
	log.Fatal(http.ListenAndServe(addressString, router))
}
