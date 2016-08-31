# Dolater.io

Dolater.io lets you execute background jobs on a remote docker server.

# How to run it

You'll need [docker-compose](https://docs.docker.com/compose/) to run the services and its dependencies easily. To run all services at once use, run:

```
docker-compose up -d --no-recreate
```

This will run the API server as well as one dolater.io worker. You can always scale the amount of workers by using `docker-compose scale worker=N` command.

Now it's all ready to use.

# How to write Worker

You can find documentation here for [writing your own worker.](https://github.com/ShoppinPal/dolaterio/blob/master/docs/write_a_worker.md)

# Simple Example

Since dolater.io is running in docker, you'll need to know your docker host IP address to access it.
* If you use boot2docker, run `boot2docker ip` to find and substitute the value for `DOCKERHOST`
  * If that option is too outdated then another option for 'nix:

    ```
    export DOCKERHOST=`docker-machine ip default`
    ```
* If you use `Docker for Mac and Windows beta` then you can substitute `DOCKERHOST` with `localhost`.

Create a worker using our parrot docker image:

```
curl http://DOCKERHOST:7000/v1/workers -H "Content-Type: application/json" -X POST -d '{"docker_image": "dolaterio/parrot"}'
```

You'll get a JSON response back with the information of the worker you just created. Use its `id` to create jobs on it:

```
curl http://DOCKERHOST:7000/v1/jobs -H "Content-Type: application/json" -X POST -d '{"worker_id": WORKER_ID, "stdin": "Hello world!"}'
```

It will return a new JSON containing, between others, the `id` of the job. You can request dolater.io for the current state of the job:

```
curl http://DOCKERHOST:7000/v1/jobs/JOB_ID
```

Passing environment variables with worker :

```
curl http://DOCKERHOST:7000/v1/workers -H "Content-Type: application/json" -X POST -d '{"docker_image": "dolaterio/parrot", "env": {"NODE_ENV": "local"}}'
```

Passing environment variables with the job :

```
curl http://127.0.0.1:7000/v1/jobs -H "Content-Type: application/json" -X POST -d '{"worker_id": "6e1935fc-328b-40d7-9957-2f10654360f1", "stdin": "Hello World!", "env": {"HI": "BYE"}}'
```
