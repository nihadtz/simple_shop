FROM golang:alpine

WORKDIR /api

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT go run server.go