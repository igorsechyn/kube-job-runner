#!/bin/bash

set -eu
set -x

GOARCH=amd64 GOOS=linux go build -o test main.go

docker build --rm --tag "igorsechyn/samplejob:1.0.0" .
docker push igorsechyn/samplejob:1.0.0
