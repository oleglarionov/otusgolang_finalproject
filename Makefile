BIN := "./bin/bannerrotation"

generate:
	go generate ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -race ./internal/...

run:
	docker-compose --file build/local/docker-compose.yaml up --build

build:
	go build -v -o $(BIN) ./cmd

.PHONY: build