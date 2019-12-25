FROM golang:1.13.4 AS build-env
ARG binary_name
ARG main_pkg
WORKDIR /service/
RUN git config --global http.postBuffer 524288000
COPY go.mod go.sum /service/
RUN go mod download
COPY cmd /service/cmd
COPY pkg /service/pkg
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o build/bin/linux.amd64/${binary_name} ${main_pkg}

FROM alpine
RUN apk --no-cache add \
    ca-certificates
WORKDIR /service/
COPY --from=build-env /service/build/bin/linux.amd64/${binary_name} /service/${binary_name}
EXPOSE 8080
ENTRYPOINT ["./kube-job-runner"]
