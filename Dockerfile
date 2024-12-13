FROM golang:1.22.10-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY app ./app
COPY internal ./internal
COPY userinterface ./userinterface
COPY Makefile ./
COPY sql ./sql

RUN CGO_ENABLED=0 GOOS=linux go build -o /application

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /application ./application
COPY --from=build-stage /app/sql/migrations /sql/migrations
