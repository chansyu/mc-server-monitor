# syntax=docker/dockerfile:1
# docker build -t mc-server-monitor:multistage .

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

EXPOSE 8080

USER nonroot:nonroot

CMD ["/mc-server-monitor"]