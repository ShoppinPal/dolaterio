package runner

import (
	"errors"

	"github.com/shoppinpal/dolaterio/db"
	"github.com/shoppinpal/dolaterio/docker"
	"github.com/shoppinpal/dolaterio/queue"
)

// JobRunner models a job runner
type JobRunner struct {
	engine       *docker.Engine
	queue        queue.Queue
	dbConnection *db.Connection
	jobs         chan *db.Job
	stop         chan bool
	concurrency  int
	stopped      bool

	Errors chan error
}

// JobRunnerOptions models the data required to initialize a JobRunner
type JobRunnerOptions struct {
	DbConnection *db.Connection
	Engine       *docker.Engine
	Concurrency  int
	Queue        queue.Queue
}

var (
	errJobRunnerStopped = errors.New("this runner is stopped")
)

// NewJobRunner build and initializes a runner
func NewJobRunner(options *JobRunnerOptions) *JobRunner {
	return &JobRunner{
		engine:       options.Engine,
		concurrency:  options.Concurrency,
		queue:        options.Queue,
		dbConnection: options.DbConnection,
		jobs:         make(chan *db.Job),
		stop:         make(chan bool),
		Errors:       make(chan error),
	}
}

// Start starts the runner, so will start consuming and processing tasks
func (runner *JobRunner) Start() {
	for i := 0; i < runner.concurrency; i++ {
		go runner.run()
	}

	go func() {
		var message *queue.Message
		var job *db.Job
		var err error
		cont := true

		for cont {
			message, _ = runner.queue.Dequeue()
			if message != nil {
				job, err = db.GetJob(runner.dbConnection, message.JobID)
				if err != nil {
					runner.Errors <- err
				} else {
					if job != nil {
						runner.jobs <- job
					}
				}
			} else {
				cont = false
			}
		}
	}()
}

// Stop stops the job runner
func (runner *JobRunner) Stop() {
	runner.queue.Close()
	runner.stopped = true
	for i := 0; i < runner.concurrency; i++ {
		runner.stop <- true
	}
	close(runner.jobs)
	close(runner.Errors)
}

func (runner *JobRunner) run() {
	var err error
	var job *db.Job

	for {
		select {
		case job = <-runner.jobs:
			job.Status = db.StatusQueued
			job.Update(runner.dbConnection)
			err = Run(job, runner.engine)
			if err != nil {
				job.Syserr = err.Error()
			}
			job.Status = db.StatusFinished
			job.Update(runner.dbConnection)
		case <-runner.stop:
			return
		}
	}
}
