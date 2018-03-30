package users

import (
	"fmt"
	"net/http"
)

func listUsers(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello, World!")
}
