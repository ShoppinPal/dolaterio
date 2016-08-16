package runner

import (
	"github.com/shoppinpal/dolaterio/db"
	"github.com/shoppinpal/dolaterio/docker"
	"github.com/shoppinpal/dolaterio/queue"
)

var (
	dbConnection *db.Connection
	engine       *docker.Engine
	q            queue.Queue
)

func setup() {
	var err error
	dbConnection, err = db.NewConnection()
	if err != nil {
		panic(err)
	}
	engine, err = docker.NewEngine()
	if err != nil {
		panic(err)
	}
	engine.SkipPull = true
	q = newFakeQueue()
}

func clean() {
	dbConnection.Close()
	q.Close()
}

func logErrors(errors chan error) {
	go func() {
		for err := range errors {
			panic(err)
		}
	}()
}
