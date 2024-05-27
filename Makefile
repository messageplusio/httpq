.PHONY: all
all: build lint

.PHONY: build
build:
	go build -o httpq ./cmd

.PHONY: lint
lint:
	golangci-lint run ./...
