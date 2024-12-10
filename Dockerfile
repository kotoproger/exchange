FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY app ./app
COPY internal ./internal
COPY userinterface ./userinterface
COPY Makefile ./
COPY sql ./sql
COPY .env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /application

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /application ./application
COPY --from=build-stage /app/sql/migrations /sql/migrations
COPY --from=build-stage /app/.env ./.env
# RUN make migration-up