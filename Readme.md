[![Build Status](https://travis-ci.org/igorsechyn/kube-job-runner.svg?branch=master)](https://travis-ci.org/igorsechyn/kube-job-runner)

# Kubernetes job runner

A spike to implement a simple control plane on top of kubernetes Job object. This is a validation of a design and not meant to be used in production.

## Design

Application consists of several components (pkg/runner), which run separately from each other on their own compute:

- `web_server` - HTTP API to submit job requests and get the execution status. On posting a job request, details are stored in data store and a queue message is sent to start processing
- `messages_processor` - a worker to poll queue and process messages to create, update and delet jobs
- `job_status_processor` - a worker to watch Job objects using k8s informer and send messges to the queue to update job statuses
- `pod_events_processor` - a worker to watche Pod objects events and send a message to mark job as failed, if Pod failed to start. Pods can fail for different reasons (e.g. image name is wrong), which will not cause the Job status to update. The idea is to handle scenarios where the reason is not recoverable and fail fast instead waiting for a timeout

Job details are stored in a datastore (curently postgres) and all updates are triggered through through a message to a queue (currently implemented with SQS). 

## Development dependencies

- go 1.12.9
- docker 18.06.1-ce (multistage build is required)
- minikube 1.6.2 for local development (https://kubernetes.io/docs/tasks/tools/install-minikube/)
- skaffold (https://github.com/GoogleContainerTools/skaffold)

## Setting up a development machine

1. Start minikube

On OSX use hyperkit (see https://minikube.sigs.k8s.io/docs/start/macos/)

```
minikube start --vm-driver=hyperkit --kubernetes-version=1.17.0
```

2. Install project dependencies

    ```
    make install
    ```

    If your IDE does not support go modules or the project is inside GOPATH, run `go mod vendor` to put dependencies into a vendor folder


3. Deploy runner with postgres and sqs to local minikube. Under the hood skaffold builds the image and applies all files under `minikube` directory.
See `skaffold.yaml` for configuration

    ```
    make deploy-local
    ```

5. Ensure your environment is working by running the pre-commit check

    ```
    make test
    ```

## Migrations

DB migrations are executed automatically as an `initContainer` of the `web-server` in `./minikube-config/deployment.yaml`. For local dev loop and CI two additional init container are added to wait for sqs and pg being available.

## During development

Commits to this codebase should follow the [conventional changelog conventions](https://github.com/bcoe/conventional-changelog-standard/blob/master/convention.md).

- `make test` - A check to be run before pushing any changes
- `make build` - builds a server executable under `build` directory
- `make docker` - builds the server and creates a docker image for release
- `make docker-publish` - publishes docker image
- `make deploy-local` - executes `skaffold run` to deploy runner to minikube (including postgres and elasticmq)
- `make watch` - runs unit tests on every change
- `skaffold dev` - will redeploy application on each change to minikube

## CI

Travis CI - https://travis-ci.org/igorsechyn/kube-job-runner