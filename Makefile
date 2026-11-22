.PHONY: all build test clean run

VERSION ?= 0.1.0-dev
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "local")
BUILD_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS = -s -w -X 'orbit/internal/buildinfo.Version=$(VERSION)' -X 'orbit/internal/buildinfo.GitCommit=$(GIT_COMMIT)' -X 'orbit/internal/buildinfo.BuildDate=$(BUILD_DATE)'

all: test build

build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/orbit-agent ./cmd/orbit-agent
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/orbit ./cmd/orbit

test:
	go test -v -race ./...

clean:
	rm -rf bin/ /tmp/orbit_*

run: build
	./bin/orbit-agent
