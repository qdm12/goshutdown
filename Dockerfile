ARG ALPINE_VERSION=3.13
ARG GO_VERSION=1.16

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk --update add git g++
ENV CGO_ENABLED=0
ARG GOLANGCI_LINT_VERSION=v1.40.1
RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we install g++ to support the race detector
ENV CGO_ENABLED=1
RUN apk --update --no-cache add g++

FROM base AS lint
RUN golangci-lint run --timeout=10m

FROM base AS tidy
RUN git init && \
  git config user.email ci@localhost && \
  git config user.name ci && \
  git add -A && git commit -m ci && \
  sed -i '/\/\/ indirect/d' go.mod && \
  go mod tidy && \
  git diff --exit-code -- go.mod
