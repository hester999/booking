SERVER_IMAGE ?= booking-server:local
WORKER_IMAGE ?= booking-worker:local
BIN_DIR ?= bin
MIGRATION_PATH=./migration

.PHONY: help build build-server build-worker run-server run-worker test vet fmt check docker-build docker-build-server docker-build-worker

help:
	@echo "Available targets:"
	@echo "  build                Build server and worker binaries"
	@echo "  run-server           Run booking HTTP server"
	@echo "  run-worker           Run outbox Kafka worker"
	@echo "  test                 Run tests"
	@echo "  vet                  Run go vet"
	@echo "  fmt                  Format Go files"
	@echo "  check                Run test, vet and build"
	@echo "  docker-build         Build server and worker Docker images"

build: build-server build-worker

build-server:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -o $(BIN_DIR)/server ./cmd/server

build-worker:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -o $(BIN_DIR)/worker ./cmd/worker

run-server:
	go run ./cmd/server

run-worker:
	go run ./cmd/worker

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

check: test vet build

docker-build: docker-build-server docker-build-worker

docker-build-server:
	docker build -f internal/infra/Dockerfile.server -t $(SERVER_IMAGE) .

docker-build-worker:
	docker build -f internal/infra/Dockerfile.worker -t $(WORKER_IMAGE) .
