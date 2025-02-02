# ------------------------------------------------------------------------------
# Build metal
# ------------------------------------------------------------------------------
FROM golang:1.19.6 AS metal

ARG METAL_VERSION

RUN git clone https://github.com/MetalBlockchain/metalgo.git \
  /go/src/github.com/!metal!blockchain/metalgo

WORKDIR /go/src/github.com/!metal!blockchain/metalgo

RUN git checkout $METAL_VERSION && \
    ./scripts/build.sh

# ------------------------------------------------------------------------------
# Build metal rosetta
# ------------------------------------------------------------------------------
FROM golang:1.19.6 AS rosetta

ARG ROSETTA_VERSION

RUN git clone https://github.com/MetalBlockchain/metal-rosetta.git \
  /go/src/github.com/!metal!blockchain/metal-rosetta

WORKDIR /go/src/github.com/!metal!blockchain/metal-rosetta

ENV CGO_ENABLED=1
ENV GOARCH=amd64
ENV GOOS=linux

RUN git checkout $ROSETTA_VERSION && \
    go mod download

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  make build

# ------------------------------------------------------------------------------
# Target container for running the node and rosetta server
# ------------------------------------------------------------------------------
FROM ubuntu:20.04

# Install dependencies
RUN apt-get update -y && \
    apt-get install -y wget

WORKDIR /app

# Install metal daemon
COPY --from=metal \
  /go/src/github.com/!metal!blockchain/metalgo/build/metalgo \
  /app/metalgo

# Install rosetta server
COPY --from=rosetta \
  /go/src/github.com/!metal!blockchain/metal-rosetta/rosetta-server \
  /app/rosetta-server

# Install rosetta runner
COPY --from=rosetta \
  /go/src/github.com/!metal!blockchain/metal-rosetta/rosetta-runner \
  /app/rosetta-runner

# Install service start script
COPY --from=rosetta \
  /go/src/github.com/!metal!blockchain/metal-rosetta/docker/entrypoint.sh \
  /app/entrypoint.sh

EXPOSE 9650
EXPOSE 9651
EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
