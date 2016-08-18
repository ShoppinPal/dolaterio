package main

import (
	"net/http"

	"github.com/shoppinpal/dolaterio/db"
	"github.com/shoppinpal/dolaterio/docker"
	"github.com/shoppinpal/dolaterio/queue"
	"github.com/gorilla/mux"
)

type apiHandler struct {
	dbConnection *db.Connection
	q            queue.Queue
	engine       *docker.Engine
}

func (api *apiHandler) rootHandler() http.Handler {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	jobs := v1.PathPrefix("/jobs").Subrouter()
	jobs.Methods("POST").HandlerFunc(api.jobsCreateHandler)
	jobs.Methods("GET").Path("/{id}").HandlerFunc(api.jobsIndexHandler)
	jobs.Methods("GET").Path("/worker/{id}").HandlerFunc(api.jobsGetAllHandler)

	workers := v1.PathPrefix("/workers").Subrouter()
	workers.Methods("POST").HandlerFunc(api.workersCreateHandler)
	workers.Methods("GET").Path("/").HandlerFunc(api.workersGetAllHandler)
	workers.Methods("GET").Path("/{id}").HandlerFunc(api.workersIndexHandler)


	return r
}
