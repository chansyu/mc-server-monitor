# syntax=docker/dockerfile:1
# https://stackoverflow.com/questions/40873165/use-docker-run-command-to-pass-arguments-to-cmd-in-dockerfile

FROM golang:1.21 AS build-stage

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY ui ./ui

RUN CGO_ENABLED=0 GOOS=linux go build -o /mc-server-monitor ./cmd/web 

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /mc-server-monitor /mc-server-monitor

EXPOSE 4000

USER nonroot:nonroot

CMD ["/mc-server-monitor"]