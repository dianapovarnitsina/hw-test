API_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/scheduler-calendar"
SENDER_BIN := "./bin/sender-calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-api:
	go build -v -o $(API_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run-api: build-api
	$(API_BIN) -config ./configs/calendar_config.toml

build-scheduler:
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

run-scheduler: build-scheduler
	$(SCHEDULER_BIN) -config ./configs/scheduler_config.yaml

build-sender:
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run-sender: build-sender
	$(SENDER_BIN) -config ./configs/sender_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...
	#go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

generate:
	rm -rf internal/server/pb
	mkdir -p internal/server/pb

	protoc \
        --proto_path=api/ \
        --go_out=internal/server/pb \
        --go-grpc_out=internal/server/pb \
        api/*.proto
