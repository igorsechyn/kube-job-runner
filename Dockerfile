FROM golang:1.12.9 AS build-env
ARG binary_name
ARG main_pkg
WORKDIR /service/
ADD go.mod /service/
ADD go.sum /service/
RUN go mod download
COPY . /service/
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o build/bin/linux.amd64/${binary_name} ${main_pkg}

FROM alpine
RUN apk --no-cache add \
    ca-certificates
WORKDIR /service/
COPY --from=build-env /service/build/bin/linux.amd64/${binary_name} /service/${binary_name}
EXPOSE 8080
ENTRYPOINT ["./kube-job-runner"]
