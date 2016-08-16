package runner

import (
	"testing"
	"time"

	"github.com/shoppinpal/dolaterio/db"
	"github.com/shoppinpal/dolaterio/docker"
	"github.com/stretchr/testify/assert"
)

func TestRunEcho(t *testing.T) {
	setup()
	defer clean()

	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"echo", "hello world"},
		},
	}
	engine, err := docker.NewEngine()
	assert.Nil(t, err)

	err = Run(job, engine)
	assert.Nil(t, err)
	assert.Equal(t, "hello world\n", string(job.Stdout))
	assert.Equal(t, "", string(job.Stderr))
}

func TestRunEnv(t *testing.T) {
	setup()
	defer clean()
	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"env"},
		},
		Env: map[string]string{"K1": "V1", "K2": "V2"},
	}
	engine, err := docker.NewEngine()
	assert.Nil(t, err)

	err = Run(job, engine)
	assert.Nil(t, err)
	assert.Contains(t, job.Stdout, "K1=V1")
	assert.Contains(t, job.Stdout, "K2=V2")
}

func TestRunStdin(t *testing.T) {
	setup()
	defer clean()
	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"cat"},
		},

		Stdin: "hello world\n",
	}
	engine, err := docker.NewEngine()
	assert.Nil(t, err)

	err = Run(job, engine)
	assert.Nil(t, err)
	assert.Equal(t, "hello world\n", job.Stdout)
}

func TestRunStderr(t *testing.T) {
	setup()
	defer clean()
	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"bash", "-c", "echo hello world >&2"},
		},
	}
	engine, err := docker.NewEngine()
	assert.Nil(t, err)

	err = Run(job, engine)
	assert.Nil(t, err)
	assert.Equal(t, "hello world\n", job.Stderr)
}

func TestRunTimeout(t *testing.T) {
	setup()
	defer clean()
	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"sleep", "2000"},
			Timeout:     1 * time.Millisecond,
		},
	}
	engine, err := docker.NewEngine()
	assert.Nil(t, err)

	err = Run(job, engine)
	assert.Equal(t, errTimeout.Error(), err.Error())
}

func TestRequiresAValidDockerImage(t *testing.T) {
	setup()
	defer clean()

	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "dolaterio/yolo",
			Cmd:         []string{"echo", "hello world"},
		},
	}
	engine, err := docker.NewEngine()
	assert.Nil(t, err)

	err = Run(job, engine)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid docker image", string(err.Error()))
	assert.Equal(t, "", string(job.Stderr))
}
