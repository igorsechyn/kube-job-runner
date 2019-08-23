.PHONY: build install fmt install lint test test-unit install-ci clean watch test-race test-integration release
MAIN_PKG := kube-job-runner/cmd/kube-job-runner
BINARY_NAME := kube-job-runner
GIT_HASH := $$(git rev-parse --short HEAD)
SERVICE_URL := $$(minikube service webserver --namespace default --url)
IMAGE_NAME := igorsechyn/kube-job-runner
PWD = $$(pwd)
DOCKER_DIGEST := $$(docker inspect --format='{{index .RepoDigests 0}}' $(IMAGE_NAME):$(GIT_HASH))
clean:
	rm -rf build/bin/*

install:
	go mod download

build: clean fmt
	env GOOS=darwin GOARCH=amd64 go build -o build/bin/darwin.amd64/$(BINARY_NAME) $(GOBUILD_VERSION_ARGS) $(MAIN_PKG)
	chmod +x build/bin/darwin.amd64/$(BINARY_NAME)

fmt:
	gofmt -w=true -s $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -w=true -d $$(find . -type f -name '*.go' -not -path "./vendor/*")

test-unit:
	go test -tags unit ./... -timeout 120s -count 1

test-integration:
	SERVICE_URL=$(SERVICE_URL) go test -tags integration ./... -timeout 120s -count 1

test: test-unit test-integration

docker: clean
	docker build --build-arg binary_name=$(BINARY_NAME) --build-arg main_pkg=$(MAIN_PKG) -t $(IMAGE_NAME):$(GIT_HASH) .

docker-publish: docker
	docker push $(IMAGE_NAME):$(GIT_HASH)

docker-sample-job:
	cd sample-job && ./publish-image.sh

deploy-local:
	skaffold run

run-local-migrations:
	curl -v -X POST $(SERVICE_URL)/migrate

generate-sql-scripts:
	cd pkg/app/migrations && go-bindata -prefix ../../../migrations/ -pkg migrations ../../../migrations

watch:
	CompileDaemon -color=true -exclude-dir=.git -build="make test-unit"