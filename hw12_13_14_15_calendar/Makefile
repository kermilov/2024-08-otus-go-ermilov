BIN := "./bin/calendar"
DOCKER_APP_IMG="calendar:develop"
DOCKER_SCHEDULER_IMG="calendar_scheduler:develop"
DOCKER_STORER_IMG="calendar_storer:develop"
DOCKER_INTEGRATION_TESTS_IMG="calendar_integration_tests:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)_scheduler -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN)_storer -ldflags "$(LDFLAGS)" ./cmd/calendar_storer

run: build
	$(BIN) -config ./configs/calendar_config.json
	$(BIN)_scheduler -config ./configs/calendar_scheduler.json
	$(BIN)_storer -config ./configs/calendar_storer.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_APP_IMG) \
		-f build/Dockerfile-calendar . \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SCHEDULER_IMG) \
		-f build/Dockerfile-scheduler . \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_STORER_IMG) \
		-f build/Dockerfile-storer . 
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_INTEGRATION_TESTS_IMG) \
		-f build/Dockerfile-integration-tests . 

run-img: build-img
	docker run $(DOCKER_APP_IMG)
	docker run $(DOCKER_SCHEDULER_IMG)
	docker run $(DOCKER_STORER_IMG)
	docker run $(DOCKER_INTEGRATION_TESTS_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb
	protoc --proto_path=api/ --go_out=internal/server/grpc/pb --go-grpc_out=internal/server/grpc/pb api/*.proto

up:
	docker-compose -f deployments/docker-compose.yaml up -d

down:
	docker-compose -f deployments/docker-compose.yaml down

integration-tests:
	docker-compose -f deployments/docker-compose.yaml up --build --abort-on-container-exit --exit-code-from integration-tests
	docker-compose -f deployments/docker-compose.yaml down	

.PHONY: build run build-img run-img version test lint integration-tests
