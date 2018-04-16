package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	DEFAULT_WORKSPACE = "/workplace"
)

var (
	WORKSPACE string
)

type Workspace struct {
	// Id string `json:"workspace_id"`
	// Name string `json:"name"`
	// PlatformOverride string `json:"platform_override"`

	// VersionSet VersionSet `json:"versionset"`
	// Packages []Package `json:"packages"`

	Base BaseInfo `parser:"'base' '=' '{' @@ '}' ';'"`

	PlatformOverride string `parser:"'platformOverride' '=' @Ident ';'"`

	Packages []Package `parser:"'packages' '=' '{' @@ '}' ';'"`
}

type BaseInfo struct {
	Workspace  string `parser:"'workspace' '=' @Ident ';'"`
	VersionSet string `parser:"'versionSet' '=' @Ident ';'"`
}

// type VersionSet struct {
//   // Id string `json:"versionset_id"`
//   // Group string `json:"group_name"`
//
//   Id string `parser:""`
//   Group string `parser:""`
// }
//
// func (v VersionSet) Name() string {
//   return fmt.Sprintf("%s/%s", v.Group, v.Id)
// }

type Package struct {
	// Id string `json:"package_id"`
	// Version string `json:"version"`

	// Id string `parser:""`
	// Version string `parser:""`
	// Location string `parser:""`
	ID      string `parser:"@Ident '-' @String '=' @String ';'"`
	Version string
}

func (p Package) Name() string {
	return fmt.Sprintf("%s-%s", p.ID, p.Version)
}

func init() {
	// Setup Configuration Variables
	// WORKSPACE = getenvOrElse("WORKSPACE", DEFAULT_WORKSPACE)
	WORKSPACE = ""
}

func ListWorkspaces(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Received ListWorkspaces request...")
	res.Header().Add("Content-Type", "application/json; charset=UTF-8")

	// Read Child Workspaces in Workspace Root
	files, err := ioutil.ReadDir(WORKSPACE)
	if err != nil {
		// TODO: Handle Error
		panic(err)
	}

	// workspaces := make([]string)
	var workspaces []string
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		filename := file.Name()
		workspaces = append(workspaces, filename)
	}

	fmt.Printf("Found Workspaces: %v\n", workspaces)
	fmt.Printf("Marshalling to Json...\n")

	err = json.NewEncoder(res).Encode(workspaces)
	if err != nil {
		// TODO: Handle Error
		panic(err)
	}

	fmt.Printf("ListWorkspaces Success\n")
}

func DescribeWorkspace(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	workspaceId := vars["workspaceId"]
	fmt.Fprintf(res, "DescribeWorkspace: %s\n", workspaceId)
}

func NewWorkspaceRouter(r *mux.Route) *mux.Router {
	router := r.Subrouter()

	router.HandleFunc("/", ListWorkspaces).Methods("GET")
	router.HandleFunc("/{workspaceId}", DescribeWorkspace).Methods("GET")

	return router
}
