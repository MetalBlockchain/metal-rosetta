# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.15 AS build

WORKDIR /go/src/github.com/figment-networks/avalanche-rosetta

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  make setup && make build

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.10 AS release

WORKDIR /app

COPY --from=build \
  /go/src/github.com/figment-networks/avalanche-rosetta/avalanche-rosetta \
  /app/avalanche-rosetta

EXPOSE 8080

ENTRYPOINT ["/app/avalanche-rosetta"]