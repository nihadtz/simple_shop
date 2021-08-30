FROM golang:1.17-alpine

WORKDIR /api

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT go run server.go --runas ${RUNAS} --migrate