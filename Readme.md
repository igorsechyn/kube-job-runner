# Kubernetes job runner

A spike to implement kubernetes job runner

## Development dependencies

- go 1.12.9
- docker 18.06.1-ce (multistage build is required)
- minikube 0.26 for local development
- skaffold (https://github.com/GoogleContainerTools/skaffold)

## Setting up a development machine

1. Start minikube

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

4. Run local database migrations, after the service is actually running (when deploy is finished, it does not mean it is running)

    ```
    make run-local-migrations
    ```

5. Ensure your environment is working by running the pre-commit check

    ```
    make test
    ```

## During development

Commits to this codebase should follow the [conventional changelog conventions](https://github.com/bcoe/conventional-changelog-standard/blob/master/convention.md).

- `make test` - A check to be run before pushing any changes
- `make build` - builds a server executable under `build` directory
- `make docker` - builds the server and creates a docker image for release
- `make docker-publish` - publishes docker image
- `make deploy-local` - executes `skaffold run` to deploy runner to minikube (including postgres and elasticmq)
- `make watch` - runs unit tests on every change
- `skaffold dev` - will redeploy application on each change to minikube