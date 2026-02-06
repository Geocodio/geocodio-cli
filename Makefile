.PHONY: build test lint clean install cover record-cassettes

BINARY := geocodio
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/geocodio/geocodio-cli/internal/cli.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/geocodio

install:
	go install $(LDFLAGS) ./cmd/geocodio

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

record-cassettes:
	@if [ -z "$(GEOCODIO_API_KEY)" ]; then \
		echo "GEOCODIO_API_KEY required"; exit 1; \
	fi
	VCR_MODE=record go test -v ./internal/api/...
