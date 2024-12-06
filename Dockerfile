FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY app ./app
COPY internal ./internal
COPY userinterface ./userinterface
COPY bin ./bin
COPY Makefile ./
COPY sql ./sql
COPY .env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /application

FROM build-stage AS run-test-stage
RUN go test -v ./...
RUN make migration-up

# Deploy the application binary into a lean image
FROM build-stage AS build-release-stage

# RUN make migration-up