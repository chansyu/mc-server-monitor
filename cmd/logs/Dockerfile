# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY ui ./ui

RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o ./mc-logs-server ./cmd/logs

# # Run the tests in the container
# # FROM build-stage AS run-test-stage
# # RUN go test -v ./...

FROM gcr.io/distroless/static
WORKDIR /app

COPY --from=build-stage /app/mc-logs-server ./mc-logs-server

USER nonroot:nonroot
ENTRYPOINT ["/app/mc-logs-server"]