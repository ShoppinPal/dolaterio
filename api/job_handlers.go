package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shoppinpal/dolaterio/db"
	"github.com/shoppinpal/dolaterio/queue"
	"github.com/gorilla/mux"
)

type jobObjectRequest struct {
	WorkerName string            `json:"worker_name"`
	Stdin    json.RawMessage   `json:"stdin"`
	Env      map[string]string `json:"env"`
}

func (api *apiHandler) jobsCreateHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var jobReq jobObjectRequest
	decoder.Decode(&jobReq)

	job := &db.Job{
		Stdin:    string(jobReq.Stdin),
		WorkerName: jobReq.WorkerName,
		Env:      jobReq.Env,
	}
	err := job.Store(api.dbConnection)
	api.q.Enqueue(&queue.Message{
		JobID: job.ID,
	})
	if err != nil {
		renderError(res, err, 500)
		return
	}
	res.WriteHeader(201)
	api.renderJob(res, job)
}

func (api *apiHandler) jobsIndexHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job, err := db.GetJob(api.dbConnection, vars["id"])

	if err != nil {
		renderError(res, err, 500)
		return
	}
	if job == nil {
		renderError(res, errors.New("Job not found"), 404)
		return
	}
	res.WriteHeader(200)
	api.renderJob(res, job)
}

func (api *apiHandler) jobsGetAllHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobs, err := db.GetAllJobs(api.dbConnection, vars["id"])

	if err != nil {
		renderError(res, err, 500)
		return
	}
	if jobs == nil {
		renderError(res, errors.New("Job not found"), 404)
		return
	}
	res.WriteHeader(200)
	for _, job := range jobs {
        api.renderAllJobs(res, &job)
  }
}

func (api *apiHandler) renderJob(res http.ResponseWriter, job *db.Job) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(job)
}

func (api *apiHandler) renderAllJobs(res http.ResponseWriter, job *interface{}) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(job)
}
