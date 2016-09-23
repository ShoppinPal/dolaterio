package runner

import (
	"errors"
	"time"

	"github.com/shoppinpal/dolaterio/db"
	"github.com/shoppinpal/dolaterio/docker"
)

var (
	errTimeout = errors.New("timeout")
)

// Run runs the job against in the container engine
func Run(job *db.Job, engine *docker.Engine) error {
	// valid, err := engine.ValidImage(job.Worker.DockerImage)
	// if err != nil {
	// 	return err
	// }

	// if !valid {
	// 	return errors.New("Invalid docker image")
	// }

	container, err := engine.BuildContainer(job)
	if err != nil {
		return err
	}
	defer container.Remove()

	done := make(chan int)
	errChn := make(chan error)
	go func() {
		err = container.AttachStdin()
		if err != nil {
			errChn <- err
		}

		err = container.Wait()
		if err != nil {
			errChn <- err
		} else {
			done <- 1
		}
	}()

	timeout := job.Worker.Timeout
	if timeout == 0 {
		timeout = engine.Timeout
	}

	select {
	case <-done:
	case err := <-errChn:
		return err

	case <-time.After(timeout):
		return errTimeout
	}

	err = container.FetchOutput()
	if err != nil {
		return err
	}

	job.Stdout = string(container.StdOut)
	job.Stderr = string(container.StdErr)
	return nil
}
