GOOSE:=GOOSE_MIGRATION_DIR=./migrations goose

.PHONY: all
all: build lint
.PHONY: build
build:
	go build -o httpq ./cmd

migrate:
	 ${GOOSE} postgres "postgresql://user:password@0.0.0.0/queue" up

.PHONY: lint
lint:
	golangci-lint run ./...
