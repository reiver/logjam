FROM golang:1.20 AS build-stage

WORKDIR /backend
COPY . .

RUN GORACE='halt_on_error=1' go test -race -v ./...
RUN CGO_ENABLED=0 GOOS=linux  go build -o ./logjam

FROM ubuntu:focal

RUN apt update && apt install wget curl htop nano tar xz-utils unzip gzip net-tools netcat -y

WORKDIR /backend
COPY --from=build-stage /backend/logjam ./logjam

EXPOSE 8080

ENTRYPOINT ["./logjam"]