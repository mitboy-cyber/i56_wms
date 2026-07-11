.PHONY: build test run clean docker

APP_NAME := i56-server
BUILD_DIR := /tmp

build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/
	@echo "✅ Build: $(BUILD_DIR)/$(APP_NAME)"

test:
	go test ./tests/unit/ -v -count=1

test-all:
	go test ./... -v -count=1

run: build
	$(BUILD_DIR)/$(APP_NAME)

dev:
	go run ./cmd/server/

clean:
	rm -f $(BUILD_DIR)/$(APP_NAME)
	go clean

docker-build:
	docker build -f deployments/docker/Dockerfile -t i56-framework:latest .

docker-up:
	docker compose -f deployments/compose/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/compose/docker-compose.yml down

lint:
	go vet ./...

fmt:
	go fmt ./...

all: build test
	@echo "✅ All checks passed"
