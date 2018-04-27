package main

import (
	"fmt"
)

// ////////////////// //
// /// Base Types /// //
// ////////////////// //

// Manifest is the type-level abstraction for the workflow's
// configuration manifest.
type Manifest struct{}

// Orchestrator orchestrates workers and distributes jobs to them to execute.
type Orchestrator interface{}

// Worker executes a job using some built-in engine.
//
// TODO:
// * Make worker type that executes a docker container
// * Make worker type that is just a function
// * Make worker type that uses processes
type Worker interface{}

// Job represents an indivisible unit of work to be executed by the orchestrator.
type Job interface{}

// //////////////////// //
// /// Docker Types /// //
// //////////////////// //

// DockerOrchestrator is an orchestrator that manages docker containers for
// executing jobs.
type DockerOrchestrator struct {
	workers []*DockerWorker
}

// DockerWorker is a worker type for representing docker containers that will
// run jobs.
type DockerWorker struct {
	// ...
}

func main() {
	fmt.Println("Starting...")

	// Use Case:
	// * create manifest containing state machine
	// * query engine to construct orchestrator from manifest
	// * write jobs for handling tasks
	// * bootstrap orchestrator with worker type
	// * execute orchestrator on input job
}
