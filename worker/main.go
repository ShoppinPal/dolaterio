package main

import (
	"github.com/shoppinpal/dolaterio/db"
	"github.com/shoppinpal/dolaterio/docker"
	"github.com/shoppinpal/dolaterio/queue"
	"github.com/shoppinpal/dolaterio/runner"
)

func main() {
	dbConnection, err := db.NewConnection()
	if err != nil {
		panic(err)
	}
	engine, err := docker.NewEngine()
	if err != nil {
		panic(err)
	}
	queue, err := queue.NewRedisQueue()
	if err != nil {
		panic(err)
	}

	runner := runner.NewJobRunner(&runner.JobRunnerOptions{
		DbConnection: dbConnection,
		Engine:       engine,
		Queue:        queue,
		Concurrency:  8,
	})
	runner.Start()
	done := make(chan bool, 1)
	select {
	case <-done:
		runner.Stop()
	}

}
